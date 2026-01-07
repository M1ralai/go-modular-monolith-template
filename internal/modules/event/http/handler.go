package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.EventService }

func NewHandler(service service.EventService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/events", h.GetAll).Methods("GET")
	router.HandleFunc("/events/calendar", h.GetByDateRange).Methods("GET")
	router.HandleFunc("/events", h.Create).Methods("POST")
	router.HandleFunc("/events/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/events/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/events/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	event, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlik oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, event, http.StatusCreated, "Etkinlik oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	event, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "NOT_FOUND", "Etkinlik bulunamadı", err.Error())
		return
	}
	utils.WriteJson(w, event, http.StatusOK, "Etkinlik getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlikler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, events, http.StatusOK, "Etkinlikler getirildi")
}

func (h *Handler) GetByDateRange(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	start, _ := time.Parse("2006-01-02", startStr)
	end, _ := time.Parse("2006-01-02", endStr)
	if start.IsZero() {
		start = time.Now().AddDate(0, -1, 0)
	}
	if end.IsZero() {
		end = time.Now().AddDate(0, 1, 0)
	}
	events, err := h.service.GetByDateRange(r.Context(), h.getUserID(r), start, end)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlikler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, events, http.StatusOK, "Takvim etkinlikleri getirildi")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	event, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlik güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, event, http.StatusOK, "Etkinlik güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlik silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Etkinlik silindi")
}
