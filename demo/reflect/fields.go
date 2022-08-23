package reflect

import (
	"errors"
	"fmt"
	"reflect"
	//"reflect"
)

func IterateFields(val any) {
	// 复杂逻辑
	res, err := iterateFields(val)

	// 简单逻辑
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range res {
		fmt.Println(k, v)
	}
}

func iterateFields(val any) (map[string]any, error) {
	if val == nil {
		return nil, errors.New("不能为 nil")
	}

	typ := reflect.TypeOf(val)
	refVal := reflect.ValueOf(val)

	// 怎么拿到指针指向的对象
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		refVal = refVal.Elem()
	}

	numField := typ.NumField()
	res := make(map[string]any, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		res[fdType.Name] = refVal.Field(i).Interface()
	}
	return res, nil
}

func SetField(entity any, field string, newVal any) error {
	val := reflect.ValueOf(entity)
	typ := val.Type()

	// 只能是一级指针，类似 *User
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return errors.New("非法类型")
	}

	typ = typ.Elem()
	val = val.Elem()

	// 这个地方判断不出来 field 在不在
	fd := val.FieldByName(field)
	// 利用 type 来判断 field 在不在
	if _, found := typ.FieldByName(field); !found {
		return errors.New("字段不存在")
	}

	if !fd.CanSet() {
		return errors.New("不可修改字段")
	}
	fd.Set(reflect.ValueOf(newVal))
	return nil
}
