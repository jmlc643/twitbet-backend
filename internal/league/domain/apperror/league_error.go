package apperror

import "errors"

var (
	ErrInvalidLeagueName     = errors.New("el nombre de la liga es inválido")
	ErrInvalidInitialBalance = errors.New("el saldo inicial es inválido")
	ErrLeagueNotFound        = errors.New("liga no encontrada")
	ErrInvalidInviteCode     = errors.New("código de invitación inválido")
	ErrUserAlreadyJoined     = errors.New("el usuario ya es participante de esta liga")
)
