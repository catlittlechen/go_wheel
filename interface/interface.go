package util

import (
	"fmt"
	"reflect"
)

// InterfaceToSimpleStruct is a simple unmarshal function
// cover the map[string]interface{} to Struct with reflect
func InterfaceToSimpleStruct(tag string, inte interface{}, structs interface{}) (err error) {
	mapInterface, ok := inte.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("cannot use %s as Map", reflect.ValueOf(inte).Kind())
		return
	}
	sv := reflect.ValueOf(structs)
	if sv.Kind() != reflect.Ptr || sv.IsNil() {
		err = fmt.Errorf("cannot use %s as Ptr or is NIL", sv.Kind())
		return
	}
	s2 := sv.Type().Elem()
	s := sv.Elem()
	if s.Kind() != reflect.Struct {
		err = fmt.Errorf("cannot use %s as Struct", s.Kind())
		return err
	}
	mapTagToValue := make(map[string]string)
	for i := 0; i < s2.NumField(); i++ {
		mapTagToValue[string(s2.Field(i).Tag)] = s2.Field(i).Name
	}
	for k, v := range mapInterface {
		f := s.FieldByName(mapTagToValue[tag+":\""+k+"\""])
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if reflect.ValueOf(v).Kind() == reflect.Float64 {
				if !f.OverflowInt(int64(v.(float64))) {
					f.SetInt(int64(v.(float64)))
				} else {
					err = fmt.Errorf("the value is over flow the Int64")
				}
			} else {
				err = fmt.Errorf("cannot cover type Int from type %s", reflect.ValueOf(v).Kind())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if reflect.ValueOf(v).Kind() == reflect.Float64 {
				if !f.OverflowUint(uint64(v.(float64))) {
					f.SetUint(uint64(v.(float64)))
				} else {
					err = fmt.Errorf("the value is over flow the Uint64")
				}
			} else {
				err = fmt.Errorf("cannot cover type Uint from type %s", reflect.ValueOf(v).Kind())
			}
		case reflect.Float32, reflect.Float64:
			if reflect.ValueOf(v).Kind() == reflect.Float64 {
				if !f.OverflowFloat(v.(float64)) {
					f.SetFloat(v.(float64))
				} else {
					err = fmt.Errorf("the value is over flow the Float64")
				}
			} else {
				err = fmt.Errorf("cannot cover type Float from type %s", reflect.ValueOf(v).Kind())
			}
		case reflect.String:
			if reflect.ValueOf(v).Kind() == reflect.String {
				f.SetString(v.(string))
			} else {
				err = fmt.Errorf("cannot cover type String from type %s", reflect.ValueOf(v).Kind())
			}
		case reflect.Bool:
			if reflect.ValueOf(v).Kind() == reflect.Bool {
				f.SetBool(v.(bool))
			} else {
				err = fmt.Errorf("cannot cover type Bool from type %s", reflect.ValueOf(v).Kind())
			}
		default:
			err = fmt.Errorf("this struct is complicated")
		}
		if err != nil {
			return
		}
	}
	return nil
}
