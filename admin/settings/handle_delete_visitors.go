package settings

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dracory/api"
)

func (controller *settingsController) handleDeleteVisitors(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	store := controller.opts.Store
	if store == nil {
		api.Respond(w, r, api.Error("Stats store not available"))
		return ""
	}

	var reqBody struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if reqBody.IP == "" {
		api.Respond(w, r, api.Error("IP address is required"))
		return ""
	}

	count, err := store.VisitorDeleteByIP(r.Context(), reqBody.IP)
	if err != nil {
		slog.Error("statsadmin settings: failed to delete visitors by IP", "ip", reqBody.IP, "error", err)
		api.Respond(w, r, api.Error(err.Error()))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData(
		fmt.Sprintf("Deleted %d visitor record(s) for IP %s", count, reqBody.IP),
		map[string]any{
			FieldDeletedCount: count,
			FieldIP:           reqBody.IP,
		},
	))
	return ""
}
