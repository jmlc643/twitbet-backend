package mapper

import (
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
)

func CreateLeagueRequestToInput(req request.CreateLeagueRequest, OwnerID uuid.UUID) input.CreateLeagueInput {
	return input.CreateLeagueInput{
		OwnerID:        OwnerID,
		Name:           req.Name,
		InitialBalance: req.InitialBalance,
		MaxRecharges:   req.MaxRecharges,
		HideStandings:  req.HideStandings,
		MinBetsToQualify: req.MinBetsToQualify,
	}
}

func CreateLeagueOutputToResponse(out *output.CreateLeagueOutput) response.CreateLeagueResponse {
	return response.CreateLeagueResponse{
		ID:         out.ID.String(),
		Slug:       out.Slug,
		InviteCode: out.InviteCode,
	}
}

func JoinLeagueRequestToInput(req request.JoinLeagueRequest, userID uuid.UUID) input.JoinLeagueInput {
	return input.JoinLeagueInput{
		UserID:     userID,
		InviteCode: req.InviteCode,
	}
}

func JoinLeagueOutputToResponse(out *output.JoinLeagueOutput) response.JoinLeagueResponse {
	return response.JoinLeagueResponse{
		LeagueID:   out.LeagueID.String(),
		LeagueName: out.LeagueName,
		Slug:       out.Slug,
		Balance:    out.Balance,
	}
}

func GetUserLeaguesOutputToResponse(out *output.GetUserLeaguesOutput) response.GetUserLeaguesResponse {
	var leagues []response.LeagueSummaryResponse
	for _, l := range out.Leagues {
		leagues = append(leagues, response.LeagueSummaryResponse{
			LeagueID:         l.LeagueID.String(),
			Slug:             l.Slug,
			Name:             l.Name,
			Role:             string(l.Role),
			ParticipantCount: l.ParticipantCount,
			Balance:          l.Balance,
		})
	}
	return response.GetUserLeaguesResponse{
		Leagues: leagues,
	}
}

func GetLeagueDetailsOutputToResponse(out output.GetLeagueDetailsOutput) response.GetLeagueDetailsResponse {
	var participants []response.ParticipantRankingResponse
	for _, p := range out.Participants {
		participants = append(participants, response.ParticipantRankingResponse{
			ParticipantID:  p.ParticipantID.String(),
			UserID:         p.UserID.String(),
			Username:       p.Username,
			ProfilePicture: p.ProfilePicture,
			Balance:        p.Balance,
			Position:       p.Position,
			Role:           p.Role,
		})
	}
	return response.GetLeagueDetailsResponse{
		LeagueID:         out.LeagueID.String(),
		Slug:             out.Slug,
		Name:             out.Name,
		OwnerID:          out.OwnerID.String(),
		InitialBalance:   out.InitialBalance,
		MaxRecharges:     out.MaxRecharges,
		IsRankingVisible: out.IsRankingVisible,
		InviteCode:       out.InviteCode,
		CreatedAt:        out.CreatedAt,
		ParticipantsCount: out.ParticipantsCount,
		Participants:     participants,
	}
}

func GetLeagueLeaderboardOutputToResponse(out *output.GetLeagueLeaderboardOutput) response.GetLeagueLeaderboardResponse {
	var leaderboard []response.LeaderboardEntryResponse
	for _, entry := range out.Leaderboard {
		leaderboard = append(leaderboard, response.LeaderboardEntryResponse{
			ParticipantID:  entry.ParticipantID.String(),
			UserID:         entry.UserID.String(),
			Username:       entry.Username,
			ProfilePicture: entry.ProfilePicture,
			Balance:        entry.Balance,
			Position:       entry.Position,
			Role:           entry.Role,
			IsUnranked:     entry.IsUnranked,
		})
	}
	if leaderboard == nil {
		leaderboard = []response.LeaderboardEntryResponse{}
	}
	return response.GetLeagueLeaderboardResponse{
		LeagueID:         out.LeagueID.String(),
		Status:           out.Status,
		HideStandings:    out.HideStandings,
		MinBetsToQualify: out.MinBetsToQualify,
		Leaderboard:      leaderboard,
	}
}