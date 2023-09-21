package tg

import (
	"path"
)

// The type implements changing screen to the underlying ScreenId
type ScreenGo struct {
	Path Path
	Args []any
}

func (sc ScreenGo) Act(c *Context) {
	err := c.Go(sc.Path, sc.Args...)
	if err != nil {
		panic(err)
	}
}

// The same as Act.
func (sc ScreenGo) Serve(c *Context) {
	sc.Act(c)
}

// Unique identifier for the screen
// and relative paths to the screen.
type Path string

// Returns true if the path is empty.
func (p Path) IsEmpty() bool {
	return p == ""
}

// Returns true if the path is absolute.
func (p Path) IsAbs() bool {
	if len(p) == 0 {
		return false
	}
	return p[0] == '/'
}

func (p Path) Dir() Path {
	return Path(path.Dir(string(p)))
}

// Clean the path deleting exceed ., .. and / .
func (p Path) Clean() Path {
	return Path(path.Clean(string(p)))
}

// Screen statement of the bot.
// Mostly what buttons to show.
type Screen struct {
	// The widget to run when reaching the screen.
	Widget Widget
}

// The first node with the "/" path.
type RootNode struct {
	Screen *Screen
	Subs []*Node
}

// The node is a simple way to represent
// tree-like structured applications.
type Node struct {
	Path Path
	Screen *Screen
	Subs []*Node
}

// Return new root node with the specified widget in the screen.
func NewRootNode(widget Widget, subs ...*Node) *RootNode {
	ret := &RootNode{}
	ret.Screen = NewScreen(widget)
	ret.Subs = subs
	return ret
}

func NewNode(relPath Path, widget Widget, subs ...*Node) *Node {
	ret := &Node{}
	ret.Path = relPath
	ret.Screen = NewScreen(widget)
	ret.Subs = subs
	return ret
}

func (n *RootNode) ScreenMap() ScreenMap {
	m := make(ScreenMap)
	var root Path = "/"
	m[root] = n.Screen
	for _, sub := range n.Subs {
		buf := sub.ScreenMap(root)
		for k, v := range buf {
			_, ok := m[k]
			if ok {
				panic("duplicate paths in node definition")
			}
			m[k] = v
		}
	}
	return m
}

func (n *Node) ScreenMap(root Path) ScreenMap {
	m := make(ScreenMap)
	pth := (root + n.Path).Clean()
	m[pth] = n.Screen
	for _, sub := range n.Subs {
		buf := sub.ScreenMap(pth + "/")
		for k, v := range buf {
			_, ok := m[k]
			if ok {
				panic("duplicate paths in node definition")
			}
			m[k] = v
		}
	}
	return m
}

// Map structure for the screens.
type ScreenMap map[Path] *Screen

// Returns the new screen with specified name and widget.
func NewScreen(widget Widget) *Screen {
	return &Screen{
		Widget: widget,
	}
}

