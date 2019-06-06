package yageconomy

// import (
// 	"github.com/jonas747/discordgo"
// 	"github.com/jonas747/yagpdb/common"
// 	"sync"
// 	"time"
// )

// type HeistProgressState int

// const (
// 	HeistProgressStateWaiting HeistProgressState = iota
// 	HeistProgressStateStarting
// 	HeistProgressStateInvading
// 	HeistProgressStateCollecting
// 	HeistProgressStateLeaving
// 	HeistProgressStateGetaway
// )

// type HeistEvent struct {
// 	Description   string
// 	Chance        float64
// 	MemberLossMin int
// 	MemberLossMax int

// 	MinMembers int
// 	MaxMembers int

// 	MoneyLossPercentage int
// }

// type HeistSession struct {
// 	sync.Mutex

// 	GuildID   int64
// 	ChannelID int64
// 	MessageID int64

// 	Author  *discordgo.User
// 	Members []*discordgo.User

// 	CreatedAt time.Time
// 	StartsAt  time.Time

// 	ProgressState  HeistProgressState
// 	StateChangedAt time.Time
// }

// var (
// 	activeHeists   []*HeistSession
// 	activeHeistsmu sync.Mutex
// )

// func NewHeist(guildID, channelID int64, author *discordgo.User, waitUntilStart time.Duration) (resp string, err error) {
// 	activeHeistsmu.Lock()
// 	defer activeHeistsmu.Unlock()

// 	for _, v := range activeHeists {
// 		if v.ChannelID == channelID {
// 			return "Already a heist going on in this channel", nil
// 		}

// 		if v.GuildID == guildID {
// 			for _, m := range v.Members {
// 				if m.ID == author.ID {
// 					return "You're already in another heist on this server", nil
// 				}
// 			}
// 		}
// 	}

// 	msg, err := common.BotSession.ChannelMessageSend(channelID, "Setting up heist...")
// 	if err != nil {
// 		return "", err
// 	}

// 	heist := &HeistSession{
// 		GuildID:   guildID,
// 		ChannelID: channelID,
// 		Author:    author,
// 		MessageID: msg.ID,
// 		Members:   []*discordgo.User{author},

// 		StartsAt:       time.Now().Add(waitUntilStart),
// 		StateChangedAt: time.Now(),
// 		CreatedAt:      time.Now(),
// 	}

// 	activeHeists = append(activeHeists, heist)
// 	go heist.Run()
// 	return "", nil
// }

// func removeHeist(h *HeistSession) {
// 	activeHeistsmu.Lock()
// 	defer activeHeistsmu.Unlock()

// 	for i, v := range activeHeists {
// 		if v == h {
// 			activeHeists = append(activeHeists[:i], activeHeists[i+1:]...)
// 			return
// 		}
// 	}
// }

// func (hs *HeistSession) init() {
// 	err := common.BotSession.MessageReactionAdd(hs.ChannelID, hs.MessageID, "✅")
// 	if err != nil {
// 		logger.WithError(err).Error("failed adding reaction")
// 	}
// }

// func (hs *HeistSession) Run() {
// 	hs.init()

// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			hs.update()
// 		}
// 	}
// }

// func (hs *HeistSession) update() {
// 	hs.Lock()
// 	defer hs.Unlock()

// 	switch hs.ProgressState {
// 	case HeistProgressStateWaiting:
// 		hs.tickWaiting()
// 		return
// 	}
// }

// func (hs *HeistSession) tickWaiting() {
// 	timeUntilStart := hs.StartsAt.Sub(time.Now())
// 	if timeUntilStart < 0 {
// 		hs.Start()
// 		return
// 	}

// 	timeBetweenUpdates := time.Minute
// 	if timeUntilStart < 10 {
// 		timeBetweenUpdates = time.Second
// 	} else if timeUntilStart < 20 {
// 		timeBetweenUpdates = time.Second * 5
// 	} else if timeUntilStart < 40 {
// 		timeBetweenUpdates = time.Second * 10
// 	} else if timeUntilStart < time.Minute {
// 		timeBetweenUpdates = time.Second * 20
// 	} else if timeUntilStart < time.Minute*2 {
// 		timeUntilStart = time.Second * 30
// 	}

