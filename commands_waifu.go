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
	"github.com/lib/pq"
	"github.com/mediocregopher/radix"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"strconv"
	"strings"
)

var (
	WaifuCmdTop = &commands.YAGCommand{
		CmdCategory: CategoryWaifu,
		Name:        "Top",
		Description: "Shows top waifus",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Page", Type: dcmd.Int, Default: 1},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			_, err := paginatedmessages.CreatePaginatedMessage(parsed.GS.ID, parsed.CS.ID, parsed.Args[0].Int(), 0,
				func(p *paginatedmessages.PaginatedMessage, newPage int) (*discordgo.MessageEmbed, error) {

					offset := (newPage - 1) * 10
					items, err := models.EconomyUsers(
						models.EconomyUserWhere.GuildID.EQ(parsed.GS.ID),
						qm.OrderBy("waifu_item_worth desc"),
						qm.Limit(10),
						qm.Offset(offset),
					).AllG(context.Background())

					if err != nil {
						return nil, err
					}

					ids := make([]int64, len(items))
					for i, v := range items {
						ids[i] = v.UserID
					}
					users := bot.GetUsersGS(parsed.GS, ids...)

					embed := SimpleEmbedResponse(ms, "")
					embed.Title = "Waifu Leaderboard"

					for i, v := range items {
						user := users[i]
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:  fmt.Sprintf("#%d %s", i+offset+1, user.Username),
							Value: fmt.Sprintf("%s%d", conf.CurrencySymbol, WaifuWorth(v)),
						})

					}

					return embed, nil
				})

			return nil, err
		},
	}
	WaifuCmdInfo = &commands.YAGCommand{
		CmdCategory: CategoryWaifu,
		Name:        "Info",
		Description: "Shows waifu stats of you or your targets",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			target := commands.ContextMS(parsed.Context())
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())

			if parsed.Args[0].Value != nil {
				target = parsed.Args[0].Value.(*dstate.MemberState)
				var err error
				account, _, err = GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
				if err != nil {
					return nil, err
				}
			}

			embed := &discordgo.MessageEmbed{
				Author: UserEmebdAuthor(target),
				Color:  ColorBlue,
				Title:  "Waifu stats",
			}

			var usersToFetch []int64

			if account.WaifudBy != 0 {
				usersToFetch = append(usersToFetch, account.WaifudBy)
			}

			if account.WaifuAffinityTowards != 0 {
				usersToFetch = append(usersToFetch, account.WaifuAffinityTowards)
			}

			usersToFetch = append(usersToFetch, account.Waifus...)

			claimedByStr := "No one :("
			affinityStr := "No one"

			var waifus []*discordgo.User
			if len(usersToFetch) > 0 {
				waifus = bot.GetUsersGS(parsed.GS, usersToFetch...)
				if account.WaifudBy != 0 {
					claimedByStr = waifus[0].Username
					waifus = waifus[1:]
				}

				if account.WaifuAffinityTowards != 0 {
					affinityStr = waifus[0].Username
					waifus = waifus[1:]
				}
			}

			var claimedBuf strings.Builder
			if len(account.Waifus) > 0 {
				for _, v := range waifus {
					claimedBuf.WriteString(v.Username + "\n")
				}
			} else {
				claimedBuf.WriteString("No one...")
			}

			embed.Fields = []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:  "Waifu Worth",
					Value: conf.CurrencySymbol + strconv.FormatInt(WaifuWorth(account), 10),
				},
				&discordgo.MessageEmbedField{
					Name:  "Affinity towards",
					Value: affinityStr,
				},
				&discordgo.MessageEmbedField{
					Name:  "Claimed By",
					Value: claimedByStr,
				},
				&discordgo.MessageEmbedField{
					Name:  "Waifus claimed",
					Value: claimedBuf.String(),
				},
			}

			return embed, nil
		},
	}
	WaifuCmdClaim = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "Claim",
		Description:  "Claims the target as your waifu, using your wallet money, if no amount is specified it will use the lowest",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "Money", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())

			target := parsed.Args[0].Value.(*dstate.MemberState)
			if target.ID == ms.ID {
				return ErrorEmbed(ms, "You can't claim yourself, silly..."), nil
			}

			// pre-generate the account since its simpler with race conditions and whatnot in the transactions
			targetAccount, _, err := GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			// safety checks
			cost := WaifuCost(account, targetAccount)
			claimAmount := cost
			if parsed.Args[1].Int() > 0 {
				claimAmount = int64(parsed.Args[1].Int())
				if claimAmount < cost {
					return ErrorEmbed(ms, "That waifu costs more than that to claim (%s%d%s)", conf.CurrencySymbol, cost), nil
				}

			}

			if account.MoneyWallet < claimAmount {
				return ErrorEmbed(ms, "You don't have that much money in your wallet"), nil
			}

			forcedErrorResp := ""
			err = common.SqlTX(func(tx *sql.Tx) error {

				numRows, err := models.EconomyUsers(qm.Where("guild_id = ? AND user_id = ? AND waifud_by = 0", parsed.GS.ID, target.ID)).UpdateAll(
					parsed.Context(), tx, models.M{"waifud_by": ms.ID})

				if err != nil {
					return err
				}

				if numRows < 1 {
					forcedErrorResp = "That waifu is already claimed by soemone else :("
					return nil
				}

				_, err = tx.Exec("UPDATE economy_users SET waifus = waifus || $4, money_wallet = money_wallet - $3 WHERE guild_id = $1 AND user_id = $2",
					parsed.GS.ID, ms.ID, claimAmount, pq.Int64Array([]int64{target.ID}))
				return errors.Wrap(err, "update_waifus")
			})

			if err != nil {
				return nil, err
			}

			if forcedErrorResp != "" {
				return ErrorEmbed(ms, forcedErrorResp), nil
			}

			return SimpleEmbedResponse(ms, "Claimed **%s** as your waifu using **%s%d**!", target.Username, conf.CurrencySymbol, claimAmount), nil
		},
	}
	WaifuCmdReset = &commands.YAGCommand{
		CmdCategory: CategoryWaifu,
		Name:        "Reset",
		Description: "Resets your waifu stats, keeping your current waifus",
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())

			account.WaifuItems = nil
			account.WaifuItemWorth = 0

			_, err := account.UpdateG(parsed.Context(), boil.Whitelist("waifu_items", "waifu_item_worth"))
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(commands.ContextMS(parsed.Context()), "Reset your waifu stats, keeping your waifus"), nil
		},
	}
	WaifuCmdTransfer = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "Transfer",
		Description:  "Transfer the ownership of one of your waifus to another user. You must pay 10% of your waifu's value.",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Waifu", Type: &commands.MemberArg{}},
			&dcmd.ArgDef{Name: "New Owner", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())
			account := CtxUser(parsed.Context())

			waifu := parsed.Args[0].Value.(*dstate.MemberState)
			newOwner := parsed.Args[1].Value.(*dstate.MemberState)

			conf := CtxConfig(parsed.Context())
			_, _, err := GetCreateAccount(parsed.Context(), newOwner.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			waifuAccount, _, err := GetCreateAccount(parsed.Context(), waifu.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			if !common.ContainsInt64Slice(account.Waifus, waifu.ID) {
				return ErrorEmbed(ms, "That person is not your waifu >:u"), nil
			}

			if newOwner.ID == waifu.ID {
				return ErrorEmbed(ms, "Can't transfer the waifu to itself!?"), nil
			}

			if newOwner.ID == ms.ID {
				return ErrorEmbed(ms, "Can't transfer the waifu to yourself!?"), nil
			}

			worth := WaifuWorth(waifuAccount)
			transferFee := int64(float64(worth) * 0.1)
			if account.MoneyWallet < transferFee {
				return ErrorEmbed(ms, "Not enough money in your wallet to transfer this waifu (costs %s%d)", conf.CurrencySymbol, transferFee), nil
			}

			err = common.SqlTX(func(tx *sql.Tx) error {
				// update old owner account
				result, err := tx.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3, waifus = array_remove(waifus, $4) WHERE guild_id = $1 AND user_id = $2 AND $5 <@ waifus",
					parsed.GS.ID, ms.ID, transferFee, waifu.ID, pq.Int64Array([]int64{waifu.ID}))
				if err != nil {
					return err
				}

				rows, err := result.RowsAffected()
				if err != nil {
					return err
				}

				if rows < 1 {
					return errors.New("failed updating tables, no rows, most likely a race condition")
				}

				// update new owner account
				_, err = tx.Exec("UPDATE economy_users SET waifus = waifus || $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, newOwner.ID, pq.Int64Array([]int64{waifu.ID}))
				if err != nil {
					return err
				}

				// update waifu
				_, err = tx.Exec("UPDATE economy_users SET waifud_by = $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, waifu.ID, newOwner.ID)
				return err
			})
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Transferred **%s** to **%s** for **%s%d**", waifu.Username, newOwner.Username, conf.CurrencySymbol, transferFee), nil
		},
	}
	WaifuCmdDivorce = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "Divorce",
		Description:  "Releases your claim on a specific waifu. You will get some of the money you've spent back unless that waifu has an affinity towards you. 6 hours cooldown.",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Waifu", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())

			waifu := parsed.Args[0].Value.(*dstate.MemberState)

			conf := CtxConfig(parsed.Context())

			waifuAccount, _, err := GetCreateAccount(parsed.Context(), waifu.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			if waifuAccount.WaifudBy != ms.ID {
				return ErrorEmbed(ms, "This person is not your waifu..."), nil
			}

			worth := WaifuWorth(waifuAccount)
			moneyBack := int64(float64(worth) * 0.5)
			if waifuAccount.WaifuAffinityTowards == ms.ID {
				// you get no money back >:u
				moneyBack = 0
			}

			err = common.SqlTX(func(tx *sql.Tx) error {
				// update old owner account
				result, err := tx.Exec("UPDATE economy_users SET money_wallet = money_wallet + $3, waifus = array_remove(waifus, $4) WHERE guild_id = $1 AND user_id = $2 AND $5 <@ waifus",
					parsed.GS.ID, ms.ID, moneyBack, waifu.ID, pq.Int64Array([]int64{waifu.ID}))
				if err != nil {
					return err
				}

				rows, err := result.RowsAffected()
				if err != nil {
					return err
				}

				if rows < 1 {
					return errors.New("failed updating tables, no rows, most likely a race condition")
				}

				// update waifu
				_, err = tx.Exec("UPDATE economy_users SET waifud_by = 0 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, waifu.ID)
				return err
			})

			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "You're now divorced with **%s**, you got back **%s%d**", waifu.Username, conf.CurrencySymbol, moneyBack), nil
		},
	}
	WaifuCmdAffinity = &commands.YAGCommand{
		CmdCategory: CategoryWaifu,
		Name:        "Affinity",
		Description: "Sets your affinity towards someone you want to be claimed by. Setting affinity will reduce their claim on you by 20%. Provide no parameters to clear your affinity.",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			// check cooldown
			var cdResp string
			err := common.RedisPool.Do(radix.Cmd(&cdResp, "SET", fmt.Sprintf("economy_affinity_cd:%d", parsed.GS.ID), "1", "EX", "1800", "NX"))
			if err != nil {
				return nil, err
			}

			if cdResp != "OK" {
				return ErrorEmbed(ms, "This command is still on cooldown"), nil
			}

			resp := ""
			if parsed.Args[0].Value == nil {
				account.WaifuAffinityTowards = 0
				resp = "Reset your affinity to no-one."
			} else {
				targetMS := parsed.Args[0].Value.(*dstate.MemberState)
				resp = "Set your affinity towards " + targetMS.Username
				account.WaifuAffinityTowards = targetMS.ID
			}

			_, err = account.UpdateG(parsed.Context(), boil.Whitelist("waifu_affinity_towards"))
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(commands.ContextMS(parsed.Context()), resp), nil
		},
	}
	WaifuCmdGift = &commands.YAGCommand{
		CmdCategory: CategoryWaifu,
		Name:        "Gift",
		Description: "Gift an item to someone. This will increase their waifu value by 50% of the gifted item's value if you are not their waifu, or 95% if you are. Provide no parameters to see a list of items that you can gift.",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Item", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Target", Type: &commands.MemberArg{}},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())
			ms := commands.ContextMS(parsed.Context())

			if parsed.Args[0].Value == nil {
				return ListWaifuItems(parsed.GS.ID, parsed.CS.ID, ms, account.MoneyWallet, conf.CurrencySymbol)
			}

			if parsed.Args[1].Value == nil {
				return ErrorEmbed(ms, "Re-run the command with the user you wanna gift it to included"), nil
			}

			itemToBuy, err := models.FindEconomyWaifuItemG(parsed.Context(), parsed.GS.ID, parsed.Args[0].Int64())
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "Unknown item"), nil
				}

				return nil, err
			}

			if int64(itemToBuy.Price) > account.MoneyWallet {
				return ErrorEmbed(ms, "You don't have enough money in your wallet to gift this item"), nil
			}

			// ensure the target has a wallet
			target := parsed.Args[1].Value.(*dstate.MemberState)
			_, _, err = GetCreateAccount(parsed.Context(), target.ID, parsed.GS.ID, conf.StartBalance)
			if err != nil {
				return nil, err
			}

			worthIncreaseModifier := 0.5
			if account.WaifudBy == target.ID {
				worthIncreaseModifier = 0.95
			}

			err = common.SqlTX(func(tx *sql.Tx) error {
				// deduct money from our account
				err = TransferMoneyWallet(parsed.Context(), tx, conf, false, ms.ID, common.BotUser.ID, int64(itemToBuy.Price), int64(itemToBuy.Price))
				// _, err := tx.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3 WHERE guild_id = $1 AND user_id = $2", parsed.GS.ID, ms.ID, itemToBuy.Price)
				if err != nil {
					return err
				}

				worthIncrease := int64(float64(itemToBuy.Price) * worthIncreaseModifier)

				// add the item and increase their worth
				_, err = tx.Exec("UPDATE economy_users SET waifu_item_worth = waifu_item_worth + $3 WHERE guild_id = $1 AND user_id = $2",
					parsed.GS.ID, target.ID, worthIncrease)
				return err
			})
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Gifted **%s** to **%s** for **%s%d**!", itemToBuy.Name, target.Username, conf.CurrencySymbol, itemToBuy.Price), nil
		},
	}

	/*
		Shop
	*/

	WaifuShopAdd = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "ItemAdd",
		Description:  "Adds an item to the waifu shop, only economy adins can use this",
		RequiredArgs: 3,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Price", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Icon", Type: dcmd.String},
			&dcmd.ArgDef{Name: "Name", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			localID, err := common.GenLocalIncrID(parsed.GS.ID, "economy_item")
			if err != nil {
				return nil, err
			}

			m := &models.EconomyWaifuItem{
				GuildID: parsed.GS.ID,
				LocalID: localID,
				Price:   parsed.Args[0].Int(),
				Icon:    parsed.Args[1].Str(),
				Name:    parsed.Args[2].Str(),
			}

			err = m.InsertG(parsed.Context(), boil.Infer())
			if err != nil {
				return nil, err
			}

			conf := CtxConfig(parsed.Context())
			return SimpleEmbedResponse(commands.ContextMS(parsed.Context()), "Added **%s** to the shop at the price of **%s%d**, it received the ID **%d**",
				m.Name, conf.CurrencySymbol, m.Price, m.LocalID), nil
		},
	}
	WaifuShopEdit = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "ItemEdit",
		Description:  "Edits an item in the waifu shop",
		RequiredArgs: 4,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Item", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Price", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Icon", Type: dcmd.String},
			&dcmd.ArgDef{Name: "Name", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {

			ms := commands.ContextMS(parsed.Context())

			item, err := models.FindEconomyWaifuItemG(parsed.Context(), parsed.GS.ID, parsed.Args[0].Int64())
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No item by that id"), nil
				}

				return nil, err
			}

			item.Price = parsed.Args[1].Int()
			item.Icon = parsed.Args[2].Str()
			item.Name = parsed.Args[3].Str()

			_, err = item.UpdateG(parsed.Context(), boil.Whitelist("price", "icon", "name"))
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Updated **%s**", item.Name), nil
		},
	}
	WaifuCmdDel = &commands.YAGCommand{
		CmdCategory:  CategoryWaifu,
		Name:         "ItemDel",
		Description:  "Removes a item from the waifu shop",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Item", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {

			numDeleted, err := models.EconomyWaifuItems(qm.Where("local_id = ? AND guild_id = ?", parsed.Args[0].Int(), parsed.GS.ID)).DeleteAll(parsed.Context(), common.PQ)
			if err != nil {
				return nil, err
			}

			ms := commands.ContextMS(parsed.Context())
			if numDeleted < 1 {
				return ErrorEmbed(ms, "No item by that ID"), nil
			}

			return SimpleEmbedResponse(ms, "Deleted item ID %d", parsed.Args[0].Int()), nil
		},
	}
)

