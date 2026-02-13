package ratelimit

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	mu        sync.Mutex
	voters    = make(map[string]int)
	lastReset = time.Now()
)

func IsLimited(ip string) bool {
	salt := os.Getenv("HASH_SALT")
	hash := sha256.Sum256([]byte(ip + salt))
	ipHash := fmt.Sprintf("%x", hash)

	mu.Lock()
	defer mu.Unlock()

	if time.Since(lastReset) > time.Minute {
		voters = make(map[string]int)
		lastReset = time.Now()
	}

	if voters[ipHash] >= 2 {
		return true
	}
	voters[ipHash]++
	return false
}
