package sqlstore

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/tizianocitro/climate-change/cc-data/server/app"
	"github.com/tizianocitro/climate-change/cc-data/server/util"
)

// eventStore is a sql store for events
// Use NewEventStore to create it
type eventStore struct {
	pluginAPI    PluginAPIClient
	store        *SQLStore
	queryBuilder sq.StatementBuilderType

	urlHashTelemetrySelect sq.SelectBuilder
}

// This is a way to implement interface explicitly
var _ app.EventStore = (*eventStore)(nil)

// NewEventStore creates a new store for events service.
func NewEventStore(pluginAPI PluginAPIClient, sqlStore *SQLStore) app.EventStore {
	urlHashTelemetrySelect := sqlStore.builder.
		Select("*").
		From("CSA_URL_HASH_TELEMETRY")

	return &eventStore{
		pluginAPI:              pluginAPI,
		store:                  sqlStore,
		queryBuilder:           sqlStore.builder,
		urlHashTelemetrySelect: urlHashTelemetrySelect,
	}
}

// SaveURLHashTelemetry saves url hash telemetry data
func (e *eventStore) SaveURLHashTelemetry(params app.URLHashTelemetryParams) error {
	tx, err := e.store.db.Beginx()
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}
	defer e.store.finalizeTransaction(tx)

	if _, err := e.store.execBuilder(tx, sq.
		Insert("CSA_URL_HASH_TELEMETRY").
		SetMap(map[string]interface{}{
			"ID":          uuid.New().String(),
			"ChannelID":   params.ChannelID,
			"ChannelName": params.ChannelName,
			"TeamID":      params.TeamID,
			"TeamName":    params.TeamName,
			"UserID":      params.UserID,
			"Username":    params.Username,
			"URLHash":     params.URLHash,
		})); err != nil {
		return errors.Wrap(err, "could not save new url hash telemetry")
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "could not commit transaction")
	}
	return nil
}

func (e *eventStore) ExportURLHashTelemetry() ([]app.URLHashTelemetryParams, error) {
	var urlHashTelemetryEntities []URLHashTelemetryEntity
	err := e.store.selectBuilder(e.store.db, &urlHashTelemetryEntities, e.urlHashTelemetrySelect)
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(app.ErrNotFound, "no url hash telemetry data was found")
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to export url hash telemetry data")
	}

	return e.toURLHashTelemetries(urlHashTelemetryEntities), nil
}

func (e *eventStore) toURLHashTelemetries(urlHashTelemetryEntities []URLHashTelemetryEntity) []app.URLHashTelemetryParams {
	if urlHashTelemetryEntities == nil {
		return nil
	}
	urlHashTelemetries := make([]app.URLHashTelemetryParams, 0, len(urlHashTelemetryEntities))
	for _, u := range urlHashTelemetryEntities {
		urlHashTelemetries = append(urlHashTelemetries, e.toURLHashTelemetry(u))
	}
	return urlHashTelemetries
}

func (e *eventStore) toURLHashTelemetry(urlHashTelemetryEntity URLHashTelemetryEntity) app.URLHashTelemetryParams {
	urlHashTelemetry := app.URLHashTelemetryParams{}
	err := util.Convert(urlHashTelemetryEntity, &urlHashTelemetry)
	if err != nil {
		return app.URLHashTelemetryParams{}
	}
	return urlHashTelemetry
}
