package queue

func (q *Queue[T]) Add(item T) {
	key := item.Key()

	if !q.visited[key] {
		q.items = append(q.items, item)
		q.visited[key] = true
	}
}

func (q *Queue[T]) Next() (T, bool) {
	if len(q.items) == 0 {
		var zeroVal T
		return zeroVal, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) IsVisited(item T) bool {
	key := item.Key()
	_, visited := q.visited[key]
	return visited
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}
