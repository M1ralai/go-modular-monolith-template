package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.ScheduleService }

func NewHandler(service service.ScheduleService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/schedule/check-conflict", h.CheckConflict).Methods("POST")
	router.HandleFunc("/schedule/free-slots", h.GetFreeSlots).Methods("GET")
	router.HandleFunc("/schedule/blocked-slots", h.GetBlockedSlots).Methods("GET")
	router.HandleFunc("/schedule/blocked-slots", h.CreateBlockedSlot).Methods("POST")
	router.HandleFunc("/schedule/blocked-slots/{id}", h.DeleteBlockedSlot).Methods("DELETE")
	router.HandleFunc("/schedule/generate-events", h.GenerateEvents).Methods("POST")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) CheckConflict(w http.ResponseWriter, r *http.Request) {
	var req dto.ConflictCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	result, err := h.service.CheckConflict(r.Context(), h.getUserID(r), req.Start, req.End)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Çakışma kontrolü başarısız", err.Error())
		return
	}
	utils.WriteJson(w, result, http.StatusOK, "Çakışma kontrolü tamamlandı")
}

func (h *Handler) GetFreeSlots(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	durationStr := r.URL.Query().Get("duration")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz tarih formatı", "YYYY-MM-DD formatında olmalı")
		return
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration < 15 {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz süre", "duration en az 15 dakika olmalı")
		return
	}

	slots, err := h.service.GetFreeSlots(r.Context(), h.getUserID(r), date, duration)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Boş zaman dilimleri alınamadı", err.Error())
		return
	}
	utils.WriteJson(w, slots, http.StatusOK, "Boş zaman dilimleri")
}

func (h *Handler) GetBlockedSlots(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz tarih formatı", "YYYY-MM-DD formatında olmalı")
		return
	}

	slots, err := h.service.GetBlockedSlots(r.Context(), h.getUserID(r), date)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Bloklu zamanlar alınamadı", err.Error())
		return
	}
	utils.WriteJson(w, slots, http.StatusOK, "Bloklu zaman dilimleri")
}

func (h *Handler) CreateBlockedSlot(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBlockedSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	slot, err := h.service.CreateBlockedSlot(r.Context(), h.getUserID(r), &req)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Bloklu zaman oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, slot, http.StatusCreated, "Bloklu zaman dilimi oluşturuldu")
}

func (h *Handler) DeleteBlockedSlot(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.DeleteBlockedSlot(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Bloklu zaman silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Bloklu zaman dilimi silindi")
}

func (h *Handler) GenerateEvents(w http.ResponseWriter, r *http.Request) {
	var req dto.GenerateEventsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	result, err := h.service.GenerateEventsForSchedule(r.Context(), h.getUserID(r), &req)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Etkinlikler oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, result, http.StatusCreated, "Dönem etkinlikleri oluşturuldu")
}
