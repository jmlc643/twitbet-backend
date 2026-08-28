package response

type UserStatsResponse struct {
	Leagues       int `json:"leagues"`
	Wins          int `json:"wins"`
	Effectiveness int `json:"effectiveness"`
}

type UserResponse struct {
	ID        string             `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	AvatarURL *string            `json:"avatar_url"`
	CreatedAt string             `json:"created_at"`
	Stats     *UserStatsResponse `json:"stats,omitempty"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
