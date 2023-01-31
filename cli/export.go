// Package cli is the commands collector for the CLI application.
package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/n25a/portal/internal/config"
)

type exportCmd struct {
	ConfigPath string `help:"Path to config file."`
	ExportPath string `help:"export path." type:"existingdir"`
}

var phoneData = template.Must(
	template.New("Phone").Parse(`{
"server" : {{ .Server }}
"password": "{{ .Password }}"
"port": {{ .Port }}
"Encryption": "AES-256-GCM"
}`),
)

var bashData = template.Must(
	template.New("Bash").Parse(`#!/bin/bash
go-shadowsocks2 -c 'ss://AES-256-GCM:{{ .Password }}@{{ .Server }}:{{ .Port }}' \
    -verbose -socks :8558
`))

func (e *exportCmd) Run() error {
	if e.ConfigPath == "" {
		return errors.New("config path is required")
	}

	err := config.LoadConfig(e.ConfigPath)
	if err != nil {
		return err
	}

	// create text file for phone client
	var data []string
	for _, user := range config.C.Users {
		var phoneBuilder strings.Builder
		s := struct {
			Server   string
			Password string
			Port     int
		}{
			Server:   config.C.Server,
			Password: user.Password,
			Port:     user.Port,
		}
		err := phoneData.Execute(&phoneBuilder, s)
		if err != nil {
			return err
		}

		data = append(data, phoneBuilder.String())
	}

	for i, user := range config.C.Users {
		err := ioutil.WriteFile(
			filepath.Join(e.ExportPath, fmt.Sprintf("%d.txt", user.Port)), []byte(data[i]), 0644)
		if err != nil {
			return err
		}
	}

	// create bash file for os client
	data = []string{}
	for _, user := range config.C.Users {
		var bashBuilder strings.Builder
		s := struct {
			Server   string
			Password string
			Port     int
		}{
			Server:   config.C.Server,
			Password: user.Password,
			Port:     user.Port,
		}
		err := bashData.Execute(&bashBuilder, s)
		if err != nil {
			return err
		}

		data = append(data, bashBuilder.String())
	}

	for i, user := range config.C.Users {
		err := ioutil.WriteFile(
			filepath.Join(e.ExportPath, fmt.Sprintf("%d.sh", user.Port)), []byte(data[i]), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
