package job

import (
	"github.com/bearqy/go-gin-templete/internal/home"
)

func init() {
	Register("aaa", aaa)
}

func aaa() {
	home.HomeJob()
}
