package config

import (
	"os"
	"time"
)

const (
	DataFile      = "data/notificaciones.json"
	CheckInterval = 30 * time.Second
)

// InitDataFile crea el archivo de datos si no existe
func InitDataFile() error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}

	if _, err := os.Stat(DataFile); os.IsNotExist(err) {
		return os.WriteFile(DataFile, []byte("[]"), 0644)
	}

	return nil
}
