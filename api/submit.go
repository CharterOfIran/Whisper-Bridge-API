package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"feedback-service/internal/github"
	"feedback-service/internal/ratelimit"
	"feedback-service/internal/security"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// ۱. استفاده از ماژول لیمیتر
	ip := r.Header.Get("X-Forwarded-For")
	if ratelimit.IsLimited(ip) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// ۲. دریافت داده
	var req struct{ Message, Honeypot string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// ۳. فیلتر ربات
	if req.Honeypot != "" {
		// برای ربات‌ها وانمود می‌کنیم که موفق بودیم!
		w.WriteHeader(http.StatusCreated)
		return
	}

	// ۴. استفاده از ماژول امنیتی
	cleanMsg := security.Clean(req.Message)

	// ۵. ارسال به گیت‌هاب
	client := github.NewClient(
		os.Getenv("GITHUB_TOKEN"),
		os.Getenv("GITHUB_OWNER"),
		os.Getenv("GITHUB_REPO"),
	)

	// دریافت نتیجه و خطا از تابع
	_, err := client.SaveComment(cleanMsg)
	if err != nil {
		// چاپ خطا در لاگ‌های ورسل برای عیب‌یابی
		fmt.Printf("GitHub Error: %v\n", err)
		http.Error(w, "Failed to save feedback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