// 	if time.Since(hs.StateChangedAt) > timeBetweenUpdates {
// 		hs.updateWaitingMessage()
// 	}
// }

// func (hs *HeistSession) Start() {

// }

// func (hs *HeistSession) updateWaitingMessage() {
// 	timeUntilStart := hs.StartsAt.Sub(time.Now())
// 	precision := common.DurationPrecisionSeconds
// 	if timeUntilStart > time.Minute*2 {
// 		precision = common.DurationPrecisionMinutes
// 	}

// 	embed := SimpleEmbedResponse(hs.Author, "A heist is being set up by **%s#%s**\nIt's scheduled to start in `%s` at `%s` UTC",
// 		hs.Author.Username, hs.Author.Discriminator, common.HumanizeDuration(precision, timeUntilStart), hs.StartsAt.UTC().Format(time.Kitchen))

// 	embed.Title = "Heist being set up"

// }

// var HeistEvents = map[HeistProgressState][]*HeistEvent{
// 	HeistProgressStateStarting: []*HeistEvent{
// 		&HeistEvent{
// 			Description: "Alright guys, check your guns. We are storming into the bank through all entrances. Let's get the cash and get out before the cops get here.",
// 			MinMembers:  2,
// 			Chance:      1,
// 		},
// 		&HeistEvent{
// 			Description: "So uh, you're charging into the bank alone huh? Well you better get ready because you're starting now!",
// 			MaxMembers:  1,
// 			Chance:      1,
// 		},
// 	},
// 	HeistProgressStateInvading: []*HeistEvent{
// 		&HeistEvent{
// 			Description: "You've entered the building and you're trying to get control and also get the location of the money.",
// 			Chance:      1,
// 		},
// 		&HeistEvent{
// 			Description:         "One of your members shot a hostage trying to play hero, yikes that's 10% off as a penalty.",
// 			Chance:              0.5,
// 			MoneyLossPercentage: 10,
// 		},
// 		&HeistEvent{
// 			Description:   "A hostage played hero and killed one of your guys.",
// 			Chance:        0.25,
// 			MemberLossMin: 1,
// 			MemberLossMax: 1,
// 		},
// 	},
// 	HeistProgressStateCollecting: []*HeistEvent{
// 		&HeistEvent{
// 			Description: "You've found the money and started collacting the dough.",
// 			Chance:      1,
// 		},
// 		&HeistEvent{
// 			Description:         "God damnit! One of your members bags ripped open and money is now everywhere, 25% penalty.",
// 			Chance:              0.5,
// 			MoneyLossPercentage: 25,
// 		},
// 	},
// 	HeistProgressStateLeaving: []*HeistEvent{
// 		&HeistEvent{
// 			Description: "Alright the cops are getting close, time to head out!",
// 			Chance:      1,
// 		},
// 		&HeistEvent{
// 			Description:         "One of your members tripped, spilling some money in the process.",
// 			MoneyLossPercentage: 10,
// 			Chance:              0.25,
// 		},
// 	},
// 	HeistProgressStateGetaway: []*HeistEvent{
// 		&HeistEvent{
// 			Description: "There's the getaway, hopefully nobody knows about it.",
// 			Chance:      1,
// 		},
// 		&HeistEvent{
// 			Description:         "A cop spots you while driving off and shoots at you, missing you all but damaging the money.",
// 			Chance:              0.25,
// 			MoneyLossPercentage: 10,
// 		},
// 		&HeistEvent{
// 			Description:         "There's a blockade up ahead, giving you all kinds of trouble!",
// 			Chance:              0.25,
// 			MemberLossMin:       0,
// 			MemberLossMax:       3,
// 			MoneyLossPercentage: 10,
// 		},
// 	},
// }
