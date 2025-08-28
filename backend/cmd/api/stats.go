package main

import (
	"net/http"
	"time"

	"github.com/A-Jama01/spaced-repetition-app/internal/store"
)


func (app *app) listStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := app.getUserIDFromContext(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	queryString := r.URL.Query()
	deckName := app.readString(queryString, "deck_name", "")
	timeZone := app.readString(queryString, "tz", "UTC")

	_, err = time.LoadLocation(timeZone)
	if err != nil {
		timeZone = "UTC"
		app.logger.Printf("Failed time zone")
	}

	queryParams := store.StatsQueryParams{
		UserID: userID,
		DeckName: deckName,
		TimeZone: timeZone,
	}

	var response struct {
		ReviewCount int64 `json:"review_count"`
		Retention float64 `json:"retention"`
		HeatMap []*store.ReviewCell `json:"heatmap"`
		Forecasts []*store.DueForecast `json:"forecasts"`
	}

	response.ReviewCount, err = app.store.Logs.GetCount(ctx, queryParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response.Retention, err = app.store.Logs.GetRetention(ctx, queryParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response.HeatMap, err = app.store.Logs.GetHeatMap(ctx, queryParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response.Forecasts, err = app.store.Cards.GetDueForecast(ctx, queryParams)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	

	err = app.writeJSON(w, http.StatusOK, envelope{"stats": response}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
