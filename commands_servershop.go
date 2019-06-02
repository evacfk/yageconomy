package yageconomy

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/yageconomy/models"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/bot/paginatedmessages"
	"github.com/jonas747/yagpdb/commands"
	"github.com/jonas747/yagpdb/common"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"strconv"
	"strings"
)

const (
	ItemTypeRole          = 0
	ItemTypeList          = 1
	ItemTypeGamblingBoost = 2
)

var ShopCommands = []*commands.YAGCommand{
	&commands.YAGCommand{
		CmdCategory: CategoryEconomy,
		Name:        "Shop",
		Description: "Shows the items available in the server shop",
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())

			_, err := paginatedmessages.CreatePaginatedMessage(parsed.GS.ID, parsed.CS.ID, 1, 0, func(p *paginatedmessages.PaginatedMessage, newPage int) (*discordgo.MessageEmbed, error) {
				offset := (newPage - 1) * 12

				items, err := models.EconomyShopItems(
					models.EconomyShopItemWhere.GuildID.EQ(parsed.GS.ID),
					qm.OrderBy("local_id asc"),
					qm.Limit(12),
					qm.Offset(offset),
				).AllG(context.Background())
				if err != nil {
					return nil, err
				}

				embed := SimpleEmbedResponse(ms, "")
				embed.Title = "Server Shop!"

				for _, v := range items {
					name := v.Name
					if v.RoleID != 0 {
						r := parsed.GS.RoleCopy(true, v.RoleID)
						if r != nil {
							name = r.Name
						} else {
							name += "(deleted-role)"
						}
					}

					canAffordStr := ""
					if account.MoneyWallet < v.Cost {
						canAffordStr = "(you can't afford)"
					}

					typStr := "role"
					if v.Type == ItemTypeList {
						typStr = "list"
					} else if v.Type == ItemTypeGamblingBoost {
						typStr = "gambling boost"
					}

					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:   fmt.Sprintf("#%d - %s", v.LocalID, name),
						Value:  fmt.Sprintf("%s - %d%s%s", typStr, v.Cost, conf.CurrencySymbol, canAffordStr),
						Inline: true,
					})
				}
				if len(items) < 1 {
					embed.Description = "(no items)"
				}

				return embed, nil
			})

			return nil, err
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "Buy",
		Description:  "Buys an item from the server shop",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "item", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())
			account := CtxUser(parsed.Context())
			conf := CtxConfig(parsed.Context())

			shopItem, err := models.FindEconomyShopItemG(parsed.Context(), parsed.GS.ID, parsed.Args[0].Int64())
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No shop item with that ID"), nil
				}

				return nil, err
			}

			if shopItem.Cost > account.MoneyWallet {
				return ErrorEmbed(ms, "Not enough money in your wallet :("), nil
			}

			forcedErrorMsg := ""

			// do this in a transaction so we can rollback the purchase if the delivery failed
			err = common.SqlTX(func(tx *sql.Tx) error {
				listValue := ""

				if shopItem.Type == ItemTypeList {

					query := `UPDATE economy_shop_list_items
SET    purchased_by = $3
WHERE  value = (
	SELECT value
	FROM   economy_shop_list_items
	WHERE  purchased_by = 0 AND guild_id = $1 AND list_id = $2
	LIMIT  1
	FOR    UPDATE
)
RETURNING value;`

					row := tx.QueryRow(query, parsed.GS.ID, shopItem.LocalID, ms.ID)

					err = row.Scan(&listValue)
					if err != nil {
						if errors.Cause(err) == sql.ErrNoRows {
							forcedErrorMsg = "No more items in that list available"
							return nil
						}
						return err
					}
				}

				err = TransferMoneyWallet(parsed.Context(), tx, conf, false, ms.ID, common.BotUser.ID, shopItem.Cost, shopItem.Cost)
				// _, err = tx.Exec("UPDATE economy_users SET money_wallet = money_wallet - $3, gambling_boost_percentage = gambling_boost_percentage + $4 WHERE guild_id = $1 AND user_id = $2",
				// parsed.GS.ID, ms.ID, shopItem.Cost, shopItem.GamblingBoostPercentage)

				if err != nil {
					return err
				}

				// deliver the item
				if shopItem.Type == ItemTypeList {
					err = bot.SendDMEmbed(ms.ID, SimpleEmbedResponse(ms, "You purhcased one of **%s**, here it is: ||%s||", shopItem.Name, listValue))
				} else if shopItem.Type == ItemTypeRole {
					err = common.AddRoleDS(ms, shopItem.RoleID)
				}

				return err

			})
			if forcedErrorMsg != "" {
				return ErrorEmbed(ms, forcedErrorMsg), err
			}

			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "You purchased **%s** for **%d%s**!", shopItem.Name, shopItem.Cost, conf.CurrencySymbol), nil
		},
	},
}

