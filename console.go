//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"notifier/config"
	"notifier/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("🔔 RGDEV Notificador - Modo Consola")
	fmt.Println("==================================")

	// Inicializar
	if err := config.InitDataFile(); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}

	fmt.Printf("📁 Archivo: %s\n", config.DataFile)
	fmt.Printf("⏱️  Verificación cada: %v\n", config.CheckInterval)
	fmt.Printf("🕐 Hora actual: %s\n", time.Now().Format("15:04"))
	fmt.Println("🚀 Servicio iniciado - Presiona Ctrl+C para salir")
	fmt.Println()

	// Iniciar servicio
	go service.StartService()

	// Esperar señal para salir
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n✅ Servicio detenido")
}
