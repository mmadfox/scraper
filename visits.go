package scraper

import (
	"sync"

	cmap "github.com/streamrail/concurrent-map"
)

type Visiter interface {
	Visit(string) bool
	ResetVisit(string) error
	Drop() error
	Close() error
}

type memoryVisits struct {
	m     cmap.ConcurrentMap
	mutex *sync.Mutex
}

func (v *memoryVisits) Visit(u string) bool {
	return !v.m.SetIfAbsent(u, 1)
}

func (v *memoryVisits) ResetVisit(u string) error {
	v.m.Remove(u)
	return nil
}

func (v *memoryVisits) Drop() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.m = cmap.New()
	return nil
}

func (v *memoryVisits) Close() error {
	return nil
}

func NewMemoryVisits() Visiter {
	return &memoryVisits{
		m:     cmap.New(),
		mutex: &sync.Mutex{},
	}
}
