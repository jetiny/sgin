package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	pub, pri, err := RsaKeyDefault.Create()
	assert.NoError(t, err)
	s := NewFileStore(NewGzipOption(&GzipStoreConfig{
		Encoder: MixTypeDefault.NewEncoder(pub),
		Decoder: MixTypeDefault.NewDecoder(pri),
	}), NewContentMixOption("ss"))
	filename := filepath.Join(os.TempDir(), "test.txt")
	{
		err := s.Create(filename)
		assert.NoError(t, err)
		err = s.WriteBool(true)
		assert.NoError(t, err)
		err = s.WriteString("")
		assert.NoError(t, err)
		s.Close()
		info, err := os.Stat(filename)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	}
	{
		err := s.Open(filename)
		assert.NoError(t, err)
		b, err := s.ReadBool()
		assert.NoError(t, err)
		str, err := s.ReadString()
		assert.NoError(t, err)
		assert.Equal(t, "", str)
		assert.Equal(t, true, b)
		s.Close()
	}
}

func TestStoreTypes(t *testing.T) {
	s := NewFileStore(NewGzipOption(nil), NewContentMixOption("ss"))
	filename := filepath.Join(os.TempDir(), "test.txt")
	{
		err := s.Create(filename)
		assert.NoError(t, err)
		err = s.WriteInt(123)
		assert.NoError(t, err)
		err = s.WriteInt8(123)
		assert.NoError(t, err)
		err = s.WriteInt16(123)
		assert.NoError(t, err)
		err = s.WriteInt32(123)
		assert.NoError(t, err)
		err = s.WriteInt64(123)
		assert.NoError(t, err)
		err = s.WriteUint(123)
		assert.NoError(t, err)
		err = s.WriteUint8(123)
		assert.NoError(t, err)
		err = s.WriteUint16(123)
		assert.NoError(t, err)
		err = s.WriteUint32(123)
		assert.NoError(t, err)
		err = s.WriteUint64(123)
		assert.NoError(t, err)
		err = s.WriteFloat32(123.456)
		assert.NoError(t, err)
		err = s.WriteFloat64(123.456)
		assert.NoError(t, err)
		err = s.WriteString("hello, world")
		assert.NoError(t, err)
		s.Close()
		info, err := os.Stat(filename)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	}
	{
		err := s.Open(filename)
		assert.NoError(t, err)
		i, err := s.ReadInt()
		assert.NoError(t, err)
		assert.Equal(t, 123, i)
		i8, err := s.ReadInt8()
		assert.NoError(t, err)
		assert.Equal(t, int8(123), i8)
		i16, err := s.ReadInt16()
		assert.NoError(t, err)
		assert.Equal(t, int16(123), i16)
		i32, err := s.ReadInt32()
		assert.NoError(t, err)
		assert.Equal(t, int32(123), i32)
		i64, err := s.ReadInt64()
		assert.NoError(t, err)
		assert.Equal(t, int64(123), i64)
		u, err := s.ReadUint()
		assert.NoError(t, err)
		assert.Equal(t, uint(123), u)
		u8, err := s.ReadUint8()
		assert.NoError(t, err)
		assert.Equal(t, uint8(123), u8)
		u16, err := s.ReadUint16()
		assert.NoError(t, err)
		assert.Equal(t, uint16(123), u16)
		u32, err := s.ReadUint32()
		assert.NoError(t, err)
		assert.Equal(t, uint32(123), u32)
		u64, err := s.ReadUint64()
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), u64)
		f32, err := s.ReadFloat32()
		assert.NoError(t, err)
		assert.NotEmpty(t, f32)
		f64, err := s.ReadFloat64()
		assert.NoError(t, err)
		assert.NotEmpty(t, f64)
		str, err := s.ReadString()
		assert.NoError(t, err)
		assert.Equal(t, "hello, world", str)
		s.Close()
	}
}
