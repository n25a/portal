package cli

import (
	"errors"
	"fmt"

	"github.com/n25a/portal/internal/config"
)

type listenCmd struct {
	ConfigPath string `help:"Path to config file."`
}

func (e *listenCmd) Run() error {
	if e.ConfigPath == "" {
		return errors.New("config path is required")
	}

	err := config.LoadConfig(e.ConfigPath)
	if err != nil {
		return err
	}

	for _, user := range config.C.Users {
		fmt.Println(
			"go-shadowsocks2", "-s",
			fmt.Sprintf("ss://AES-256-GCM:%s@:%d",
				user.Password,
				user.Port))

	}

	return nil
}
