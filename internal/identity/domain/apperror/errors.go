package apperror

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("El correo electrónico ya está registrado")
	ErrUserNotFound       = errors.New("Usuario no encontrado")
	ErrInvalidCredentials = errors.New("Credenciales inválidas")
	ErrUnauthorized       = errors.New("No autorizado")
	ErrInternal           = errors.New("Error interno del servidor")
	ErrAccountNotVerified = errors.New("Por favor verifica tu cuenta")
	ErrAlreadyVerified    = errors.New("La cuenta ya está verificada")
	ErrOTPExpired         = errors.New("El código ha expirado o no existe")
	ErrOTPInvalid         = errors.New("El código es incorrecto")
	ErrPasswordsDoNotMatch = errors.New("Las contraseñas no coinciden")
	ErrInvalidOldPassword = errors.New("La contraseña actual es incorrecta")
)
