package utils

import (
	"crypto/rand"
	"math/big"
	"nnyd-back/config"
	"time"
)

func GenerateRandomDelayTime() (delay_time time.Duration) {
	duration := int64(config.POST_MAX_DELAY_TIME - config.POST_MIN_DELAY_TIME)

	randomValue := 60
	randomBigValue, err := rand.Int(rand.Reader, big.NewInt(duration))
	if err == nil {
		randomValue = int(randomBigValue.Int64())
	}
	return time.Duration(randomValue + 1 + config.POST_MIN_DELAY_TIME)
}
