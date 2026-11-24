package config

import (
	"TheNeonCenturyBot2/storage"
	"strings"

	"gorm.io/gorm"
)

func IsTryingToJailBreak(db *gorm.DB, requests []string, userId string) bool {
	var JailBreakRequestCount = 0
	for _, request := range requests {

		if strings.HasPrefix(request, "429") {
			JailBreakRequestCount++
		}
	}
	var NotJailBreakRequestCount = len(requests) - JailBreakRequestCount
	if JailBreakRequestCount > NotJailBreakRequestCount {
		storage.BanUser(db, userId)
		return true
	}
	return false
}
