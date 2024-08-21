package twoq

import (
	"level-zero/internal/dto"
	"level-zero/internal/storage/cache"
	"log"
)

type InQueue struct {
	size     int
	index    int
	head     *cache.LinkedList
	last     *cache.LinkedList
	elements map[string]*cache.LinkedList
}

func NewInQueue(size int) *InQueue {
	return &InQueue{
		size:     size,
		index:    0,
		elements: make(map[string]*cache.LinkedList),
	}
}

func (i *InQueue) Add(order dto.Order) *cache.LinkedList {
	//Returns node that was deleted

	var deleted *cache.LinkedList

	if _, ok := i.elements[order.UID]; ok {
		//this queue should not contain duplicates
		return nil
	}

	if i.index == 0 {
		i.head = cache.NewLinkedList(order)
		i.last = i.head
		i.index++
		i.elements[order.UID] = i.head
		return nil
	}

	if i.index < i.size {
		i.index++
		i.head = i.head.NewHead(order)
		i.elements[order.UID] = i.head
		return nil
	}
	if i.index == i.size {
		deleted = i.last
		if i.last != nil {
			i.last = i.last.Prev
			if i.last != nil {
				i.last.Next = nil
			}
			delete(i.elements, deleted.Order.UID)
		}
		i.head = i.head.NewHead(order)
		i.elements[order.UID] = i.head
	}
	return deleted
}

func (i *InQueue) OrderByID(id string) (order dto.Order) {
	const op = `storage.2Q.in.OrderByID`

	order.UID = id
	order.Content = i.elements[id].Order.Content

	log.Printf("%s content for %s was found in cache", op, id)

	return order
}
