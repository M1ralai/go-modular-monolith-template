package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.PersonService }

func NewHandler(service service.PersonService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/people", h.GetAll).Methods("GET")
	router.HandleFunc("/people/search", h.Search).Methods("GET")
	router.HandleFunc("/people/tag/{tag}", h.SearchByTag).Methods("GET")
	router.HandleFunc("/people", h.Create).Methods("POST")
	router.HandleFunc("/people/{id}", h.GetByID).Methods("GET")
	router.HandleFunc("/people/{id}", h.Update).Methods("PUT", "PATCH")
	router.HandleFunc("/people/{id}", h.Delete).Methods("DELETE")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	person, err := h.service.Create(r.Context(), &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Kişi oluşturulamadı", err.Error())
		return
	}
	utils.WriteJson(w, person, http.StatusCreated, "Kişi oluşturuldu")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	person, err := h.service.GetByID(r.Context(), id, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "NOT_FOUND", "Kişi bulunamadı", err.Error())
		return
	}
	utils.WriteJson(w, person, http.StatusOK, "Kişi getirildi")
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	people, err := h.service.GetAll(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Kişiler getirilemedi", err.Error())
		return
	}
	utils.WriteJson(w, people, http.StatusOK, "Kişiler getirildi")
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	people, err := h.service.Search(r.Context(), h.getUserID(r), query)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Arama yapılamadı", err.Error())
		return
	}
	utils.WriteJson(w, people, http.StatusOK, "Arama sonuçları")
}

func (h *Handler) SearchByTag(w http.ResponseWriter, r *http.Request) {
	tag := mux.Vars(r)["tag"]
	people, err := h.service.SearchByTag(r.Context(), h.getUserID(r), tag)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Tag araması yapılamadı", err.Error())
		return
	}
	utils.WriteJson(w, people, http.StatusOK, "Tag sonuçları")
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req dto.UpdatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	person, err := h.service.Update(r.Context(), id, &req, h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Kişi güncellenemedi", err.Error())
		return
	}
	utils.WriteJson(w, person, http.StatusOK, "Kişi güncellendi")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.Delete(r.Context(), id, h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Kişi silinemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Kişi silindi")
}
