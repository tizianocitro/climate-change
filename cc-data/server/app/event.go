package app

type UserAddedParams struct {
	TeamID string `json:"teamId"`
	UserID string `json:"userId"`
}

type URLHashTelemetryParams struct {
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	TeamID      string `json:"teamId"`
	TeamName    string `json:"teamName"`
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	URLHash     string `json:"urlHash"`
}
