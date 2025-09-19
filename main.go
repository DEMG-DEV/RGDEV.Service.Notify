package main

import (
	"context"
	"log"
	"notifier/config"
	"notifier/service"
	"notifier/ui"
	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
)

var (
	cfg          *config.Config
	notifService *service.Service
	ctx          context.Context
	cancelFunc   context.CancelFunc
)

func init() {
	// Configurar logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Crear contexto global
	ctx, cancelFunc = context.WithCancel(context.Background())

	// Inicializar configuración
	cfg = config.Default()

	// Asegurar que el directorio de datos existe
	if err := cfg.InitDataFile(); err != nil {
		log.Fatalf("Error inicializando configuración: %v", err)
	}

	log.Println("Aplicación inicializada correctamente")
}

func onReady() {
	// Configurar icono de la bandeja del sistema
	systray.SetTitle("Notificador RGDEV")
	systray.SetTooltip("Servicio de Notificaciones Programadas")

	// Crear elementos del menú
	mOpen := systray.AddMenuItem("🔔 Agregar Notificación", "Abrir interfaz para crear nueva notificación")
	mStatus := systray.AddMenuItem("📊 Estado del Servicio", "Ver estado actual del servicio")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("❌ Salir", "Cerrar aplicación completamente")

	// Inicializar y iniciar el servicio de notificaciones
	notifService = service.NewService(cfg)
	go notifService.Start()

	log.Println("Servicio de notificaciones iniciado")

	// Configurar manejo de señales del sistema
	go handleSystemSignals()

	// Manejo de eventos del menú en un goroutine separado
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				log.Println("Abriendo interfaz de usuario...")
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Error en UI: %v", r)
						}
					}()
					ui.ShowUI()
				}()

			case <-mStatus.ClickedCh:
				log.Println("Mostrando estado del servicio...")
				// Aquí se podría implementar una ventana de estado
				showServiceStatus()

			case <-mQuit.ClickedCh:
				log.Println("Cerrando aplicación...")
				cleanup()
				systray.Quit()
				return

			case <-ctx.Done():
				log.Println("Contexto cancelado, cerrando menú...")
				return
			}
		}
	}()
}

func onExit() {
	log.Println("Ejecutando limpieza al salir...")
	cleanup()
}

func cleanup() {
	// Cancelar contexto para detener todas las goroutines
	if cancelFunc != nil {
		cancelFunc()
	}

	// Detener el servicio de notificaciones
	if notifService != nil {
		notifService.Stop()
		log.Println("Servicio de notificaciones detenido")
	}

	log.Println("Limpieza completada")
}

// handleSystemSignals maneja las señales del sistema para una salida elegante
func handleSystemSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Señal recibida: %v. Cerrando aplicación...", sig)
		cleanup()
		systray.Quit()
	case <-ctx.Done():
		return
	}
}

// showServiceStatus muestra información sobre el estado del servicio
func showServiceStatus() {
	// Esta función se puede expandir para mostrar una ventana con estadísticas
	log.Println("=== ESTADO DEL SERVICIO ===")
	log.Printf("Archivo de configuración: %s", cfg.DataFile)
	log.Printf("Intervalo de verificación: %v", cfg.CheckInterval)
	log.Printf("Archivo de datos existe: %v", cfg.DataFileExists())

	if stat, err := os.Stat(cfg.DataFile); err == nil {
		log.Printf("Tamaño del archivo: %d bytes", stat.Size())
		log.Printf("Última modificación: %v", stat.ModTime())
	}
	log.Println("===========================")
}

func main() {
	log.Println("Iniciando RGDEV.Service.Notify...")

	// Ejecutar la aplicación de bandeja del sistema
	systray.Run(onReady, onExit)

	log.Println("Aplicación terminada")
}
