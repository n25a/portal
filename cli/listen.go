package cli

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"text/template"

	"github.com/n25a/portal/internal/config"
)

type listenCmd struct {
	ConfigPath string `help:"Path to config file."`
}

var serverCommand = template.Must(
	template.New("server").
		Parse(`go-shadowsocks2 -s 'ss://AES-256-GCM:{{ .Password }}@:{{ .Port }}'`),
)

func (e *listenCmd) Run() error {
	if e.ConfigPath == "" {
		return errors.New("config path is required")
	}

	var data []string
	for _, user := range config.C.Users {
		var serverCommandBuilder strings.Builder
		s := struct {
			Password string
			Port     int
		}{
			Password: user.Password,
			Port:     user.Port,
		}
		err := serverCommand.Execute(&serverCommandBuilder, s)
		if err != nil {
			return err
		}

		data = append(data, serverCommandBuilder.String())
	}

	for i, _ := range config.C.Users {
		go func(idx int) {
			cmd := exec.Command(data[idx])
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}(i)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}
