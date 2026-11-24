package storage

import (
	"errors"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type AllowedChannels struct {
	ChannelId string `gorm:"primaryKey"`
}

type User struct {
	UserId       string `gorm:"primaryKey"`
	IsBanned     bool
	IsAdmin      bool
	LastRequests []string `gorm:"serializer:json"`
}

func DbStart() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage/bot.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Database error: ", err)
	}

	if err := db.AutoMigrate(&AllowedChannels{}, &User{}); err != nil {
		log.Fatal("Migration Error: ", err)
	}

	return db
}

func AllowChannel(db *gorm.DB, channelId string) error {
	return db.FirstOrCreate(&AllowedChannels{ChannelId: channelId}).Error
}

func IsChannelAllowed(db *gorm.DB, channelId string) bool {
	var count int64
	if err := db.Model(&AllowedChannels{}).Where("channel_id = ?", channelId).Count(&count).Error; err != nil {
		log.Println("IsChannelAllowed error:", err)
		return false
	}
	return count > 0
}

func DisAllowChannel(db *gorm.DB, channelId string) error {
	return db.Where("channel_id = ?", channelId).Delete(&AllowedChannels{}).Error
}

func IsUserExist(db *gorm.DB, userId string) bool {
	var count int64
	if err := db.Model(&User{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		log.Println("IsUserExist error:", err)
		return false
	}
	return count > 0
}

func GetUserRequests(db *gorm.DB, userId string) ([]string, error) {
	var user User
	if err := db.Select("last_requests").First(&user, "user_id = ?", userId).Error; err != nil {
		log.Println("GetUserRequests error:", err)
		return nil, err
	}
	return user.LastRequests, nil
}

func IsUserBanned(db *gorm.DB, userId string) bool {
	var user User
	// Select("is_banned") не тянет и не парсит JSON из last_requests на каждое сообщение
	err := db.Select("is_banned").First(&user, "user_id = ?", userId).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("IsUserBanned error:", err)
		}
		return false
	}
	return user.IsBanned
}

func AddUser(db *gorm.DB, userId string, isBanned, isAdmin bool) error {
	user := User{
		UserId:       userId,
		IsBanned:     isBanned,
		IsAdmin:      isAdmin,
		LastRequests: []string{},
	}
	return db.Create(&user).Error
}

func AddRequestToUser(db *gorm.DB, userId, request string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, "user_id = ?", userId).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("AddRequestToUser error:", err)
			}
			return err
		}

		user.LastRequests = append(user.LastRequests, request)

		return tx.Model(&user).Select("LastRequests").Updates(&user).Error
	})
}

func BanUser(db *gorm.DB, userId string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, "user_id = ?", userId).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("BanUser error:", err)
			}
			return err
		}

		if user.IsBanned {
			return errors.New("user is already banned")
		}
		if user.IsAdmin {
			return errors.New("you can't ban admin")
		}

		return tx.Model(&user).Update("is_banned", true).Error
	})
}

func UnBanUser(db *gorm.DB, userId string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, "user_id = ?", userId).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("UnBanUser error:", err)
			}
			return err
		}

		if !user.IsBanned {
			return errors.New("user is not banned")
		}
		if user.IsAdmin {
			return errors.New("you can't unban admin")
		}

		return tx.Model(&user).Update("is_banned", false).Error
	})
}

func CleanUserQuestionsHistory(db *gorm.DB, userId string) error {
	res := db.Model(&User{}).
		Where("user_id = ?", userId).
		Select("LastRequests").
		Updates(&User{LastRequests: []string{}})

	if res.Error != nil {
		log.Println("CleanUserQuestionsHistory error:", res.Error)
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
