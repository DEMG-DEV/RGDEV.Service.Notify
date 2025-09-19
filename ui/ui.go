package ui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"notifier/config"
	"notifier/types"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	server    *http.Server
	isRunning bool
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>RGDEV Notificador</title>
    <link rel="icon" type="image/jpeg" href="/favicon">
    <link rel="shortcut icon" type="image/jpeg" href="/favicon">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #f0f2f5; 
            padding: 20px;
        }
        .container { 
            max-width: 600px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 12px; 
            box-shadow: 0 4px 20px rgba(0,0,0,0.1); 
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white; 
            padding: 25px 20px; 
            text-align: center;
            position: relative;
        }
        .header img {
            border-radius: 12px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
            transition: transform 0.3s ease;
        }
        .header img:hover {
            transform: scale(1.05);
        }
        .header h1 { font-size: 24px; margin-bottom: 5px; }
        .header p { opacity: 0.9; }
        .content { padding: 30px; }
        .status {
            background: #e8f5e8; 
            border: 1px solid #4caf50; 
            border-radius: 8px; 
            padding: 15px; 
            margin-bottom: 30px;
            text-align: center;
        }
        .form-section { margin-bottom: 30px; }
        .form-section h2 { margin-bottom: 20px; color: #333; }
        .form-group { margin-bottom: 20px; }
        .form-group label { 
            display: block; 
            margin-bottom: 8px; 
            font-weight: 600; 
            color: #555;
        }
        .form-group input, .form-group textarea { 
            width: 100%; 
            padding: 12px; 
            border: 2px solid #e0e0e0; 
            border-radius: 8px; 
            font-size: 14px;
            transition: border-color 0.3s;
        }
        .form-group input:focus, .form-group textarea:focus { 
            outline: none; 
            border-color: #667eea; 
        }
        .checkbox-group {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 20px;
        }
        .checkbox-group input[type="checkbox"] {
            width: auto;
            transform: scale(1.2);
        }
        .btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s;
            width: 100%;
        }
        .btn:hover { transform: translateY(-2px); }
        .btn:active { transform: translateY(0); }
        .notification-list { margin-top: 30px; }
        .notification-item {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 0 8px 8px 0;
        }
        .notification-time { font-weight: bold; color: #667eea; }
        .notification-title { font-size: 16px; margin: 5px 0; }
        .notification-message { color: #666; font-size: 14px; }
        .notification-type { 
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 12px;
            margin-top: 5px;
        }
        .close-btn {
            background: #dc3545;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 6px;
            cursor: pointer;
            margin-top: 20px;
        }
        .close-btn:hover { background: #c82333; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="/logo" alt="RGDEV Logo" style="height: 60px; margin-bottom: 10px; border-radius: 8px;">
            <h1>🔔 RGDEV Notificador</h1>
            <p>Sistema de Notificaciones Programadas</p>
        </div>
        
        <div class="content">
            <div class="status">
                <strong>📊 Estado del Sistema</strong><br>
                ⏰ Hora actual: <strong>{{.CurrentTime}}</strong><br>
                📋 Total notificaciones: <strong>{{.TotalNotifications}}</strong><br>
                ✅ Servicio activo
            </div>

            <div class="form-section">
                <h2>➕ Nueva Notificación</h2>
                <form method="POST" action="/add">
                    <div class="form-group">
                        <label for="hora">⏰ Hora (formato 24h):</label>
                        <input type="time" id="hora" name="hora" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="titulo">📝 Título:</label>
                        <input type="text" id="titulo" name="titulo" placeholder="Ingrese el título de la notificación" required>
                    </div>
                    
                    <div class="form-group">
                        <label for="mensaje">💬 Mensaje:</label>
                        <textarea id="mensaje" name="mensaje" rows="3" placeholder="Ingrese el contenido del mensaje" required></textarea>
                    </div>
                    
                    <div class="checkbox-group">
                        <input type="checkbox" id="recurrente" name="recurrente">
                        <label for="recurrente">🔄 Repetir todos los días</label>
                    </div>
                    
                    <button type="submit" class="btn">💾 Guardar Notificación</button>
                </form>
            </div>

            <div class="notification-list">
                <h2>📋 Notificaciones Programadas</h2>
                {{if .Notifications}}
                    {{range .Notifications}}
                    <div class="notification-item">
                        <div class="notification-time">🕐 {{.Hora}}</div>
                        <div class="notification-title">{{.Titulo}}</div>
                        <div class="notification-message">{{.Mensaje}}</div>
                        <span class="notification-type">
                            {{if .Recurrente}}🔄 Recurrente{{else}}📅 Una vez{{end}}
                        </span>
                    </div>
                    {{end}}
                {{else}}
                    <p style="text-align: center; color: #666; padding: 20px;">
                        📭 No hay notificaciones programadas
                    </p>
                {{end}}
            </div>

            <button onclick="window.close()" class="close-btn">❌ Cerrar</button>
        </div>
    </div>

    <script>
        // Auto-refresh cada 30 segundos
        setTimeout(() => location.reload(), 30000);
        
        // Confirmar al guardar
        document.querySelector('form').addEventListener('submit', function(e) {
            const hora = document.getElementById('hora').value;
            const titulo = document.getElementById('titulo').value;
            
            if (!confirm('¿Guardar notificación "' + titulo + '" para las ' + hora + '?')) {
                e.preventDefault();
            }
        });
    </script>
</body>
</html>
`

type PageData struct {
	CurrentTime        string
	TotalNotifications int
	Notifications      []types.Notificacion
}

// ShowUI inicia la interfaz web
func ShowUI() {
	port := ":8080"
	url := "http://localhost" + port

	// Si el servidor ya está ejecutándose, solo abrir el navegador
	if isRunning {
		log.Println("Servidor web ya está ejecutándose, abriendo navegador...")
		openBrowser(url)
		return
	}

	// Configurar rutas (solo la primera vez)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/add", handleAdd)
	mux.HandleFunc("/favicon", handleFavicon)
	mux.HandleFunc("/logo", handleLogo)

	// Crear servidor
	server = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Marcar como ejecutándose
	isRunning = true

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Iniciando interfaz web en %s", url)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error en servidor web: %v", err)
			isRunning = false
		}
	}()

	// Esperar un momento para que el servidor inicie
	time.Sleep(500 * time.Millisecond)

	// Abrir navegador automáticamente
	openBrowser(url)
}

// handleHome maneja la página principal
func handleHome(w http.ResponseWriter, r *http.Request) {
	notifications := loadNotifications()

	data := PageData{
		CurrentTime:        time.Now().Format("15:04:05"),
		TotalNotifications: len(notifications),
		Notifications:      notifications,
	}

	tmpl := template.Must(template.New("home").Parse(htmlTemplate))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error ejecutando template: %v", err)
	}
}

// handleAdd maneja el formulario de agregar notificación
func handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Crear notificación desde formulario
	notification := types.Notificacion{
		Hora:       r.FormValue("hora"),
		Titulo:     r.FormValue("titulo"),
		Mensaje:    r.FormValue("mensaje"),
		Recurrente: r.FormValue("recurrente") == "on",
	}

	// Validar
	if err := notification.Validate(); err != nil {
		http.Error(w, "❌ Error de validación: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Cargar existentes
	notifications := loadNotifications()

	// Verificar duplicados
	for _, existing := range notifications {
		if existing.Hora == notification.Hora && existing.Titulo == notification.Titulo {
			http.Error(w, "❌ Ya existe una notificación con la misma hora y título", http.StatusBadRequest)
			return
		}
	}

	// Agregar y guardar
	notifications = append(notifications, notification)
	if err := saveNotifications(notifications); err != nil {
		http.Error(w, "❌ Error guardando: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Notificación guardada: %s a las %s", notification.Titulo, notification.Hora)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// loadNotifications carga las notificaciones del archivo
func loadNotifications() []types.Notificacion {
	var notifications []types.Notificacion

	file, err := os.ReadFile(config.DataFile)
	if err != nil {
		return notifications
	}

	if len(file) > 0 {
		json.Unmarshal(file, &notifications)
	}

	return notifications
}

// saveNotifications guarda las notificaciones al archivo
func saveNotifications(notifications []types.Notificacion) error {
	data, err := json.MarshalIndent(notifications, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(config.DataFile, data, 0644)
}

// handleFavicon sirve el favicon
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, config.DataFile[:strings.LastIndex(config.DataFile, "/")]+"/../assets/images/logo.jpeg")
}

// handleLogo sirve el logo principal
func handleLogo(w http.ResponseWriter, r *http.Request) {
	logoPath := "assets/images/logo.jpeg"
	http.ServeFile(w, r, logoPath)
}

// openBrowser abre el navegador web
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("plataforma no soportada")
	}

	if err != nil {
		log.Printf("No se pudo abrir el navegador: %v", err)
		log.Printf("Abra manualmente: %s", url)
	}
}
