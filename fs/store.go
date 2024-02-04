package fs

import (
	"encoding/binary"
	"io"
	"math"
	"os"

	"github.com/pkg/errors"
)

type FileStore interface {
	Open(file string) error
	Create(file string) error
	Close() error

	IsEOF(err error) bool

	ReadBuffer() ([]byte, error)
	ReadString() (string, error)
	ReadBool() (bool, error)
	ReadFloat32() (float32, error)
	ReadFloat64() (float64, error)
	ReadInt() (int, error)
	ReadInt64() (int64, error)
	ReadInt32() (int32, error)
	ReadInt16() (int16, error)
	ReadInt8() (int8, error)
	ReadUint() (uint, error)
	ReadUint64() (uint64, error)
	ReadUint32() (uint32, error)
	ReadUint16() (uint16, error)
	ReadUint8() (uint8, error)

	WriteBuffer(buf []byte) error
	WriteString(msg string) error
	WriteBool(value bool) error
	WriteFloat32(value float32) error
	WriteFloat64(value float64) error
	WriteInt(value int) error
	WriteInt64(value int64) error
	WriteInt32(value int32) error
	WriteInt16(value int16) error
	WriteInt8(value int8) error
	WriteUint(value uint64) error
	WriteUint64(value uint64) error
	WriteUint32(value uint32) error
	WriteUint16(value uint16) error
	WriteUint8(value uint8) error
}

type mode uint8

const (
	modeNone  mode = 0
	modeRead  mode = 1
	modeWrite mode = 2
)

type fileStore struct {
	f    *os.File
	m    mode
	o    FileStoreOption
	path string
}

type OnReadBuffer func(buf []byte) ([]byte, error)
type OnWriteBuffer func(buf []byte) ([]byte, error)
type OnBeforeOpenFile func(file string) (string, error)
type OnBeforeCreateFile func(file string) (string, error)
type OnAfterCloseFile func(file string) error

type FileStoreOption struct {
	OnReadBuffer       OnReadBuffer
	OnWriteBuffer      OnWriteBuffer
	OnBeforeOpenFile   OnBeforeOpenFile
	OnBeforeCreateFile OnBeforeCreateFile
	OnAfterCloseFile   OnAfterCloseFile
}

type mergeStore struct {
	FileStoreOption
	readQueues         []OnReadBuffer
	writeQueues        []OnWriteBuffer
	beforeOpenQueues   []OnBeforeOpenFile
	beforeCreateQueues []OnBeforeCreateFile
	afterCloseQueues   []OnAfterCloseFile
}

func newMergeStore(opts []FileStoreOption) *mergeStore {
	s := &mergeStore{
		readQueues:         make([]OnReadBuffer, 0),
		writeQueues:        make([]OnWriteBuffer, 0),
		beforeOpenQueues:   make([]OnBeforeOpenFile, 0),
		beforeCreateQueues: make([]OnBeforeCreateFile, 0),
		afterCloseQueues:   make([]OnAfterCloseFile, 0),
	}
	s.OnReadBuffer = func(buf []byte) ([]byte, error) {
		for _, v := range s.readQueues {
			ref, err := v(buf)
			if err != nil {
				return nil, err
			}
			buf = ref
		}
		return buf, nil
	}
	s.OnWriteBuffer = func(buf []byte) ([]byte, error) {
		for _, v := range s.writeQueues {
			ref, err := v(buf)
			if err != nil {
				return nil, err
			}
			buf = ref
		}
		return buf, nil
	}
	s.OnBeforeOpenFile = func(file string) (string, error) {
		for _, v := range s.beforeOpenQueues {
			ref, err := v(file)
			if err != nil {
				return "", err
			}
			file = ref
		}
		return file, nil
	}
	s.OnAfterCloseFile = func(file string) error {
		for _, v := range s.afterCloseQueues {
			err := v(file)
			if err != nil {
				return err
			}
		}
		return nil
	}
	s.OnBeforeCreateFile = func(file string) (string, error) {
		for _, v := range s.beforeCreateQueues {
			ref, err := v(file)
			if err != nil {
				return "", err
			}
			file = ref
		}
		return file, nil
	}
	s.merge(opts)
	return s
}

func (s *mergeStore) merge(opts []FileStoreOption) {
	for _, opt := range opts {
		if opt.OnReadBuffer != nil {
			s.readQueues = append(s.readQueues, opt.OnReadBuffer)
		}
		if opt.OnWriteBuffer != nil {
			s.writeQueues = append(s.writeQueues, opt.OnWriteBuffer)
		}
		if opt.OnBeforeOpenFile != nil {
			s.beforeOpenQueues = append(s.beforeOpenQueues, opt.OnBeforeOpenFile)
		}
		if opt.OnAfterCloseFile != nil {
			s.afterCloseQueues = append(s.afterCloseQueues, opt.OnAfterCloseFile)
		}
		if opt.OnBeforeCreateFile != nil {
			s.beforeCreateQueues = append(s.beforeCreateQueues, opt.OnBeforeCreateFile)
		}
	}
}

func NewFileStore(options ...FileStoreOption) FileStore {
	return &fileStore{
		f: nil,
		m: modeNone,
		o: newMergeStore(options).FileStoreOption,
	}
}

