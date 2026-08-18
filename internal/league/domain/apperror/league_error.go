package apperror

import "errors"

var (
	ErrInvalidLeagueName     = errors.New("El nombre de la liga es inválido")
	ErrInvalidInitialBalance = errors.New("El saldo inicial es inválido")
	ErrLeagueNotFound        = errors.New("Liga no encontrada")
	ErrInvalidInviteCode     = errors.New("Código de invitación inválido")
	ErrUserAlreadyJoined     = errors.New("El usuario ya es participante de esta liga")
	ErrInvalidMatchTitle     = errors.New("El título del partido es inválido")
	ErrInvalidMarketName     = errors.New("El nombre del mercado es inválido")
	ErrInvalidMarketOptions  = errors.New("El mercado debe tener al menos dos opciones")
	ErrInsufficientBalance   = errors.New("Saldo insuficiente para realizar la apuesta")
	ErrMarketNotActive       = errors.New("El mercado no está activo o está suspendido")
	ErrInvalidBetAmount      = errors.New("El monto de la apuesta es inválido")
	ErrMarketOptionNotFound  = errors.New("Opción de mercado no encontrada")
	ErrMarketNotFound        = errors.New("Mercado no encontrado")
	ErrMarketOptionBlocked   = errors.New("Opción de mercado bloqueada")
	ErrInvalidMarketOptionStatus = errors.New("Estado de opción de mercado inválido")
	ErrInvalidMarketType     = errors.New("Tipo de mercado inválido")
	ErrMatchNotFound         = errors.New("Partido no encontrado")
	ErrUnauthorized          = errors.New("Usuario no autorizado para realizar esta acción")
	ErrBetNotFound           = errors.New("Apuesta no encontrada")
	ErrArbitrageMarket   = errors.New("Mercado de arbitraje detectado (overround < 1)")
	ErrOddsOutOfBounds   = errors.New("El rebalanceo deja cuotas fuera del rango permitido")
)

type RebalanceError struct {
	TouchedMass    float64
	MaxTouchedMass float64
	Hint           string
}

func (e *RebalanceError) Error() string {
	return "Rebalanceo imposible"
}

type OddsChangedError struct {
	CurrentOdds float64
}

func (e *OddsChangedError) Error() string {
	return "ODDS_CHANGED"
}
