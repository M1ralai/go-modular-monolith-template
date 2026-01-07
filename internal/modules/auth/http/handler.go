package http

import (
	"encoding/json"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service service.AuthService
}

func NewHandler(service service.AuthService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", h.Register).Methods("POST")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	response, err := h.service.Login(r.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			utils.ReturnError(w, "UNAUTHORIZED", "Geçersiz e-posta veya şifre", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Giriş yapılamadı", err.Error())
		return
	}

	utils.WriteJson(w, response, http.StatusOK, "Giriş başarılı")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek formatı", err.Error())
		return
	}

	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}

	response, err := h.service.Register(r.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.ReturnError(w, "BAD_REQUEST", "Bu e-posta adresi zaten kullanımda", err.Error())
			return
		}
		utils.ReturnError(w, "INTERNAL_ERROR", "Kayıt oluşturulamadı", err.Error())
		return
	}

	utils.WriteJson(w, response, http.StatusCreated, "Kayıt başarılı")
}
