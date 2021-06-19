package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:		 "theater",
		Description: "Theatre Pack - As the Bard once said, “All the world’s a stage, and all the men and women merely members of our target demographic.” Hark! It’s the Theatre Pack, 30 cards about those things that are like movies but they happen right in front of you.\nWe’d like to thank our parents, our amazing director, and of course, God.\nLin-Manuel Miranda’s favorite pack.\nAll profits fund a grant program for innovative theater and comedy run by the Theater Communications Group.\n\nReleased: June 18, 2018 Icon: Theatre masks",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `Alright everybody, HOLD!
Kelly, why is there %s on my stage?`},
			&PromptCard{Prompt: `Comedy = Tragedy + %s.`},
			&PromptCard{Prompt: `Let's take it from the top, and remember, you are %s. Show me (SAME CARD AGAIN) .`},
			&PromptCard{Prompt: `Match-maker,
match-maker,
make me a match.
Find me %s.`},
			&PromptCard{Prompt: `This season at Manhattan Theatre Club: "Who's afraid of %s?"`},
		},
		Responses: []ResponseCard{
			`A dead salesman.`,
			`A Drama Desk Award for Outstanding Sound Design in a Play.`,
			`A fog machine.`,
			`A play about people I don't like doing things I don't care about.`,
			`Absolutely butchering Sondheim.`,
			`All 59 inches of Kristin Chenoweth.`,
			`An autographed headshot of Nathan Lane.`,
			`Being crushed to death by a stage light.`,
			`Brief male nudity.`,
			`Five miso soups, four seaweed salads, three soy burger dinners, two tofu dog platters, and one pasta with meatless meatballs.`,
			`Forgetting your lines, shitting your pants, and your pants falling down.`,
			`Improv comedy.`,
			`Killing Dad and fucking Mom.`,
			`Linda, 18 but wise beyond her years, achingly beautiful.`,
			`My whole family watching.`,
			`Narcissistic Personality Disorder.`,
			`Problematic depictions of Asian characters.`,
			`Rampant misogyny and sexual harassment.`,
			`Taking a year off to study Japanese puppet theatre.`,
			`The Phantom of the Opera.`,
			`The wickedly talented, one and only, Adele Dazeem.`,
			`This old lady next to me who won't stop farting.`,
			`Two contrasting monologues-- one classical, one contemporary.`,
			`Two men in a horse costume.`,
			`Two tickets to Hamilton.`,
		},
	}

	AddPack(pack)
}
