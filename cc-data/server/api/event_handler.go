package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tizianocitro/climate-change/cc-data/server/app"
)

const TelemetryBasePath = "/url_hash_telemetry"

// EventHandler is the API handler.
type EventHandler struct {
	*ErrorHandler
	eventService *app.EventService
}

// EventHandler returns a new event api handler
func NewEventHandler(router *mux.Router, eventService *app.EventService) *EventHandler {
	handler := &EventHandler{
		ErrorHandler: &ErrorHandler{},
		eventService: eventService,
	}

	platformRouter := router.PathPrefix("/events").Subrouter()
	platformRouter.HandleFunc("/user_added", withContext(handler.userAdded)).Methods(http.MethodPost)
	platformRouter.HandleFunc("/url_hash_telemetry", withContext(handler.exportURLHashTelemetry)).Methods(http.MethodGet)
	platformRouter.HandleFunc("/url_hash_telemetry", withContext(handler.saveURLHashTelemetry)).Methods(http.MethodPost)

	return handler
}

func (h *EventHandler) userAdded(c *Context, w http.ResponseWriter, r *http.Request) {
	var params app.UserAddedParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode user added payload", err)
		return
	}
	if err := h.eventService.UserAdded(params); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to handle user added", err)
		return
	}
	ReturnJSON(w, "", http.StatusOK)
}

func (h *EventHandler) saveURLHashTelemetry(c *Context, w http.ResponseWriter, r *http.Request) {
	var params app.URLHashTelemetryParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to decode url hash telemetry payload", err)
		return
	}
	if err := h.eventService.SaveURLHashTelemetry(params); err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to save url hash telemetry", err)
		return
	}
	ReturnJSON(w, "", http.StatusOK)
}

// e.g. http://www.isislab.it:8065/plugins/climate-change-data/api/v0/events/url_hash_telemetry
func (h *EventHandler) exportURLHashTelemetry(c *Context, w http.ResponseWriter, r *http.Request) {
	data, err := h.eventService.ExportURLHashTelemetry()
	if err != nil {
		h.HandleErrorWithCode(w, c.logger, http.StatusBadRequest, "unable to export url hash telemetry", err)
		return
	}
	ReturnAttachment(w, data, "text/csv", "hyperlink_telemetry.csv", h.writeAttachment)
}

func (h *EventHandler) writeAttachment(w http.ResponseWriter, data app.WritableData) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	csvData := data.(app.CSVData)
	if err := writer.Write(csvData.Header); err != nil {
		http.Error(w, fmt.Sprintf("Error writing CSV header: %v", err), http.StatusInternalServerError)
	}
	for _, record := range csvData.Rows {
		if err := writer.Write(record); err != nil {
			http.Error(w, fmt.Sprintf("Error writing CSV row: %v", err), http.StatusInternalServerError)
		}
	}
}
