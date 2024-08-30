package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/withoutsecondd/kamibooking/internal"
	"github.com/withoutsecondd/kamibooking/model"
	"github.com/withoutsecondd/kamibooking/service"
	"net/http"
	"strconv"
)

type Handler struct {
	ReservationS service.ReservationService
}

func (h *Handler) SetupRoutes(r chi.Router) {
	r.Get("/reservations/{roomId}", h.getReservations)
	r.Post("/reservations/", h.createReservation)
}

func (h *Handler) getReservations(w http.ResponseWriter, r *http.Request) {
	roomIdStr := chi.URLParam(r, "roomId")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	reservations, err := h.ReservationS.GetReservations(int64(roomId))
	if err != nil {
		var hErr internal.HttpError
		ok := errors.As(err, &hErr)
		if ok {
			http.Error(w, hErr.Err.Error(), hErr.Code)
			return
		} else {
			http.Error(w, "unexpected error occurred", http.StatusInternalServerError)
			return
		}
	}

	response, err := json.Marshal(reservations)
	if err != nil {
		http.Error(w, "failed to form json response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createReservation(w http.ResponseWriter, r *http.Request) {
	var reservation model.Reservation

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.ReservationS.PostReservation(&reservation)
	if err != nil {
		var hErr internal.HttpError
		ok := errors.As(err, &hErr)
		if ok {
			http.Error(w, hErr.Err.Error(), hErr.Code)
			return
		} else {
			http.Error(w, "unexpected error occurred", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}
