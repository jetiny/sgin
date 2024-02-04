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

func NewGzipOption() FileStoreOption {
	r := gzipStoreOption{}
	r.FileStoreOption.OnAfterCloseFile = func(file string) error {
		if r.opath != "" {
			r.opath = ""
			return os.Remove(r.opath)
		} else if r.cpath != "" {
			cpath := r.cpath
			r.cpath = ""
			defer os.Remove(cpath)
			err := CompressFile(cpath, file)
			return err
		}
		return nil
	}
	r.FileStoreOption.OnBeforeCreateFile = func(file string) (string, error) {
		newFile := filepath.Join(os.TempDir(), utils.Uuid())
		r.cpath = newFile
		return newFile, nil
	}
	r.FileStoreOption.OnBeforeOpenFile = func(file string) (string, error) {
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
