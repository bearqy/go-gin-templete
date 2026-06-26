package job

import (
	"go-gin-templete/internal/home"
)

func RegisterDefaults() {
	Register("demo", demo)
}

func demo() {
	home.HomeJob()
}
