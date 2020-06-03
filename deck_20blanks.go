package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "20-blanks",
		Description: "Deck of 20 blank response cards",
		// Can I omit "Prompts" like this?
		// Prompts: []*PromptCard{
		// 	&PromptCard{Prompt: `Prompt text here, with %s as a blank`},
		// },
		Responses: []ResponseCard{
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
			`%blank`,
		},
	}

	AddPack(pack)
}
