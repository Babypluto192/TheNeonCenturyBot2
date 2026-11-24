package config

import (
	"TheNeonCenturyBot2/storage"

	"gorm.io/gorm"
)

func IsBanned(db *gorm.DB, userid string) bool {
	return storage.IsUserBanned(db, userid)
}
