package df_item_menu

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Item struct {
	ID         string `json:"id"`
	Meta       int    `json:"meta"`
	Slot       int    `json:"slot"`
	CustomName string `json:"custom_name"`
	Action     string `json:"action"`
}

func (i Item) Stack() (item.Stack, int) {
	if i.ID == "" {
		return item.Stack{}, i.Slot
	}
	it, ok := world.ItemByName(i.ID, int16(i.Meta))
	if !ok {
		panic(fmt.Errorf("item %s with meta %d not found", i.ID, i.Meta))
	}
	if i.Slot > 8 || i.Slot < 0 {
		panic(fmt.Errorf("invalid slot %d", i.Slot))
	}

	stack := item.NewStack(it, 1).WithCustomName(i.CustomName)
	if i.Action != "" {
		stack = stack.WithValue(actionKey, i.Action)
	}
	return stack, i.Slot
}

func (i Item) WithCustomName(name string) Item {
	i.CustomName = name
	return i
}

func (i Item) WithAction(action string) Item {
	i.Action = action
	return i
}

func (i Item) WithSlot(slot int) Item {
	i.Slot = slot
	return i
}

func ItemFromStack(st item.Stack) Item {
	if st.Empty() {
		return Item{}
	}
	id, meta := st.Item().EncodeItem()
	var action string
	if actionVal, ok := st.Value(actionKey); ok {
		action = actionVal.(string)
	}

	return Item{
		ID:         id,
		Meta:       int(meta),
		CustomName: st.CustomName(),
		Action:     action,
	}
}
