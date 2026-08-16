package usecase

import (
	"context"
	"errors"

	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetMatchDetailsUseCase struct {
	matchRepo repository.MatchRepository
}

func NewGetMatchDetailsUseCase(matchRepo repository.MatchRepository) *GetMatchDetailsUseCase {
	return &GetMatchDetailsUseCase{
		matchRepo: matchRepo,
	}
}

func (uc *GetMatchDetailsUseCase) Execute(ctx context.Context, slug string) (output.GetMatchDetailsOutput, error) {
	match, err := uc.matchRepo.GetMatchBySlug(ctx, slug)
	if err != nil {
		return output.GetMatchDetailsOutput{}, err
	}
	if match == nil {
		return output.GetMatchDetailsOutput{}, errors.New("Partido no encontrado")
	}

	markets, err := uc.matchRepo.GetMarketsByMatchID(ctx, match.ID)
	if err != nil {
		return output.GetMatchDetailsOutput{}, err
	}

	return output.GetMatchDetailsOutput{
		ID:        match.ID,
		LeagueID:  match.LeagueID,
		Slug:      match.Slug,
		Title:     match.Title,
		StartTime: match.StartTime,
		Status:    match.Status,
		CreatedAt: match.CreatedAt,
		UpdatedAt: match.UpdatedAt,
		Markets:   markets,
	}, nil
}
