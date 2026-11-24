package cogs

import (
	"TheNeonCenturyBot2/config"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Help(d *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == d.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandHelp) {

		if config.IsAdmin(d, m.Author.ID, m.ChannelID) {

			d.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title:       d.State.User.Username,
				Description: "Prefix for all commands is " + config.CommandPrefix + "\n" + " Current commands : " + "`" + config.CommandHelp + "`" + ", `" + config.CommandAsk + "`" + "\n Admin commands: " + "`" + config.CommandAllowChannelNew + "`" + ", " + "`" + config.CommandDisallowChannelNew + "`" + ", `" + config.CommandUnbanUser + "`" + ", `" + config.CommandBanUser + "`" + ", `" + config.CommandCleanHistory + "`",
				Color:       0x5b209a,
			})
			return
		}

		d.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       d.State.User.Username,
			Description: "Prefix for all commands is " + config.CommandPrefix + " Current commands : `" + config.CommandHelp + "`, `" + config.CommandAsk + "`",
			Color:       0x5b209a,
		})
	}

}
