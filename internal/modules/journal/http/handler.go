package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.JournalService }

func NewHandler(service service.JournalService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/journal", h.GetAll).Methods("GET")
	router.HandleFunc("/journal", h.Create).Methods("POST")
	router.HandleFunc("/journal/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/journal/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/journal/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	entry, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Günlük oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, entry, http.StatusCreated, "Günlük oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	entry, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "NOT_FOUND", "Günlük bulunamadı", err.Error())
		return
	}
	utils.WriteJson(w, entry, http.StatusOK, "Günlük getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Günlükler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, entries, http.StatusOK, "Günlükler getirildi")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	entry, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Günlük güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, entry, http.StatusOK, "Günlük güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Günlük silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Günlük silindi")
}
