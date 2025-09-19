# RGDEV.Service.Notify

## 📋 Descripción

RGDEV.Service.Notify es un servicio de notificaciones programadas para Windows que funciona como una aplicación de bandeja del sistema. Permite a los usuarios programar notificaciones que se mostrarán en momentos específicos del día, con soporte para notificaciones recurrentes y no recurrentes.

## 🏗️ Arquitectura del Sistema

El sistema está compuesto por tres módulos principales:

### 1. **Main Application** (`main.go`)
- **Función**: Punto de entrada principal de la aplicación
- **Responsabilidades**:
  - Inicializa el icono de bandeja del sistema usando `systray`
  - Maneja el menú contextual con opciones "Abrir UI" y "Salir"
  - Inicia el servicio de notificaciones en segundo plano
  - Gestiona los eventos del usuario

### 2. **Service Module** (`service/service.go`)
- **Función**: Motor de notificaciones que funciona en segundo plano
- **Responsabilidades**:
  - Lee el archivo de notificaciones cada 30 segundos
  - Compara la hora actual con las notificaciones programadas
  - Envía notificaciones usando la librería `beeep`
  - Maneja la lógica de notificaciones recurrentes vs no recurrentes
  - Evita duplicados en notificaciones no recurrentes

### 3. **UI Module** (`ui/ui.go`)
- **Función**: Interfaz gráfica para agregar nuevas notificaciones
- **Responsabilidades**:
  - Proporciona formulario para crear notificaciones
  - Guarda nuevas notificaciones en el archivo JSON
  - Interfaz intuitiva desarrollada con Fyne

### 4. **Data Storage** (`data/notificaciones.json`)
- **Función**: Almacenamiento persistente de notificaciones
- **Formato**: Array JSON con objetos de notificación

## 📊 Estructura de Datos

```go
type Notificacion struct {
    Hora      string `json:"hora"`      // Formato "HH:MM" (15:04)
    Titulo    string `json:"titulo"`    // Título de la notificación
    Mensaje   string `json:"mensaje"`   // Contenido del mensaje
    Recurrente bool   `json:"recurrente"` // Si se repite diariamente
}
```

## 🔧 Dependencias

El proyecto utiliza las siguientes librerías externas:

- **[systray](https://github.com/getlantern/systray)**: Para crear el icono de bandeja del sistema
- **[beeep](https://github.com/gen2brain/beeep)**: Para mostrar notificaciones nativas del sistema
- **[fyne](https://fyne.io/)**: Para la interfaz gráfica de usuario

## 📁 Estructura del Proyecto

```
RGDEV.Service.Notify/
├── main.go                 # Aplicación principal
├── README.md              # Documentación
├── data/
│   └── notificaciones.json # Almacén de notificaciones
├── service/
│   └── service.go         # Servicio de notificaciones
└── ui/
    └── ui.go             # Interfaz de usuario
```

## 🚀 Instalación y Configuración

### Prerrequisitos

- Go 1.19 o superior
- Windows 10/11

### Pasos de Instalación

1. **Clonar el repositorio**:
   ```bash
   git clone <repository-url>
   cd RGDEV.Service.Notify
   ```

2. **Inicializar el módulo Go** (si no existe go.mod):
   ```bash
   go mod init notifier
   ```

3. **Instalar dependencias**:
   ```bash
   go get github.com/getlantern/systray
   go get github.com/gen2brain/beeep
   go get fyne.io/fyne/v2/app
   go get fyne.io/fyne/v2/widget
   ```

4. **Crear el directorio de datos**:
   ```bash
   mkdir data
   echo [] > data/notificaciones.json
   ```

5. **Compilar la aplicación**:
   ```bash
   go build -o notifier.exe
   ```

## 💻 Uso

### Ejecución
```bash
./notifier.exe
```

### Funcionalidades

1. **Icono de Bandeja**: Al ejecutar, aparece un icono en la bandeja del sistema
2. **Menú Contextual**:
   - **"Abrir UI"**: Abre la ventana para agregar nuevas notificaciones
   - **"Salir"**: Cierra la aplicación completamente

3. **Agregar Notificaciones**:
   - Clic derecho en el icono → "Abrir UI"
   - Completar los campos:
     - **Hora**: Formato HH:MM (24 horas)
     - **Título**: Título de la notificación
     - **Mensaje**: Contenido descriptivo
     - **Recurrente**: Marcar si se debe repetir diariamente
   - Clic en "Guardar"

### Tipos de Notificaciones

- **No Recurrentes**: Se muestran una sola vez en la hora especificada
- **Recurrentes**: Se muestran todos los días a la misma hora

## 🔍 Funcionamiento Interno

1. **Inicio**: La aplicación se ejecuta y crea un icono en la bandeja del sistema
2. **Servicio en Segundo Plano**: Un goroutine lee el archivo JSON cada 30 segundos
3. **Verificación de Hora**: Compara la hora actual (formato HH:MM) con las notificaciones programadas
4. **Envío de Notificaciones**: Cuando coincide la hora, envía la notificación usando el sistema nativo
5. **Control de Duplicados**: Las notificaciones no recurrentes se marcan como enviadas para evitar repeticiones
6. **Gestión de UI**: La interfaz se abre bajo demanda para agregar nuevas notificaciones

## 📝 Ejemplo de Archivo de Datos

```json
[
  {
    "hora": "09:00",
    "titulo": "Reunión Matutina",
    "mensaje": "Reunión de equipo en 5 minutos",
    "recurrente": true
  },
  {
    "hora": "14:30",
    "titulo": "Recordatorio",
    "mensaje": "Llamar al cliente importante",
    "recurrente": false
  }
]
```

## 🛠️ Desarrollo y Extensión

### Posibles Mejoras

- [ ] Agregar función para editar/eliminar notificaciones existentes
- [ ] Implementar diferentes tipos de sonidos para notificaciones
- [ ] Añadir soporte para notificaciones con fecha específica
- [ ] Crear interfaz para ver historial de notificaciones enviadas
- [ ] Implementar configuración personalizable (intervalo de verificación, etc.)

### Estructura Modular

El código está diseñado de manera modular, facilitando la extensión:
- **service/**: Lógica de negocio independiente
- **ui/**: Interfaz separada que puede reemplazarse fácilmente
- **main.go**: Coordinador simple entre módulos

## 🔧 Troubleshooting

### Problemas Comunes

1. **La aplicación no muestra notificaciones**:
   - Verificar que el archivo `data/notificaciones.json` existe
   - Comprobar el formato de hora (debe ser HH:MM en formato 24h)

2. **Error al abrir la UI**:
   - Verificar que las dependencias de Fyne están instaladas correctamente
   - Comprobar que el sistema soporta interfaces gráficas

3. **El icono no aparece en la bandeja**:
   - Verificar que la librería systray es compatible con el sistema operativo
   - Comprobar configuración de la bandeja del sistema Windows

## 📄 Licencia

MIT

## 👨‍💻 Desarrollador

Desarrollado por RGDEV - Sistema de Notificaciones v1.0
