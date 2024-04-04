package base

import (
	"path/filepath"

	"github.com/jetiny/sgin/utils"
	"github.com/oschwald/geoip2-golang"
)

func initIp() (*geoip2.Reader, error) {
	cwd := utils.GetExecutablePath(EnvIPUseCwd.Bool())
	db, err := geoip2.Open(filepath.Join(cwd, EnvIPDatabase.String()))
	if err != nil {
		return nil, err
	}
	return db, nil
}
