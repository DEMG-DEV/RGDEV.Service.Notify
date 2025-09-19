package types

import (
	"testing"
)

func TestNotificacion_Validate(t *testing.T) {
	tests := []struct {
		name         string
		notificacion Notificacion
		wantErr      bool
		errorMsg     string
	}{
		{
			name: "Notificación válida básica",
			notificacion: Notificacion{
				Hora:       "14:30",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr: false,
		},
		{
			name: "Notificación válida recurrente",
			notificacion: Notificacion{
				Hora:       "09:00",
				Titulo:     "Reunión diaria",
				Mensaje:    "Stand-up meeting",
				Recurrente: true,
			},
			wantErr: false,
		},
		{
			name: "Hora inválida - formato incorrecto",
			notificacion: Notificacion{
				Hora:       "25:00",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "formato de hora inválido",
		},
		{
			name: "Hora inválida - minutos incorrectos",
			notificacion: Notificacion{
				Hora:       "14:60",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "formato de hora inválido",
		},
		{
			name: "Título vacío",
			notificacion: Notificacion{
				Hora:       "14:30",
				Titulo:     "",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "el título es requerido",
		},
		{
			name: "Título solo espacios",
			notificacion: Notificacion{
				Hora:       "14:30",
				Titulo:     "   ",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "el título es requerido",
		},
		{
			name: "Mensaje vacío",
			notificacion: Notificacion{
				Hora:       "14:30",
				Titulo:     "Test",
				Mensaje:    "",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "el mensaje es requerido",
		},
		{
			name: "Mensaje solo espacios",
			notificacion: Notificacion{
				Hora:       "14:30",
				Titulo:     "Test",
				Mensaje:    "   ",
				Recurrente: false,
			},
			wantErr:  true,
			errorMsg: "el mensaje es requerido",
		},
		{
			name: "Hora válida con un dígito",
			notificacion: Notificacion{
				Hora:       "9:05",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr: false,
		},
		{
			name: "Hora límite superior válida",
			notificacion: Notificacion{
				Hora:       "23:59",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr: false,
		},
		{
			name: "Hora límite inferior válida",
			notificacion: Notificacion{
				Hora:       "00:00",
				Titulo:     "Test",
				Mensaje:    "Mensaje de prueba",
				Recurrente: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notificacion.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() esperaba error pero no lo obtuvo")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Validate() error = %v, esperaba que contuviera %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error inesperado = %v", err)
				}
			}
		})
	}
}

func TestNotificacion_ValidateFormats(t *testing.T) {
	validHours := []string{
		"00:00", "01:30", "09:15", "12:00", "18:45", "23:59",
		"9:05", "8:30", // Formatos con un dígito
	}

	for _, hora := range validHours {
		t.Run("hora_válida_"+hora, func(t *testing.T) {
			notif := Notificacion{
				Hora:    hora,
				Titulo:  "Test",
				Mensaje: "Test message",
			}

			if err := notif.Validate(); err != nil {
				t.Errorf("Hora %s debería ser válida, pero obtuvo error: %v", hora, err)
			}
		})
	}

	invalidHours := []string{
		"24:00", "25:30", "12:60", "ab:cd", "12:", ":30", "12", "1200",
		"12:1", "-1:30", "12:-1",
	}

	for _, hora := range invalidHours {
		t.Run("hora_inválida_"+hora, func(t *testing.T) {
			notif := Notificacion{
				Hora:    hora,
				Titulo:  "Test",
				Mensaje: "Test message",
			}

			if err := notif.Validate(); err == nil {
				t.Errorf("Hora %s debería ser inválida, pero pasó la validación", hora)
			}
		})
	}
}

func BenchmarkNotificacion_Validate(b *testing.B) {
	notif := Notificacion{
		Hora:       "14:30",
		Titulo:     "Benchmark Test",
		Mensaje:    "Este es un mensaje de prueba para el benchmark",
		Recurrente: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		notif.Validate()
	}
}

// Helper function para verificar si una cadena contiene otra
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
