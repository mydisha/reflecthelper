package reflecthelper

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	Hello string
}

func TestCast(t *testing.T) {
	type args struct {
		val reflect.Value
	}
	tests := []struct {
		name     string
		args     args
		wantKind reflect.Kind
	}{
		{
			name: "invalid nil value",
			args: args{
				val: reflect.ValueOf(nil),
			},
			wantKind: reflect.Invalid,
		},
		{
			name: "valid slice value",
			args: args{
				val: reflect.ValueOf([]int{1, 2, 3}),
			},
			wantKind: reflect.Slice,
		},
		{
			name: "valid struct value",
			args: args{
				val: reflect.ValueOf(test{"Hi!"}),
			},
			wantKind: reflect.Struct,
		},
		{
			name: "valid ptr struct value",
			args: args{
				val: reflect.ValueOf(test{"Hi!"}),
			},
			wantKind: reflect.Struct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := Cast(tt.args.val)
			assert.Equal(t, tt.wantKind, GetKind(gotRes.Value))
		})
	}
}

func TestValue_IterateStruct(t *testing.T) {
	t.Run("kind is not struct", func(t *testing.T) {
		var hello int
		val := Cast(reflect.ValueOf(hello))
		assert.Nil(t, val.IterateStruct().Error())
	})
	t.Run("iterate example function", func(t *testing.T) {
		type test struct {
			Hello string
		}

		val := Cast(reflect.ValueOf(test{"Hi!"}))
		val.IterateStruct(nil, func(val reflect.Value, field reflect.Value) (err error) {
			fmt.Println(val.String())
			fmt.Println(field.String())
			return
		})
		assert.Nil(t, val.Error())
	})
	t.Run("error in the iteration", func(t *testing.T) {
		type test struct {
			Hello string
		}

		val := Cast(reflect.ValueOf(test{"Hi!"}))
		val.IterateStruct(nil, func(val reflect.Value, field reflect.Value) (err error) {
			return errors.New("random error")
		})
		assert.NotNil(t, val.Error())
	})
	t.Run("error in the iteration ignored", func(t *testing.T) {
		type test struct {
			Hello string
		}

		val := Cast(reflect.ValueOf(test{"Hi!"}), EnableIgnoreError())
		val.IterateStruct(nil, func(val reflect.Value, field reflect.Value) (err error) {
			return errors.New("random error")
		})
		assert.Nil(t, val.Error())
	})
	t.Run("panic with recoverer - error type", func(t *testing.T) {
		type test struct {
			Hello string
		}
		errTest := errors.New("test")

		val := Cast(reflect.ValueOf(test{"Hi!"}), EnablePanicRecoverer())
		val.IterateStruct(nil, func(val reflect.Value, field reflect.Value) (err error) {
			panic(errTest)
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, errTest, val.Error())
	})
	t.Run("panic with recoverer - any type", func(t *testing.T) {
		type test struct {
			Hello string
		}
		val := Cast(reflect.ValueOf(test{"Hi!"}), EnablePanicRecoverer())
		val.IterateStruct(nil, func(val reflect.Value, field reflect.Value) (err error) {
			panic("hello")
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, "hello", val.Error().Error())
	})
}

func TestValue_iterateArraySlice(t *testing.T) {
	t.Run("kind is not slice or array", func(t *testing.T) {
		var hello int
		val := Cast(reflect.ValueOf(hello))
		assert.Nil(t, val.IterateArraySlice().Error())
	})
	t.Run("iterate example function", func(t *testing.T) {
		var hello = []int{1, 2, 3, 4, 5}

		val := Cast(reflect.ValueOf(hello))
		val.IterateArraySlice(nil, func(parent reflect.Value, index int, elem reflect.Value) (err error) {
			fmt.Println("index: ", index, "elem: ", elem.Interface())
			return
		})
		assert.Nil(t, val.Error())

		var helloArray = [5]int{1, 2, 3, 4, 5}

		val = Cast(reflect.ValueOf(helloArray))
		val.IterateArraySlice(nil, func(parent reflect.Value, index int, elem reflect.Value) (err error) {
			fmt.Println("index: ", index, "elem: ", elem.Interface())
			return
		})
		assert.Nil(t, val.Error())
	})
	t.Run("error in the iteration", func(t *testing.T) {
		var hello = []int{1, 2, 3, 4, 5}

		val := Cast(reflect.ValueOf(hello))
		val.IterateArraySlice(nil, func(parent reflect.Value, index int, elem reflect.Value) (err error) {
			err = errors.New("any error")
			return
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, "any error", val.Error().Error())
	})
	t.Run("error in the iteration ignored", func(t *testing.T) {
		var hello = [5]int{1, 2, 3, 4, 5}
		val := Cast(reflect.ValueOf(hello), EnableIgnoreError())
		val.IterateArraySlice(nil, func(parent reflect.Value, index int, elem reflect.Value) (err error) {
			err = errors.New("any error")
			return
		})
		assert.Nil(t, val.Error())
	})
	t.Run("panic with recoverer", func(t *testing.T) {
		var hello = [5]int{1, 2, 3, 4, 5}
		val := Cast(reflect.ValueOf(hello), EnablePanicRecoverer())
		val.IterateArraySlice(nil, func(parent reflect.Value, index int, elem reflect.Value) (err error) {
			panic(errors.New("panic error"))
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, "panic error", val.Error().Error())
	})
}

func TestValue_IterateMap(t *testing.T) {
	t.Run("kind is not map", func(t *testing.T) {
		var hello int
		val := Cast(reflect.ValueOf(hello))
		assert.Nil(t, val.IterateMap().Error())
	})
	t.Run("iterate example function", func(t *testing.T) {
		var test = map[string]string{
			"hello": "hi",
			"hi":    "hello",
		}

		val := Cast(reflect.ValueOf(test))
		val.IterateMap(nil, func(parent reflect.Value, key reflect.Value, elem reflect.Value) (err error) {
			fmt.Println("key: ", key, "elem: ", elem.Interface())
			return
		})
		assert.Nil(t, val.Error())
	})
	t.Run("error in the iteration", func(t *testing.T) {
		var test = map[string]string{
			"hello": "hi",
			"hi":    "hello",
		}

		val := Cast(reflect.ValueOf(test))
		val.IterateMap(nil, func(parent reflect.Value, key reflect.Value, elem reflect.Value) (err error) {
			err = errors.New("random error")
			return
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, "random error", val.Error().Error())
	})
	t.Run("error in the iteration ignored", func(t *testing.T) {
		var test = map[string]string{
			"hello": "hi",
			"hi":    "hello",
		}

		val := Cast(reflect.ValueOf(test), EnableIgnoreError())
		val.IterateMap(nil, func(parent reflect.Value, key reflect.Value, elem reflect.Value) (err error) {
			err = errors.New("random error")
			return
		})
		assert.Nil(t, val.Error())
	})
	t.Run("panic with recoverer", func(t *testing.T) {
		var test = map[string]string{
			"hello": "hi",
			"hi":    "hello",
		}

		val := Cast(reflect.ValueOf(test), EnablePanicRecoverer())
		val.IterateMap(nil, func(parent reflect.Value, key reflect.Value, elem reflect.Value) (err error) {
			panic(errors.New("panic error"))
		})
		assert.NotNil(t, val.Error())
		assert.Equal(t, "panic error", val.Error().Error())
	})
}

func TestValue_IterateChan(t *testing.T) {
	t.Run("kind is not map", func(t *testing.T) {
		var hello int
		val := Cast(reflect.ValueOf(hello))
		assert.Nil(t, val.IterateChan().Error())
	})
	t.Run("iterate example function", func(t *testing.T) {
	})
	t.Run("error in the iteration", func(t *testing.T) {
	})
	t.Run("error in the iteration ignored", func(t *testing.T) {
	})
	t.Run("panic with recoverer", func(t *testing.T) {
	})
}
