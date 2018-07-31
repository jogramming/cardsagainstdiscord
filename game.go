package cardsagainstdiscord

import (
	"fmt"
	"github.com/jonas747/discordgo"
	"log"
	"math/rand"
	"sort"
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
	PreRoundDelayDuration  = time.Second * 15
	PickResponseDuration   = time.Second * 60
	PickWinnerDuration     = time.Second * 90
	GameExpireAfter        = time.Second * 300
	GameExpireAfterPregame = time.Minute * 30
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

	JoinEmoji      = "➕"
	LeaveEmoji     = "➖"
	PlayPauseEmoji = "⏯"
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

	PlayerLimit        int
	Packs              []string
	availablePrompts   []*PromptCard
	availableResponses []ResponseCard

	Players []*Player

	State        GameState
	StateEntered time.Time

	// The time the most recent action was taken, if we go too long without a user action we expire the game
	LastAction time.Time

	CurrentPropmpt *PromptCard

	LastMenuMessage int64

	Responses []*PickedResonse

	stopped bool
	stopch  chan bool
}

type PickedResonse struct {
	Player     *Player
	Selections []ResponseCard
}

func (g *Game) Created() {
	g.LastAction = time.Now()

	embed := &discordgo.MessageEmbed{
		Title:       "Game created!",
		Description: fmt.Sprintf("React with %s to join and %s to leave, the game master can start/stop the game with %s", JoinEmoji, LeaveEmoji, PlayPauseEmoji),
	}

	msg, err := g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)

	if err != nil {
		return
	}

	g.LastMenuMessage = msg.ID

	g.stopch = make(chan bool)

	g.loadPackPrompts()
	g.loadPackResponses()

	go g.runTicker()

	go g.addCommonMenuReactions(msg.ID)
}

func (g *Game) loadPackResponses() {
	for _, v := range g.Packs {
		pack := Packs[v]
		g.availableResponses = append(g.availableResponses, pack.Responses...)
	}
}
func (g *Game) loadPackPrompts() {
	for _, v := range g.Packs {
		pack := Packs[v]
		g.availablePrompts = append(g.availablePrompts, pack.Prompts...)
	}
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
	if len(g.availableResponses) < 1 {
		g.loadPackResponses() // re-shuffle basically, TODO: exclude current hands
	}

	i := rand.Intn(len(g.availableResponses))
	card := g.availableResponses[i]
	g.availableResponses = append(g.availableResponses[:i], g.availableResponses[i+1:]...)
	return card
}

func (g *Game) getRandomPlayerCards(num int) []ResponseCard {
	result := make([]ResponseCard, 0, num)

	if len(g.availableResponses) < 1 {
		return result
	}

	for len(result) < num {
		card := g.getRandomResponseCard()
		result = append(result, card)
	}

	return result
}

func (g *Game) sendAnnouncment(msg string, allPlayers bool) {
	embed := &discordgo.MessageEmbed{
		Description: msg,
	}

	if allPlayers {
		for _, v := range g.Players {
			go func(channel int64) {
				g.Session.ChannelMessageSendEmbed(channel, embed)
			}(v.Channel)
		}
	}

	g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)
}

func (g *Game) sendAnnouncmentMenu(msg string) {
	embed := &discordgo.MessageEmbed{
		Description: msg,
	}

	m, err := g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)
	if err == nil {
		g.LastMenuMessage = m.ID
		go g.addCommonMenuReactions(m.ID)
	}
}

func (g *Game) Stop() {
	g.Lock()
	if g.stopped {
		g.Unlock()
		return // Already stopped
	}

	close(g.stopch)
	g.Unlock()
}

func (g *Game) runTicker() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-g.stopch:
			return
		case <-ticker.C:
			g.Tick()
		}
	}
}

