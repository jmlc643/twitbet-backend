package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/apperror"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GrantLeagueBonusUseCase struct {
	leagueRepo repository.LeagueRepository
}

func NewGrantLeagueBonusUseCase(leagueRepo repository.LeagueRepository) *GrantLeagueBonusUseCase {
	return &GrantLeagueBonusUseCase{
		leagueRepo: leagueRepo,
	}
}

func (uc *GrantLeagueBonusUseCase) Execute(ctx context.Context, userID, leagueID uuid.UUID, amount float64) error {
	league, err := uc.leagueRepo.GetLeagueByID(ctx, leagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return apperror.ErrLeagueNotFound
	}
	if league.OwnerID != userID {
		participant, err := uc.leagueRepo.GetParticipant(ctx, leagueID, userID)
		if err != nil || !participant.IsAdmin {
			return errors.New("Solo el administrador de la liga puede otorgar bonos")
		}
	}

	participants, err := uc.leagueRepo.GetParticipantsByLeagueID(ctx, leagueID)
	if err != nil {
		return err
	}

	var bonuses []*entity.ParticipantBonus
	for _, p := range participants {
		bonuses = append(bonuses, entity.NewParticipantBonus(leagueID, p.ParticipantID, amount))
	}

	return uc.leagueRepo.CreateParticipantBonuses(ctx, bonuses)
}
