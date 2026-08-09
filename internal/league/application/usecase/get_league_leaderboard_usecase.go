package usecase

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/entity"
	"github.com/jmlc643/twitbet-backend/internal/league/domain/repository"
)

type GetLeagueLeaderboardUseCase struct {
	repo repository.LeagueRepository
}

func NewGetLeagueLeaderboardUseCase(repo repository.LeagueRepository) *GetLeagueLeaderboardUseCase {
	return &GetLeagueLeaderboardUseCase{repo: repo}
}

func (uc *GetLeagueLeaderboardUseCase) Execute(ctx context.Context, leagueID, userID uuid.UUID) (*output.GetLeagueLeaderboardOutput, error) {
	league, err := uc.repo.GetLeagueByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}
	if league == nil {
		return nil, errors.New("liga no encontrada")
	}

	hideBalances := league.HideStandings && league.Status != "FINALIZED"
	var outputLeaderboard []output.LeaderboardEntryOutput

	if hideBalances && userID != league.OwnerID {
		return &output.GetLeagueLeaderboardOutput{
			LeagueID:         league.ID,
			Status:           league.Status,
			HideStandings:    league.HideStandings,
			MinBetsToQualify: league.MinBetsToQualify,
			Leaderboard:      []output.LeaderboardEntryOutput{},
		}, nil
	}

	rawLeaderboard, err := uc.repo.GetLeaderboard(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	var ranked []entity.LeaderboardEntry
	var unranked []entity.LeaderboardEntry

	for _, entry := range rawLeaderboard {
		if entry.BetsCount < league.MinBetsToQualify {
			unranked = append(unranked, entry)
		} else {
			ranked = append(ranked, entry)
		}
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Balance > ranked[j].Balance
	})

	position := 1
	for _, entry := range ranked {
		var balance *float64
		if !hideBalances {
			bal := entry.Balance
			balance = &bal
		}

		pos := position
		outputLeaderboard = append(outputLeaderboard, output.LeaderboardEntryOutput{
			ParticipantID:  entry.ParticipantID,
			UserID:         entry.UserID,
			Username:       entry.Username,
			ProfilePicture: entry.ProfilePicture,
			Balance:        balance,
			Position:       &pos,
			Role:           entry.Role,
			IsUnranked:     false,
		})
		position++
	}

	for _, entry := range unranked {
		var balance *float64
		if !hideBalances {
			bal := entry.Balance
			balance = &bal
		}

		outputLeaderboard = append(outputLeaderboard, output.LeaderboardEntryOutput{
			ParticipantID:  entry.ParticipantID,
			UserID:         entry.UserID,
			Username:       entry.Username,
			ProfilePicture: entry.ProfilePicture,
			Balance:        balance,
			Position:       nil,
			Role:           entry.Role,
			IsUnranked:     true,
		})
	}

	return &output.GetLeagueLeaderboardOutput{
		LeagueID:         league.ID,
		Status:           league.Status,
		HideStandings:    league.HideStandings,
		MinBetsToQualify: league.MinBetsToQualify,
		Leaderboard:      outputLeaderboard,
	}, nil
}
