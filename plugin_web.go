package yageconomy

import (
	"bytes"
	"database/sql"
	"github.com/ericlagergren/decimal"
	"github.com/jonas747/yageconomy/models"
	"github.com/jonas747/yagpdb/common"
	"github.com/jonas747/yagpdb/web"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/types"
	"goji.io"
	"goji.io/pat"
	"image"
	"io"
	"net/http"
)

type PostConfigForm struct {
	Enabled             bool
	Admins              []int64 `valid:"role,true"`
	CurrencyName        string  `valid:",1,50"`
	CurrencyNamePlural  string  `valid:",1,50"`
	CurrencySymbol      string  `valid:",1,50"`
	DailyFrequency      int64
	DailyAmount         int64
	ChatmoneyFrequency  int64
	ChatmoneyAmountMin  int64
	ChatmoneyAmountMax  int64
	StartBalance        int64
	AutoPlantChannels   []int64 `valid:"channel,true"`
	AutoPlantMin        int64
	AutoPlantMax        int64
	AutoPlantChance     float64
	RobFine             int
	FishingCooldown     int
	FishingMaxwinAmount int64
	FishingMinWinAmount int64
}

func (p PostConfigForm) DBModel() *models.EconomyConfig {
	return &models.EconomyConfig{
		Enabled:             p.Enabled,
		Admins:              p.Admins,
		CurrencyName:        p.CurrencyName,
		CurrencyNamePlural:  p.CurrencyNamePlural,
		CurrencySymbol:      p.CurrencySymbol,
		DailyFrequency:      p.DailyFrequency,
		DailyAmount:         p.DailyAmount,
		ChatmoneyFrequency:  p.ChatmoneyFrequency,
		ChatmoneyAmountMin:  p.ChatmoneyAmountMin,
		ChatmoneyAmountMax:  p.ChatmoneyAmountMax,
		StartBalance:        p.StartBalance,
		AutoPlantChannels:   p.AutoPlantChannels,
		AutoPlantMin:        p.AutoPlantMin,
		AutoPlantMax:        p.AutoPlantMax,
		AutoPlantChance:     types.NewDecimal(decimal.New(int64(p.AutoPlantChance*100), 4)),
		RobFine:             p.RobFine,
		FishingCooldown:     p.FishingCooldown,
		FishingMaxWinAmount: p.FishingMaxwinAmount,
		FishingMinWinAmount: p.FishingMinWinAmount,
	}
}

func (p *Plugin) InitWeb() {
	web.LoadHTMLTemplate("../../../yageconomy/assets/economy.html", "templates/plugins/economy.html")
	web.AddSidebarItem(web.SidebarCategoryFun, &web.SidebarItem{
		Name: "Economy",
		URL:  "economy",
	})

	subMux := goji.SubMux()

	web.CPMux.Handle(pat.New("/economy"), subMux)
	web.CPMux.Handle(pat.New("/economy/*"), subMux)

	subMux.Use(web.RequireGuildChannelsMiddleware)

	mainGetHandler := web.ControllerHandler(handleGetEconomy, "cp_economy_settings")

	subMux.Handle(pat.Get(""), mainGetHandler)
	subMux.Handle(pat.Get("/"), mainGetHandler)
	subMux.Handle(pat.Get("/pick_image"), http.HandlerFunc(handleGetPickImage))
	subMux.Handle(pat.Post("/pick_image"), web.ControllerPostHandler(HandleSetImage, mainGetHandler, nil, "Updated economy pick image"))
	subMux.Handle(pat.Post(""), web.ControllerPostHandler(handlePostEconomy, mainGetHandler, PostConfigForm{}, "Updated economy config"))
	subMux.Handle(pat.Post("/"), web.ControllerPostHandler(handlePostEconomy, mainGetHandler, PostConfigForm{}, "Updated economy config"))

}

func tmplFormatPercentage(in *decimal.Big) string {
	result := in.Mul(in, decimal.New(100, 0))
	return result.String()
}

func handleGetEconomy(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	g, templateData := web.GetBaseCPContextData(r.Context())

	templateData["fmtDecimalPercentage"] = tmplFormatPercentage

	if templateData["PluginSettings"] == nil {
		conf, err := models.FindEconomyConfigG(r.Context(), g.ID)
		if err != nil {
			if errors.Cause(err) == sql.ErrNoRows {
				conf = DefaultConfig(g.ID)
			} else {
				return templateData, err
			}
		}

		templateData["PluginSettings"] = conf
	}

	return templateData, nil
}

func handlePostEconomy(w http.ResponseWriter, r *http.Request) (templateData web.TemplateData, err error) {
	g, templateData := web.GetBaseCPContextData(r.Context())

	form := r.Context().Value(common.ContextKeyParsedForm).(*PostConfigForm)
	conf := form.DBModel()
	conf.GuildID = g.ID

	templateData["PluginSettings"] = conf

	err = conf.UpsertG(r.Context(), true, []string{"guild_id"}, boil.Infer(), boil.Infer())
	return templateData, nil
}

func handleGetPickImage(w http.ResponseWriter, r *http.Request) {
	g, _ := web.GetBaseCPContextData(r.Context())

	row, err := models.FindEconomyPickImageG(r.Context(), g.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}

		web.CtxLogger(r.Context()).WithError(err).Error("failed retrieving econ pick image")
		w.WriteHeader(500)
		return
	}

	w.Write(row.Image)
}

func HandleSetImage(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	g, tmpl := web.GetBaseCPContextData(ctx)

	file, header, err := r.FormFile("image")
	if err != nil {
		return tmpl, err
	}

	if header.Size > 250000 {
		return tmpl.AddAlerts(web.ErrorAlert("Max image size is 250KB")), nil
	}

	buf := make([]byte, int(header.Size))
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return tmpl, err
	}

	imgHeader, _, err := image.DecodeConfig(bytes.NewReader(buf))
	if err != nil {
		return tmpl, err
	}

	if imgHeader.Width > 1080 || imgHeader.Height > 1920 {
		return tmpl.AddAlerts(web.ErrorAlert("Max image size is 1080x1920")), nil
	}

	m := models.EconomyPickImage{
		GuildID: g.ID,
		Image:   buf,
	}

	err = m.UpsertG(r.Context(), true, []string{"guild_id"}, boil.Whitelist("image"), boil.Whitelist("guild_id", "image"))
	return tmpl, err
}