var ShopAdminCommands = []*commands.YAGCommand{
	&commands.YAGCommand{
		CmdCategory:     CategoryEconomy,
		Name:            "ShopAdd",
		Description:     "Adds an item to the shop, only economy admins can use this command",
		LongDescription: "Types are 'role', 'list' and 'gamblingboostx[percentage]' where percentage is the gambling boost percentage\n\nExample: -shopadd gamblingboostx10 1000 10% gambling income increase\nThat will add a item with 10% gambling income increase and with the name '10% gambling income increase'",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "Type", Type: dcmd.String},
			&dcmd.ArgDef{Name: "Price", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Name", Type: dcmd.String},
		},
		RequiredArgs: 3,
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())

			tStr := strings.ToLower(parsed.Args[0].Str())

			t := ItemTypeRole
			gamblingBoostPercentage := int64(0)
			if tStr == "list" {
				t = ItemTypeList
			} else if strings.HasPrefix(tStr, "gamblingboostx") {
				t = ItemTypeGamblingBoost
				split := strings.Split(tStr, "x")
				if len(split) < 2 {
					return ErrorEmbed(ms, "No boost percentage specified, example: `gamblingboostx10` for 10%%"), nil
				}

				var err error
				gamblingBoostPercentage, err = strconv.ParseInt(split[1], 10, 32)
				if err != nil {
					return nil, err
				}
			}

			lID, err := common.GenLocalIncrID(parsed.GS.ID, "economy_shop_item")
			if err != nil {
				return nil, err
			}

			roleID := int64(0)
			name := parsed.Args[2].Str()

			// this is a role
			if t == ItemTypeRole {
				parsed.GS.RLock()
				for _, v := range parsed.GS.Guild.Roles {
					if strings.EqualFold(v.Name, name) {
						roleID = v.ID
						name = v.Name
						break
					}
				}
				parsed.GS.RUnlock()
				if roleID == 0 {
					return ErrorEmbed(ms, "Unknown role %q", name), nil
				}
			}

			m := &models.EconomyShopItem{
				GuildID: parsed.GS.ID,
				LocalID: lID,

				Type: int16(t),

				Cost:                    int64(parsed.Args[1].Int()),
				Name:                    name,
				RoleID:                  roleID,
				GamblingBoostPercentage: int(gamblingBoostPercentage),
			}

			err = m.InsertG(parsed.Context(), boil.Infer())
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Added **%s** to the shop, it was given the ID **%d**", name, lID), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "ShopListAdd",
		Description:  "Adds a item to a shop list, only economy admin can use this command",
		RequiredArgs: 2,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "List ID", Type: dcmd.Int},
			&dcmd.ArgDef{Name: "Item", Type: dcmd.String},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())

			shopItem, err := models.FindEconomyShopItemG(parsed.Context(), parsed.GS.ID, parsed.Args[0].Int64())
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No shop item with that ID"), nil
				}

				return nil, err
			}

			if shopItem.Type != 1 {
				return ErrorEmbed(ms, "That shop item is not a list"), nil
			}

			lID, err := common.GenLocalIncrID(parsed.GS.ID, "economy_shop_list_item")
			if err != nil {
				return nil, err
			}

			m := &models.EconomyShopListItem{
				GuildID: parsed.GS.ID,
				LocalID: lID,
				ListID:  shopItem.LocalID,
				Value:   parsed.Args[1].Str(),
			}

			err = m.InsertG(parsed.Context(), boil.Infer())
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Added to the list **%s**", shopItem.Name), nil
		},
	},
	&commands.YAGCommand{
		CmdCategory:  CategoryEconomy,
		Name:         "ShopRem",
		Aliases:      []string{"ShopDel"},
		Description:  "Removes a item from the shop, only economy admins can use this command",
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "ID", Type: dcmd.Int},
		},
		RunFunc: func(parsed *dcmd.Data) (interface{}, error) {
			ms := commands.ContextMS(parsed.Context())

			shopItem, err := models.FindEconomyShopItemG(parsed.Context(), parsed.GS.ID, parsed.Args[0].Int64())
			if err != nil {
				if errors.Cause(err) == sql.ErrNoRows {
					return ErrorEmbed(ms, "No shop item with that ID"), nil
				}

				return nil, err
			}
			_, err = shopItem.DeleteG(parsed.Context())
			if err != nil {
				return nil, err
			}

			return SimpleEmbedResponse(ms, "Deleted **%s** from the server shop", shopItem.Name), nil
		},
	},
}
