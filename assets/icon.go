package assets

import (
	_ "embed"
)

// Icon contiene los datos del icono para la bandeja del sistema
//
//go:embed images/tray_icon.ico
var Icon []byte

// GetIcon devuelve los datos del icono
func GetIcon() []byte {
	return Icon
}
