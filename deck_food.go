package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:		 "food",
		Description: "Food Pack - Originally released as three seperate promo packs sold with popsicles (cherry, coconut, and mango) during Pax Prime 2015. 30 cards about brekky gruffles and syrupy friend chortles. Co-written with Lucky Peach Magazine. Released to be purchased directly from CAH’s website on November 19, 2015, with one card changed.\nIcon: Crisscrossed knife and spoon.",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `Aw babe, your burbs smell like %s.`},
			&PromptCard{Prompt: `Don’t miss Rachel Ray’s hit new show, Cooking with %s.`},
			&PromptCard{Prompt: `Don't miss Rachael Ray's hit new show, Cooking with %s.`},
			&PromptCard{Prompt: `Excuse me, waiter. Could you take this back? This soup tastes like %s.`},
			&PromptCard{Prompt: `I’m Bobby Flay, and if you can’t stand %s. get out of the kitchen!`},
			&PromptCard{Prompt: `It’s not delivery. It’s %s.`},
			&PromptCard{Prompt: `Now on Netflix: Jiro Dreams Of %s.`},
		},
		Responses: []ResponseCard{
			`A belly full of hard-boiled eggs.`,
			`A joyless vegan patty.`,
			`A sobering quantity of chili cheese fries.`,
			`A table for one at The Cheesecake Factory.`,
			`Being emotionally and physically dominated by Gordon Ramsay.`,
			`Clamping down on a gazelle’s jugular and tasting its warm life waters.`,
			`Committing suicide at the Old Country Buffet.`,
			`Father’s forbidden chocolates.`,
			`Going vegetarian and feeling so great all the time.`,
			`Jizz Twinkies.`,
			`Kale farts.`,
			`Kevin Bacon Bits.`,
			`Licking the cake batter off of grandma’s fingers.`,
			`Not knowing what to believe anymore about butter.`,
			`Oreos for dinner.`,
			`Real cheese flavor.`,
			`Soup that’s better than pussy.`,
			`Sucking down thousands of pounds of krill every day.`,
			`Swishing the wine around and sniffing it like a big fancy man.`,
			`The Dial-A-Slice Apple Divider from Williams-Sonoma.`,
			`The Hellman’s Mayonnaise Corporation.`,
			`The hot dog I put in my vagina ten days ago.`,
			`The inaudible screams of carrots.`,
			`What to do with all of this chocolate on my penis.`,
		},
	}

	AddPack(pack)
}
