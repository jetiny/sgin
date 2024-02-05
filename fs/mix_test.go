package fs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMix(t *testing.T) {
	pub, pri, err := RsaKeySmall.NewFile()
	assert.NoError(t, err)
	fmt.Println(pub)
	fmt.Println(pri)
	enc, err := MixTypeDefault.NewEncoderFile(pub)
	assert.NoError(t, err)
	dec, err := MixTypeDefault.NewDecoderFile(pri)
	assert.NoError(t, err)
	{
		files := []string{"mix.go", "mix_test.go"}
		for _, f := range files {
			dist, err := enc.Encode(f)
			assert.NoError(t, err)
			fmt.Println(dist)
			dist, err = dec.Decode(dist)
			assert.NoError(t, err)
			fmt.Println(dist)
		}
	}
}

func TestMixAll(t *testing.T) {
	pub, pri, err := RsaKeySmall.NewFile()
	assert.NoError(t, err)
	fmt.Println(pub)
	fmt.Println(pri)
	types := []MixType{
		MixTypeDefault,
		MixTypeCfb,
		MixTypeOfb,
		MixTypeCtr,
		MixTypeGcm,
		MixTypeCbc,
	}
	for _, v := range types {
		enc, err := v.NewEncoderFile(pub)
		assert.NoError(t, err)
		dec, err := v.NewDecoderFile(pri)
		assert.NoError(t, err)
		{
			files := []string{"mix.go", "mix_test.go"}
			for _, f := range files {
				dist, err := enc.Encode(f)
				assert.NoError(t, err)
				fmt.Println(dist)
				dist, err = dec.Decode(dist)
				assert.NoError(t, err)
				fmt.Println(dist)
			}
		}
	}
}
