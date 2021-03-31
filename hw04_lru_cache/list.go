package hw04lrucache

import "fmt"

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
	first, last *ListItem
	length      int
}

// Len возвращает количество ListItem в list (длинну списка).
func (l *list) Len() int { return l.length }

// Front возвращает первый элемент списка.
func (l *list) Front() *ListItem {
	if l.Len() == 0 {
		return nil
	}
	return l.first
}

// Back возвращает последний элемент списка.
func (l *list) Back() *ListItem {
	if l.Len() == 0 {
		return nil
	}
	return l.last
}

// PushFront добавляет новый элемент в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	if listItem, ok := l.listIsEmpty(v); ok { // если список пуст и добавляется первый элемент
		return listItem
	}
	itemToBind := &ListItem{Value: v, Next: nil, Prev: nil}
	l.length++
	l.bindToFront(itemToBind)
	return l.Front()
}

// PushBack добавляет новый элемент в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	if listItem, ok := l.listIsEmpty(v); ok { // если список пуст и добавляется первый элемент
		return listItem
	}
	itemToBind := &ListItem{Value: v, Next: nil, Prev: nil}
	l.length++
	l.bindToBack(itemToBind)
	return l.Back()
}

// Remove удаляет элемент из списка.
func (l *list) Remove(i *ListItem) {
	l.unbind(i)
	*i = ListItem{}
	l.length--
}

// MoveToFront перемещает элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	l.unbind(i)
	l.bindToFront(i)
}

// listIsEmpty проверяет пуст ли список, если да, то это готовит первый элемент для списка.
func (l *list) listIsEmpty(v interface{}) (*ListItem, bool) {
	if l.Len() == 0 {
		l.first = &ListItem{v, nil, nil}
		l.last = l.first
		l.length = 1
		return l.first, true
	}
	return nil, false
}

// bindToFront привязывает элемент в начале списка, симметричен bindToBack.
func (l *list) bindToFront(i *ListItem) {
	saveFirst := l.Front()
	i.Prev = nil
	i.Next = saveFirst
	l.first = i
	if l.Len() > 1 {
		saveFirst.Prev = l.Front() // когда в списке много элементов
	} else {
		l.last = i // случай, когда в списке 1 элемент и он перемещается методом MoveToFront
	}
}

// bindToBack привязывает элемент в начале списка, симметричен bindToFront.
func (l *list) bindToBack(i *ListItem) {
	saveLast := l.Back()
	i.Next = nil
	i.Prev = saveLast
	l.last = i
	if l.Len() > 1 {
		saveLast.Next = l.Back() // когда в списке много элементов
	} else {
		l.first = i // данный случай невозможен, так как нет метода "MoveToBack" (аналогичный MoveToFront)
	}
}

// unbind отвязывает элемент для Remove или MoveToFront.
func (l *list) unbind(i *ListItem) {
	switch {
	case l.Len() == 1:
		l.last = nil
		l.first = nil
	case i == l.Front():
		{
			i.Next.Prev = nil
			l.first = i.Next
		}
	case i == l.Back():
		{
			i.Prev.Next = nil
			l.last = i.Prev
		}
	default: // не единственный, не первый и не последний элемент - это средний
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
}

// String выводит значения списка.
func (l *list) String() string {
	for i := l.Front(); i != nil; i = i.Next {
		fmt.Print(i.Value, " ")
	}
	return ""
}

// NewList создает новый *list в интерфейсе List.
func NewList() List {
	return new(list)
}
