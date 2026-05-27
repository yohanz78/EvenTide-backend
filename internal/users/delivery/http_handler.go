package delivery

import (
	"encoding/json"
	"net/http"

	"EventTide-backend/internal/domain"
)

// Handler hanya memegang kontrak Usecase
type UserHandler struct {
	UserUC domain.UserUsecase
}

// Struct khusus untuk memetakan input JSON dari klien
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Fungsi ini yang akan didaftarkan ke Router
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Format JSON salah"})
		return
	}

	// Lempar ke layer Usecase
	user, err := h.UserUC.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusConflict) // 409 Conflict (Misal: email sudah ada)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Berhasil
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registrasi berhasil",
		"data":    user, // Field password otomatis hilang karena tag json:"-" di Domain
	})
}
