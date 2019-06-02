package yageconomy

import (
	"fmt"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/dstate"
	"github.com/jonas747/yageconomy/models"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"math/rand"
	"strings"
	"time"
)

var GameCommands = []*commands.YAGCommand{

	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "BetFlip",
		Description:  "Bet on heads or tail, if you guess correct you win 2x your bet",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Side", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[0].Int()
			moneyIn := amount
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if int64(amount) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have that amount in your wallet"), nil
			}

			guessedHeads := strings.HasPrefix(strings.ToLower(parsed.Args[1].Str()), "h")

			isHeads := rand.Intn(2) == 0

			won := false
			winningsLosses := int64(amount)
			var err error

			if (isHeads && guessedHeads) || (!isHeads && !guessedHeads) {
				won = true
				winningsLosses = ApplyGamblingBoost(account, int64(amount))
				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, 0, ms.ID, 0, winningsLosses)
			} else {
				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, ms.ID, common.BotUser.ID, winningsLosses, winningsLosses)
			}

			if err != nil {
				return nil, err
			}

			strResult := "heads"
			if !isHeads {
				strResult = "tails"
			}

			msg := ""
			if won {
				msg = fmt.Sprintf("Result is... **%s**: You won! Awarded with **%d%s**", strResult, int64(amount)+winningsLosses, conf.CurrencySymbol)
			} else {
				msg = fmt.Sprintf("Result is... **%s**: You lost... you're now **%d%s** poorer...", strResult, moneyIn, conf.CurrencySymbol)
			}

			return SimpleEmbedResponse(ms, msg), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "BetRoll",
		Description:  "Rolls 1-100, Rolling over 66 yields x2 of your bet, over 90 -> x4 and 100 -> x10.",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Bet", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[0].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if int64(amount) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have that amount in your wallet"), nil
			}

			walletMod := -amount

			roll := rand.Intn(100) + 1

			won := roll > 66
			if roll == 100 {
				walletMod = amount * 9
			} else if roll > 90 {
				walletMod = amount * 3
			} else if roll > 66 {
				walletMod = amount
			}

			var err error
			if won {
				// transfer winnings into our account
				walletMod = int(ApplyGamblingBoost(account, int64(walletMod)))
				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, 0, ms.ID, 0, int64(walletMod))
			} else {
				// transfer losses into bot account
				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, ms.ID, common.BotUser.ID, int64(amount), int64(amount))
			}
			if err != nil {
				return nil, err
			}

			msg := ""
			if won {
				msg = fmt.Sprintf("Rolled **%d** and won! You have been awarded with **%d%s**", roll, walletMod+amount, conf.CurrencySymbol)
			} else {
				msg = fmt.Sprintf("Rolled **%d** and lost... you're now **%d%s** poorer...", roll, amount, conf.CurrencySymbol)
			}

			return SimpleEmbedResponse(ms, msg), nil
		},
	}, &commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Rob",
		Description:  "Steals money from someone, the chance of suceeding = your networth / (their cash + your networth)",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := parsed.Args[0].Value.(*dstate.MemberState)

			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			if conf.RobFine < 1 {
				return ErrorEmbed(ms, "No fine as been set, as a result the rob command has been disabled"), nil
			}

			targetAccount, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			if targetAccount.MoneyWallet < 1 {
				return ErrorEmbed(ms, "This person has no money left in their wallet :("), nil
			}

			if account.MoneyWallet < int64(conf.RobFine) {
				return ErrorEmbed(ms, "You don't have enough money in your wallet to pay the fine if you fail"), nil
			}

			sucessChance := float64(account.MoneyWallet+account.MoneyBank) / float64(targetAccount.MoneyWallet+account.MoneyWallet+account.MoneyBank)
			if rand.Float64() < sucessChance {
				// sucessfully robbed them

				amount := targetAccount.MoneyWallet

				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, target.ID, ms.ID, amount, ApplyGamblingBoost(account, amount))
				if err != nil {
					return nil, err
				}

				return SimpleEmbedResponse(ms, "You sucessfully robbed **%s** for **%d%s**!", target.Username, ApplyGamblingBoost(account, amount), conf.CurrencySymbol), nil
			} else {
				fine := int64(float64(conf.RobFine/100) * float64(account.MoneyWallet))

				err = TransferMoneyWallet(parsed.Context(), nil, conf, false, ms.ID, common.BotUser.ID, fine, fine)

				if err != nil {
					return nil, err
				}

				return ErrorEmbed(ms, "You failed robbing **%s**, you were fined **%d%s** as a result, hopefully you have learned your lesson now.",
					target.Username, conf.RobFine, conf.CurrencySymbol), nil
			}

		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Fish",
		Description:  "Attempts to fish for some easy money",
		RequiredArgs: 1,
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			if conf.FishingMaxWinAmount < 1 {
				return ErrorEmbed(ms, "Fishing not set up on this server"), nil
			}

			wonAmount := int64(0)
			fishAmount := 0
			if rand.Float64() > 0.25 {
				// 75% chance to catch a fish
				wonAmount = rand.Int63n(conf.FishingMaxWinAmount-conf.FishingMinWinAmount) + conf.FishingMinWinAmount
				wonAmount = ApplyGamblingBoost(account, wonAmount)
				fishAmount = 1
			}

			result, err := common.PQ.Exec(`UPDATE economy_users SET last_fishing = now(), money_wallet = money_wallet + $4, fish_caugth = fish_caugth + $5
			WHERE guild_id = $1 AND user_id = $2 AND EXTRACT(EPOCH FROM (now() - last_fishing)) > $3`, parsed.GS.ID, ms.ID, conf.FishingCooldown*60, wonAmount, fishAmount)
			if err != nil {
				return nil, err
			}

			rows, err := result.RowsAffected()
			if err != nil {
				return nil, err
			}

			if rows < 1 {
				timeToWait := account.LastFishing.Add(time.Duration(conf.FishingCooldown) * time.Minute).Sub(time.Now())
				return ErrorEmbed(ms, "You can't fish again yet, please wait another %s.", common.HumanizeDuration(common.DurationPrecisionSeconds, timeToWait)), nil
			}

			if wonAmount == 0 {
				return SimpleEmbedResponse(ms, "Aww man, you let your fish slip away..."), nil
			}

			return SimpleEmbedResponse(ms, "Nice! You caught your fish worth **%d%s**!", wonAmount, conf.CurrencySymbol), nil

		},
	},
}

func ApplyGamblingBoost(account *models.EconomyUser, winnings int64) int64 {
	return int64(float64(winnings) * ((float64(account.GamblingBoostPercentage) / 100) + 1))
}
