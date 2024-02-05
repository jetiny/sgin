package fs

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"

	"github.com/jetiny/sgin/utils"
)

const (
	bufLen     = 32 * 1024
	aesKeySize = 32 // AES-256
)

type (
	EncCipher func() ([]byte, any, error)
	DecCipher func([]byte) (any, error)
)

// encodeStream
//
//	@Description:  encrypted data
//	@param r       data source read stream
//	@param w       encrypted data is written to the stream
//	@param encMode public key or encryption function
//	@param h       specify the hash verification method
//	@param gen     EncCipher, generate key and specify encryption algorithm
//	@return error
func encodeStream(r io.Reader, w io.Writer, encMode any, h hash.Hash, gen EncCipher) error {
	tmp := make([]byte, bufLen)

	// ensure that the data and cipher algorithm objects are legal
	data, cp, err := gen()
	if err != nil {
		return err
	}

	var encKey []byte
	switch mode := encMode.(type) {
	case []byte:
		encKey, err = rsaEncrypt(mode, data)
	case string:
		encKey, err = rsaEncrypt([]byte(mode), data)
	case func([]byte) ([]byte, error):
		encKey, err = mode(data)
	default:
		return errors.New("encMode not support")
	}
	if err != nil {
		return err
	}

	n := len(encKey)
	tmp[0], tmp[1] = byte(n), byte(n>>8)
	head := append(tmp[:2], encKey...)

	_, err = w.Write(head)
	if err != nil {
		return err
	}

	h.Write(head)

	var cw io.WriteCloser

	switch cc := cp.(type) {
	case cipher.AEAD:
		ac := &aeadCipher{
			h: h, w: w,
			aead:  cc,
			oSize: cc.Overhead(),
			tmp:   make([]byte, bufLen),
			buf:   bytes.NewBuffer(make([]byte, 0, bufLen)),
		}
		if err = ac.readNonce(data); err != nil {
			return err
		}

		cw = ac
		// limit the maximum length of plaintext
		tmp = tmp[:bufLen-cc.Overhead()]
	case cipher.Stream:
		cw = &streamCipher{
			h: h, w: w,
			stream: cc,
		}
	case cipher.BlockMode:
		cw = &blockCipher{
			h: h, w: w,
			block: cc,
			bSize: cc.BlockSize(),
			tmp:   make([]byte, bufLen),
			buf:   bytes.NewBuffer(make([]byte, 0, bufLen)),
		}
	default:
		return errors.New("gen return parameter error")
	}

	_, err = io.CopyBuffer(cw, &onlyRW{r: r}, tmp)
	if err != nil {
		return err
	}

	if err = cw.Close(); err != nil {
		return err
	}

	_, err = w.Write(h.Sum(tmp[:0]))

	return err
}

