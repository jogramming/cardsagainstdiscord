package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "40-blanks",
		Description: "Deck of 40 blank response cards",
		// Can I define an empty "Prompts" set like this?
		Prompts: []*PromptCard{},
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
