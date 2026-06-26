package cli

import (
	"flag"
	"os"
)

// ConfigFilePath 默认读取的配置文件路径
var ConfigFilePath = "config/default.yaml"

func Init() {
	configFilePath, err := Parse(os.Args[1:])
	if err != nil {
		flag.CommandLine.Parse(os.Args[1:])
		return
	}
	ConfigFilePath = configFilePath
}

func Parse(args []string) (string, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	configFilePath := ConfigFilePath
	fs.StringVar(&configFilePath, "config", ConfigFilePath, "配置文件路径")
	if err := fs.Parse(args); err != nil {
		return "", err
	}
	return configFilePath, nil
}
