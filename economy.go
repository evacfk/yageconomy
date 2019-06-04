package yageconomy

//go:generate sqlboiler --no-hooks psql

import (
	"github.com/ericlagergren/decimal"
	"github.com/jonas747/yageconomy/models"
	"github.com/jonas747/yagpdb/common"
	"github.com/volatiletech/sqlboiler/types"
)

var logger = common.GetPluginLogger(&Plugin{})

func RegisterPlugin() {
	plugin := &Plugin{}
	common.InitSchema(DBSchema, "economy")
	common.RegisterPlugin(plugin)
}

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {

	return &common.PluginInfo{
		Name:     "Economy",
		SysName:  "economy",
		Category: common.PluginCategoryMisc,
	}
}

const (
	DefaultCurrencyName   = "YAGBuck"
	DefaultCurrencySymbol = "$"
)

func DefaultConfig(g int64) *models.EconomyConfig {
	return &models.EconomyConfig{
		GuildID:            g,
		CurrencyName:       DefaultCurrencyName,
		CurrencyNamePlural: DefaultCurrencyName + "s",
		CurrencySymbol:     DefaultCurrencySymbol,
		StartBalance:       1000,

		DailyFrequency: 1440,
		DailyAmount:    250,

		ChatmoneyFrequency: 100,
		ChatmoneyAmountMin: 10,
		ChatmoneyAmountMax: 50,

		FishingMaxWinAmount: 200,
		FishingMinWinAmount: 50,
		FishingCooldown:     30,

		AutoPlantMin:    50,
		AutoPlantMax:    200,
		AutoPlantChance: types.NewDecimal(decimal.New(2, 2)),
	}
}
