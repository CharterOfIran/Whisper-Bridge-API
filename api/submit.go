package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"feedback-service/internal/github"
	"feedback-service/internal/ratelimit"
	"feedback-service/internal/security"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// ... سایر تنظیمات CORS

	// ۱. استفاده از ماژول لیمیتر
	ip := r.Header.Get("X-Forwarded-For")
	if ratelimit.IsLimited(ip) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// ۲. دریافت داده
	var req struct{ Message, Honeypot string }
	json.NewDecoder(r.Body).Decode(&req)

	// ۳. فیلتر ربات
	if req.Honeypot != "" {
		return
	}

	// ۴. استفاده از ماژول امنیتی
	cleanMsg := security.Clean(req.Message)

	// ۵. ارسال به گیت‌هاب
	client := github.NewClient(os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_OWNER"), os.Getenv("GITHUB_REPO"))
	client.SaveComment(cleanMsg)

	w.WriteHeader(http.StatusCreated)
}
