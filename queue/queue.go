package queue

func (q *Queue[T]) Add(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := item.Key()

	if _, visited := q.visited[key]; !visited {
		q.visited[key] = struct{}{}
		q.items = append(q.items, item)
	}
}

func (q *Queue[T]) Next() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zeroVal T
		return zeroVal, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue[T]) IsVisited(item T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	key := item.Key()
	_, visited := q.visited[key]
	return visited
}

func (q *Queue[T]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}

func (q *Queue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
