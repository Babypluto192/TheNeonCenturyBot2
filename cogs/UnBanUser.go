package cogs

import (
	"TheNeonCenturyBot2/Filters"
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func ParseUserId(arg string) (string, bool) {
	id := arg

	if strings.HasPrefix(id, "<@") && strings.HasSuffix(id, ">") {
		id = strings.TrimPrefix(id, "<@")
		id = strings.TrimSuffix(id, "!")
		id = strings.TrimSuffix(id, ">")
	}

	if len(id) < 17 || len(id) > 20 {
		return "", false
	}
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return "", false
	}

	return id, true
}

func UnBanUser(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandUnbanUser) {
			args := strings.Fields(m.Content)
			if len(args) < 2 {
				d.ChannelMessageSend(m.ChannelID, "Wrong number of arguments use >unban userid or >unban @UserPing")
				return
			}
			if !Filters.IsAdminFilter(m.ChannelID, d, m.Author.ID) {
				return
			}

			if !Filters.IsChannelAllowedFilter(db, m.ChannelID, d) {
				return
			}

			targetID, ok := ParseUserId(args[1])
			if !ok {
				d.ChannelMessageSend(m.ChannelID, "Wrong user id.")
				return
			}
			if targetID == m.Author.ID {
				d.ChannelMessageSend(m.ChannelID, "You can't unban yourself [baka](https://klipy.com/gifs/nanami-aoyama-nanami).")
				return
			}
			if targetID == d.State.User.ID {
				d.ChannelMessageSend(m.ChannelID, "You can't unban the bot.")
				return
			}

			if _, err := d.User(targetID); err != nil {
				d.ChannelMessageSend(m.ChannelID, "This user probably doesn't exists lol")
				return
			}
			if !storage.IsUserExist(db, targetID) {
				d.ChannelMessageSend(m.ChannelID, "You are not banned")
				return
			}

			err := storage.UnBanUser(db, targetID)
			if err != nil {
				log.Println(err)
				return
			}

			d.ChannelMessageSend(m.ChannelID, "User has been unbanned.")
		}
	}

}
