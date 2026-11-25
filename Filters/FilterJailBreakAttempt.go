package Filters

import (
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func FilterJailBreakAttempt(db *gorm.DB, channelId string, d *discordgo.Session, filteredQuestion string, userId string) bool {
	if strings.Contains(filteredQuestion, "429") {
		err := storage.AddRequestToUser(db, userId, filteredQuestion)
		d.ChannelMessageSend(channelId, "Your question is not related to lore")
		if err != nil {
			log.Println("AddRequestToUser error:", err)
			return false
		}
		return false
	}
	return true
}
