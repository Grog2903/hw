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
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.length
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

	if l.front != nil {
		l.front.Prev = newItem
	} else {
		l.back = newItem
	}

	l.front = newItem
	l.length++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Prev:  l.back,
	}

	if l.back != nil {
		l.back.Next = newItem
	} else {
		l.front = newItem
	}

	l.back = newItem
	l.length++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if nil == i.Prev {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if nil == i.Next {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}

	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}
