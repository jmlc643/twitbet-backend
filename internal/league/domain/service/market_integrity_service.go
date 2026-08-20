package service

import (
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
)

type MarketIntegrityService struct{}

func NewMarketIntegrityService() *MarketIntegrityService {
	return &MarketIntegrityService{}
}

func (s *MarketIntegrityService) CalculateOverround(options []entity.MarketOption) float64 {
	var overround float64
	for _, opt := range options {
		if opt.CurrentOdds > 0 {
			overround += 1.0 / opt.CurrentOdds
		}
	}
	return overround
}

func (s *MarketIntegrityService) CalculateRebalance(
	currentOptions []entity.MarketOption,
	updatedOdds map[uuid.UUID]float64,
	targetMargin float64,
) (map[uuid.UUID]float64, error) {
	
	targetOverround := 1.0 + targetMargin
	var touchedMass float64
	var untouchedMass float64
	touchedCount := 0

	currentProbs := make(map[uuid.UUID]float64)

	for _, opt := range currentOptions {
		p := 1.0 / opt.CurrentOdds
		currentProbs[opt.ID] = p

		if newOdd, isTouched := updatedOdds[opt.ID]; isTouched {
			touchedMass += 1.0 / newOdd
			touchedCount++
		} else {
			untouchedMass += p
		}
	}

	if touchedCount == len(currentOptions) {
		if touchedMass < targetOverround {
			return nil, &apperror.RebalanceError{
				TouchedMass:    touchedMass,
				MaxTouchedMass: targetOverround,
				Hint:           "Las cuotas ingresadas son demasiado altas y no permiten asegurar un margen de ganancia para la casa.",
			}
		}
		return updatedOdds, nil
	}

	if touchedMass >= targetOverround {
		return nil, &apperror.RebalanceError{
			TouchedMass:    touchedMass,
			MaxTouchedMass: targetOverround,
			Hint:           "Las cuotas editadas son demasiado altas y no dejan margen para recalcular las demás opciones de forma segura.",
		}
	}

	resultOdds := make(map[uuid.UUID]float64)
	massToDistribute := targetOverround - touchedMass

	for _, opt := range currentOptions {
		if newOdd, isTouched := updatedOdds[opt.ID]; isTouched {
			resultOdds[opt.ID] = newOdd
		} else {
			pCurrent := currentProbs[opt.ID]
			pNew := massToDistribute * (pCurrent / untouchedMass)
			newOdd := 1.0 / pNew
			
			if newOdd < 1.01 || newOdd > 100.0 {
				return nil, apperror.ErrOddsOutOfBounds
			}
			resultOdds[opt.ID] = newOdd
		}
	}

	return resultOdds, nil
}
