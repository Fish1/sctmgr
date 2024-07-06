package mainscreen

import (
	"github.com/charmbracelet/bubbles/list"
)

func New() Model {

	items := []list.Item{}

	items = append(items, item(ScreenItem{
		name: "Glorious Eggroll Manager",
	}))

	items = append(items, item(ScreenItem{
		name: "Prefix Manager",
	}))

	list := list.New(items, itemDelegate{}, 40, 20)

	return Model{
		list: list,
	}
}
