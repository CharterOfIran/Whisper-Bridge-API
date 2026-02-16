package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"feedback-service/internal/ratelimit"
	"feedback-service/internal/security"
	"feedback-service/internal/storage"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// ۱. لیمیتر (Rate Limit)
	ip := r.Header.Get("X-Forwarded-For")
	if ratelimit.IsLimited(ip) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// ۲. دریافت داده
	var req struct {
		Message  string `json:"message"`
		Honeypot string `json:"hp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// ۳. فیلتر ربات
	if req.Honeypot != "" {
		w.WriteHeader(http.StatusCreated)
		return
	}

	// ۴. پاک‌سازی قبل از رمزنگاری (طبق استاندارد پروژه)
	cleanMsg := security.Clean(req.Message)

	// ۵. رمزنگاری ترکیبی (Hybrid Encryption)

	rawKey := os.Getenv("RSA_PUBLIC_KEY")
	// تبدیل کاراکترهای متنی n\ به خط جدید واقعی برای درک توسط بلوک PEM
	fixedKey := strings.ReplaceAll(rawKey, "\\n", "\n")
	payload, err := security.EncryptHybrid(cleanMsg, fixedKey)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	// تبدیل ساختار رمزنگاری شده به یک رشته واحد (Blob) برای گمنامی کامل
	finalBlob, _ := json.Marshal(payload)

	// ۶. ذخیره‌سازی در Redis
	err = storage.SaveToQueue(string(finalBlob))
	if err != nil {
		http.Error(w, "Failed to queue message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
