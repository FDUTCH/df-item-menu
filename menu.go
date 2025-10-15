package df_item_menu

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
)

const actionKey = "action_key"

type Menu struct {
	callbacks map[int]func(*player.Player)
	items     []Item
}

func (m *Menu) Set(pl *player.Player) {
	inv := pl.Inventory()
	inv.Clear()
	for _, it := range m.items {
		st, slot := it.Stack()
		_ = inv.SetItem(slot, st)
	}
}

func NewMenu(items ...Item) *Menu {
	return &Menu{callbacks: map[int]func(*player.Player){}, items: items}
}

func (m *Menu) HandleItemUse(ctx *player.Context) {
	pl := ctx.Val()
	main, _ := pl.HeldItems()
	if main.Empty() {
		return
	}
	slot, _ := pl.Inventory().First(main)
	if callback, ok := m.callbacks[slot]; ok {
		ctx.Cancel()
		callback(pl)
	}

	val, ok := main.Value(actionKey)
	if ok {
		ctx.Cancel()
		reflect.ValueOf(pl.Handler()).
			MethodByName(val.(string)).
			Call([]reflect.Value{reflect.ValueOf(pl)})
	}
}

func (m *Menu) WithItem(item Item) *MenuWithSlot {
	m.items = append(m.items, item)
	return &MenuWithSlot{
		Menu:     m,
		lastSlot: item.Slot,
	}
}

func (m *Menu) Slot(slot int) *MenuWithSlot {
	return &MenuWithSlot{
		Menu:     m,
		lastSlot: slot,
	}
}

type MenuWithSlot struct {
	*Menu
	lastSlot int
}

func (m *MenuWithSlot) WithCallback(callback func(*player.Player)) *Menu {
	m.callbacks[m.lastSlot] = callback
	return m.Menu
}
