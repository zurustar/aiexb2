// backend/internal/handler/reservation_handler.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/service"
)

// ReservationServiceInterface は予約サービスのインターフェース
type ReservationServiceInterface interface {
	CreateReservation(ctx context.Context, req *service.CreateReservationRequest) (*domain.Reservation, error)
	CancelReservation(ctx context.Context, id uuid.UUID, startAt time.Time, userID uuid.UUID) error
}

// ApprovalServiceInterface は承認サービスのインターフェース
type ApprovalServiceInterface interface {
	ApproveReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID) error
	RejectReservation(ctx context.Context, reservationID uuid.UUID, startAt time.Time, approverID uuid.UUID, reason string) error
}

// ReservationHandler は予約関連のHTTPハンドラー
type ReservationHandler struct {
	reservationService ReservationServiceInterface
	approvalService    ApprovalServiceInterface
}

// NewReservationHandler は新しいReservationHandlerを作成します
func NewReservationHandler(
	reservationService ReservationServiceInterface,
	approvalService ApprovalServiceInterface,
) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		approvalService:    approvalService,
	}
}

// RegisterRoutes はルートを登録します
func (h *ReservationHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/events", h.CreateReservation).Methods("POST")
	r.HandleFunc("/api/v1/events/{id}", h.GetReservation).Methods("GET")
	r.HandleFunc("/api/v1/events/{id}", h.CancelReservation).Methods("DELETE")
	r.HandleFunc("/api/v1/events/{id}/approve", h.ApproveReservation).Methods("POST")
	r.HandleFunc("/api/v1/events/{id}/reject", h.RejectReservation).Methods("POST")
}

// CreateReservationRequest は予約作成リクエスト
type CreateReservationRequest struct {
	ResourceIDs []string  `json:"resource_ids"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Timezone    string    `json:"timezone"`
}

// CreateReservation は予約を作成します
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	var req CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// バリデーション
	if req.Timezone == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_TIMEZONE", "Timezone is required")
		return
	}
	if req.StartAt.IsZero() || req.EndAt.IsZero() {
		WriteError(w, http.StatusBadRequest, "INVALID_TIME_RANGE", "StartAt and EndAt are required")
		return
	}
	if req.EndAt.Before(req.StartAt) {
		WriteError(w, http.StatusBadRequest, "INVALID_TIME_RANGE", "EndAt must be after StartAt")
		return
	}

	// UUIDに変換
	resourceIDs := make([]uuid.UUID, len(req.ResourceIDs))
	for i, id := range req.ResourceIDs {
		parsed, err := uuid.Parse(id)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "INVALID_RESOURCE_ID", "Invalid resource ID")
			return
		}
		resourceIDs[i] = parsed
	}

	serviceReq := &service.CreateReservationRequest{
		OrganizerID: session.UserID,
		ResourceIDs: resourceIDs,
		Title:       req.Title,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Timezone:    req.Timezone,
	}

	reservation, err := h.reservationService.CreateReservation(r.Context(), serviceReq)
	if err != nil {
		if err == service.ErrResourceNotAvailable {
			WriteError(w, http.StatusConflict, "RESOURCE_CONFLICT", "One or more resources are not available")
			return
		}
		WriteError(w, http.StatusBadRequest, "CREATE_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, reservation)
}

// GetReservation は予約を取得します
func (h *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid reservation ID")
		return
	}

	// TODO: リポジトリから取得
	// startAtはクエリパラメータから取得して使用
	WriteJSON(w, http.StatusOK, map[string]string{
		"id":      id.String(),
		"message": "Reservation details",
	})
}

// CancelReservation は予約をキャンセルします
func (h *ReservationHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid reservation ID")
		return
	}

	startAtStr := r.URL.Query().Get("start_at")
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_START_AT", "Invalid start_at parameter")
		return
	}

	err = h.reservationService.CancelReservation(r.Context(), id, startAt, session.UserID)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "CANCEL_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Reservation cancelled successfully",
	})
}

// ApproveReservation は予約を承認します
func (h *ReservationHandler) ApproveReservation(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid reservation ID")
		return
	}

	startAtStr := r.URL.Query().Get("start_at")
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_START_AT", "Invalid start_at parameter")
		return
	}

	err = h.approvalService.ApproveReservation(r.Context(), id, startAt, session.UserID)
	if err != nil {
		if err == service.ErrNotApprover {
			WriteError(w, http.StatusForbidden, "FORBIDDEN", "User is not an approver")
			return
		}
		WriteError(w, http.StatusBadRequest, "APPROVE_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Reservation approved successfully",
	})
}

// RejectReservationRequest は予約却下リクエスト
type RejectReservationRequest struct {
	Reason string `json:"reason"`
}

// RejectReservation は予約を却下します
func (h *ReservationHandler) RejectReservation(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(ContextKeySession).(*service.Session)
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid reservation ID")
		return
	}

	startAtStr := r.URL.Query().Get("start_at")
	startAt, err := time.Parse(time.RFC3339, startAtStr)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_START_AT", "Invalid start_at parameter")
		return
	}

	var req RejectReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	err = h.approvalService.RejectReservation(r.Context(), id, startAt, session.UserID, req.Reason)
	if err != nil {
		if err == service.ErrNotApprover {
			WriteError(w, http.StatusForbidden, "FORBIDDEN", "User is not an approver")
			return
		}
		WriteError(w, http.StatusBadRequest, "REJECT_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Reservation rejected successfully",
	})
}
