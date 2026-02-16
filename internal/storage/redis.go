package storage

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func SaveToQueue(encryptedBlob string) error {
	// استفاده از REDIS_URL که در پنل Upstash/Vercel ست شده است
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return err
	}

	client := redis.NewClient(opt)
	defer client.Close()

	// ذخیره‌سازی طبق بند ۳ نقشه راه: فقط LPUSH و بدون متادیتا
	return client.LPush(ctx, "whisper_queue", encryptedBlob).Err()
}
