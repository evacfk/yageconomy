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
						Value:  fmt.Sprint(targetAccount.MoneyBank) + conf.CurrencySymbol,
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Wallet",
						Value:  fmt.Sprint(targetAccount.MoneyWallet) + conf.CurrencySymbol,
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Gambling profit boost %",
						Value:  fmt.Sprintf("%d%%", targetAccount.GamblingBoostPercentage),
					},
					&discordgo.MessageEmbedField{
						Inline: true,
						Name:   "Fish Caugth",
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

			return SimpleEmbedResponse(ms, "Withdrew **%d%s** from your bank, your wallet now has **%d%s**", amount, conf.CurrencySymbol, account.MoneyWallet, conf.CurrencySymbol), nil
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

			return SimpleEmbedResponse(ms, "Deposited **%d%s** Into your bank account, your bank now contains **%d%s**", amount, conf.CurrencySymbol, account.MoneyBank, conf.CurrencySymbol), nil
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

			return SimpleEmbedResponse(ms, "Sent **%d%s** to **%s**%s", amount, conf.CurrencySymbol, target.Username, extraStr), nil
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

			return SimpleEmbedResponse(ms, "Awarded **%s** with %d%s%s", target.Username, amount, conf.CurrencySymbol, extraStr), nil
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

			return SimpleEmbedResponse(ms, "Took away %d%s from **%s**%s", amount, conf.CurrencySymbol, target.Username, extraStr), nil
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

			return SimpleEmbedResponse(ms, "Sent %d%s to %s%s", amount, conf.CurrencySymbol, target.Username, extraStr), nil
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
				return SimpleEmbedResponse(ms, "Claimed your daily of **%d%s**", conf.DailyAmount, conf.CurrencySymbol), nil
			}

			timeToWait := account.LastDailyClaim.Add(time.Duration(conf.DailyFrequency) * time.Minute).Sub(time.Now())
			return ErrorEmbed(ms, "You can't claim your daily yet again! Please wait another %s.", common.HumanizeDuration(common.DurationPrecisionSeconds, timeToWait)), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory: CategoryEconomy,
		Name:        "TopMoney",
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

				offset := (newPage - 1) * 20
				if offset < 0 {
					offset = 0
				}

				result, err := models.EconomyUsers(
					models.EconomyUserWhere.GuildID.EQ(parsed.GS.ID),
					qm.OrderBy("money_wallet + money_bank desc"),
					qm.Limit(20),
					qm.Offset(offset)).AllG(context.Background())
				if err != nil {
					return nil, err
				}

				buf := strings.Builder{}
				buf.WriteString("Economy leaderboard:```\n")

				userIDs := make([]int64, len(result))
				for i, v := range result {
					userIDs[i] = v.UserID
				}

				users := bot.GetUsersGS(parsed.GS, userIDs...)

				for i, v := range result {
					user := users[i]
					buf.WriteString(fmt.Sprintf("#%2d: %-20s : %d%s\n", i+offset+1, user.Username, v.MoneyBank+v.MoneyWallet, conf.CurrencySymbol))
				}

				buf.WriteString("```")

				return SimpleEmbedResponse(ms, buf.String()), nil
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

			err = PlantMoney(parsed.Context(), conf, parsed.CS.ID, ms.ID, amount, parsed.Args[1].Str())
			if err != nil {
				return nil, err
			}

			_, err = common.PQ.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, ms.ID, amount)
			if err != nil {
				return nil, err
			}

			if errors.Cause(err) != sql.ErrNoRows {
				return nil, err
			}

			m := &models.EconomyPlant{
				ChannelID: parsed.CS.ID,
				GuildID:   parsed.GS.ID,
				AuthorID:  ms.ID,
				Amount:    int64(amount),
				Password:  parsed.Args[1].Str(),
			}

			err = m.InsertG(parsed.Context(), boil.Infer())
			if err != nil {
				return nil, err
			}

			extraStr := ""
			if parsed.Args[1].Str() != "" {
				extraStr = " with the passowrd `" + parsed.Args[1].Str() + "`"
			}

			return SimpleEmbedResponse(commands.ContextMS(parsed.Context()),
				fmt.Sprintf("%s planted **%d%s**%s in the channel!\nUse `pick [password]` to pick it up", ms.Username, amount, conf.CurrencySymbol, extraStr)), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Pick",
		Description:  "Picks up money planted in the channel previously using plant",
		RequiredArgs: 0,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Password", Type: dcmd.String, Default: ""},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			p, err := models.FindEconomyPlantG(parsed.Context(), parsed.CS.ID)
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No plant in this channel :("), nil
				}

				return nil, err
			}

			noPlant := false
			noMoneyLeft := false
			wrongPassword := false
			pmAmount := int64(0)
			err = common.SqlTX(func(tx *sql.Tx) error {
				pm, err := models.EconomyPlants(qm.Where("channel_id = ?", parsed.CS.ID), qm.For("update")).One(parsed.Context(), tx)
				if err != nil {
					if errors.Cause(err) == sql.ErrNoRows {
						noPlant = true
					}
					return err
				}
				pmAmount = pm.Amount

				if pm.Password != "" {
					if !strings.EqualFold(pm.Password, parsed.Args[0].Str()) {
						wrongPassword = true
						return nil
					}
				}

				_, err = tx.Exec("UPDATE economy_users SET money_wallet = money_wallet + $3 WHERE user_id = $2 AND guild_id = $1", parsed.GS.ID, ms.ID, pm.Amount)
				if err != nil {
					return err
				}

				_, err = pm.Delete(parsed.Context(), tx)
				return err
			})

			if noPlant {
				return ErrorEmbed(ms, "No plant in this channel :("), nil
			}

			if err != nil {
				return nil, err
			}

			if noMoneyLeft {
				return ErrorEmbed(ms, "The person that planted this didn't manage their money correctly and no longer has enough money..."), nil
			}

			if wrongPassword {
				return ErrorEmbed(ms, "Incorrect passowrd >:u"), nil
			}

			common.BotSession.ChannelMessageDelete(parsed.CS.ID, p.MessageID)

			return SimpleEmbedResponse(ms, fmt.Sprintf("Picked up **%d%s**!", pmAmount, conf.CurrencySymbol)), nil
		},
	},
}
