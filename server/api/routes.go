package api

import (
	"net/http"

	"github.com/zihao-liu-qs/qs-blog/server/handlers"
)

// RegisterRoutes 挂载所有 API 路由。
// 按功能分组，未实现的路由留有注释，方便后续扩展。
func RegisterRoutes(mux *http.ServeMux) {

	// ── 基础 ──────────────────────────────────────────────
	mux.HandleFunc("GET /api/health", handlers.Health)

	// ── 支付（待实现） ─────────────────────────────────────
	// mux.HandleFunc("POST /api/checkout", handlers.CreateCheckout)
	// mux.HandleFunc("POST /api/webhook",  handlers.StripeWebhook)

	// ── License（待实现） ──────────────────────────────────
	// mux.HandleFunc("GET  /api/license/verify", handlers.VerifyLicense)

	// ── 下载（待实现） ─────────────────────────────────────
	// mux.HandleFunc("GET  /api/download/:product", handlers.ServeDownload)
}
