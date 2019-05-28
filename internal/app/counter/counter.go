package counter

type Node struct {
	Time int64
}

func Counter(nodes []Node, cutOff int64) int {
	var from int

	for i, n := range nodes {
		if from == 0 {
			if n.Time >= cutOff {
				from = i
				break
			}
		}
	}

	return from
}
