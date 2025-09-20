# RGDEV Notificador - Sitio Web Promocional

Este directorio contiene el sitio web promocional para RGDEV Notificador, desplegado automáticamente en GitHub Pages.

## 🌐 Sitio Web

El sitio web está disponible en: `https://DEMG-DEV.github.io/RGDEV.Service.Notify`

## 📁 Estructura

```
web/
├── index.html      # Página principal
├── styles.css      # Estilos CSS
├── script.js       # JavaScript interactivo
└── README.md       # Este archivo
```

## ✨ Características

- **Diseño Responsivo**: Optimizado para desktop, tablet y móvil
- **Animaciones Suaves**: Efectos de scroll y transiciones elegantes
- **SEO Optimizado**: Meta tags y estructura semántica
- **Carga Rápida**: CSS y JS optimizados
- **Accesibilidad**: Cumple estándares de accesibilidad web

## 🚀 Despliegue Automático

El sitio se despliega automáticamente a GitHub Pages cuando:

- Se hace push a la rama `main`
- Los archivos en la carpeta `web/` son modificados

## 🛠️ Desarrollo Local

Para probar el sitio localmente:

1. Navega a la carpeta `web/`
2. Abre `index.html` en tu navegador
3. O usa un servidor local:

   ```bash
   # Con Python
   python -m http.server 8000
   
   # Con Node.js
   npx serve .
   ```

## 📱 Secciones del Sitio

1. **Hero**: Presentación principal con call-to-action
2. **Características**: Funcionalidades principales del sistema
3. **Cómo Funciona**: Pasos para usar el sistema
4. **Descarga**: Enlaces de descarga y GitHub
5. **Footer**: Enlaces adicionales y información

## 🎨 Personalización

Para personalizar el sitio:

1. **Colores**: Modifica las variables CSS en `:root` en `styles.css`
2. **Contenido**: Edita el texto en `index.html`
3. **Animaciones**: Ajusta los efectos en `script.js`

## 📊 Analytics

Para agregar analytics (Google Analytics, etc.):

1. Agrega el código de tracking en `index.html`
2. Configura los eventos de conversión en `script.js`
