package service

import (
	"encoding/json"
	"fmt"
	"log"
	"notifier/config"
	"notifier/types"
	"os"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
)

// Service maneja el servicio de notificaciones
type Service struct {
	config   *config.Config
	enviadas map[string]time.Time // Almacena cuándo se envió cada notificación no recurrente
	mutex    sync.RWMutex         // Protege el acceso concurrente al map
	stopCh   chan struct{}        // Canal para detener el servicio
}

// NewService crea una nueva instancia del servicio
func NewService(cfg *config.Config) *Service {
	return &Service{
		config:   cfg,
		enviadas: make(map[string]time.Time),
		stopCh:   make(chan struct{}),
	}
}

// StartService inicia el servicio de notificaciones (función legacy para compatibilidad)
func StartService() {
	cfg := config.Default()
	if err := cfg.InitDataFile(); err != nil {
		log.Printf("Error inicializando archivo de datos: %v", err)
		return
	}

	service := NewService(cfg)
	service.Start()
}

// Start inicia el bucle principal del servicio
func (s *Service) Start() {
	log.Println("Iniciando servicio de notificaciones...")

	// Limpiar notificaciones enviadas una vez al día
	go s.cleanupRoutine()

	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.checkNotifications(); err != nil {
				log.Printf("Error verificando notificaciones: %v", err)
			}
		case <-s.stopCh:
			log.Println("Deteniendo servicio de notificaciones...")
			return
		}
	}
}

// Stop detiene el servicio de notificaciones
func (s *Service) Stop() {
	close(s.stopCh)
}

// checkNotifications verifica y envía las notificaciones pendientes
func (s *Service) checkNotifications() error {
	// Verificar tamaño del archivo antes de leerlo
	if stat, err := os.Stat(s.config.DataFile); err == nil {
		if stat.Size() > s.config.MaxFileSize {
			return fmt.Errorf("archivo de notificaciones demasiado grande: %d bytes", stat.Size())
		}
	}

	file, err := os.ReadFile(s.config.DataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Crear archivo si no existe
			if initErr := s.config.InitDataFile(); initErr != nil {
				return fmt.Errorf("error creando archivo de datos: %w", initErr)
			}
			return nil
		}
		return fmt.Errorf("error leyendo archivo de notificaciones: %w", err)
	}

	// Validar que el archivo no esté vacío
	if len(file) == 0 {
		return s.config.InitDataFile()
	}

	var notificaciones []types.Notificacion
	if err := json.Unmarshal(file, &notificaciones); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	ahora := time.Now()
	horaActual := ahora.Format("15:04")
	fechaActual := ahora.Format("2006-01-02")

	for _, n := range notificaciones {
		// Validar notificación antes de procesarla
		if err := n.Validate(); err != nil {
			log.Printf("Notificación inválida ignorada: %v", err)
			continue
		}

		if n.Hora == horaActual {
			notifKey := fmt.Sprintf("%s_%s_%s", fechaActual, n.Hora, n.Titulo)

			s.mutex.RLock()
			_, yaEnviada := s.enviadas[notifKey]
			s.mutex.RUnlock()

			// Si no es recurrente y ya fue enviada hoy, omitir
			if !n.Recurrente && yaEnviada {
				continue
			}

			// Enviar notificación
			if err := s.sendNotification(n); err != nil {
				log.Printf("Error enviando notificación '%s': %v", n.Titulo, err)
				continue
			}

			log.Printf("Notificación enviada: %s - %s", n.Titulo, n.Mensaje)

			// Marcar como enviada si no es recurrente
			if !n.Recurrente {
				s.mutex.Lock()
				s.enviadas[notifKey] = ahora
				s.mutex.Unlock()
			}
		}
	}

	return nil
}

// sendNotification envía una notificación al sistema
func (s *Service) sendNotification(n types.Notificacion) error {
	if err := beeep.Notify(n.Titulo, n.Mensaje, ""); err != nil {
		return fmt.Errorf("error enviando notificación del sistema: %w", err)
	}
	return nil
}

// cleanupRoutine limpia las notificaciones enviadas una vez al día
func (s *Service) cleanupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupOldNotifications()
		case <-s.stopCh:
			return
		}
	}
}

// cleanupOldNotifications elimina notificaciones enviadas de días anteriores
func (s *Service) cleanupOldNotifications() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ahora := time.Now()
	fechaActual := ahora.Format("2006-01-02")

	for key, fechaEnvio := range s.enviadas {
		fechaEnvioStr := fechaEnvio.Format("2006-01-02")
		if fechaEnvioStr != fechaActual {
			delete(s.enviadas, key)
		}
	}

	log.Printf("Limpieza completada. Notificaciones en memoria: %d", len(s.enviadas))
}
