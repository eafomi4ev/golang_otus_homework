package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/eafomi4ev/golang_otus_homework/hw12_13_14_15_calendar/internal/storage"
)

type PeriodRequest struct {
	period   string
	startDay time.Time
}

type API struct {
	App  *app.App
	Logg app.Logger
}

func (api *API) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	_, _ = w.Write([]byte(`{"hello": "world"}`))
}

func (api *API) AddEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	event := storage.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		api.Logg.Error(fmt.Errorf("api call error: cannot add event: %w", err).Error())
	}

	event, err = api.App.Storage.Add(ctx, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *API) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	event := storage.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		api.Logg.Error(fmt.Errorf("api call error: cannot update event: %w", err).Error())
	}

	event, err = api.App.Storage.Update(ctx, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *API) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	event := storage.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)

	if err != nil {
		api.Logg.Error(fmt.Errorf("api call error: cannot delete event: %w", err).Error())
	}

	err = api.App.Storage.Delete(ctx, event.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok": "true"}`))
}

func (api *API) ListPerPeriod(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var err error

	qp := r.URL.Query()

	pr := PeriodRequest{}
	// err := json.NewDecoder(r.Body).Decode(&pr)
	pr.period = qp["period"][0]
	pr.startDay, err = time.Parse("2006-01-02T15:04:05Z07", qp["startDay"][0])

	if err != nil {
		api.Logg.Error(fmt.Errorf("api call error: cannot delete event: %w", err).Error())
	}

	var events []storage.Event
	switch pr.period {
	case "day":
		events, err = api.App.Storage.ListPerDay(ctx, pr.startDay)
	case "week":
		events, err = api.App.Storage.ListPerWeek(ctx, pr.startDay)
	case "month":
		events, err = api.App.Storage.ListPerMonth(ctx, pr.startDay)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	j, _ := json.Marshal(events)
	err = json.NewEncoder(w).Encode(j)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
