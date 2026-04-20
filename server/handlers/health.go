package handlers

import (
	"encoding/json"
	"net/http"
)

// Health 返回服务状态，用于确认后端正常运行。
// GET /api/health
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
