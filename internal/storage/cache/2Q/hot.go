package twoq

import (
	"level-zero/internal/dto"
	"level-zero/internal/storage/cache"
	"log"
)

type HotQueue struct {
	size     int
	index    int
	head     *cache.LinkedList
	last     *cache.LinkedList
	elements map[string]*cache.LinkedList
}

func NewHotQueue(size int) *HotQueue {
	return &HotQueue{size: size, elements: make(map[string]*cache.LinkedList, size)}
}

func (h *HotQueue) Add(order dto.Order) {
	//if order is not in cache yet
	if node, ok := h.elements[order.UID]; !ok {
		if h.index == 0 {
			h.head = cache.NewLinkedList(order)
			h.last = h.head
			h.elements[order.UID] = h.head
			h.index++
		} else {
			//if queue is not full yet
			if h.size == h.index {
				//last's previous node will be the last node now, the old one will be deleted
				delete(h.elements, h.head.Order.UID)
				h.last = h.last.Prev
				if h.last != nil {
					h.last.Next = nil
				}
			} else {
				h.index++
			}
			h.head = h.head.NewHead(order)
			h.elements[order.UID] = h.head
		}
	} else {
		//send node to the top if it's somewhere in the queue
		h.head = node.PlaceBeforeHead(h.head)
	}
}

func (h *HotQueue) NodeWithOrderByID(id string) (node *cache.LinkedList) {
	const op = `storage.cache.2Q.LRU.OrderByID`

	log.Printf("%s content for %s was found in cache", op, id)

	node = h.elements[id]

	return node
}

func (h *HotQueue) PlaceBeforeHead(newHead *cache.LinkedList) {
	h.head = h.head.PlaceBeforeHead(newHead)
}
