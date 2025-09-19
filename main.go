// Package main implementa el sistema de notificaciones programadas RGDEV.Service.Notify.
//
// Este sistema permite crear y gestionar notificaciones que se ejecutan en momentos
// específicos del día, con soporte para notificaciones recurrentes y no recurrentes.
//
// Características principales:
//   - Interfaz web moderna para gestión de notificaciones
//   - Icono de bandeja del sistema para acceso rápido
//   - Notificaciones nativas del sistema operativo
//   - Almacenamiento persistente en JSON
//   - Sin dependencias gráficas complejas
//   - Ejecución en background sin ventana de consola
//
// Uso básico:
//
//	go run main.go
//
// O después de compilar:
//
//	./notifier.exe
//
// Para más información consulte el README.md

//go:build windows
// +build windows

package main

import (
	"log"
	"notifier/assets"
	"notifier/config"
	"notifier/service"
	"notifier/ui"
	"syscall"

	"github.com/getlantern/systray"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procFreeConsole = kernel32.NewProc("FreeConsole")
)

func main() {
	// Ocultar ventana de consola en Windows
	hideConsole()

	// Inicializar configuración
	if err := config.InitDataFile(); err != nil {
		log.Fatalf("Error inicializando: %v", err)
	}

	// Ejecutar aplicación de bandeja
	systray.Run(onReady, onExit)
}

// hideConsole oculta la ventana de consola en Windows
func hideConsole() {
	procFreeConsole.Call()
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
