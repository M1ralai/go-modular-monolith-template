package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.GoalService }

func NewHandler(service service.GoalService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/goals", h.GetAll).Methods("GET")
	router.HandleFunc("/goals", h.Create).Methods("POST")
	router.HandleFunc("/goals/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/goals/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/goals/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	goal, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Hedef oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, goal, http.StatusCreated, "Hedef oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	goal, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "NOT_FOUND", "Hedef bulunamadı", err.Error())
		return
	}
	utils.WriteJson(w, goal, http.StatusOK, "Hedef getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	goals, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Hedefler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, goals, http.StatusOK, "Hedefler getirildi")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdateGoalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	goal, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Hedef güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, goal, http.StatusOK, "Hedef güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Hedef silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Hedef silindi")
}
