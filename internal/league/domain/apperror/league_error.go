package apperror

import "errors"

var (
	ErrInvalidLeagueName     = errors.New("el nombre de la liga es inválido")
	ErrInvalidInitialBalance = errors.New("el saldo inicial es inválido")
	ErrLeagueNotFound        = errors.New("liga no encontrada")
	ErrInvalidInviteCode     = errors.New("código de invitación inválido")
	ErrUserAlreadyJoined     = errors.New("el usuario ya es participante de esta liga")
	ErrInvalidMatchTitle     = errors.New("el título del partido es inválido")
	ErrInvalidMarketName     = errors.New("el nombre del mercado es inválido")
	ErrInvalidMarketOptions  = errors.New("el mercado debe tener al menos dos opciones")
)
