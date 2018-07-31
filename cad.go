package cardsagainstdiscord

import (
	"fmt"
	"github.com/jonas747/discordgo"
	"strings"
)

var Packs = make(map[string]*CardPack)

func AddPack(name string, pack *CardPack) {
	// Count picks
	for _, v := range pack.Prompts {
		numPicks := strings.Count(v.Prompt, "%s")
		if numPicks == 0 {
			v.Prompt += " %s"
			v.NumPick = 1
		} else {
			v.NumPick = numPicks
		}
	}

	Packs[name] = pack
}

type CardPack struct {
	Name      string
	Prompts   []*PromptCard
	Responses []ResponseCard
}

type PromptCard struct {
	Prompt  string
	NumPick int
}

var (
	EscaperReplacer = strings.NewReplacer("*", "\\*", "_", "\\_")
)

func (p *PromptCard) PlaceHolder() string {
	s := strings.Replace(p.Prompt, "%s", "_____", -1)
	s = strings.Replace(s, "%%", `%`, -1)

	s = EscaperReplacer.Replace(s)

	return s
}

func (p *PromptCard) WithCards(cards interface{}) string {
	args := make([]interface{}, p.NumPick)
	switch t := cards.(type) {
	case []string:
		for i, v := range t {
			args[i] = "**" + v + "**"
		}
	case []ResponseCard:
		for i, v := range t {
			args[i] = "**" + v + "**"
		}
	}

	s := fmt.Sprintf(p.Prompt, args...)
	// s = EscaperReplacer.Replace(s)
	return s
}

type ResponseCard string

type SessionProvider interface {
	SessionForGuild(guildID int64) *discordgo.Session
}

type StaticSessionProvider struct {
	Session *discordgo.Session
}

func (sp *StaticSessionProvider) SessionForGuild(guildID int64) *discordgo.Session {
	return sp.Session
}
