package types

import (
	"fmt"
	"regexp"
	"strings"
)

// Notificacion representa una notificación programada
type Notificacion struct {
	Hora       string `json:"hora"`       // Formato "HH:MM" (15:04)
	Titulo     string `json:"titulo"`     // Título de la notificación
	Mensaje    string `json:"mensaje"`    // Contenido del mensaje
	Recurrente bool   `json:"recurrente"` // Si se repite diariamente
}

// ValidateNotificacion valida los campos de una notificación
func (n *Notificacion) Validate() error {
	// Validar formato de hora
	horaRegex := regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`)
	if !horaRegex.MatchString(n.Hora) {
		return fmt.Errorf("formato de hora inválido: %s. Use HH:MM (00:00-23:59)", n.Hora)
	}

	// Validar campos requeridos
	if strings.TrimSpace(n.Titulo) == "" {
		return fmt.Errorf("el título es requerido")
	}

	if strings.TrimSpace(n.Mensaje) == "" {
		return fmt.Errorf("el mensaje es requerido")
	}

	// Normalizar la hora (agregar cero inicial si es necesario)
	n.Hora = normalizeTime(n.Hora)
	n.Titulo = strings.TrimSpace(n.Titulo)
	n.Mensaje = strings.TrimSpace(n.Mensaje)

	return nil
}

// normalizeTime normaliza el formato de hora a HH:MM
func normalizeTime(hora string) string {
	parts := strings.Split(hora, ":")
	if len(parts) != 2 {
		return hora
	}

	// Asegurar que la hora tenga 2 dígitos
	if len(parts[0]) == 1 {
		parts[0] = "0" + parts[0]
	}

	// Asegurar que los minutos tengan 2 dígitos
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}

	return parts[0] + ":" + parts[1]
}
