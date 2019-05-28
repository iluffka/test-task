package counter

import (
	"testing"
	"time"
)

const (
	countNodes            = 20
	positiveFromWithSleep = 10
	positiveFrom          = 0
)

func TestCounter(t *testing.T) {
	t.Run("positive case with sleep", func(t *testing.T) {
		nodes := make([]Node, 0, countNodes)
		for i := 0; i < countNodes; i++ {
			if i == countNodes/2 {
				time.Sleep(6 * time.Second)
			}
			nodes = append(nodes, Node{
				Time: time.Now().Unix(),
			})
		}
		node := Node{
			Time: time.Now().Unix(),
		}
		nodes = append(nodes)
		cutOff := node.Time - 5

		from := Counter(nodes, cutOff)
		if from != positiveFromWithSleep {
			t.Fatalf("ожидаемый результат %d, полученный %d", positiveFromWithSleep, from)
		}
	})
	t.Run("positive case", func(t *testing.T) {
		nodes := make([]Node, 0, countNodes+1)
		for i := 0; i < countNodes; i++ {
			nodes = append(nodes, Node{
				Time: time.Now().Unix(),
			})
		}
		node := Node{
			Time: time.Now().Unix(),
		}
		nodes = append(nodes, node)
		cutOff := node.Time - 5

		from := Counter(nodes, cutOff)
		if from != positiveFrom {
			t.Fatalf("ожидаемый результат %d, полученный %d", positiveFrom, from)
		}
	})
}
