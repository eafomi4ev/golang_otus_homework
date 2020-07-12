package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                      // длина списка
	Front() *Item                  // первый Item
	Back() *Item                   // последний Item
	PushFront(v interface{}) *Item // добавить значение в начало
	PushBack(v interface{}) *Item  // добавить значение в конец
	Remove(i *Item)                // удалить элемент
	MoveToFront(i *Item)           // переместить элемент в начало

}

type Item struct {
	Value interface{} // значение
	Next  *Item       // следующий элемент
	Prev  *Item       // предыдущий элемент

}

type list struct {
	length int
	front  *Item
	back   *Item
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *Item {
	return l.front
}

func (l *list) Back() *Item {
	return l.back
}

func (l *list) PushFront(v interface{}) *Item {
	newFrontItem := Item{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.length == 0 {
		l.back = &newFrontItem
	} else {
		l.front.Prev = &newFrontItem
	}

	l.front = &newFrontItem
	l.length++

	return &newFrontItem
}

func (l *list) PushBack(v interface{}) *Item {
	newBackItem := Item{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.length == 0 {
		l.front = &newBackItem
	} else {
		l.back.Next = &newBackItem
	}

	l.back = &newBackItem
	l.length++

	return &newBackItem
}

func (l *list) Remove(i *Item) {
	defer func() {
		l.length--
	}()

	if l.length == 1 {
		l.front = nil
		l.back = nil
		return
	}

	if prevItem := i.Prev; prevItem != nil {
		prevItem.Next = i.Next
	} else {
		i.Next.Prev = nil
		l.front = i.Next
	}

	if nextItem := i.Next; nextItem != nil {
		nextItem.Prev = i.Prev
	} else {
		i.Prev.Next = nil
	}
}

func (l *list) MoveToFront(i *Item) {
	if l.length == 1 || i == l.front {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next

	if nextItem != nil {
		nextItem.Prev = prevItem
	} else {
		l.back = prevItem
	}
	prevItem.Next = nextItem

	l.front.Prev = i
	i.Next = l.front
	i.Prev = nil
	l.front = i
}

func NewList() List {
	return &list{}
}
