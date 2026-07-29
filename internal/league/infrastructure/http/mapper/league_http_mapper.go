package mapper

import (
	"github.com/google/uuid"
	"github.com/jmlc643/twitbet-backend/internal/league/application/input"
	"github.com/jmlc643/twitbet-backend/internal/league/application/output"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/request"
	"github.com/jmlc643/twitbet-backend/internal/league/infrastructure/http/dto/response"
)

func CreateLeagueRequestToInput(req request.CreateLeagueRequest, adminID uuid.UUID) input.CreateLeagueInput {
	return input.CreateLeagueInput{
		AdminID:        adminID,
		Name:           req.Name,
		InitialBalance: req.InitialBalance,
		MaxRecharges:   req.MaxRecharges,
		HideStandings:  req.HideStandings,
	}
}

func CreateLeagueOutputToResponse(out *output.CreateLeagueOutput) response.CreateLeagueResponse {
	return response.CreateLeagueResponse{
		ID:         out.ID.String(),
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
		Balance:    out.Balance,
	}
}
