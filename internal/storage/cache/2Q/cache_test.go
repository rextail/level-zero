package twoq

import (
	"fmt"
	"level-zero/internal/dto"
	"strconv"
	"testing"
)

func TestCache_AddOrder(t *testing.T) {
	t.Run("From Queue to Queue", func(t *testing.T) {
		//InQueue And HotQueue of size 3, OutQueue of size 6

		//First element must reach OutQueue
		cache := New(12)
		counter := 0
		for counter <= 3 {
			order := dto.Order{
				UID:     strconv.Itoa(counter),
				Content: "",
			}
			cache.AddOrder(order)
			counter++
		}
		fmt.Println(cache.qLocator)
		if cache.qLocator["0"] != 1 {
			t.Errorf("got %d, expected 1", cache.qLocator["3"])
		}
		//Now we have to reach hot queue
		order := dto.Order{
			UID:     "0",
			Content: "",
		}
		cache.AddOrder(order)
		fmt.Println(cache.qLocator)
		if cache.qLocator["0"] != 2 {
			t.Errorf("got %d, expected 2", cache.qLocator["0"])
		}
	})
	t.Run("Pull odd element out from OutQueue", func(t *testing.T) {
		//If the OutQueue is full and one more unique element
		//come to the queue then last element will be
		//deleted, unique element become head

		//InQueue And HotQueue of size 2, OutQueue of size 4
		cache := New(8)
		counter := 0

		for counter <= 6 {
			order := dto.Order{
				UID:     strconv.Itoa(counter),
				Content: "",
			}
			cache.AddOrder(order)
			counter++
		}
		if id, ok := cache.qLocator["0"]; ok {
			t.Errorf("got %d, unexpected", id)
		}
	})

}
