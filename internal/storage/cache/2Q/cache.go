package twoq

import (
	"level-zero/internal/dto"
)

const (
	inQ = iota
	outQ
	hotQ
)

type Cache struct {
	//2Q caching
	inQueue  *InQueue
	outQueue *OutQueue
	hotQueue *HotQueue
	qLocator map[string]int
}

func New(size int) *Cache {
	outSize := size / 2
	inSize := size / 4
	hotSize := size / 4
	return &Cache{
		inQueue:  NewInQueue(inSize),
		outQueue: NewOutQueue(outSize),
		hotQueue: NewHotQueue(hotSize),
		qLocator: make(map[string]int, size),
	}
}

func (c *Cache) AddOrder(order dto.Order) {
	const op = `storage.cache.AddOrder`

	qID, ok := c.qLocator[order.UID]
	if !ok {
		c.qLocator[order.UID] = inQ
	}

	switch qID {
	case inQ:
		//if order is in inQ
		if evacuated := c.inQueue.Add(order); evacuated != nil {
			//if queue is full we move element to the next queue
			c.qLocator[evacuated.Order.UID] = outQ
			res := c.outQueue.Add(evacuated.Order)
			if res.Evacuated != nil {
				//if after we've
				//added the element to the outQ and some non-nil element was
				//evacuated from outQ
				if res.IsHot {
					c.hotQueue.Add(order)
					c.qLocator[order.UID] = hotQ
				} else {
					delete(c.qLocator, res.Evacuated.Order.UID)
				}
			}
		}
	case outQ:
		//if
		res := c.outQueue.Add(order)
		if res.Evacuated != nil {
			if res.IsHot {
				c.hotQueue.Add(order)
				c.qLocator[order.UID] = hotQ
			} else {
				delete(c.qLocator, res.Evacuated.Order.UID)
			}
		}
	case hotQ:
		c.hotQueue.Add(order)
	}
}

func (c *Cache) OrderByID(id string) string {
	const op = `storage.cache.OrderByID`
	if qID, ok := c.qLocator[id]; ok {

		switch qID {

		case inQ:
			//nothing happens if we read from this queue
			order := c.inQueue.OrderByID(id)
			return order.Content
		case outQ:
			//if element in this queue we move it to the hot queue
			order := c.outQueue.OrderByID(id)

			//exclude node from out queue
			node := c.outQueue.RemoveByID(id)

			c.hotQueue.Add(node.Order)

			c.qLocator[id] = hotQ

			return order.Content
		case hotQ:
			//if element in this queue we move it to the head
			order := c.hotQueue.NodeWithOrderByID(id)
			c.hotQueue.PlaceBeforeHead(order)
			return order.Order.Content
		}
	}
	//if element wasn't found in the qLocator
	return ""
}
