package fs

import "crypto/rc4"

func NewContentMixOption(key string) FileStoreOption {
	res := FileStoreOption{}
	if key != "" {
		process := func(encodeBytes []byte) ([]byte, error) {
			chiper, err := rc4.NewCipher(([]byte)(key))
			if err != nil {
				return nil, err
			}
			dst := encodeBytes
			chiper.XORKeyStream(dst, encodeBytes)
			return dst, nil
		}
		res.OnReadBuffer = process
		res.OnWriteBuffer = process
	}
	return res
}