func (g *Game) Tick() {
	g.Lock()
	defer g.Unlock()

	expireAfter := GameExpireAfter
	if g.State == GameStatePreGame {
		expireAfter = GameExpireAfterPregame
	}
	if time.Since(g.LastAction) > expireAfter || len(g.Players) < 1 {
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
	if len(g.Players) < 2 {
		g.setState(GameStatePreGame)
		g.sendAnnouncment("Not enough players...", false)
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

	lastPick := 1
	if g.CurrentPropmpt != nil {
		lastPick = g.CurrentPropmpt.NumPick
	}

	// Pick random propmpt
	g.CurrentPropmpt = g.randomPrompt()

	// Give each player a random card (if they're below 10 cards)
	g.giveEveryoneCards(lastPick)

	// Pick next cardzar
	g.CurrentCardCzar = NextCardCzar(g.Players, g.CurrentCardCzar)

	// Present the board
	g.presentStartRound()

	g.setState(GameStatePickingResponses)

}

func NextCardCzar(players []*Player, current int64) int64 {
	var next int64
	var lowest int64
	for _, v := range players {
		if v.ID == current || !v.Playing {
			continue
		}

		if v.ID > current && (v.ID < next || next == 0) {
			next = v.ID
		}

		if lowest == 0 || v.ID < lowest {
			lowest = v.ID
		}
	}

	if next == 0 {
		next = lowest
	}

	return next
}

func (g *Game) randomPrompt() *PromptCard {
	if len(g.availablePrompts) < 1 {
		g.loadPackPrompts() // ran out of cards, just relaod the packs
	}

	i := rand.Intn(len(g.availablePrompts))
	prompt := g.availablePrompts[i]
	g.availablePrompts = append(g.availablePrompts[:i], g.availablePrompts[i+1:]...)

	return prompt
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

	for _, player := range g.Players {
		go func(p *Player) {
			p.PresentBoard(g.Session, g.CurrentPropmpt, g.CurrentCardCzar)
		}(player)
	}

	playerInstructions := fmt.Sprintf("Check your dm for your cards and make your selections there, then return here, you have %d seconds", int(PickResponseDuration.Seconds()))
	embed := &discordgo.MessageEmbed{
		Title: "Next round started!",
		Color: 7001855,
		// Description: msg,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Prompt",
				Value: g.CurrentPropmpt.PlaceHolder(),
			},
			&discordgo.MessageEmbedField{
				Name:  "CardCzar",
				Value: fmt.Sprintf("<@%d>", g.CurrentCardCzar),
			},
			&discordgo.MessageEmbedField{
				Name:  "Instructions",
				Value: "Players: " + playerInstructions + "\nCardCzar: Wait until all players have picked cards(s) then select the best one(s)",
			},
		},
	}
	m, err := g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)
	if err == nil {
		g.LastMenuMessage = m.ID
		go g.addCommonMenuReactions(m.ID)
	}
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

	// Shuffle them
	perm := rand.Perm(len(g.Responses))
	newResponses := make([]*PickedResonse, len(g.Responses))
	for i, v := range perm {
		newResponses[i] = g.Responses[v]
	}
	g.Responses = newResponses

	// Shows all the picks in both dm's and main channel
	g.presentPickedResponseCards()
	g.setState(GameStatePickingWinner)
}

func (g *Game) presentPickedResponseCards() {

	embed := &discordgo.MessageEmbed{
		Title:       "Pick the winner",
		Description: fmt.Sprintf("Cards have been picked, pick the best one(s) <@%d>! you have %d seconds.", g.CurrentCardCzar, int(PickWinnerDuration.Seconds())),
		Color:       5659830,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "Prompt",
				Value: g.CurrentPropmpt.PlaceHolder(),
			},
			&discordgo.MessageEmbedField{
				Name: "Candidates",
			},
		},
	}

	for i, v := range g.Responses {
		filledPrompt := g.CurrentPropmpt.WithCards(v.Selections)

		embed.Fields[1].Value += CardSelectionEmojis[i] + ": " + filledPrompt + "\n\n"
	}

	msg, err := g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)
	if err != nil {
		return
	}

	numOptions := len(g.Responses)
	go func() {
		for i := 0; i < numOptions; i++ {
			g.Session.MessageReactionAdd(g.MasterChannel, msg.ID, CardSelectionEmojis[i])
		}

		g.addCommonMenuReactions(msg.ID)
	}()

	g.LastMenuMessage = msg.ID
}

func (g *Game) cardzarExpired() {
	msg, err := g.Session.ChannelMessageSend(g.MasterChannel, fmt.Sprintf("<%d> didn't pick a winner in %d seconds, skipping round...", g.CurrentCardCzar, int(PickWinnerDuration.Seconds())))
	if err == nil {
		go g.addCommonMenuReactions(msg.ID)
	}

	g.setState(GameStatePreRoundDelay)
}

