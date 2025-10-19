package df_item_menu

import (
	"reflect"

	"github.com/FDUTCH/patterns"
	"github.com/df-mc/dragonfly/server/player"
)

const (
	useActionKey   = "use_action_key"
	clickActionKey = "click_action_key"
)

type Menu[T player.Handler] struct {
	useCallbacks   map[int]func(*player.Player)
	clickCallbacks map[int]func(*player.Player)
	items          []Item
}

func (m *Menu[T]) Set(pl *player.Player) {
	inv := pl.Inventory()
	inv.Clear()
	for _, it := range m.items {
		st, slot := it.Stack()
		_ = inv.SetItem(slot, st)
	}
}

func NewMenu[T player.Handler](items ...Item) *Menu[T] {
	return &Menu[T]{useCallbacks: make(map[int]func(*player.Player)), clickCallbacks: make(map[int]func(*player.Player)), items: items}
}

func (m *Menu[T]) HandleItemUse(ctx *player.Context) {
	m.trigger(ctx, false)
}

func (m *Menu[T]) HandleClick(ctx *player.Context) {
	m.trigger(ctx, true)
}

func (m *Menu[T]) trigger(ctx *player.Context, click bool) {
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
		handler, ok := patterns.UnwrapUntilCast[player.Handler, T](pl.Handler())
		if ok {
			reflect.ValueOf(handler).
				MethodByName(val.(string)).
				Call([]reflect.Value{reflect.ValueOf(pl)})
		}
	}
}

func (m *Menu[T]) WithItem(item Item) *MenuWithSlot[T] {
	m.items = append(m.items, item)
	return &MenuWithSlot[T]{
		Menu:     m,
		lastSlot: item.Slot,
	}
}

func (m *Menu[T]) Slot(slot int) *MenuWithSlot[T] {
	return &MenuWithSlot[T]{
		Menu:     m,
		lastSlot: slot,
	}
}

type MenuWithSlot[T player.Handler] struct {
	*Menu[T]
	lastSlot int
}

func (m *MenuWithSlot[T]) WithUseCallback(callback func(*player.Player)) *Menu[T] {
	m.useCallbacks[m.lastSlot] = callback
	return m.Menu
}

func (m *MenuWithSlot[T]) WithClickCallback(callback func(*player.Player)) *Menu[T] {
	m.clickCallbacks[m.lastSlot] = callback
	return m.Menu
}
