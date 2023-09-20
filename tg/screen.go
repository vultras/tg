package tg

// Unique identifier for the screen.
type ScreenId string

// Screen statement of the bot.
// Mostly what buttons to show.
type Screen struct {
	// Unique identifer to change to the screen
	// via Context.ChangeScreen method.
	Id ScreenId
	// The widget to run when reaching the screen.
	Widget Widget
}

// The node is a simple way to represent
// tree-like structured applications.
type Node struct {
	Screen *Screen
	Subs []*Node
}

func NewNode(id ScreenId, widget Widget, subs ...*Node) *Node {
	ret := &Node{}
	ret.Screen = NewScreen(id, widget)
	ret.Subs = subs
	return ret
}

func (n *Node) ScreenMap() ScreenMap {
	m := make(ScreenMap)
	id := n.Screen.Id
	m[id] = n.Screen
	n.Screen.Id = id
	var root ScreenId
	if id == "/" {
		root = ""
	} else {
		root = id
	}
	for _, sub := range n.Subs {
		buf := sub.screenMap(root + "/")
		for k, v := range buf {
			m[k] = v
		}
	}
	return m
}

func (n *Node) screenMap(root ScreenId) ScreenMap {
	m := make(ScreenMap)
	id := root+ n.Screen.Id
	m[id] = n.Screen
	n.Screen.Id = id
	for _, sub := range n.Subs {
		buf := sub.screenMap(id + "/")
		for k, v := range buf {
			m[k] = v
		}
	}
	return m
}

// Map structure for the screens.
type ScreenMap map[ScreenId]*Screen

// Returns the new screen with specified name and widget.
func NewScreen(id ScreenId, widget Widget) *Screen {
	return &Screen{
		Id: id,
		Widget: widget,
	}
}

