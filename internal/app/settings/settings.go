package settings

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

var paths = [2]string{"./app.yml", "./configs/app.yml"}

type log struct {
	FileEnabled bool `yaml:"file-enabled"`
	MaxSize     int  `yaml:"max-size"`
	MaxAge      int  `yaml:"max-age"`
}

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
	Log           log           `yaml:"log"`
	Server        server        `yaml:"server"`
	Authorization authorization `yaml:"authorization"`
}

var loadedSettings *Settings

func defaultSettings() *Settings {
	return &Settings{
		Log: log{
			FileEnabled: false,
			MaxSize:     10,
			MaxAge:      30,
		},
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

func (s *Settings) configureLogger() {
	var handler *slog.TextHandler

	if s.Log.FileEnabled {
		logFile := &lumberjack.Logger{
			Filename:  "./logs/app.log",
			MaxSize:   s.Log.MaxSize,
			MaxAge:    s.Log.MaxAge,
			LocalTime: true,
		}
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		handler = slog.NewTextHandler(multiWriter, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
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
			break
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

	loadedSettings.configureLogger()

	return *loadedSettings
}
