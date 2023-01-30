package cli

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

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
		go func(user config.User) {
			cmd := exec.Command(
				"go-shadowsocks2", "-s",
				fmt.Sprintf("ss://AES-256-GCM:%s@:%d",
					user.Password,
					user.Port),
				"-verbose")
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}(user)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
