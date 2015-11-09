package main

import "container/list"

// VisitCache stores a cache of recent visits.
type VisitCache struct {
	cache    *list.List // The cache of recent visits.
	capacity int        // The capacity of the cache.
}

// NewVisitCache returns an initialised visit cache of the given capacity.
func NewVisitCache(capacity int) *VisitCache {
	return &VisitCache{list.New(), capacity}
}

// Add a visit to the cache, removing older entries if the cache is full.
func (vc *VisitCache) Add(v Visit) {
	vc.cache.PushBack(v)
	if vc.cache.Len() > vc.capacity {
		vc.cache.Remove(vc.cache.Front())
	}
}

// Iterator returns a channel that receives the contents of the cache.
func (vc *VisitCache) Iterator() <-chan Visit {
	ch := make(chan Visit)
	go func() {
		for e := vc.cache.Front(); e != nil; e = e.Next() {
			ch <- e.Value.(Visit)
		}
		close(ch)
	}()
	return ch
}
