# Guía de Contribución

## 🤝 ¿Cómo Contribuir?

¡Gracias por tu interés en contribuir a RGDEV.Service.Notify! Este documento te guiará a través del proceso.

## 📋 Formas de Contribuir

### 🐛 Reportar Bugs

1. Verifica que el bug no haya sido reportado ya
2. Usa el template de issues para bugs
3. Incluye información del sistema y pasos para reproducir

### ✨ Sugerir Nuevas Características

1. Revisa las características planificadas en el README
2. Abre un issue con la etiqueta "enhancement"
3. Describe claramente el caso de uso

### 🔧 Contribuir Código

1. Fork el repositorio
2. Crea una rama para tu feature: `git checkout -b feature/nueva-caracteristica`
3. Haz commits con mensajes descriptivos
4. Abre un Pull Request

## 📝 Estándares de Código

### Go Style Guide

- Sigue las [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Usa `gofmt` para formatear el código
- Documenta funciones públicas con comentarios

### Estructura de Commits

```
tipo(alcance): descripción breve

Descripción más detallada si es necesaria.

- Cambio específico 1
- Cambio específico 2
```

Tipos válidos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

### Testing

- Agrega tests para nuevas funcionalidades
- Ejecuta `go test ./...` antes de enviar PR
- Mantén cobertura de código > 80%

## 🔍 Proceso de Review

1. **Automated Checks**: CI debe pasar
2. **Code Review**: Al menos un maintainer debe aprobar
3. **Testing**: Verificación manual si es necesario
4. **Merge**: Squash and merge preferido

## 📚 Recursos Útiles

- [Documentación de Go](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

## ❓ ¿Preguntas?

Si tienes dudas, no dudes en:

- Abrir un issue con la etiqueta "question"
- Contactar a los maintainers

¡Esperamos tus contribuciones! 🎉
