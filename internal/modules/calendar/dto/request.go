package dto

type GoogleConnectRequest struct {
	RedirectURI string `json:"redirect_uri,omitempty"`
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state,omitempty"`
}

type SyncRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google apple"`
}

type QueueSyncRequest struct {
	EventID int    `json:"event_id" validate:"required"`
	Action  string `json:"action" validate:"required,oneof=create update delete"`
}
