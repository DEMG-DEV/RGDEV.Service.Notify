package service

import (
	"encoding/json"
	"log"
	"notifier/config"
	"notifier/types"
	"os"
	"path/filepath"
	"time"

	"github.com/gen2brain/beeep"
)

var enviadas = make(map[string]bool)

// StartService inicia el servicio de notificaciones
func StartService() {
	log.Println("Iniciando servicio de notificaciones...")

	for {
		checkNotifications()
		time.Sleep(config.CheckInterval)
	}
}

// checkNotifications verifica y envía notificaciones
func checkNotifications() {
	file, err := os.ReadFile(config.DataFile)
	if err != nil {
		return
	}

	var notificaciones []types.Notificacion
	if err := json.Unmarshal(file, &notificaciones); err != nil {
		return
	}

	horaActual := time.Now().Format("15:04")
	fechaActual := time.Now().Format("2006-01-02")

	for _, n := range notificaciones {
		if n.Hora == horaActual {
			key := fechaActual + "_" + n.Hora + "_" + n.Titulo

			// Si no es recurrente y ya fue enviada hoy, omitir
			if !n.Recurrente && enviadas[key] {
				continue
			}

			// Enviar notificación con logo
			logoPath, _ := filepath.Abs("assets/images/logo.jpeg")
			if err := beeep.Notify(n.Titulo, n.Mensaje, logoPath); err != nil {
				log.Printf("Error enviando notificación: %v", err)
				continue
			}

			log.Printf("Notificación enviada: %s - %s", n.Titulo, n.Mensaje)

			// Marcar como enviada si no es recurrente
			if !n.Recurrente {
				enviadas[key] = true
			}
		}
	}

	// Limpiar notificaciones de días anteriores
	if len(enviadas) > 100 {
		enviadas = make(map[string]bool)
	}
}
