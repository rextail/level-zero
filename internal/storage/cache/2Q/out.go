package twoq

import (
	"level-zero/internal/dto"
	"level-zero/internal/storage/cache"
	"log"
)

type AddOutResult struct {
	IsHot     bool
	Evacuated *cache.LinkedList
}

type OutQueue struct {
	size     int
	index    int
	head     *cache.LinkedList
	last     *cache.LinkedList
	elements map[string]*cache.LinkedList
}

func NewOutQueue(size int) *OutQueue {
	return &OutQueue{
		size:     size,
		index:    0,
		elements: make(map[string]*cache.LinkedList, size),
	}
}

func (o *OutQueue) Add(order dto.Order) AddOutResult {

	res := AddOutResult{}

	if o.index == 0 {
		//first element in the queue
		o.head = cache.NewLinkedList(order)
		o.last = o.head
		o.elements[order.UID] = o.head
		o.index++
		return res
	}

	if elem, ok := o.elements[order.UID]; !ok {
		//if element is not in queue
		if o.index == o.size {

			toEvacuate := o.last

			o.last = o.last.Prev
			if o.last != nil {
				o.last.Next = nil
			}

			toEvacuate.Prev = nil
			toEvacuate.Next = nil

			res.Evacuated = toEvacuate
			res.IsHot = false

			delete(o.elements, toEvacuate.Order.UID)

		} else {
			o.index++
		}
		o.head = o.head.NewHead(order)
		o.elements[order.UID] = o.head

	} else {
		//if element already in queue, delete from queue (will be added to hot queue after)
		o.RemoveByID(elem.Order.UID)

		res.Evacuated = elem
		res.IsHot = true
	}

	return res
}

func (o *OutQueue) RemoveByID(id string) *cache.LinkedList {
	const op = `storage.2Q.out.OrderByID`
	log.Printf("%s: deleting %s from Out Queue", op, id)

	node := o.elements[id]

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if node == o.last {
		o.last = node.Prev
	}

	delete(o.elements, node.Order.UID)

	o.index--

	return node
}

func (o *OutQueue) OrderByID(id string) dto.Order {
	const op = `storage.2Q.out.OrderByID`

	node := o.elements[id]

	log.Printf("%s content for %s was found in cache", op, id)

	return node.Order
}
