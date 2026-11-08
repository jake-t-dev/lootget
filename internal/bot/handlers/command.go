package handlers

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (h *handler) command(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message == nil || m.Author == nil || m.Author.Bot || m.GuildID == "" {
		return
	}
	fields := strings.Fields(m.Content)
	if len(fields) == 0 || !strings.EqualFold(fields[0], "!lootget") {
		return
	}
	if len(fields) != 2 {
		h.reply(s, m, "use `!lootget enable` to set this channel as the output, or `!lootget status` to check the output.")
		return
	}

	switch strings.ToLower(fields[1]) {
	case "status":
		output := h.store.output(m.GuildID)
		if output == "" {
			h.reply(s, m, "no output channel is set. run `!lootget enable` in the desired output channel.")
		} else {
			h.reply(s, m, "the output channel is <#"+output+">.")
		}
	case "enable":
		_, err := s.State.Channel(m.ChannelID)
		if err != nil {
			_, err = s.Channel(m.ChannelID)
		}
		if err := h.store.set(m.GuildID, m.ChannelID); err != nil {
			log.Printf("save channel settings: %v", err)
			h.reply(s, m, "error")
			return
		}
		h.reply(s, m, "yay! this channel is now the output channel.")
	default:
		h.reply(s, m, "use `!lootget enable` to set this channel as the output, or `!lootget status` to check the output.")
	}
}
