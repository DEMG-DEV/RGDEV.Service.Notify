package types

import (
	"fmt"
	"regexp"
	"strings"
)

// Notificacion representa una notificación programada
type Notificacion struct {
	Hora       string `json:"hora"`
	Titulo     string `json:"titulo"`
	Mensaje    string `json:"mensaje"`
	Recurrente bool   `json:"recurrente"`
}

// Validate valida los campos de una notificación
func (n *Notificacion) Validate() error {
	// Validar formato de hora
	if !regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]$`).MatchString(n.Hora) {
		return fmt.Errorf("formato de hora inválido: %s", n.Hora)
	}

	// Validar campos requeridos
	if strings.TrimSpace(n.Titulo) == "" {
		return fmt.Errorf("el título es requerido")
	}

	if strings.TrimSpace(n.Mensaje) == "" {
		return fmt.Errorf("el mensaje es requerido")
	}

	return nil
}
