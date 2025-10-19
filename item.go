package df_item_menu

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

type Item struct {
	ID          string `json:"id"`
	Meta        int    `json:"meta"`
	Slot        int    `json:"slot"`
	CustomName  string `json:"custom_name"`
	UseAction   string `json:"use_action"`
	ClickAction string `json:"click_action"`
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
	if i.UseAction != "" {
		stack = stack.WithValue(useActionKey, i.UseAction)
	}
	if i.ClickAction != "" {
		stack = stack.WithValue(clickActionKey, i.ClickAction)
	}
	return stack, i.Slot
}

func (i Item) WithCustomName(name string) Item {
	i.CustomName = name
	return i
}

func (i Item) WithUseAction(action string) Item {
	i.UseAction = action
	return i
}

func (i Item) WithClickAction(action string) Item {
	i.ClickAction = action
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
	var useAction string
	if actionVal, ok := st.Value(useActionKey); ok {
		useAction = actionVal.(string)
	}
	var clickAction string
	if actionVal, ok := st.Value(clickActionKey); ok {
		clickAction = actionVal.(string)
	}

	return Item{
		ID:          id,
		Meta:        int(meta),
		CustomName:  st.CustomName(),
		UseAction:   useAction,
		ClickAction: clickAction,
	}
}
