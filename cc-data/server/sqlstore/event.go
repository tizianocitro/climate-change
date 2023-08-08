package sqlstore

type URLHashTelemetryEntity struct {
	ID          string `json:"id"`
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	TeamID      string `json:"teamId"`
	TeamName    string `json:"teamName"`
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	URLHash     string `json:"urlHash"`
}
