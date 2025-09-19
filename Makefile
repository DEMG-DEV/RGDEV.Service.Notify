# RGDEV.Service.Notify Makefile
# Automatización de tareas de desarrollo y build

# Variables
BINARY_NAME=notifier
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
VERSION=1.0.0
BUILD_DIR=build
DIST_DIR=dist

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build clean test lint fmt vet deps run console install dist package

# Comando por defecto
all: clean fmt vet test build

# Ayuda
help: ## Mostrar esta ayuda
	@echo "$(GREEN)RGDEV.Service.Notify - Comandos Disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

# Desarrollo
fmt: ## Formatear código Go
	@echo "$(GREEN)Formateando código...$(NC)"
	go fmt ./...

vet: ## Verificar código con go vet
	@echo "$(GREEN)Ejecutando go vet...$(NC)"
	go vet ./...

lint: ## Ejecutar linter (requiere golangci-lint)
	@echo "$(GREEN)Ejecutando linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint no instalado. Ejecuta: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

# Dependencias
deps: ## Instalar dependencias
	@echo "$(GREEN)Instalando dependencias...$(NC)"
	go mod tidy
	go mod download

deps-update: ## Actualizar dependencias
	@echo "$(GREEN)Actualizando dependencias...$(NC)"
	go get -u ./...
	go mod tidy

# Testing
test: ## Ejecutar tests
	@echo "$(GREEN)Ejecutando tests...$(NC)"
	go test -v ./...

test-coverage: ## Ejecutar tests con cobertura
	@echo "$(GREEN)Ejecutando tests con cobertura...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Cobertura generada en coverage.html$(NC)"

benchmark: ## Ejecutar benchmarks
	@echo "$(GREEN)Ejecutando benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

# Build
build: ## Compilar aplicación
	@echo "$(GREEN)Compilando aplicación...$(NC)"
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_WINDOWS) .

build-linux: ## Compilar para Linux
	@echo "$(GREEN)Compilando para Linux...$(NC)"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_UNIX) .

build-all: ## Compilar para todas las plataformas
	@echo "$(GREEN)Compilando para todas las plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "$(GREEN)Binarios generados en $(BUILD_DIR)/$(NC)"

# Ejecución
run: ## Ejecutar aplicación principal
	@echo "$(GREEN)Ejecutando aplicación...$(NC)"
	go run main.go

console: ## Ejecutar versión consola
	@echo "$(GREEN)Ejecutando versión consola...$(NC)"
	go run console.go

# Instalación
install: ## Instalar en GOPATH
	@echo "$(GREEN)Instalando aplicación...$(NC)"
	go install .

# Distribución
dist: build-all ## Crear distribución completa
	@echo "$(GREEN)Creando distribución...$(NC)"
	@mkdir -p $(DIST_DIR)
	@cp -r $(BUILD_DIR)/* $(DIST_DIR)/
	@cp README.md $(DIST_DIR)/
	@cp LICENSE $(DIST_DIR)/
	@cp CHANGELOG.md $(DIST_DIR)/
	@echo "$(GREEN)Distribución creada en $(DIST_DIR)/$(NC)"

package: dist ## Crear archivos de distribución comprimidos
	@echo "$(GREEN)Creando paquetes...$(NC)"
	@cd $(DIST_DIR) && \
	tar -czf $(BINARY_NAME)-$(VERSION)-windows-amd64.tar.gz $(BINARY_NAME)-windows-amd64.exe README.md LICENSE CHANGELOG.md && \
	tar -czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 README.md LICENSE CHANGELOG.md && \
	tar -czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 README.md LICENSE CHANGELOG.md
	@echo "$(GREEN)Paquetes creados en $(DIST_DIR)/$(NC)"

# Release
release: clean fmt vet test build-all package ## Proceso completo de release
	@echo "$(GREEN)Release $(VERSION) completado!$(NC)"

# Limpieza
clean: ## Limpiar archivos generados
	@echo "$(GREEN)Limpiando archivos generados...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)
	@rm -f $(BINARY_WINDOWS)
	@rm -f coverage.out coverage.html
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)

# Docker (futuro)
docker-build: ## Construir imagen Docker
	@echo "$(GREEN)Construyendo imagen Docker...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .

# Desarrollo
dev: ## Modo desarrollo con hot reload (requiere air)
	@echo "$(GREEN)Iniciando modo desarrollo...$(NC)"
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)Air no instalado. Ejecuta: go install github.com/cosmtrek/air@latest$(NC)"; \
		echo "$(YELLOW)Usando go run en su lugar...$(NC)"; \
		go run main.go; \
	fi

# Información del proyecto
info: ## Mostrar información del proyecto
	@echo "$(GREEN)RGDEV.Service.Notify$(NC)"
	@echo "Versión: $(VERSION)"
	@echo "Arquitectura Go: $(shell go env GOOS)/$(shell go env GOARCH)"
	@echo "Versión Go: $(shell go version)"
	@echo "Módulos:"
	@go list -m all

# Security
security: ## Verificar vulnerabilidades de seguridad
	@echo "$(GREEN)Verificando seguridad...$(NC)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "$(YELLOW)govulncheck no instalado. Ejecuta: go install golang.org/x/vuln/cmd/govulncheck@latest$(NC)"; \
	fi
