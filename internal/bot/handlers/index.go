package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var itemPattern = regexp.MustCompile(`\[([^\[\]\n]{1,100})\]`)

type handler struct {
	store  *channelStore
	client *http.Client
}

func AddHandlers(dg *discordgo.Session) (*discordgo.Session, error) {
	path := os.Getenv("OUTPUT_CHANNELS_FILE")
	if path == "" {
		path = filepath.Join("data", "output-channels.json")
	}
	store, err := newChannelStore(path)
	if err != nil {
		return nil, fmt.Errorf("load channel settings: %w", err)
	}
	h := &handler{store: store, client: &http.Client{Timeout: 10 * time.Second}}
	dg.AddHandler(h.command)
	dg.AddHandler(h.message)
	return dg, nil
}

func (h *handler) reply(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:         content,
		Reference:       m.Reference(),
		AllowedMentions: &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{}, RepliedUser: false},
	})
	if err != nil {
		log.Printf("send reply in channel %s: %v", m.ChannelID, err)
	}
}

func (h *handler) message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message == nil || m.Author == nil || m.Author.Bot || m.GuildID == "" {
		return
	}
	fields := strings.Fields(m.Content)
	if len(fields) > 0 && strings.EqualFold(fields[0], "!lootget") {
		return
	}
	matches := itemPattern.FindAllStringSubmatch(m.Content, -1)
	if len(matches) == 0 {
		return
	}
	output := h.store.output(m.GuildID)
	if output == "" {
		log.Printf("no output channel configured for guild %s", m.GuildID)
		return
	}
	var lines []string
	seen := make(map[string]bool)
	for _, match := range matches {
		name := strings.TrimSpace(match[1])
		key := strings.ToLower(name)
		if name == "" || seen[key] {
			continue
		}
		seen[key] = true
		items, err := lookupItems(h.client, name)
		if err != nil {
			log.Printf("Error from WowHead: %v", err)
			continue
		}
		for _, item := range items {
			lines = append(lines, item.url())
		}
	}
	var batch string
	for _, line := range lines {
		if len(batch)+len(line)+1 > 2000 {
			h.sendOutput(s, output, batch)
			batch = ""
		}
		if batch != "" {
			batch += "\n"
		}
		batch += line
	}
	if batch != "" {
		h.sendOutput(s, output, batch)
	}
}

func (h *handler) sendOutput(s *discordgo.Session, channel, content string) {
	_, err := s.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
		Content:         content,
		AllowedMentions: &discordgo.MessageAllowedMentions{Parse: []discordgo.AllowedMentionType{}},
	})
	if err != nil {
		log.Printf("send output in channel %s: %v", channel, err)
	}
}
