package reflecthelper

import (
	"fmt"
	"reflect"
)

func checkAssigner(assigner reflect.Value) (err error) {
	err = getErrIsValid(assigner)
	if err != nil {
		return
	}
	if !assigner.CanSet() {
		err = ErrAssignerCantSet
	}
	return
}

// AssignReflect assigns the val of the reflect.Value to the assigner.
// This function asserts that the assigner Kind is same as the val Kind.
func AssignReflect(assigner reflect.Value, val reflect.Value) (err error) {
	assigner = GetChildElem(assigner)
	val = GetChildElem(val)
	err = checkAssigner(assigner)
	if err != nil {
		return
	}
	err = getErrIsValid(val)
	if err != nil {
		return
	}
	err = tryAssign(assigner, val)
	return
}

func tryAssign(assigner reflect.Value, val reflect.Value) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			switch val := rec.(type) {
			case error:
				err = val
			default:
				err = fmt.Errorf("%v", val)
			}
		}
	}()

	assignerKind := GetKind(assigner)
	switch assignerKind {
	case reflect.Bool:
		var result bool
		result, err = ExtractBool(val)
		if err != nil {
			return
		}
		assigner.SetBool(result)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var result int64
		result, err = ExtractInt(val)
		if err != nil {
			return
		}
		if assigner.OverflowInt(result) {
			err = getErrOverflow(assigner)
			return
		}
		assigner.SetInt(result)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var result uint64
		result, err = ExtractUint(val)
		if err != nil {
			return
		}
		if assigner.OverflowUint(result) {
			err = getErrOverflow(assigner)
			return
		}
		assigner.SetUint(result)
	case reflect.Float32, reflect.Float64:
		var result float64
		result, err = ExtractFloat(val)
		if err != nil {
			return
		}
		if assigner.OverflowFloat(result) {
			err = getErrOverflow(assigner)
			return
		}
		assigner.SetFloat(result)
	case reflect.Complex64, reflect.Complex128:
		var result complex128
		result, err = ExtractComplex(val)
		if err != nil {
			return
		}
		if assigner.OverflowComplex(result) {
			err = getErrOverflow(assigner)
			return
		}
		assigner.SetComplex(result)
	case reflect.Array, reflect.Slice:
		switch GetKind(val) {
		case reflect.Array, reflect.Slice:
			isSlice := assignerKind == reflect.Slice
			if !isSlice {
				err = checkOverLength(assigner, val)
				if err != nil {
					return
				}
			}
			err = iterateAndAssign(assigner, val, isSlice)
		case reflect.String:
			err = iterateAndAssignString(assigner, val)
		default:
			err = getErrUnimplementedAssign(assigner, val)
		}
	case reflect.Chan:
		// TODO: Implement assigner for channel type
	case reflect.Func:
		// TODO: Implement assigner for func
	case reflect.Map:
		// TODO: Implement assigner for map
	case reflect.String:
		var result string
		result, err = ExtractString(val)
		if err != nil {
			return
		}
		assigner.SetString(result)
	case reflect.Struct:
		// TODO: Implement assigner for struct
		// TODO: Implement time.Time struct type
	default:
		err = getErrUnimplementedAssign(assigner, val)
	}
	return
}

func checkOverLength(assigner reflect.Value, val reflect.Value) (err error) {
	if assigner.Len() < val.Len() {
		err = getErrOverflowedLength(assigner, val)
	}
	return
}

func iterateAndAssign(assigner reflect.Value, val reflect.Value, isSlice bool) (err error) {
	if isSlice {
		emptySlice := reflect.MakeSlice(assigner.Type(), 0, val.Len())
		for index := 0; index < val.Len(); index++ {
			elemVal := reflect.New(GetElemType(assigner)).Elem()
			err = AssignReflect(elemVal, val.Index(index))
			if err != nil {
				return
			}
			emptySlice.Set(reflect.AppendSlice(emptySlice, elemVal))
		}
		assigner.Set(emptySlice)
	} else {
		typeArr := reflect.ArrayOf(assigner.Len(), GetElemType(assigner))
		emptyArray := reflect.New(typeArr).Elem()
		for index := 0; index < val.Len(); index++ {
			err = AssignReflect(emptyArray.Index(index), val.Index(index))
			if err != nil {
				return
			}
		}
		assigner.Set(emptyArray)
	}
	return
}

func iterateAndAssignString(assigner reflect.Value, val reflect.Value) (err error) {
	switch GetElemKind(assigner) {
	case reflect.Uint8:
		reflect.Copy(assigner, val)
	case reflect.Int32:
		// valStr := val.String()
		// runeSlice := []rune(valStr)
		// TODO: Implement this manually
	default:
		err = getErrUnimplementedAssign(assigner, val)
	}
	return
}
