package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"notifier/config"
	"notifier/types"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var (
	// Mutex para proteger acceso concurrente al archivo
	fileMutex sync.Mutex
)

// NotificationUI maneja la interfaz de usuario para notificaciones
type NotificationUI struct {
	config *config.Config
	app    fyne.App
	window fyne.Window
}

// NewNotificationUI crea una nueva instancia de la UI
func NewNotificationUI(cfg *config.Config) *NotificationUI {
	return &NotificationUI{
		config: cfg,
		app:    app.New(),
	}
}

// ShowUI muestra la interfaz de usuario (función legacy para compatibilidad)
func ShowUI() {
	cfg := config.Default()
	ui := NewNotificationUI(cfg)
	ui.Show()
}

// Show muestra la ventana de la interfaz de usuario
func (ui *NotificationUI) Show() {
	ui.window = ui.app.NewWindow("Agregar Notificación")
	ui.setupWindow()
	ui.window.ShowAndRun()
}

// setupWindow configura la ventana y sus componentes
func (ui *NotificationUI) setupWindow() {
	ui.window.Resize(fyne.NewSize(float32(ui.config.WindowWidth), float32(ui.config.WindowHeight)))
	ui.window.SetFixedSize(true)

	// Crear campos de entrada
	horaEntry := widget.NewEntry()
	horaEntry.SetPlaceHolder("Ejemplo: 09:30")

	tituloEntry := widget.NewEntry()
	tituloEntry.SetPlaceHolder("Ingrese el título")

	mensajeEntry := widget.NewMultiLineEntry()
	mensajeEntry.SetPlaceHolder("Ingrese el mensaje de la notificación")
	mensajeEntry.Resize(fyne.NewSize(300, 80))

	recurrenteCheck := widget.NewCheck("Notificación recurrente (diaria)", func(bool) {})

	// Etiquetas informativas
	horaLabel := widget.NewLabel("Hora (formato 24h):")
	tituloLabel := widget.NewLabel("Título:")
	mensajeLabel := widget.NewLabel("Mensaje:")

	// Crear botones
	guardarBtn := widget.NewButton("Guardar", func() {
		ui.handleSave(horaEntry, tituloEntry, mensajeEntry, recurrenteCheck)
	})
	guardarBtn.Importance = widget.HighImportance

	cancelarBtn := widget.NewButton("Cancelar", func() {
		ui.window.Close()
	})

	// Crear layout
	content := container.NewVBox(
		widget.NewCard("Nueva Notificación", "", container.NewVBox(
			horaLabel,
			horaEntry,
			widget.NewSeparator(),
			tituloLabel,
			tituloEntry,
			widget.NewSeparator(),
			mensajeLabel,
			mensajeEntry,
			widget.NewSeparator(),
			recurrenteCheck,
		)),
		container.NewBorder(nil, nil, nil, nil,
			container.NewHBox(
				cancelarBtn,
				widget.NewSeparator(),
				guardarBtn,
			),
		),
	)

	ui.window.SetContent(content)
}

// handleSave maneja el guardado de una nueva notificación
func (ui *NotificationUI) handleSave(horaEntry, tituloEntry *widget.Entry, mensajeEntry *widget.Entry, recurrenteCheck *widget.Check) {
	// Crear notificación con los datos del formulario
	notif := types.Notificacion{
		Hora:       horaEntry.Text,
		Titulo:     tituloEntry.Text,
		Mensaje:    mensajeEntry.Text,
		Recurrente: recurrenteCheck.Checked,
	}

	// Validar la notificación
	if err := notif.Validate(); err != nil {
		ui.showError("Error de Validación", err.Error())
		return
	}

	// Guardar la notificación
	if err := ui.saveNotification(notif); err != nil {
		ui.showError("Error al Guardar", fmt.Sprintf("No se pudo guardar la notificación: %v", err))
		return
	}

	// Limpiar formulario
	ui.clearForm(horaEntry, tituloEntry, mensajeEntry, recurrenteCheck)

	// Mostrar confirmación
	ui.showSuccess("Éxito", "Notificación guardada correctamente")
}

// saveNotification guarda una notificación en el archivo JSON
func (ui *NotificationUI) saveNotification(notif types.Notificacion) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// Asegurar que el directorio existe
	if err := ui.config.EnsureDataDir(); err != nil {
		return fmt.Errorf("error creando directorio: %w", err)
	}

	// Leer notificaciones existentes
	var notificaciones []types.Notificacion
	if ui.config.DataFileExists() {
		file, err := os.ReadFile(ui.config.DataFile)
		if err != nil {
			return fmt.Errorf("error leyendo archivo: %w", err)
		}

		if len(file) > 0 {
			if err := json.Unmarshal(file, &notificaciones); err != nil {
				return fmt.Errorf("error parsing JSON existente: %w", err)
			}
		}
	}

	// Verificar duplicados
	for _, existing := range notificaciones {
		if existing.Hora == notif.Hora && existing.Titulo == notif.Titulo {
			return fmt.Errorf("ya existe una notificación con la misma hora (%s) y título (%s)", notif.Hora, notif.Titulo)
		}
	}

	// Agregar nueva notificación
	notificaciones = append(notificaciones, notif)

	// Guardar archivo
	data, err := json.MarshalIndent(notificaciones, "", "  ")
	if err != nil {
		return fmt.Errorf("error generando JSON: %w", err)
	}

	if err := os.WriteFile(ui.config.DataFile, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	log.Printf("Notificación guardada: %s a las %s", notif.Titulo, notif.Hora)
	return nil
}

// clearForm limpia todos los campos del formulario
func (ui *NotificationUI) clearForm(horaEntry, tituloEntry, mensajeEntry *widget.Entry, recurrenteCheck *widget.Check) {
	horaEntry.SetText("")
	tituloEntry.SetText("")
	mensajeEntry.SetText("")
	recurrenteCheck.SetChecked(false)
}

// showError muestra un diálogo de error
func (ui *NotificationUI) showError(title, message string) {
	dialog.ShowError(fmt.Errorf(message), ui.window)
}

// showSuccess muestra un diálogo de éxito
func (ui *NotificationUI) showSuccess(title, message string) {
	dialog.ShowInformation(title, message, ui.window)
}
