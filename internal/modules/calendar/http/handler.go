package http

import (
	"encoding/json"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/service"
	"github.com/gorilla/mux"
)

type Handler struct{ service service.CalendarService }

func NewHandler(service service.CalendarService) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/calendar/google/connect", h.GetGoogleAuthURL).Methods("POST")
	router.HandleFunc("/calendar/google/callback", h.HandleGoogleCallback).Methods("GET")
	router.HandleFunc("/calendar/google/disconnect", h.DisconnectGoogle).Methods("POST")
	router.HandleFunc("/calendar/google/sync", h.SyncGoogle).Methods("POST")
	router.HandleFunc("/calendar/status", h.GetSyncStatus).Methods("GET")
	router.HandleFunc("/calendar/integrations", h.GetIntegrations).Methods("GET")
}

func (h *Handler) getUserID(r *http.Request) int {
	return utils.GetUserIDFromContext(r.Context())
}

func (h *Handler) GetGoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	authURL, err := h.service.GetGoogleAuthURL(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Google OAuth yapılandırılmamış", err.Error())
		return
	}
	utils.WriteJson(w, dto.AuthURLResponse{AuthURL: authURL}, http.StatusOK, "Yetkilendirme URL'i oluşturuldu")
}

func (h *Handler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.ReturnError(w, "BAD_REQUEST", "Kod gerekli", "code parametresi eksik")
		return
	}

	integration, err := h.service.HandleGoogleCallback(r.Context(), h.getUserID(r), code)
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Google bağlantısı başarısız", err.Error())
		return
	}
	utils.WriteJson(w, integration, http.StatusOK, "Google Calendar bağlandı")
}

func (h *Handler) DisconnectGoogle(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DisconnectGoogle(r.Context(), h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Bağlantı kesilemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Google Calendar bağlantısı kesildi")
}

func (h *Handler) SyncGoogle(w http.ResponseWriter, r *http.Request) {
	if err := h.service.SyncGoogle(r.Context(), h.getUserID(r)); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Senkronizasyon başarısız", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Google Calendar senkronize edildi")
}

func (h *Handler) GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.service.GetSyncStatus(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Durum alınamadı", err.Error())
		return
	}
	utils.WriteJson(w, status, http.StatusOK, "Senkronizasyon durumu")
}

func (h *Handler) GetIntegrations(w http.ResponseWriter, r *http.Request) {
	integrations, err := h.service.GetIntegrations(r.Context(), h.getUserID(r))
	if err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Entegrasyonlar alınamadı", err.Error())
		return
	}
	utils.WriteJson(w, integrations, http.StatusOK, "Takvim entegrasyonları")
}

func (h *Handler) QueueSync(w http.ResponseWriter, r *http.Request) {
	var req dto.QueueSyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ReturnError(w, "BAD_REQUEST", "Geçersiz istek", err.Error())
		return
	}
	if err := validation.Get().Struct(req); err != nil {
		utils.ReturnError(w, "VALIDATION_ERROR", "Doğrulama hatası", validation.FormatErr(err))
		return
	}
	if err := h.service.QueueSync(r.Context(), h.getUserID(r), req.EventID, req.Action); err != nil {
		utils.ReturnError(w, "INTERNAL_ERROR", "Senkronizasyon kuyruğa eklenemedi", err.Error())
		return
	}
	utils.WriteJson(w, nil, http.StatusOK, "Senkronizasyon kuyruğa eklendi")
}
