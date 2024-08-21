package cache

import (
	"level-zero/internal/dto"
)

type LinkedList struct {
	Order dto.Order
	Next  *LinkedList
	Prev  *LinkedList
}

func NewLinkedList(order dto.Order) *LinkedList {
	return &LinkedList{Order: order}
}

func (l *LinkedList) NewHead(order dto.Order) *LinkedList {
	node := &LinkedList{
		Order: order,
		Next:  l,
	}
	if l != nil {
		l.Prev = node
	}
	return node
}

func (l *LinkedList) PlaceBeforeHead(newHead *LinkedList) *LinkedList {
	if l == newHead {
		return l
	}

	node := l

	prev := node.Prev
	next := node.Next

	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}

	node.Next = newHead

	node.Prev = nil

	if newHead != nil {
		newHead.Prev = node
	}

	return node
}