// decodeStream
//
//	@Description:  decrypt data
//	@param r       ciphertext data read stream
//	@param w       the decrypted data is written to the stream
//	@param devMode private key data or decryption function
//	@param h       specify the hash verification method
//	@param gen     DecCipher, verify and generate a decryption algorithm based on the key
//	@return error
func decodeStream(r io.Reader, w io.Writer, decMode any, h hash.Hash, gen DecCipher) error {
	tmp := make([]byte, bufLen)

	_, err := io.ReadFull(r, tmp[:2])
	if err != nil {
		return err
	}

	h.Write(tmp[:2])

	n := int(tmp[0]) | int(tmp[1])<<8
	if n > bufLen {
		return errors.New("len(rsa) out of index")
	}

	_, err = io.ReadFull(r, tmp[:n])
	if err != nil {
		return err
	}

	h.Write(tmp[:n])

	var data []byte
	switch mode := decMode.(type) {
	case []byte:
		data, err = rsaDecrypt(mode, tmp[:n])
	case string:
		data, err = rsaDecrypt([]byte(mode), tmp[:n])
	case func([]byte) ([]byte, error):
		data, err = mode(tmp[:n])
	default:
		return errors.New("decMode not support")
	}
	if err != nil {
		return err
	}

	// it is necessary to verify the legality of data and generate cipher
	cp, err := gen(data)
	if err != nil {
		return err
	}

	var (
		cr  io.Reader
		sum = make([]byte, h.Size())
	)

	switch cc := cp.(type) {
	case cipher.AEAD:
		ac := &aeadCipher{
			h: h, r: r,
			aead:  cc,
			sum:   sum,
			oSize: cc.Overhead(),
			tmp:   make([]byte, bufLen),
			buf:   bytes.NewBuffer(make([]byte, 0, bufLen)),
		}
		if err = ac.readNonce(data); err != nil {
			return err
		}

		cr = ac
		// limit the maximum length of plaintext
		tmp = tmp[:bufLen-cc.Overhead()]
	case cipher.Stream:
		cr = &streamCipher{
			h: h, r: bufio.NewReader(r),
			stream: cc,
			sum:    sum,
		}
	case cipher.BlockMode:
		cr = &blockCipher{
			h: h, r: r,
			block: cc,
			bSize: cc.BlockSize(),
			sum:   sum,
			tmp:   make([]byte, bufLen),
			buf:   bytes.NewBuffer(make([]byte, 0, bufLen)),
		}
	default:
		return errors.New("gen return parameter error")
	}

	_, err = io.CopyBuffer(&onlyRW{w: w}, cr, tmp)
	if err != nil {
		return err
	}

	if bytes.Equal(h.Sum(tmp[:0]), sum) {
		return nil
	}

	return errors.New("h.Sum not match")
}

// -----------------------------------------------------------------------------

func genEncCipher(enc any) EncCipher {
	return func() ([]byte, any, error) {
		key := make([]byte, aesKeySize)

		_, err := rand.Read(key)
		if err != nil {
			return nil, nil, err
		}

		for i, v := range key {
			if v == 0 {
				key[i] = byte(i%0xff) + 1 // ensure that the key is not 0
			}
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, nil, err
		}

		var (
			gen any
			add = func(s int) []byte {
				info := make([]byte, s)
				_, _ = rand.Read(info)

				// generate data [key + 0 + info]
				key = append(key, 0)
				key = append(key, info...)
				return info
			}
		)

		switch e := enc.(type) {
		case func(cipher.Block) (cipher.AEAD, error):
			ad, err := e(block)
			if err != nil {
				return nil, nil, err
			}

			gen = ad
			add(ad.NonceSize()) // nonce
		case func(cipher.Block, []byte) cipher.Stream:
			gen = e(block, add(block.BlockSize())) // iv
		case func(cipher.Block, []byte) cipher.BlockMode:
			gen = e(block, add(block.BlockSize())) // iv
		default:
			return nil, nil, errors.New("enc func type error")
		}

		return key, gen, nil
	}
}

