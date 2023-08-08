package app

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v6/plugin"

	"github.com/tizianocitro/climate-change/cc-data/server/util"
)

type EventService struct {
	api   plugin.API
	store EventStore
}

// NewEventService returns a new platform config service
func NewEventService(api plugin.API, store EventStore) *EventService {
	return &EventService{
		api:   api,
		store: store,
	}
}

func (s *EventService) UserAdded(params UserAddedParams) error {
	s.api.LogInfo("Params on user added", "params", params)
	channels, err := s.api.GetPublicChannelsForTeam(params.TeamID, 0, 200)
	if err != nil {
		return fmt.Errorf("couldn't get public channels for team %s", params.TeamID)
	}
	for _, channel := range channels {
		if _, err := s.api.AddChannelMember(channel.Id, params.UserID); err != nil {
			return fmt.Errorf("couldn't add channel %s to user %s", channel.Id, params.UserID)
		}
	}
	return nil
}

func (s *EventService) SaveURLHashTelemetry(params URLHashTelemetryParams) error {
	s.api.LogInfo("Url hash telemetry params", "params", params)
	if err := s.store.SaveURLHashTelemetry(params); err != nil {
		return fmt.Errorf("couldn't save url hash telemetry due to %s", err.Error())
	}
	return nil
}

func (s *EventService) ExportURLHashTelemetry() (WritableData, error) {
	telemetries, err := s.store.ExportURLHashTelemetry()
	if err != nil {
		return CSVData{}, fmt.Errorf("couldn't export url hash telemetry due to %s", err.Error())
	}
	telemetriesRows := [][]string{}
	for _, telemetry := range telemetries {
		telemetriesRows = append(telemetriesRows, []string{
			telemetry.ChannelID,
			telemetry.ChannelName,
			telemetry.TeamID,
			telemetry.TeamName,
			telemetry.UserID,
			telemetry.Username,
			telemetry.URLHash,
		})
	}
	return CSVData{
		Header: util.ConvertStructToKeys(telemetries[0]),
		Rows:   telemetriesRows,
	}, nil
}
