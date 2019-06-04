package yageconomy

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dstate"
	"github.com/jonas747/yageconomy/models"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/paginatedmessages"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"strings"
	"time"
)

var CoreCommands = []*commands.YAGCommand{
	&commands.YAGCommand{
		CmdCategory: CategoryEconomy,
		Name:        "$",
		Aliases:     []string{"balance", "wallet"},
		Description: "Shows you balance",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			var targetAccount *models.EconomyUser
			var targetMS *dstate.MemberState
			conf := CtxConfig(parsed.Context())

			if parsed.Args[0].Value != nil {
				target := parsed.Args[0].Value.(*dstate.MemberState)

				var err error
				targetAccount, _, err = GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
				if err != nil {
					return nil, err
				}

				targetMS = target
			} else {
				targetMS = commands.ContextMS(parsed.Context())
				targetAccount = CtxUser(parsed.Context())
			}

			embed := &discordgo.MessageEmbed{
				Author:      UserEmebdAuthor(targetMS),
				Description: "Account of " + targetMS.Username,
				Color:       ColorBlue,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Bank Balance",
						Value:  conf.CurrencySymbol + fmt.Sprint(targetAccount.MoneyBank),
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Wallet",
						Value:  conf.CurrencySymbol + fmt.Sprint(targetAccount.MoneyWallet),
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Gambling profit boost %",
						Value:  fmt.Sprintf("%d%%", targetAccount.GamblingBoostPercentage),
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Fish caught",
						Value:  fmt.Sprintf("%d", targetAccount.FishCaugth),
					},
				},
			}
			return embed, nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Withdraw",
		Description:  "Withdraws money from your bank account into your wallet",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[0].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if int64(amount) > account.MoneyBank {
				return ErrorEmbed(ms, "You don't have that amount in your bank"), nil
			}

			account.MoneyBank -= int64(amount)
			account.MoneyWallet += int64(amount)
			_, err := account.UpdateG(parsed.Context(), boil.Whitelist("money_bank", "money_wallet"))
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Withdrew **%s%d** from your bank, your wallet now has **%s%d**", conf.CurrencySymbol, amount, conf.CurrencySymbol, account.MoneyWallet), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Deposit",
		Description:  "Deposits money into your bank account from your wallet",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
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

			account.MoneyBank += int64(amount)
			account.MoneyWallet -= int64(amount)
			_, err := account.UpdateG(parsed.Context(), boil.Whitelist("money_bank", "money_wallet"))
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Deposited **%s%d** Into your bank account, your bank now contains **%s%d**", conf.CurrencySymbol, amount, conf.CurrencySymbol, account.MoneyBank), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Give",
		Description:  "Give someone money from your wallet",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Reason", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := parsed.Args[0].Value.(*dstate.MemberState)

			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[1].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if int64(amount) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have that amount in your wallet"), nil
			}

			if ms.ID == target.ID {
				return ErrorEmbed(ms, "Can't send money to yourself, silly..."), nil
			}

			targetAccount, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			targetAccount.MoneyWallet += int64(amount)
			account.MoneyWallet -= int64(amount)

			// update the acconts
			err = common.SqlTX(func(tx *sql.Tx) error {

				_, err := tx.Exec("UPDATE economy_users SET money_wallet = money_wallet + $3 WHERE user_id = $2 AND guild_id = $1", parsed.GS.ID, target.ID, amount)
				if err != nil {
					return err
				}

				_, err = tx.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3 WHERE user_id = $2 AND guild_id = $1", parsed.GS.ID, ms.ID, amount)
				return err
			})

			if err != nil {
				return nil, err
			}

			extraStr := ""
			if parsed.Args[2].Str() != "" {
				extraStr = " with the message: **" + parsed.Args[2].Str() + "**"
			}

			return SimpleEmbedResponse(ms, "Sent **%s%d** to **%s**%s", conf.CurrencySymbol, amount, target.Username, extraStr), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Award",
		Description:  "Award a member of the server some money (admins only)",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Reason", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := parsed.Args[0].Value.(*dstate.MemberState)

			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[1].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if !common.ContainsInt64SliceOneOf(conf.Admins, ms.Roles) {
				return ErrorEmbed(ms, "You're not a economy admin"), nil
			}

			// esnure that the account exists
			_, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			_, err = common.PQ.Exec("UPDATE economy_users SET money_wallet = money_wallet + $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, target.ID, amount)
			if err != nil {
				return nil, err
			}

			extraStr := ""
			if parsed.Args[2].Str() != "" {
				extraStr = " with the message: **" + parsed.Args[2].Str() + "**"
			}

			return SimpleEmbedResponse(ms, "Awarded **%s** with %s%d%s", target.Username, conf.CurrencySymbol, amount, extraStr), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Take",
		Description:  "Takes away money from someone (admins only)",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Reason", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := parsed.Args[0].Value.(*dstate.MemberState)

			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[1].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if !common.ContainsInt64SliceOneOf(conf.Admins, ms.Roles) {
				return ErrorEmbed(ms, "You're not a economy admin"), nil
			}

			// esnure that the account exists
			_, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			_, err = common.PQ.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, target.ID, amount)
			if err != nil {
				return nil, err
			}

			extraStr := ""
			if parsed.Args[2].Str() != "" {
				extraStr = " with the message: **" + parsed.Args[2].Str() + "**"
			}

			return SimpleEmbedResponse(ms, "Took away %s%d from **%s**%s", conf.CurrencySymbol, amount, target.Username, extraStr), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Give",
		Description:  "Give someone money from your wallet",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "Amount", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Reason", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := parsed.Args[0].Value.(*dstate.MemberState)

			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[1].Int()
			if amount < 1 {
				return ErrorEmbed(ms, "Amount too small"), nil
			}

			if int64(amount) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have that amount in your wallet"), nil
			}

			targetAccount, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			targetAccount.MoneyWallet += int64(amount)
			account.MoneyWallet -= int64(amount)

			// update the acconts
			err = TransferMoneyWallet(parsed.Context(), nil, conf, false, ms.ID, target.ID, int64(amount), int64(amount))
			if err != nil {
				return nil, err
			}

			extraStr := ""
			if parsed.Args[2].Str() != "" {
				extraStr = " with the message: **" + parsed.Args[2].Str() + "**"
			}

			return SimpleEmbedResponse(ms, "Sent %s%d to %s%s", conf.CurrencySymbol, amount, target.Username, extraStr), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory: CategoryEconomy,
		Name:        "Daily",
		Description: "Claim your daily free cash",
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			if conf.DailyAmount < 1 {
				return ErrorEmbed(ms, "Daily not set up on this server"), nil
			}

			result, err := common.PQ.Exec(`UPDATE economy_users SET last_daily_claim = now(), money_wallet = money_wallet + $4
			WHERE guild_id = $1 AND user_id = $2 AND EXTRACT(EPOCH FROM (now() - last_daily_claim))  > $3`, parsed.GS.ID, ms.ID, conf.DailyFrequency*60, conf.DailyAmount)
			if err != nil {
				return nil, err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return nil, err
			}

			if rows > 0 {
				return SimpleEmbedResponse(ms, "Claimed your daily of **%s%d**", conf.CurrencySymbol, conf.DailyAmount), nil
			}

			timeToWait := account.LastDailyClaim.Add(time.Duration(conf.DailyFrequency) * time.Minute).Sub(time.Now())
			return ErrorEmbed(ms, "You can't claim your daily yet again! Please wait another %s.", common.HumanizeDuration(common.DurationPrecisionSeconds, timeToWait)), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory: CategoryEconomy,
		Name:        "TopMoney",
		Aliases:     []string{"LB"},
		Description: "Economy leaderboard, optionally specify a page",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Page", Type: dcmd.Int, Default: 1},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			conf := CtxConfig(parsed.Context())

			page := parsed.Args[0].Int()
			if page < 0 {
				page = 1
			}

			ms := commands.ContextMS(parsed.Context())
			_, err := paginatedmessages.CreatePaginatedMessage(parsed.GS.ID, parsed.CS.ID, page, 0, func(p *paginatedmessages.PaginatedMessage, newPage int) (*discordgo.MessageEmbed, error) {

				offset := (newPage - 1) * 10
				if offset < 0 {
					offset = 0
				}

				result, err := models.EconomyUsers(
					models.EconomyUserWhere.GuildID.EQ(parsed.GS.ID),
					qm.OrderBy("money_wallet + money_bank desc"),
					qm.Limit(10),
					qm.Offset(offset)).AllG(context.Background())
				if err != nil {
					return nil, err
				}

				embed := SimpleEmbedResponse(ms, "")
				embed.Title = conf.CurrencySymbol + " Leaderboard"

				userIDs := make([]int64, len(result))
				for i, v := range result {
					userIDs[i] = v.UserID
				}

				users := bot.GetUsersGS(parsed.GS, userIDs...)

				for i, v := range result {
					user := users[i]
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  fmt.Sprintf("#%d %s", i+offset+1, user.Username),
						Value: fmt.Sprintf("%s%d", conf.CurrencySymbol, v.MoneyBank+v.MoneyWallet),
					})

				}

				return embed, nil
				// return SimpleEmbedResponse(ms, buf.String()), nil
			})

			return nil, err
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Plant",
		Description:  "Plants a certain amount of currency in the channel, optionally with a password, use Pick to pick it",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Money", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Password", Type: dcmd.String, Default: ""},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			amount := parsed.Args[0].Int()

			if int64(amount) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have that amount in your wallet"), nil
			}

			if amount < 1 {
				return ErrorEmbed(ms, "Too low amount"), nil
			}

			_, err := models.FindEconomyPlantG(parsed.Context(), parsed.CS.ID)
			if err == nil {
				return ErrorEmbed(ms, "There's already money planted in this channel"), nil
			}

			cmdPrefix, _ := commands.GetCommandPrefix(conf.GuildID)
			msgContent := fmt.Sprintf("%s planted **%s%d** in the channel!\nUse `%spick (code-here)` to pick it up", ms.Username, conf.CurrencySymbol, amount, cmdPrefix)

			err = PlantMoney(parsed.Context(), conf, parsed.CS.ID, ms.ID, amount, parsed.Args[1].Str(), msgContent)
			if err != nil {
				return nil, err
			}

			err = TransferMoneyWallet(parsed.Context(), nil, conf, false, ms.ID, common.BotUser.ID, int64(amount), int64(amount))
			if err != nil {
				return nil, err
			}

			bot.MessageDeleteQueue.DeleteMessages(parsed.CS.ID, parsed.Msg.ID)

			return nil, nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Pick",
		Description:  "Picks up money planted in the channel previously using plant",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Password", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			bot.MessageDeleteQueue.DeleteMessages(parsed.CS.ID, parsed.Msg.ID)

			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			p, err := models.EconomyPlants(
				models.EconomyPlantWhere.ChannelID.EQ(parsed.CS.ID),
				models.EconomyPlantWhere.Password.EQ(strings.ToLower(parsed.Args[0].Str())),
				qm.OrderBy("message_id desc"),
			).OneG(parsed.Context())

			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No plant in this channel or incorrect passowrd :("), nil
				}

				return nil, err
			}

			noPlant := false
			pmAmount := int64(0)
			err = common.SqlTX(func(tx *sql.Tx) error {
				pm, err := models.EconomyPlants(models.EconomyPlantWhere.MessageID.EQ(p.MessageID), qm.For("update")).One(parsed.Context(), tx)
				if err != nil {
					if errors.Cause(err) == sql.ErrNoRows {
						noPlant = true
					}
					return err
				}
				if pm.MessageID != p.MessageID {
					noPlant = true
					return nil
				}

				pmAmount = pm.Amount

				_, err = tx.Exec("UPDATE economy_users SET money_wallet = money_wallet + $3 WHERE user_id = $2 AND guild_id = $1", parsed.GS.ID, ms.ID, pm.Amount)
				if err != nil {
					return err
				}

				_, err = pm.Delete(parsed.Context(), tx)
				return err
			})

			if noPlant {
				return ErrorEmbed(ms, "Yikes, someone snatched it before you."), nil
			}

			if err != nil {
				return nil, err
			}

			common.BotSession.ChannelMessageDelete(parsed.CS.ID, p.MessageID)

			return SimpleEmbedResponse(ms, fmt.Sprintf("Picked up **%s%d**!", conf.CurrencySymbol, pmAmount)), nil
		},
	},
}