func genDecCipher(dec any) DecCipher {
	return func(data []byte) (any, error) {
		// verify data validity, data like this [key + 0 + iv/nonce]
		key, info, ok := bytes.Cut(data, []byte{0})
		if !ok || len(key) != aesKeySize {
			return nil, errors.New("key error")
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		var gen any

		switch d := dec.(type) {
		case func(cipher.Block) (cipher.AEAD, error):
			ad, err := d(block)
			if err != nil {
				return nil, err
			}

			gen = ad // subsequent verification of nonce by read nonce
		case func(cipher.Block, []byte) cipher.Stream:
			if len(info) != block.BlockSize() {
				return nil, errors.New("len(iv) error")
			}

			gen = d(block, info)
		case func(cipher.Block, []byte) cipher.BlockMode:
			if len(info) != block.BlockSize() {
				return nil, errors.New("len(iv) error")
			}

			gen = d(block, info)
		default:
			return nil, errors.New("dec func type error")
		}

		return gen, nil
	}
}

// -----------------------------------------------------------------------------

type streamCipher struct {
	w      io.Writer
	h      hash.Hash
	stream cipher.Stream
	r      *bufio.Reader
	sum    []byte
}

func (sc *streamCipher) Write(p []byte) (n int, err error) {
	sc.stream.XORKeyStream(p, p)

	n, err = writeFull(sc.w, p)
	if err == nil {
		sc.h.Write(p[:n])
	}

	return
}

func (sc *streamCipher) Read(p []byte) (n int, err error) {
	n, err = sc.r.Read(p)
	if err == nil {
		var this []byte

		this, err = sc.r.Peek(len(sc.sum))
		switch err {
		case nil:
		case io.EOF:
			// last Peek data: len(last) = len(sc.sum)
			// this Peek data: len(this) < len(sc.sum)
			// there is no len(this) = len(sc.sum), in this case it will Peek again
			// |     p[:n]    | this | the last data obtained
			// |  xxx  | last | this | len(p) > len(sc.sum), so last is in p[:n]
			// |  xxx + yyy |   sum  | extract the sum from this data
			lp := n - len(sc.sum) + len(this)
			copy(sc.sum, append(p[lp:n], this...))
			n = lp
		default:
			return
		}

		if n > 0 {
			sc.h.Write(p[:n])
			sc.stream.XORKeyStream(p, p[:n])
		}
	}

	return
}

func (sc *streamCipher) Close() error {
	return nil
}

// -----------------------------------------------------------------------------

type blockCipher struct {
	r     io.Reader
	w     io.Writer
	h     hash.Hash
	block cipher.BlockMode
	buf   *bytes.Buffer
	sum   []byte
	tmp   []byte
	bSize int
}

func (bc *blockCipher) Write(p []byte) (n int, err error) {
	n, err = bc.buf.Write(p)

	if bn := bc.buf.Len() / bc.bSize; bn >= 1 {
		// encrypt data in blocks of (bn * BlockSize)
		// guaranteed to encrypt at least (1 * BlockSize) block data
		bn *= bc.bSize

		// bc.tmp[0] = 0, ordinary data, data length stored in encrypted block
		bc.tmp[0], bc.tmp[1], bc.tmp[2] = 0, byte(bn), byte(bn>>8)

		bc.block.CryptBlocks(bc.tmp, bc.tmp[:bc.bSize])

		_, err = writeFull(bc.w, bc.tmp[:bc.bSize])
		if err != nil {
			return
		}

		bc.h.Write(bc.tmp[:bc.bSize])

		_, err = io.ReadFull(bc.buf, bc.tmp[:bn])
		if err != nil {
			return
		}

		bc.block.CryptBlocks(bc.tmp, bc.tmp[:bn])

		_, err = writeFull(bc.w, bc.tmp[:bn])
		if err != nil {
			return
		}

		bc.h.Write(bc.tmp[:bn])
	}

	return
}

func (bc *blockCipher) Read(p []byte) (n int, err error) {
	if bc.buf.Len() == 0 {
		// read (1 * BlockSize) block data
		_, err = io.ReadFull(bc.r, bc.tmp[:bc.bSize])
		if err != nil {
			return
		}

		bc.h.Write(bc.tmp[:bc.bSize])
		bc.block.CryptBlocks(bc.tmp, bc.tmp[:bc.bSize])

		if n = int(bc.tmp[0]); n > 0 {
			if n == 1 {
				n = 0 // no data left to read
			} else {
				n = copy(p, bc.tmp[1:n]) // read remaining data
			}

			_, err = io.ReadFull(bc.r, bc.sum)
			if err == nil {
				err = io.EOF // after reading the crc, return the data to read
			}

			return
		}

		n = int(bc.tmp[1]) | int(bc.tmp[2])<<8

		_, err = io.ReadFull(bc.r, bc.tmp[:n])
		if err != nil {
			return
		}

		bc.h.Write(bc.tmp[:n])
		bc.block.CryptBlocks(bc.tmp, bc.tmp[:n])
		bc.buf.Write(bc.tmp[:n])
	}

	return bc.buf.Read(p)
}

func (bc *blockCipher) Close() error {
	if bc.buf.Len() > 0 {
		// remaining data up to (BlockSize - 1) byte
		// storage length of bc.tmp[0] is at most just to fill the BlockSize
		bc.tmp[0] = byte(bc.buf.Len() + 1)
		copy(bc.tmp[1:], bc.buf.Bytes())
	} else {
		bc.tmp[0] = 1 // mark no data left
	}

	bc.block.CryptBlocks(bc.tmp, bc.tmp[:bc.bSize])

	_, err := writeFull(bc.w, bc.tmp[:bc.bSize])
	if err != nil {
		return err
	}

	bc.h.Write(bc.tmp[:bc.bSize])

	return nil
}

// -----------------------------------------------------------------------------

//goland:noinspection SpellCheckingInspection
type aeadCipher struct {
	r     io.Reader
	w     io.Writer
	h     hash.Hash
	aead  cipher.AEAD
	buf   *bytes.Buffer
	tmp   []byte
	sum   []byte
	nonce []byte
	oSize int
	eof   bool
}

func (ac *aeadCipher) readNonce(data []byte) error {
	// data must be [key + 0 + nonce], key content does not contain 0
	if i := bytes.IndexByte(data, 0); i > 0 {
		ac.nonce = data[i+1:]
	}

	if len(ac.nonce) != ac.aead.NonceSize() {
		return errors.New("can not read nonce")
	}

	return nil
}

func (ac *aeadCipher) Write(p []byte) (n int, err error) {
	n, err = ac.buf.Write(p)

	if bn := ac.buf.Len() / ac.oSize; bn >= 1 {
		bn *= ac.oSize // logic similar to blockCipher algorithm

		ac.tmp[0], ac.tmp[1], ac.tmp[2] = 0, byte(bn), byte(bn>>8)

		tmp := ac.aead.Seal(ac.tmp[:0], ac.nonce, ac.tmp[:3], nil)

		_, err = writeFull(ac.w, tmp)
		if err != nil {
			return
		}

		increment(ac.nonce)
		ac.h.Write(tmp)

		_, err = io.ReadFull(ac.buf, ac.tmp[:bn])
		if err != nil {
			return
		}

		tmp = ac.aead.Seal(ac.tmp[:0], ac.nonce, ac.tmp[:bn], nil)

		_, err = writeFull(ac.w, tmp)
		if err != nil {
			return
		}

		increment(ac.nonce)
		ac.h.Write(tmp)
	}

	return
}

func (ac *aeadCipher) Read(p []byte) (n int, err error) {
	if ac.buf.Len() == 0 {
		if ac.eof {
			return 0, io.EOF
		}

		n = 3 + ac.oSize // read 3 bytes of length information

		_, err = io.ReadFull(ac.r, ac.tmp[:n])
		if err != nil {
			return
		}

		ac.h.Write(ac.tmp[:n])

		var tmp []byte

		tmp, err = ac.aead.Open(ac.tmp[:0], ac.nonce, ac.tmp[:n], nil)
		if err != nil {
			return
		}

		increment(ac.nonce)

		eof := tmp[0] == 1

		if n = int(tmp[1]) | int(tmp[2])<<8; n > 0 {
			n += ac.oSize

			_, err = io.ReadFull(ac.r, ac.tmp[:n])
			if err != nil {
				return
			}

			ac.h.Write(ac.tmp[:n])

			tmp, err = ac.aead.Open(ac.tmp[:0], ac.nonce, ac.tmp[:n], nil)
			if err != nil {
				return
			}

			increment(ac.nonce)
			ac.buf.Write(tmp)
		}

		if eof {
			_, err = io.ReadFull(ac.r, ac.sum)
			if err != nil {
				return
			}

			ac.eof = true // last piece of data
		}
	}

	return ac.buf.Read(p)
}

func (ac *aeadCipher) Close() error {
	n := ac.buf.Len() // encrypt the last data length information
	ac.tmp[0], ac.tmp[1], ac.tmp[2] = 1, byte(n), byte(n>>8)

	tmp := ac.aead.Seal(ac.tmp[:0], ac.nonce, ac.tmp[:3], nil)

	_, err := writeFull(ac.w, tmp)
	if err != nil {
		return err
	}

	increment(ac.nonce)
	ac.h.Write(tmp)

	if n > 0 { // encrypt the last remaining data
		tmp = ac.aead.Seal(ac.tmp[:0], ac.nonce, ac.buf.Bytes(), nil)

		_, err = writeFull(ac.w, tmp)
		if err != nil {
			return err
		}

		increment(ac.nonce)
		ac.h.Write(tmp)
	}

	return nil
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}

// -----------------------------------------------------------------------------
type MixType string

const (
	MixTypeDefault MixType = ""
	MixTypeCfb     MixType = "CFB"
	MixTypeOfb     MixType = "OFB"
	MixTypeCtr     MixType = "CTR"
	MixTypeGcm     MixType = "GCM"
	MixTypeCbc     MixType = "CBC"
)

func (s MixType) NewEncoder(pub []byte) Encoder {
	return &innerEncoder{
		enc: s.getEncoder(),
		h:   md5.New(),
		pub: pub,
	}
}

func (s MixType) NewEncoderFile(pubFile string) (Encoder, error) {
	pub, err := os.ReadFile(pubFile)
	if err != nil {
		return nil, err
	}
	return s.NewEncoder(pub), nil
}

func (s MixType) NewDecoder(pri []byte) Decoder {
	return &innerDecoder{
		dec: s.getDecoder(),
		h:   md5.New(),
		pri: pri,
	}
}

func (s MixType) NewDecoderFile(priFile string) (Decoder, error) {
	pri, err := os.ReadFile(priFile)
	if err != nil {
		return nil, err
	}
	return s.NewDecoder(pri), nil
}

type Encoder interface {
	EncodeStream(r io.Reader, w io.Writer) error
	EncodeBytes(data []byte) ([]byte, error)
	EncodeFile(src, dst string) error
	Encode(src string) (string, error)
}

type Decoder interface {
	DecodeStream(r io.Reader, w io.Writer) error
	DecodeBytes(data []byte) ([]byte, error)
	DecodeFile(src, dst string) error
	Decode(src string) (string, error)
}

type innerDecoder struct {
	dec DecCipher
	h   hash.Hash
	pri []byte
}

func (s *innerDecoder) DecodeStream(r io.Reader, w io.Writer) error {
	return decodeStream(r, w, s.pri, s.h, s.dec)
}

func (s *innerDecoder) DecodeBytes(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	w := bytes.NewBuffer([]byte{})
	err := decodeStream(r, w, s.pri, s.h, s.dec)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *innerDecoder) DecodeFile(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	g, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer g.Close()
	return decodeStream(f, g, s.pri, s.h, s.dec)
}

func (s *innerDecoder) Decode(src string) (string, error) {
	dstFile := filepath.Join(os.TempDir(), utils.Uuid())
	err := s.DecodeFile(src, dstFile)
	return dstFile, err
}

type innerEncoder struct {
	enc EncCipher
	h   hash.Hash
	pub []byte
}

func (s *innerEncoder) EncodeStream(r io.Reader, w io.Writer) error {
	return encodeStream(r, w, s.pub, s.h, s.enc)
}

func (s *innerEncoder) EncodeBytes(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	w := bytes.NewBuffer([]byte{})
	err := encodeStream(r, w, s.pub, s.h, s.enc)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (s *innerEncoder) EncodeFile(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()
	g, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer g.Close()
	return encodeStream(f, g, s.pub, s.h, s.enc)
}

func (s *innerEncoder) Encode(src string) (string, error) {
	dstFile := filepath.Join(os.TempDir(), utils.Uuid())
	err := s.EncodeFile(src, dstFile)
	return dstFile, err
}

func (s MixType) getEncoder() EncCipher {
	switch s {
	case MixTypeCtr:
		return genEncCipher(cipher.NewCTR)
	case MixTypeOfb:
		return genEncCipher(cipher.NewOFB)
	case MixTypeGcm:
		return genEncCipher(cipher.NewGCM)
	case MixTypeCbc:
		return genEncCipher(cipher.NewCBCEncrypter)
	default:
		return genEncCipher(cipher.NewCFBEncrypter)
	}
}

func (s MixType) getDecoder() DecCipher {
	switch s {
	case MixTypeCtr:
		return genDecCipher(cipher.NewCTR)
	case MixTypeOfb:
		return genDecCipher(cipher.NewOFB)
	case MixTypeGcm:
		return genDecCipher(cipher.NewGCM)
	case MixTypeCbc:
		return genDecCipher(cipher.NewCBCDecrypter)
	default:
		return genDecCipher(cipher.NewCFBDecrypter)
	}
}

type RsaKey int

const (
	RsaKeyDefault RsaKey = 4096
	RsaKeySmall   RsaKey = 2048
	RsaKeyMiddle  RsaKey = 3072
	RsaKeyLarge   RsaKey = 4096
)

func (s RsaKey) Create() (pub, pri []byte, err error) {
	return genRsaKeys(int(s))
}

func (s RsaKey) CreateFile(pub, pri string) error {
	return genRsaKeyFiles(int(s), pub, pri)
}

func (s RsaKey) NewFile() (string, string, error) {
	return genRsaTmpFiles(int(s))
}

func genRsaKeys(bits int) (pub, pri []byte, err error) {
	pubW := bytes.NewBuffer([]byte{})
	priW := bytes.NewBuffer([]byte{})
	err = genRsaKey(bits, pubW, priW)
	if err != nil {
		return
	}
	pub = pubW.Bytes()
	pri = priW.Bytes()
	return
}

func genRsaKeyFiles(bits int, pubFile, priFile string) error {
	pub, pri, err := genRsaKeys(bits)
	if err != nil {
		return err
	}
	err = os.WriteFile(pubFile, pub, 0600)
	if err != nil {
		return err
	}
	err = os.WriteFile(priFile, pri, 0600)
	if err != nil {
		return err
	}
	return nil
}

func genRsaTmpFiles(bits int) (pubFile, priFile string, err error) {
	keyName := utils.Uuid()
	pubFile = filepath.Join(os.TempDir(), fmt.Sprintf("rsa_pri_%s.pem", keyName))
	priFile = filepath.Join(os.TempDir(), fmt.Sprintf("rsa_pub_%s.pem", keyName))
	err = genRsaKeyFiles(bits, pubFile, priFile)
	return
}

// genRsaKey
//
//	@Description: generate rsa public-private key pair
//	@param bits   generated digits
//	@param pub    public key write stream
//	@param pri    private key write stream
//	@return error
func genRsaKey(bits int, pub, pri io.Writer) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	err = pem.Encode(pri, block)
	if err != nil {
		return err
	}

	derPkix, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	return pem.Encode(pub, block)
}

// rsaEncrypt
//
//	@Description:   rsa encryption logic
//	@param pubKey   public key data
//	@param origData data to be encrypted
//	@return []byte  return encrypted data
//	@return error
func rsaEncrypt(pubKey, origData []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("block error")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pk, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pk, origData)
}

// rsaDecrypt
//
//	@Description:     rsa decryption logic
//	@param priKey     private key data
//	@param cipherText cipher text
//	@return []byte    decrypted data
//	@return error
func rsaDecrypt(priKey, cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	prIv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, prIv, cipherText)
}

func writeFull(w io.Writer, b []byte) (n int, err error) {
	var now int

	for n < len(b) {
		now, err = w.Write(b[n:])
		if err != nil {
			return
		}

		n += now
	}

	return
}

// eliminate the judgment of io.WriterTo and io.ReaderFrom in io.CopyBuffer
// make sure to only copy data using io.Reader and io.Writer
type onlyRW struct {
	r io.Reader
	w io.Writer
}

func (rw *onlyRW) Write(b []byte) (int, error) {
	return rw.w.Write(b)
}

func (rw *onlyRW) Read(b []byte) (int, error) {
	return rw.r.Read(b)
}
