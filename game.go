package cardsagainstdiscord

import (
	"fmt"
	"github.com/jonas747/discordgo"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

type GameState int

/*

Game Loop:      If a player has enough wins to win the game
pregame ------------------<|
   |                       |
PreRoundDelay ------------<|
   |                       |
PickingResponses           |
   |                       |
Picking winner             |
   |>----------------------^

*/

const (
	GameStatePreGame          GameState = 0 // Before the game starts
	GameStatePreRoundDelay    GameState = 1 // Countdown before a roundn starts
	GameStatePickingResponses GameState = 2 // Players are picking responses for the prompt card
	GameStatePickingWinner    GameState = 3 // Cardzar is picking the winning response
)

const (
	PreRoundDelayDuration = time.Second * 10
	PickResponseDuration  = time.Second * 45
	PickWinnerDuration    = time.Second * 60
	GameExpireAfter       = time.Second * 180
)

var (
	CardSelectionEmojis = []string{
		"🇦", // A
		"🇧", // B
		"🇨", // C
		"🇩", // D
		"🇪", // E
		"🇫", // F
		"🇬", // G
		"🇭", // H
		"🇮", // I
		"🇯", // J
		"🇰", // K
	}

	JoinEmoji  = "➕"
	LeaveEmoji = "➖"
	PlayEmoji  = "▶"
)

type Game struct {
	sync.RWMutex
	// Never chaged
	Manager *GameManager
	Session *discordgo.Session

	// The main channel this game resides in, never changes
	MasterChannel int64
	// The server the game resides in, never changes
	GuildID int64

	// The user that created this game
	GameMaster int64

	// The current cardzar
	CurrentCardCzar int64

	Packs   []string
	Players []*Player

	PlayerLimit int

	State        GameState
	StateEntered time.Time

	LastAction time.Time

	CurrentPropmpt *PromptCard

	LastMenuMessage int64

	Responses []*PickedResonse
}

type PickedResonse struct {
	Player     *Player
	Selections []ResponseCard
}

func (g *Game) Created() {
	g.LastAction = time.Now()

	msg, err := g.Session.ChannelMessageSend(g.MasterChannel, fmt.Sprintf("Game is created! React with %s to join and %s to leave, the game master can start the game with %s",
		JoinEmoji, LeaveEmoji, PlayEmoji))
	if err != nil {
		return
	}

	g.LastMenuMessage = msg.ID

	go g.addCommonMenuReactions(msg.ID, true)
}

// AddPlayer attempts to add a player to the game, if it fails (hit the limit for example) then it returns false
func (g *Game) AddPlayer(id int64, username string) bool {
	g.Lock()
	defer g.Unlock()

	return g.addPlayer(id, username)
}

func (g *Game) addPlayer(id int64, username string) bool {
	if g.PlayerLimit <= len(g.Players) {
		return false
	}

	// Create a userchannel and cache it for use later
	channel, err := g.Session.UserChannelCreate(id)
	if err != nil {
		return false
	}

	p := &Player{
		ID:       id,
		Username: username,
		Channel:  channel.ID,
		Cards:    g.getRandomPlayerCards(8),
	}

	g.Players = append(g.Players, p)

	go g.sendAnnouncment(fmt.Sprintf("<@%d> Joined the game! (%d/%d)", id, len(g.Players), g.PlayerLimit), false)
	return true
}

func (g *Game) RemovePlayer(id int64) {
	g.Lock()
	defer g.Unlock()
	g.removePlayer(id)
}

func (g *Game) removePlayer(id int64) {
	found := false
	for i, v := range g.Players {
		if v.ID == id {
			g.Players = append(g.Players[:i], g.Players[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return
	}

	go g.sendAnnouncment(fmt.Sprintf("<@%d> Left the game (%d/%d)", id, len(g.Players), g.PlayerLimit), false)

	if g.CurrentCardCzar == id && g.State != GameStatePreGame && g.State != GameStatePreRoundDelay {
		g.nextRound()
	}
}

func (g *Game) setState(state GameState) {
	g.State = state
	g.StateEntered = time.Now()
	log.Println("Set ", g.MasterChannel, " state to ", state)
}

func (g *Game) nextRound() {
	g.setState(GameStatePreRoundDelay)
}

func (g *Game) getRandomResponseCard() ResponseCard {
	totalAvailableCards := 0
	for _, v := range g.Packs {
		totalAvailableCards += len(Packs[v].Responses)
	}

	cardIndex := rand.Intn(totalAvailableCards)
	for _, v := range g.Packs {
		packResponses := Packs[v].Responses
		if len(packResponses) > cardIndex {
			return packResponses[cardIndex]
		}

		cardIndex -= len(packResponses)
	}

	panic("Should never get here")
	return ""
}

func (g *Game) getRandomPlayerCards(num int) []ResponseCard {
	result := make([]ResponseCard, 0, num)

	if len(g.Packs) < 1 {
		return result
	}

	for len(result) < num {
		card := g.getRandomResponseCard()

		// Duplicate
		for _, existing := range result {
			if existing == card {
				continue
			}
		}

		result = append(result, card)
	}

	return result
}

func (g *Game) sendAnnouncment(msg string, allPlayers bool) {
	session := g.Manager.SessionProvider.SessionForGuild(g.GuildID)
	if session == nil {
		return
	}

	if allPlayers {
		for _, v := range g.Players {
			go func(channel int64) {
				g.Session.ChannelMessageSend(channel, msg)
			}(v.Channel)
		}
	}

	g.Session.ChannelMessageSend(g.MasterChannel, msg)
}

func (g *Game) Tick() {
	g.Lock()
	defer g.Unlock()

	log.Println(time.Since(g.LastAction))
	if time.Since(g.LastAction) > GameExpireAfter || len(g.Players) < 1 {
		g.gameExpired()
		return
	}

	switch g.State {
	case GameStatePreGame:
		return
	case GameStatePreRoundDelay:
		if time.Since(g.StateEntered) > PreRoundDelayDuration {
			g.startRound()
			return
		}
	case GameStatePickingResponses:
		allPlayersDone := true
		oneResponsePicked := false
		for _, v := range g.Players {
			if v.ID == g.CurrentCardCzar || !v.Playing {
				continue
			}

			if len(v.SelectedCards) < g.CurrentPropmpt.NumPick {
				allPlayersDone = false
			} else {
				oneResponsePicked = true
			}
		}

		if allPlayersDone || time.Since(g.StateEntered) >= PickResponseDuration {
			if oneResponsePicked {
				g.donePickingResponses()
			} else {
				// No one picked any cards...?
				g.sendAnnouncment("No one picked any cards, going to next round", false)
				g.nextRound()
			}
		}
	case GameStatePickingWinner:
		if time.Since(g.StateEntered) >= PickWinnerDuration {
			g.cardzarExpired()
		}
	}

}

func (g *Game) startRound() {
	if len(g.Players) < 1 {
		g.setState(GameStatePreGame)
		return
	}

	// Remove previous selected cards from players decks
	for _, v := range g.Responses {
		for _, sel := range v.Selections {
			for i, c := range v.Player.Cards {
				if c == sel {
					v.Player.Cards = append(v.Player.Cards[:i], v.Player.Cards[i+1:]...)
					break
				}
			}
		}
	}

	g.Responses = nil

	for _, v := range g.Players {
		v.Playing = true
		v.SelectedCards = nil
	}

	// Pick next cardzar
	g.pickNextCardzar()

	lastPick := 1
	if g.CurrentPropmpt != nil {
		lastPick = g.CurrentPropmpt.NumPick
	}

	// Pick random propmpt
	g.CurrentPropmpt = g.randomPrompt()

	// Give each player a random card (if they're below 10 cards)
	g.giveEveryoneCards(lastPick)

	// Present the board
	g.presentStartRound()

	g.setState(GameStatePickingResponses)

}

func (g *Game) pickNextCardzar() {
	var next int64
	for _, v := range g.Players {
		if v.ID == g.CurrentCardCzar || !v.Playing {
			continue
		}

		if v.ID > g.CurrentCardCzar && v.ID > next {
			next = v.ID
		}
	}

	if next == 0 {
		next = g.Players[0].ID
	}

	g.CurrentCardCzar = next
}

func (g *Game) randomPrompt() *PromptCard {
	totalAvailablePromps := 0
	for _, v := range g.Packs {
		totalAvailablePromps += len(Packs[v].Prompts)
	}

	cardIndex := rand.Intn(totalAvailablePromps)
	for _, v := range g.Packs {
		packPrompts := Packs[v].Prompts
		if len(packPrompts) > cardIndex {
			return &packPrompts[cardIndex]
		}

		cardIndex -= len(packPrompts)
	}

	panic("Should never get here")
	return nil
}

func (g *Game) giveEveryoneCards(num int) {
	for _, p := range g.Players {
		if p.ID == g.CurrentCardCzar || !p.Playing || len(p.Cards) >= 10 {
			continue
		}

		if num+len(p.Cards) > 10 {
			num = 10 - len(p.Cards)
		}

		cards := g.getRandomPlayerCards(num)
		p.Cards = append(p.Cards, cards...)
	}

}

func (g *Game) presentStartRound() {

	cardCzarUsername := ""
	for _, player := range g.Players {
		if player.ID == g.CurrentCardCzar {
			cardCzarUsername = player.Username
		}

		go func(p *Player) {
			p.PresentBoard(g.Session, g.CurrentPropmpt, g.CurrentCardCzar)
		}(player)
	}

	// Present the main board
	msg := fmt.Sprintf("---------------------\nNext round started! **%s** is the Card Czar!\n\n**%s**\n\nCheck your dm for your card and pick your cards there, return here afterwards, you got 45 seconds.",
		cardCzarUsername, strings.Replace(g.CurrentPropmpt.Prompt, "%s", "____", -1))
	g.Session.ChannelMessageSend(g.MasterChannel, msg)
}

func (g *Game) donePickingResponses() {
	// Send a message to players that missed the round
	for _, v := range g.Players {
		if !v.Playing || v.ID == g.CurrentCardCzar {
			continue
		}

		if len(v.SelectedCards) < g.CurrentPropmpt.NumPick {
			go g.Session.ChannelMessageSend(v.Channel, fmt.Sprintf("You didn't respond in time... winner is being picked in <#%d>", g.MasterChannel))
			v.SelectedCards = nil
			continue
		}

		selections := make([]ResponseCard, 0, len(v.SelectedCards))
		for _, sel := range v.SelectedCards {
			selections = append(selections, v.Cards[sel])
		}

		g.Responses = append(g.Responses, &PickedResonse{
			Player:     v,
			Selections: selections,
		})
	}

	// Shows all the picks in both dm's and main channel
	g.presentPickedResponseCards()
	g.setState(GameStatePickingWinner)
}

func (g *Game) presentPickedResponseCards() {
	content := fmt.Sprintf("Cards have been picked, pick the best one(s) <@%d>!\n\n**%s**:\n\n", g.CurrentCardCzar, strings.Replace(g.CurrentPropmpt.Prompt, "%s", "____", -1))

	for i, v := range g.Responses {
		strCards := make([]interface{}, len(v.Selections))
		for j, resp := range v.Selections {
			strCards[j] = "**" + string(resp) + "**"
		}

		content += CardSelectionEmojis[i] + ": " + fmt.Sprintf(g.CurrentPropmpt.Prompt, strCards...) + "\n\n"
	}

	msg, err := g.Session.ChannelMessageSend(g.MasterChannel, content)
	if err != nil {
		return
	}

	numOptions := len(g.Responses)
	go func() {
		for i := 0; i < numOptions; i++ {
			g.Session.MessageReactionAdd(g.MasterChannel, msg.ID, CardSelectionEmojis[i])
		}

		g.addCommonMenuReactions(msg.ID, false)
	}()

	g.LastMenuMessage = msg.ID
}

func (g *Game) cardzarExpired() {
	msg, err := g.Session.ChannelMessageSend(g.MasterChannel, fmt.Sprintf("<%d> didn't pick a winner in %d seconds, skipping round...", g.CurrentCardCzar, int(PickWinnerDuration.Seconds())))
	if err == nil {
		go g.addCommonMenuReactions(msg.ID, false)
	}

	g.setState(GameStatePreRoundDelay)
}

func (g *Game) gameExpired() {
	g.Session.ChannelMessageSend(g.MasterChannel, "CAH Game expired")
	go g.Manager.RemoveGame(g.MasterChannel)
}

func (g *Game) addCommonMenuReactions(mID int64, play bool) {
	g.Session.MessageReactionAdd(g.MasterChannel, mID, JoinEmoji)
	g.Session.MessageReactionAdd(g.MasterChannel, mID, LeaveEmoji)
	if play {
		g.Session.MessageReactionAdd(g.MasterChannel, mID, PlayEmoji)
	}
}

func (g *Game) HandleRectionAdd(ra *discordgo.MessageReactionAdd) {
	g.Lock()
	defer g.Unlock()

	var player *Player
	for _, v := range g.Players {
		if v.ID == ra.UserID {
			player = v
			break
		}
	}

	if ra.MessageID == g.LastMenuMessage {
		switch ra.Emoji.Name {
		case JoinEmoji:
			if player != nil {
				return
			}

			member, err := g.Session.GuildMember(g.GuildID, ra.UserID)
			if member.User.Bot {
				return
			}

			if err == nil {
				g.addPlayer(ra.UserID, member.User.Username)
			}
			g.LastAction = time.Now()
			return
		case LeaveEmoji:
			g.removePlayer(ra.UserID)
			g.LastAction = time.Now()
			return
		case PlayEmoji:
			log.Println("Pressed play")
			g.LastAction = time.Now()
			if g.State == GameStatePreGame && g.GameMaster == ra.UserID {
				g.setState(GameStatePreRoundDelay)
				go g.sendAnnouncment("Starting in 10 seconds", false)
			}

			return
		default:
			log.Println("Unknown: ", ra.Emoji.Name)
		}

	}

	// From here on out only players can take actions
	if player == nil {
		return
	}

	switch g.State {
	case GameStatePickingResponses:
		if ra.MessageID != player.LastReactionMenu {
			return
		}

		g.LastAction = time.Now()
		g.playerPickedResponseReaction(player, ra)
	case GameStatePickingWinner:
		if ra.MessageID != g.LastMenuMessage || player.ID != g.CurrentCardCzar {
			return
		}
		emojiIndex := -1
		for i, v := range CardSelectionEmojis {
			if v == ra.Emoji.Name {
				emojiIndex = i
				break
			}
		}

		if emojiIndex == -1 || emojiIndex >= len(g.Responses) {
			return
		}

		winner := g.Responses[emojiIndex]
		winner.Player.Wins++
		g.presentWinner(winner)
		g.setState(GameStatePreGame)
	}

}

func (g *Game) presentWinner(winningPick *PickedResonse) {

	// Sort the players by the number of wins
	// note: this wont change the cardzar order as thats done as lowest -> highest user ids
	sort.Slice(g.Players, func(i int, j int) bool {
		return g.Players[i].Wins > g.Players[j].Wins
	})

	standings := ""
	for _, v := range g.Players {
		standings += fmt.Sprintf("%20s: %d\n", v.Username, v.Wins)
	}

	args := make([]interface{}, len(winningPick.Selections))
	for i, v := range winningPick.Selections {
		args[i] = "**" + v + "**"
	}
	winnerCard := fmt.Sprintf(g.CurrentPropmpt.Prompt, args...)

	content := fmt.Sprintf("picked <@%d>'s card(s)!\n%s\n\n**Standings:**\n%s\n\nNext round in 10 seconds...", winningPick.Player.ID, winnerCard, standings)
	msg, err := g.Session.ChannelMessageSend(g.MasterChannel, content)
	if err != nil {
		return
	}

	g.addCommonMenuReactions(msg.ID, false)
}

func (g *Game) playerPickedResponseReaction(player *Player, ra *discordgo.MessageReactionAdd) {
	if len(player.SelectedCards) >= g.CurrentPropmpt.NumPick {
		return
	}

	emojiIndex := -1
	for i, v := range CardSelectionEmojis {
		if v == ra.Emoji.Name {
			emojiIndex = i
			break
		}
	}

	if emojiIndex < 0 {
		// Unknown reaction
		return
	}

	if emojiIndex >= len(player.Cards) {
		// Somehow picked a reaction that they cant (probably added the reaction themselv to mess with the bot)
		return
	}

	for _, selection := range player.SelectedCards {
		if selection == emojiIndex {
			// Already selected this card
			return
		}
	}

	player.SelectedCards = append(player.SelectedCards, emojiIndex)
	card := player.Cards[emojiIndex]

	respMsg := fmt.Sprintf("Selected %s", card)
	if len(player.SelectedCards) >= g.CurrentPropmpt.NumPick {
		respMsg += fmt.Sprintf(", go to <#%d> and wait for the other players to finish their selections, the winner will be picked there", g.MasterChannel)
	} else {
		respMsg += fmt.Sprintf(", select %d more cards", g.CurrentPropmpt.NumPick-len(player.SelectedCards))
	}

	go g.Session.ChannelMessageSend(player.Channel, respMsg)
}

type Player struct {
	ID            int64
	Username      string
	Cards         []ResponseCard
	SelectedCards []int
	Wins          int

	Channel int64

	// Wether this user is playing this round, if the user joined in the middle of a round this will be false
	Playing bool

	LastReactionMenu int64
}

func (p *Player) PresentBoard(session *discordgo.Session, currentPrompt *PromptCard, currentCardCzar int64) {
	if currentCardCzar == p.ID {
		return
	}

	msg := "-----------\n**Next round:**\n" + strings.Replace(currentPrompt.Prompt, "%s", "____", -1)
	msg += "\n\n"
	msg += fmt.Sprintf("Pick %d cards:\n", currentPrompt.NumPick)

	for i, v := range p.Cards {
		msg += CardSelectionEmojis[i] + ": " + string(v) + "\n"
	}

	resp, err := session.ChannelMessageSend(p.Channel, msg)
	if err != nil {
		return
	}

	p.LastReactionMenu = resp.ID

	if currentCardCzar != p.ID {
		for i, _ := range p.Cards {
			session.MessageReactionAdd(p.Channel, resp.ID, CardSelectionEmojis[i])
		}
	}
}
