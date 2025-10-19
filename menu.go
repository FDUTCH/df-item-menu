package df_item_menu

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
)

const (
	useActionKey   = "use_action_key"
	clickActionKey = "click_action_key"
)

type Menu struct {
	useCallbacks   map[int]func(*player.Player)
	clickCallbacks map[int]func(*player.Player)
	items          []Item
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
	return &Menu{useCallbacks: make(map[int]func(*player.Player)), clickCallbacks: make(map[int]func(*player.Player)), items: items}
}

func (m *Menu) HandleItemUse(ctx *player.Context) {
	m.trigger(ctx, false)
}

func (m *Menu) HandleClick(ctx *player.Context) {
	m.trigger(ctx, true)
}

func (m *Menu) trigger(ctx *player.Context, click bool) {
	pl := ctx.Val()
	main, _ := pl.HeldItems()
	if main.Empty() {
		return
	}
	slot, _ := pl.Inventory().First(main)

	if click {
		if callback, ok := m.clickCallbacks[slot]; ok {
			ctx.Cancel()
			callback(pl)
		}
	} else {
		if callback, ok := m.useCallbacks[slot]; ok {
			ctx.Cancel()
			callback(pl)
		}
	}

	var action = useActionKey
	if click {
		action = clickActionKey
	}

	val, ok := main.Value(action)
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

func (m *MenuWithSlot) WithUseCallback(callback func(*player.Player)) *Menu {
	m.useCallbacks[m.lastSlot] = callback
	return m.Menu
}

func (m *MenuWithSlot) WithClickCallback(callback func(*player.Player)) *Menu {
	m.clickCallbacks[m.lastSlot] = callback
	return m.Menu
}
