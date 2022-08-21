package list

import "fmt"

// ArrayList 基于切片的简单封装
type ArrayList[E any] struct {
	elements []E
}

// New 初始化一个len为0，cap为cap的ArrayList
func New[E any](cap int) *ArrayList[E] {
	return &ArrayList[E]{elements: make([]E, 0, cap)}
}
func NewWithE[E any](elements []E) *ArrayList[E] {
	return &ArrayList[E]{elements: elements}
}

// Append 往ArrayList里追加数据
func (a *ArrayList[E]) Append(t E) error {
	a.elements = append(a.elements, t)
	return nil
}

// Add 在ArrayList下标为index的位置插入一个元素
// 当index等于ArrayList长度等同于append
func (a *ArrayList[E]) Add(index int, t E) error {
	if index < 0 || index > len(a.elements) {
		return newErrIndexOutOfRange(len(a.elements), index)
	}
	a.elements = append(a.elements, t)
	copy(a.elements[index+1:], a.elements[index:])
	a.elements[index] = t
	return nil
}

// newErrIndexOutOfRange 创建一个代表
func newErrIndexOutOfRange(length int, index int) error {
	return fmt.Errorf("kit: 下标超出范围，长度 %d, 下标 %d", length, index)
}