func WaifuWorth(target *models.EconomyUser) int64 {
	const base = 100

	worth := base + target.WaifuItemWorth
	return worth
}

func WaifuCost(from, target *models.EconomyUser) int64 {
	worth := WaifuWorth(target)

	cost := int64(float64(worth) * 1.1)

	if from != nil {
		if target.WaifuAffinityTowards == from.UserID {
			cost = int64(float64(cost) * 0.8)
		}
	}

	return cost

}

func ListWaifuItems(guildID, channelID int64, ms *dstate.MemberState, currentMoney int64, currencySymbol string) (*discordgo.MessageEmbed, error) {
	_, err := paginatedmessages.CreatePaginatedMessage(guildID, channelID, 1, 0, func(p *paginatedmessages.PaginatedMessage, newPage int) (*discordgo.MessageEmbed, error) {

		offset := (newPage - 1) * 12

		items, err := models.EconomyWaifuItems(models.EconomyWaifuItemWhere.GuildID.EQ(guildID), qm.OrderBy("local_id asc"), qm.Limit(12), qm.Offset(offset)).AllG(context.Background())
		if err != nil {
			return nil, err
		}

		embed := SimpleEmbedResponse(ms, "")
		if len(items) < 1 {
			embed.Description = "No items :("
		}

		embed.Title = "Waifu gift shot!"

		for _, v := range items {
			extraVal := ""
			if int64(v.Price) > currentMoney {
				extraVal = " *(you can't afford)*"
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("#%d - %s %s", v.LocalID, v.Icon, v.Name),
				Value:  fmt.Sprintf("%s%d%s", currencySymbol, v.Price, extraVal),
				Inline: true,
			})
		}

		return embed, nil
	})

	return nil, err

}
