// backend/internal/handler/resource_handler.go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/your-org/esms/internal/domain"
	"github.com/your-org/esms/internal/repository"
)

// ResourceHandler はリソース関連のHTTPハンドラー
type ResourceHandler struct {
	resourceRepo repository.ResourceRepository
}

// NewResourceHandler は新しいResourceHandlerを作成します
func NewResourceHandler(resourceRepo repository.ResourceRepository) *ResourceHandler {
	return &ResourceHandler{
		resourceRepo: resourceRepo,
	}
}

// RegisterRoutes はルートを登録します
func (h *ResourceHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/resources", h.ListResources).Methods("GET")
	r.HandleFunc("/api/v1/resources/{id}", h.GetResource).Methods("GET")
	r.HandleFunc("/api/v1/resources", h.CreateResource).Methods("POST")
	r.HandleFunc("/api/v1/resources/{id}", h.UpdateResource).Methods("PUT")
	r.HandleFunc("/api/v1/resources/{id}", h.DeleteResource).Methods("DELETE")
}

// ListResources はリソース一覧を取得します
func (h *ResourceHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	// TODO: リポジトリにListメソッドを追加するか、別の方法で一覧取得
	WriteJSON(w, http.StatusOK, []interface{}{})
}

// GetResource はリソースを取得します
func (h *ResourceHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid resource ID")
		return
	}

	resource, err := h.resourceRepo.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
		return
	}

	WriteJSON(w, http.StatusOK, resource)
}

// CreateResourceRequest はリソース作成リクエスト
type CreateResourceRequest struct {
	Name        string                 `json:"name"`
	Type        domain.ResourceType    `json:"type"`
	Description string                 `json:"description"`
	Location    string                 `json:"location"`
	Capacity    *int                   `json:"capacity"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// CreateResource はリソースを作成します
func (h *ResourceHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	resource := &domain.Resource{
		ID:       uuid.New(),
		Name:     req.Name,
		Type:     req.Type,
		Location: &req.Location,
		Capacity: req.Capacity,
		IsActive: true,
	}

	if err := h.resourceRepo.Create(r.Context(), resource); err != nil {
		WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, resource)
}

// UpdateResource はリソースを更新します
func (h *ResourceHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid resource ID")
		return
	}

	resource, err := h.resourceRepo.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusNotFound, "NOT_FOUND", "Resource not found")
		return
	}

	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	resource.Name = req.Name
	if req.Location != "" {
		resource.Location = &req.Location
	}
	resource.Capacity = req.Capacity

	if err := h.resourceRepo.Update(r.Context(), resource); err != nil {
		WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, resource)
}

// DeleteResource はリソースを削除します
func (h *ResourceHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid resource ID")
		return
	}

	if err := h.resourceRepo.Delete(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Resource deleted successfully",
	})
}
