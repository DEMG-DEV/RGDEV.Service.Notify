//go:build !windows
// +build !windows

package main

import (
	"log"
	"notifier/assets"
	"notifier/config"
	"notifier/service"
	"notifier/ui"

	"github.com/getlantern/systray"
)

func main() {
	// Inicializar configuración
	if err := config.InitDataFile(); err != nil {
		log.Fatalf("Error inicializando: %v", err)
	}

	// Ejecutar aplicación de bandeja
	systray.Run(onReady, onExit)
}

func onReady() {
	// Configurar icono de la bandeja con nuestro logo
	systray.SetIcon(assets.GetIcon())
	systray.SetTitle("RGDEV Notificador")
	systray.SetTooltip("Sistema de Notificaciones Programadas")

	// Menú
	mOpen := systray.AddMenuItem("Agregar Notificación", "Abrir interfaz")
	mQuit := systray.AddMenuItem("Salir", "Cerrar aplicación")

	// Iniciar servicio
	go service.StartService()

	// Manejar eventos
	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				go ui.ShowUI()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Println("Aplicación cerrada")
}
