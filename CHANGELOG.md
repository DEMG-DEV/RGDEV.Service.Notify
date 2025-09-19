# Changelog

Todos los cambios notables de este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto se adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Sin Publicar]

### Planificado

- Función para eliminar notificaciones desde la interfaz web
- Edición de notificaciones existentes
- Soporte para sonidos personalizados
- Notificaciones con fecha específica

## [1.0.0] - 2025-01-19

### 🎉 Lanzamiento Inicial

#### Agregado

- **Sistema de notificaciones programadas** para Windows
- **Interfaz web moderna** (HTML/CSS/JavaScript)
- **Icono de bandeja del sistema** con menú contextual
- **Soporte para notificaciones recurrentes** y no recurrentes
- **Validación de entrada** con formato de hora HH:MM
- **Almacenamiento persistente** en JSON
- **Servidor HTTP integrado** en puerto 8080
- **Auto-apertura del navegador** al usar la interfaz
- **Prevención de duplicados** para notificaciones no recurrentes
- **Logs informativos** para debugging
- **Compilación sin CGO** (sin dependencias gráficas complejas)

#### Características Técnicas

- Verificación automática cada 30 segundos
- Interfaz responsive con gradientes modernos
- Control de concurrencia con mutexes
- Gestión elegante de servidor HTTP reutilizable
- Validación robusta con regex para formato de hora
- Limpieza automática de memoria

#### Archivos del Sistema

- `main.go` - Aplicación principal con bandeja del sistema
- `console.go` - Versión solo consola para desarrollo
- `service/service.go` - Motor de notificaciones en segundo plano
- `ui/ui.go` - Interfaz web con servidor HTTP integrado
- `types/types.go` - Definiciones y validación de datos
- `config/config.go` - Configuración del sistema
- `README.md` - Documentación completa
- `LICENSE` - Licencia MIT
- `.gitignore` - Configuración Git profesional

#### Dependencias

- `github.com/getlantern/systray` - Icono de bandeja del sistema
- `github.com/gen2brain/beeep` - Notificaciones nativas del SO

### Notas de la Versión

- **Total**: ~600 líneas de código
- **Arquitectura**: Modular y escalable
- **Compatibilidad**: Windows 10/11
- **Requerimientos**: Go 1.19+

---

### Formato de Versiones

Este proyecto usa [Semantic Versioning](https://semver.org/):

- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Nueva funcionalidad compatible con versiones anteriores
- **PATCH**: Correcciones de bugs compatibles

### Tipos de Cambios

- `Agregado` para nuevas características
- `Cambiado` para cambios en funcionalidad existente
- `Obsoleto` para características que serán removidas
- `Removido` para características removidas
- `Arreglado` para corrección de bugs
- `Seguridad` para vulnerabilidades
