package reflect

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFuncs(t *testing.T) {
	testCases := []struct {
		name string

		input any

		wantRes map[string]*FuncInfo
		wantErr error
	}{
		{
			// 普通结构体
			name:  "normal struct",
			input: User{},
			wantRes: map[string]*FuncInfo{
				"GetAge": {
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(User{})},
					Out:    []reflect.Type{reflect.TypeOf(0)},
					Result: []any{0},
				},
			},
		},
		{
			// 指针
			name:  "pointer",
			input: &User{},
			wantRes: map[string]*FuncInfo{
				"GetAge": {
					Name:   "GetAge",
					In:     []reflect.Type{reflect.TypeOf(&User{})},
					Out:    []reflect.Type{reflect.TypeOf(0)},
					Result: []any{0},
				},
				"ChangeName": {
					Name:   "ChangeName",
					In:     []reflect.Type{reflect.TypeOf(&User{}), reflect.TypeOf("")},
					Out:    []reflect.Type{},
					Result: []any{},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFuncs(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func Test(t *testing.T) {
	typ := reflect.TypeOf(Method)
	fmt.Println(typ.Kind())
	numIn := typ.NumIn()
	for i := 0; i < numIn; i++ {
		f := typ.In(i)
		fmt.Println(f.Kind())
		//fmt.Println(reflect.Zero(f).Interface())
		numField := f.NumField()
		for i := 0; i < numField; i++ {
			field := f.Field(i)
			fmt.Println(field.Type, reflect.Zero(field.Type).Interface())

		}
	}

}

func Method(user User) {

}
