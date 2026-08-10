package apperror

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("el correo electrónico ya está registrado")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUnauthorized       = errors.New("no autorizado")
	ErrInternal           = errors.New("error interno del servidor")
	ErrAccountNotVerified = errors.New("por favor verifica tu cuenta")
	ErrAlreadyVerified    = errors.New("la cuenta ya está verificada")
	ErrOTPExpired         = errors.New("el código ha expirado o no existe")
	ErrOTPInvalid         = errors.New("el código es incorrecto")
	ErrPasswordsDoNotMatch = errors.New("las contraseñas no coinciden")
	ErrInvalidOldPassword = errors.New("la contraseña actual es incorrecta")
)
