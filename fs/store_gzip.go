package fs

import (
	"os"
	"path/filepath"

	"github.com/jetiny/sgin/utils"
)

type gzipStoreOption struct {
	FileStoreOption
	cpath string
	opath string
}

type GzipStoreConfig struct {
	Encoder Encoder
	Decoder Decoder
}

func NewGzipOption(option *GzipStoreConfig) FileStoreOption {
	var enc Encoder = nil
	var dec Decoder = nil
	if option != nil {
		enc = option.Encoder
		dec = option.Decoder
	}
	r := &gzipStoreOption{}
	r.FileStoreOption.OnAfterCloseFile = func(file string) error {
		if r.opath != "" {
			opath := r.opath
			r.opath = ""
			return os.Remove(opath)
		} else if r.cpath != "" {
			cpath := r.cpath
			r.cpath = ""
			defer os.Remove(cpath)
			err := CompressFile(cpath, file)
			if err != nil {
				return err
			}
			if enc != nil {
				newFile, e := enc.Encode(file)
				if e != nil {
					return e
				}
				err = os.Rename(newFile, file)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	r.FileStoreOption.OnBeforeCreateFile = func(file string) (string, error) {
		newFile := filepath.Join(os.TempDir(), utils.Uuid())
		r.cpath = newFile
		return newFile, nil
	}
	r.FileStoreOption.OnBeforeOpenFile = func(file string) (string, error) {
		if dec != nil {
			decFile, e := dec.Decode(file)
			if e != nil {
				return "", e
			}
			file = decFile
		}
		newFile := filepath.Join(os.TempDir(), utils.Uuid())
		err := DecompressFile(file, newFile)
		if err != nil {
			return "", err
		}
		r.opath = newFile
		return newFile, nil
	}
	return r.FileStoreOption
}
