package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "pride",
		Description: "Pride Pack - 30 cards that completely encapsulate the queer experience",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `Excuse me, straight man, but %s isn't for you, STRAIGHT MAN.`},
			&PromptCard{Prompt: `GOD HATES %s!`},
			&PromptCard{Prompt: `If you can't love yourself, how the hell you gonna love %s?`},
			&PromptCard{Prompt: `We're here! We're %s! Get used to it!`},
			&PromptCard{Prompt: `YAAAAAAS! You are serving me %s realness!`},
		},
		Responses: []ResponseCard{
			`30 shirtless bears emerging from the fog.`,
			`A 6-hour conversation on gender and queer theory.`,
			`A big black dick strapped to a frail white body.`,
			`A genderless hole.`,
			`A messy bitch who lives for drama.`,
			`A Subaru.`,
			`A twink in a bounce house.`,
			`All the different kinds of lesbians.`,
			`Black, white, Puerto Rican, and Chinese boys.`,
			`Britney, bitch!`,
			`Getting your ass ate.`,
			`Having your titties sucked while sucking on titties.`,
			`Licking that pussy right.`,
			`Marsha P. Johnson, the trans woman of color who may have thrown the first brick at Stonewall.`,
			`Older fitness gays.`,
			`Peeing in a bathroom.`,
			`Poppers and lube.`,
			`PrEP.`,
			`Repeatedly coming out as bisexual.`,
			`Talking, laughing, loving, breathing, fighting, fucking, crying, drinking, riding, winning, losing, cheating, kissing, thinking, dreaming.`,
			`Telling Heather she can't pull off that top.`,
			`The careless cunt who left a water ring on my credenza.`,
			`The pan-ethnic, gender-fluid children of the future.`,
			`Those cheekbones, honey.`,
			`Whatever straight people do for fun.`,
		},
	}

	AddPack(pack)
}
