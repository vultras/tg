package tg

import (
	"reflect"
)

var (
	actionMapByReflect  = make(map[reflect.Type]string)
	actionMapByTypeName = make(map[string]reflect.Type)
	Init                func()
)

func initEncoding() {
	actions := map[string]Action{
		"action-func":   ActionFunc(nil),
		"screen-change": ScreenGo{},
	}
	for k, action := range actions {
		DefineAction(k, action)
	}
}

// Define interface to make it marshalable to JSON etc.
// Like in GOB. Must be done both on client and server
// if one is provided.
func DefineAction(typeName string, a Action) error {
	t := reflect.TypeOf(a)

	actionMapByReflect[t] = typeName
	actionMapByTypeName[typeName] = t

	return nil
}

func DefineGroupAction(typ string, a GroupAction) error {
	return nil
}
