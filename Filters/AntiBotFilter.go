package Filters

import (
	"github.com/bwmarrin/discordgo"
)

func AntiBotFilter(d *discordgo.Session, userId string) bool {
	if userId == d.State.User.ID {
		return false
	}
	return true
}