func (g *Game) gameExpired() {
	g.Session.ChannelMessageSend(g.MasterChannel, "CAH Game expired, too long without any actions or no players.")
	go g.Manager.RemoveGame(g.MasterChannel)
}

func (g *Game) addCommonMenuReactions(mID int64) {
	g.Session.MessageReactionAdd(g.MasterChannel, mID, JoinEmoji)
	g.Session.MessageReactionAdd(g.MasterChannel, mID, LeaveEmoji)
	g.Session.MessageReactionAdd(g.MasterChannel, mID, PlayPauseEmoji)
}

func (g *Game) HandleRectionAdd(ra *discordgo.MessageReactionAdd) {
	g.Lock()
	defer g.Unlock()

	log.Println("Handling RA in game: ", ra.Emoji.Name, ", ", ra.UserID)

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

			go func() {
				member, err := g.Session.GuildMember(g.GuildID, ra.UserID)
				if err != nil || member.User.Bot {
					return
				}

				if err = g.Manager.PlayerTryJoinGame(g.MasterChannel, member.User.ID, member.User.Username); err == nil {
					g.Lock()
					g.LastAction = time.Now()
					g.Unlock()
				} else {
					log.Println("Failed adding", ra.UserID, "to", g.MasterChannel, ":", err.Error())
				}
			}()

			return
		case LeaveEmoji:
			go g.Manager.PlayerTryLeaveGame(ra.UserID)
			g.LastAction = time.Now()
			return
		case PlayPauseEmoji:
			log.Println("Pressed play/pause")
			g.LastAction = time.Now()
			if g.State == GameStatePreGame && g.GameMaster == ra.UserID {
				g.setState(GameStatePreRoundDelay)
				go g.sendAnnouncment(fmt.Sprintf("Starting in %d seconds", int(PreRoundDelayDuration.Seconds())), false)
			} else if g.GameMaster == ra.UserID {
				for _, v := range g.Players {
					v.SelectedCards = nil
				}

				g.Responses = nil
				g.setState(GameStatePreGame)

				go g.sendAnnouncmentMenu(fmt.Sprintf("Paused, react with %s to continue, the game can be paused for max 30 minutes before it expires.", PlayPauseEmoji))
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
		g.setState(GameStatePreRoundDelay)
		g.LastAction = time.Now()
	}
}

func (g *Game) presentWinner(winningPick *PickedResonse) {

	// Sort the players by the number of wins
	// note: this wont change the cardzar order as thats done as lowest -> highest user ids
	sort.Slice(g.Players, func(i int, j int) bool {
		return g.Players[i].Wins > g.Players[j].Wins
	})

	standings := "```\n"
	for _, v := range g.Players {
		standings += fmt.Sprintf("%-20s: %d\n", v.Username, v.Wins)
	}
	standings += "```"

	winnerCard := g.CurrentPropmpt.WithCards(winningPick.Selections)

	title := fmt.Sprintf("%s Won!", winningPick.Player.Username)
	content := fmt.Sprintf("%s\n\n**Standings:**\n%s\n\nNext round in %d seconds...", winnerCard, standings, int(PreRoundDelayDuration.Seconds()))
	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: content,
		Color:       15276265,
	}

	msg, err := g.Session.ChannelMessageSendEmbed(g.MasterChannel, embed)
	// msg, err := g.Session.ChannelMessageSend(g.MasterChannel, content)
	if err != nil {
		return
	}

	g.addCommonMenuReactions(msg.ID)
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

	respMsg := fmt.Sprintf("Selected **%s**", card)
	if len(player.SelectedCards) >= g.CurrentPropmpt.NumPick {
		respMsg += fmt.Sprintf(", go to <#%d> and wait for the other players to finish their selections, the winner will be picked there", g.MasterChannel)
	} else {
		respMsg += fmt.Sprintf(", select %d more cards", g.CurrentPropmpt.NumPick-len(player.SelectedCards))
	}

	embed := &discordgo.MessageEmbed{
		Description: respMsg,
	}
	go g.Session.ChannelMessageSendEmbed(player.Channel, embed)
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

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Pick %d cards!", currentPrompt.NumPick),
		Description: currentPrompt.PlaceHolder(),
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Options",
			},
		},
	}

	for i, v := range p.Cards {
		embed.Fields[0].Value += CardSelectionEmojis[i] + ": " + string(v) + "\n"
	}

	resp, err := session.ChannelMessageSendEmbed(p.Channel, embed)
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
