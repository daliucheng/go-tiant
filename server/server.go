package server

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type ServerConfig struct {
	Address   string `yaml:"address"`
	CloseChan chan struct{}
}

func (conf *ServerConfig) check() {
	if strings.Trim(conf.Address, " ") == "" {
		conf.Address = ":8080"
	}
}

func Run(engine *gin.Engine, conf ServerConfig) (err error) {
	conf.check()
	return engine.Run(conf.Address)
}
