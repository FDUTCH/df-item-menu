package main

import (
	"fmt"
	"log/slog"
	"os"

	df_item_menu "github.com/FDUTCH/df-item-menu"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/restartfu/gophig"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for p := range srv.Accept() {
		p.Handle(&Handler{})
		LobbyItemsMenu.Set(p)
		_ = p
	}
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}

type Handler struct {
	player.NopHandler
}

func (h *Handler) HandleItemUse(ctx *player.Context) {
	LobbyItemsMenu.HandleItemUse(ctx)
}

func (h *Handler) OpenDuelsMenu(pl *player.Player) {
	pl.Message("OpenDuelsMenu")
}

func (h *Handler) OpenFFAMenu(pl *player.Player) {
	pl.Message("OpenFFAMenu")
}

func (h *Handler) OpenFriendsMenu(pl *player.Player) {
	pl.Message("OpenFriendsMenu")
}

func (h *Handler) OpenSpectatorMenu(pl *player.Player) {
	pl.Message("OpenSpectatorMenu")
}

var (
	LobbyMenuConfig *gophig.Gophig[[]df_item_menu.Item]
	LobbyItemsMenu  *df_item_menu.Menu
)

func init() {
	LobbyMenuConfig = gophig.NewGophig[[]df_item_menu.Item]("./lobby_items.json", gophig.JSONMarshaler{}, os.ModePerm)
	val, err := LobbyMenuConfig.LoadConf()
	if err != nil {
		_ = LobbyMenuConfig.SaveConf([]df_item_menu.Item{})
	}
	LobbyItemsMenu = df_item_menu.NewMenu(val...)
}
