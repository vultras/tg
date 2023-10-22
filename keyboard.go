package tg

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The general keyboard type used both in Reply and Inline.
type Keyboard struct {
	// The action is called if there is no
	// defined action for the button.
	Action Action
	Rows []ButtonRow
	buttonMap ButtonMap
}

// Returns the new keyboard with specified rows.
func NewKeyboard(rows ...ButtonRow) *Keyboard {
	ret := &Keyboard{}
	for _, row := range rows {
		if row != nil && len(row) > 0 {
			ret.Rows = append(ret.Rows, row)
		}
	}
	return ret
}

// Adds a new button row to the current keyboard.
func (kbd *Keyboard) Row(btns ...*Button) *Keyboard {
	// For empty row. We do not need that.
	if len(btns) < 1 {
		return kbd
	}
	retBtns := []*Button{}
	for _, btn := range btns {
		if btn == nil {
			continue
		}
		retBtns = append(retBtns, btn)
	}
	// Add only if there is something to add.
	if len(retBtns) > 0 {
		kbd.Rows = append(kbd.Rows, retBtns)
	}
	return kbd
}

// Adds buttons as one column list.
func (kbd *Keyboard) List(btns ...*Button) *Keyboard {
	for _, btn := range btns {
		if btn == nil {
			continue
		}
		kbd.Rows = append(kbd.Rows, ButtonRow{btn})
	}
	return kbd
}

// Set the default action when no button provides
// key to the data we got.
func (kbd *Keyboard) WithAction(a Action) *Keyboard {
	kbd.Action = a
	return kbd
}

// Alias to WithAction but better typing when setting
// a specific function
func (kbd *Keyboard) ActionFunc(fn ActionFunc) *Keyboard {
	return kbd.WithAction(fn)
}

// Returns the map of buttons. Used to define the Action.
func (kbd Keyboard) ButtonMap() ButtonMap {
	if kbd.buttonMap != nil {
		return kbd.buttonMap
	}
	ret := make(ButtonMap)
	for _, vi := range kbd.Rows {
		for _, vj := range vi {
			ret[vj.Key()] = vj
		}
	}
	kbd.buttonMap = ret

	return ret
}

// Convert the keyboard to the more specific inline one.
func (kbd *Keyboard) Inline() *Inline {
	ret := &Inline{}
	ret.Keyboard = kbd
	return ret
}

// Convert the keyboard to the more specific reply one.
func (kbd *Keyboard) Reply() *Reply {
	ret := &Reply{}
	ret.Keyboard = kbd
	// it is used more often than not once.
	ret.OneTime = true
	return ret
}

