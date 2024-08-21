package settings

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var paths = [2]string{"./app.yml", "./configs/app.yml"}

type server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type authorization struct {
	Enabled bool   `yaml:"enabled"`
	Header  string `yaml:"header"`
	Key     string `yaml:"key"`
}

type Settings struct {
	Server        server        `yaml:"server"`
	Authorization authorization `yaml:"authorization"`
}

var loadedSettings *Settings

func defaultSettings() *Settings {
	return &Settings{
		Server: server{
			Host: "0.0.0.0",
			Port: "8080",
		},
		Authorization: authorization{
			Enabled: false,
			Header:  "Api-Key",
			Key:     "",
		},
	}
}

func LoadSettings() Settings {
	if loadedSettings != nil {
		slog.Error("Settings already loaded")
		os.Exit(1)
	}

	loadedSettings = defaultSettings()
	var file *os.File
	var fileErr error

	for _, path := range paths {
		file, fileErr = os.Open(path)

		if fileErr != nil {
			slog.Warn("Settings not found", "path", path)
		} else {
			slog.Info("Settings found", "path", path)
		}
	}

	if file != nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)

		if err := decoder.Decode(&loadedSettings); err != nil {
			slog.Error("Decoding YAML settings", "error", err)
			os.Exit(1)
		}
	}

	return *loadedSettings
}
