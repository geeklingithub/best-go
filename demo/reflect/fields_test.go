package reflect

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {

	u1 := &User{
		Name: "大明",
	}
	u2 := &u1

	tests := []struct {
		// 名字
		name string

		// 输入部分
		val any

		// 输出部分
		wantRes map[string]any
		wantErr error
	}{
		{
			name:    "nil",
			val:     nil,
			wantErr: errors.New("不能为 nil"),
		},
		{
			name:    "user",
			val:     User{Name: "Tom"},
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Tom",
			},
		},
		{
			// 指针
			name: "pointer",
			val:  &User{Name: "Jerry"},
			// 要支持指针
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "Jerry",
			},
		},
		{
			// 多重指针
			name: "multiple pointer",
			val:  u2,
			// 要支持指针
			wantErr: nil,
			wantRes: map[string]any{
				"Name": "大明",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := iterateFields(tt.val)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestSetField(t *testing.T) {
	testCases := []struct {
		name string

		field  string
		entity any
		newVal any

		wantErr error
	}{
		{
			name:    "struct",
			entity:  User{},
			field:   "Name",
			wantErr: errors.New("非法类型"),
		},
		{
			name:    "private field",
			entity:  &User{},
			field:   "age",
			wantErr: errors.New("不可修改字段"),
		},
		{
			name:    "invalid field",
			entity:  &User{},
			field:   "invalid_field",
			wantErr: errors.New("字段不存在"),
		},
		{
			name: "pass",
			entity: &User{
				Name: "",
			},
			field:  "Name",
			newVal: "Tom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := SetField(tc.entity, tc.field, tc.newVal)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

type User struct {
	Name string
	// 因为同属一个包，所以 age 还可以被测试访问到
	// 如果是不同包，就访问不到了
	age int
}

func (u User) GetAge() int {
	return u.age
}

func (u *User) ChangeName(newName string) {
	u.Name = newName
}

func (u User) private() {
	fmt.Println("private")
}
