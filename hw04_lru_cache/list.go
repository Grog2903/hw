package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.front,
	}

	if l.Back() == nil {
		l.back = newItem
	}

	return &ListItem{v, nil, l.front}
}

func (l *list) PushBack(v interface{}) *ListItem {
	return &ListItem{v, nil, l.front}
}

func (l *list) Remove(i *ListItem) {
}

func (l *list) MoveToFront(i *ListItem) {
}

func NewList() List {
	return new(list)
}
