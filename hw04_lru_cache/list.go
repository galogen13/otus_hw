package hw04_lru_cache //nolint:golint,stylecheck

// List - интерфейс - двухсвязный список.
type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый Item
	Back() *ListItem                   // последний Item
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

// ListItem - элемент списка.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	first *ListItem
	last  *ListItem
}

func NewList() List {
	return &list{}
}

// Len - возвращает длину списка.
func (l list) Len() int {
	return l.len
}

// Front - возвращает указатель на первый элемент списка.
func (l list) Front() *ListItem {
	return l.first
}

// Back - возвращает указатель на последний элемент списка.
func (l list) Back() *ListItem {
	return l.last
}

// PushFront - добавляет в начало списка элемент.
func (l *list) PushFront(v interface{}) *ListItem {
	newFirst := ListItem{Value: v}
	if l.first != nil {
		newFirst.Next = l.first
		l.first.Prev = &newFirst
	}
	if l.Len() == 0 {
		l.last = &newFirst
	}
	l.first = &newFirst
	l.len++

	return &newFirst
}

// PushBack - добавляет в конец списка элемент.
func (l *list) PushBack(v interface{}) *ListItem {
	newLast := ListItem{Value: v}
	if l.last != nil {
		newLast.Prev = l.last
		l.last.Next = &newLast
	}

	if l.Len() == 0 {
		l.first = &newLast
	}
	l.last = &newLast
	l.len++

	return &newLast
}

// Remove - удаляет элемент из списка.
func (l *list) Remove(i *ListItem) {
	switch {
	// Если известен предудыщий и следующий - сошьем их
	case i.Next != nil && i.Prev != nil:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	// Если удаляется крайний элемент, то предкрайний должен стать крайним
	case i.Next == nil && i.Prev != nil:
		i.Prev.Next = nil
		l.last = i.Prev
	// Если удаляется первый элемент, то второй должен стать первым
	case i.Next != nil && i.Prev == nil:
		i.Next.Prev = nil
		l.first = i.Next
	// Если это последний элемент в списке, очищаем первого и последнего в списке
	case i.Next == nil && i.Prev == nil:
		l.first = nil
		l.last = nil
	}

	l.len--
}

// MoveToFront - перемещает элемент в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}
