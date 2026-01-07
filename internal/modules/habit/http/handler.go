package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service     service.HabitService
	repo        repository.HabitRepository
	broadcaster *notifService.Broadcaster
	logger      *logger.ZapLogger
}

func NewHandler(service service.HabitService, repo repository.HabitRepository, broadcaster *notifService.Broadcaster, logger *logger.ZapLogger) *Handler {
	return &Handler{
		service:     service,
		repo:        repo,
		broadcaster: broadcaster,
		logger:      logger,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/habits", h.GetAll).Methods("GET")
	router.HandleFunc("/habits/active", h.GetActive).Methods("GET")
	router.HandleFunc("/habits", h.Create).Methods("POST")
	router.HandleFunc("/habits/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/habits/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/habits/{id}", h.Delete).Methods("DELETE")
	router.HandleFunc("/habits/{id}/log", h.LogHabit).Methods("POST")
	router.HandleFunc("/habits/{id}/complete", h.Complete).Methods("POST")
	router.HandleFunc("/habits/{id}/skip", h.Skip).Methods("POST")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	habit, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, habit, http.StatusCreated, "Alışkanlık oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	habit, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		if err.Error() == "habit not found" {
			utils.ReturnError(w, "NOT_FOUND", "Alışkanlık bulunamadı", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, habit, http.StatusOK, "Alışkanlık getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	habits, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlıklar getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, habits, http.StatusOK, "Alışkanlıklar getirildi")
}

func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	habits, err := h.service.GetActive(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Aktif alışkanlıklar getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, habits, http.StatusOK, "Aktif alışkanlıklar getirildi")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdateHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	habit, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, habit, http.StatusOK, "Alışkanlık güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Alışkanlık silindi")
}

func (h *Handler) LogHabit(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.LogHabitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := h.service.LogHabit(r.Context(), id, &req, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık kaydedilemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Alışkanlık kaydedildi")
}

func (h *Handler) Skip(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	userID := h.getUserID(r)

	// CRITICAL FIX: Use service layer directly - service handles WebSocket and jobs internally
	// DO NOT trigger jobs from HTTP handler - this causes duplicate executions
	if err := h.service.SkipHabit(r.Context(), id, userID); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık atlanamadı", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Alışkanlık atlandı")
}

func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.LogHabitRequest
	json.NewDecoder(r.Body).Decode(&req) // Ignore error if empty body

	userID := h.getUserID(r)

	// CRITICAL FIX: Use service layer directly - service handles WebSocket and jobs internally
	// DO NOT trigger jobs from HTTP handler - this causes duplicate executions
	if err := h.service.Complete(r.Context(), id, &req, userID); err != nil {
		if err.Error() == "habit already completed today" {
			utils.ReturnError(w, "BAD_REQUEST", "Bu alışkanlık bugün zaten tamamlandı", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Alışkanlık tamamlanamadı", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Alışkanlık tamamlandı")
}
