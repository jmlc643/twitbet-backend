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
	ErrInsufficientBalance   = errors.New("saldo insuficiente para realizar la apuesta")
	ErrMarketNotActive       = errors.New("el mercado no está activo o está suspendido")
	ErrInvalidBetAmount      = errors.New("el monto de la apuesta es inválido")
	ErrMarketOptionNotFound  = errors.New("opción de mercado no encontrada")
	ErrMarketNotFound        = errors.New("mercado no encontrado")
	ErrMarketOptionBlocked   = errors.New("opción de mercado bloqueada")
	ErrInvalidMarketOptionStatus = errors.New("estado de opción de mercado inválido")
	ErrInvalidMarketType     = errors.New("tipo de mercado inválido")
	ErrMatchNotFound         = errors.New("partido no encontrado")
	ErrUnauthorized          = errors.New("usuario no autorizado para realizar esta acción")
)
