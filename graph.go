package dj

import "context"

// BFS is a breadth-first search starting at the root node
// and traversing the graph using the next function.
func BFS[T any](root T, next func(T) []T, fn func(T)) {
	BFSCtx(context.Background(), root, next, fn)
}

// BFSCtx is a breadth-first search starting at the root node
// and traversing the graph using the next function.
// The context can be used to cancel the search.
func BFSCtx[T any](ctx context.Context, root T, next func(T) []T, fn func(T)) {
	queue := NewQueue(root)

	for {
		if v, ok := queue.PopFront(); !ok {
			return
		} else if fn(v); ctx.Err() != nil {
			return
		} else {
			for _, v := range next(v) {
				queue.PushBack(v)
			}
		}
	}
}

// DFS is a depth-first search starting at the root node
// and traversing the graph using the next function.
func DFS[T any](root T, next func(T) []T, fn func(T)) {
	DFSCtx(context.Background(), root, next, fn)
}

// DFSCtx is a depth-first search starting at the root node
// and traversing the graph using the next function.
// The context can be used to cancel the search.
func DFSCtx[T any](ctx context.Context, root T, next func(T) []T, fn func(T)) {
	queue := NewQueue(root)

	for {
		if v, ok := queue.PopBack(); !ok {
			return
		} else if fn(v); ctx.Err() != nil {
			return
		} else {
			for _, v := range next(v) {
				queue.PushBack(v)
			}
		}
	}
}
