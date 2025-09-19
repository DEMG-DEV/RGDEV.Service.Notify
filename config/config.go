package config

import (
	"os"
	"path/filepath"
	"time"
)

// Config contiene la configuración de la aplicación
type Config struct {
	DataFile      string        // Ruta al archivo de notificaciones
	CheckInterval time.Duration // Intervalo de verificación
	WindowWidth   int           // Ancho de la ventana UI
	WindowHeight  int           // Alto de la ventana UI
	MaxFileSize   int64         // Tamaño máximo del archivo JSON en bytes
}

// Default devuelve la configuración por defecto
func Default() *Config {
	return &Config{
		DataFile:      filepath.Join("data", "notificaciones.json"),
		CheckInterval: 30 * time.Second,
		WindowWidth:   350,
		WindowHeight:  400,
		MaxFileSize:   1024 * 1024, // 1MB
	}
}

// EnsureDataDir crea el directorio de datos si no existe
func (c *Config) EnsureDataDir() error {
	dir := filepath.Dir(c.DataFile)
	return os.MkdirAll(dir, 0755)
}

// DataFileExists verifica si el archivo de datos existe
func (c *Config) DataFileExists() bool {
	_, err := os.Stat(c.DataFile)
	return err == nil
}

// InitDataFile inicializa el archivo de datos con un array JSON vacío
func (c *Config) InitDataFile() error {
	if err := c.EnsureDataDir(); err != nil {
		return err
	}

	if !c.DataFileExists() {
		return os.WriteFile(c.DataFile, []byte("[]"), 0644)
	}

	return nil
}
