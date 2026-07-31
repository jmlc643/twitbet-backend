package output

import "github.com/jmlc643/twitbet-backend/internal/league/domain/entity"

type GetUserLeaguesOutput struct {
	Leagues []entity.LeagueSummary
}
