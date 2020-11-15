package reflecthelper

import (
	"reflect"
	"time"
)

// List of reflect.Type used in this package
var (
	TypeRuneSlice = reflect.TypeOf([]rune{})
	TypeByteSlice = reflect.TypeOf([]byte{})
	TypeTimePtr   = reflect.TypeOf(&time.Time{})
	TypeTime      = reflect.TypeOf(time.Time{})
)

// IsTypeValueElemable checks if the type of the reflect.Value can call Elem
func IsTypeValueElemable(val reflect.Value) bool {
	return IsKindTypeElemable(GetKind(val))
}

// IsTypeElemable checks wether the typ of reflect.Type can call Elem method.
func IsTypeElemable(typ reflect.Type) (res bool) {
	if typ == nil {
		return
	}

	res = IsKindTypeElemable(typ.Kind())
	return
}

// GetElemType returns the elem type of a val of reflect.Value
func GetElemType(val reflect.Value) (typ reflect.Type) {
	if !val.IsValid() {
		return
	}

	typ = val.Type()
	if IsTypeValueElemable(val) {
		typ = typ.Elem()
	}
	return
}

// GetChildElemType returns the child elems' (root child) type of the val of reflect.Value.
func GetChildElemType(val reflect.Value) (typ reflect.Type) {
	if !val.IsValid() {
		return
	}

	typ = val.Type()
	for IsTypeElemable(typ) {
		typ = typ.Elem()
	}
	return
}
