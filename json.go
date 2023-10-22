package tg

import (
	//"encoding/json"
	//"reflect"
)

/*func (a *action) UnmarshalJSON(data []byte) error {
	var err error
	m := make(map[string]any)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	bts, err := json.Marshal(m["Action"])
	if err != nil {
		return err
	}

	a.Type = m["Type"].(string)
	typ := actionMapByTypeName[a.Type].(reflect.Type)
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	vr := reflect.New(typ).Interface().(Action)
	err = json.Unmarshal(bts, vr)
	if err != nil {
		return err
	}

	a.Action = vr
	return nil
}*/