func (p *fileStore) Open(file string) error {
	p.path = file
	file, err := p.o.OnBeforeOpenFile(file)
	if err != nil {
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	p.f = f
	p.m = modeRead
	return nil
}

func (p *fileStore) Create(file string) error {
	p.path = file
	file, err := p.o.OnBeforeCreateFile(file)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	p.f = f
	p.m = modeWrite
	return nil
}

func (p *fileStore) Close() error {
	err := p.f.Close()
	p.f = nil
	p.m = modeNone
	path := p.path
	p.path = ""
	if p.o.OnAfterCloseFile != nil {
		return p.o.OnAfterCloseFile(path)
	}
	return err
}

func (p *fileStore) IsEOF(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF)
}

func (p *fileStore) ReadBuffer() ([]byte, error) {
	var value int32
	err := binary.Read(p.f, binary.LittleEndian, &value)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, value)
	_, err = p.f.Read(buf)
	if err != nil {
		return nil, err
	}
	if p.o.OnReadBuffer != nil {
		return p.o.OnReadBuffer(buf)
	}
	return buf, nil
}

func (p *fileStore) ReadString() (string, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// read bool
func (p *fileStore) ReadBool() (bool, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return false, err
	}
	if len(buf) == 0 {
		return false, nil
	}
	return buf[0] == 1, nil
}

// ReadFloat32
func (p *fileStore) ReadFloat32() (float32, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return float32(binary.LittleEndian.Uint32(buf)), nil
}

// ReadFloat64
func (p *fileStore) ReadFloat64() (float64, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return float64(binary.LittleEndian.Uint64(buf)), nil
}

// ReadInt
func (p *fileStore) ReadInt() (int, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint64(buf)), nil
}

// ReadInt64
func (p *fileStore) ReadInt64() (int64, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return int64(binary.LittleEndian.Uint64(buf)), nil
}

// ReadInt32
func (p *fileStore) ReadInt32() (int32, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return int32(binary.LittleEndian.Uint32(buf)), nil
}

// ReadInt16
func (p *fileStore) ReadInt16() (int16, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return int16(binary.LittleEndian.Uint16(buf)), nil
}

// ReadInt8
func (p *fileStore) ReadInt8() (int8, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return int8(buf[0]), nil
}

// ReadUint
func (p *fileStore) ReadUint() (uint, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return uint(binary.LittleEndian.Uint64(buf)), nil
}

// ReadUint64
func (p *fileStore) ReadUint64() (uint64, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// ReadUint32
func (p *fileStore) ReadUint32() (uint32, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// ReadUint16
func (p *fileStore) ReadUint16() (uint16, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// ReadUint8
func (p *fileStore) ReadUint8() (uint8, error) {
	buf, err := p.ReadBuffer()
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

// WriteInt
func (p *fileStore) WriteInt(value int) error {
	// 写入value
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return p.WriteBuffer(buf)
}

// WriteInt64
func (p *fileStore) WriteInt64(value int64) error {
	// 写入value
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return p.WriteBuffer(buf)
}

// WriteInt32
func (p *fileStore) WriteInt32(value int32) error {
	// 写入value
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))
	return p.WriteBuffer(buf)
}

// WriteInt16
func (p *fileStore) WriteInt16(value int16) error {
	// 写入value
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(value))
	return p.WriteBuffer(buf)
}

// WriteInt8
func (p *fileStore) WriteInt8(value int8) error {
	// 写入value
	return p.WriteBuffer([]byte{uint8(value)})
}

// WriteUint
func (p *fileStore) WriteUint(value uint64) error {
	// 写入value
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return p.WriteBuffer(buf)
}

// WriteUint64
func (p *fileStore) WriteUint64(value uint64) error {
	// 写入value
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, value)
	return p.WriteBuffer(buf)
}

// WriteUint32
func (p *fileStore) WriteUint32(value uint32) error {
	// 写入value
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, value)
	return p.WriteBuffer(buf)
}

// WriteUint16
func (p *fileStore) WriteUint16(value uint16) error {
	// 写入value
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, value)
	return p.WriteBuffer(buf)
}

// WriteUint8
func (p *fileStore) WriteUint8(value uint8) error {
	// 写入value
	return p.WriteBuffer([]byte{value})
}

// write bool
func (p *fileStore) WriteBool(value bool) error {
	// 写入value
	if value {
		return p.WriteBuffer([]byte{1})
	}
	return p.WriteBuffer([]byte{0})
}

// WriteFloat32
func (p *fileStore) WriteFloat32(value float32) error {
	// 写入value
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(math.Float32bits(value)))
	return p.WriteBuffer(buf)
}

// WriteFloat64
func (p *fileStore) WriteFloat64(value float64) error {
	// 写入value
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(value))
	return p.WriteBuffer(buf)
}

func (p *fileStore) WriteBuffer(buf []byte) error {
	if p.o.OnWriteBuffer != nil {
		ref, err := p.o.OnWriteBuffer(buf)
		if err != nil {
			return err
		}
		buf = ref
	}
	err := binary.Write(p.f, binary.LittleEndian, int32(len(buf)))
	if err != nil {
		return err
	}
	_, err = p.f.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (p *fileStore) WriteString(msg string) error {
	return p.WriteBuffer([]byte(msg))
}
