package apperror

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("el correo electrónico ya está registrado")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUnauthorized       = errors.New("no autorizado")
	ErrInternal           = errors.New("error interno del servidor")
)