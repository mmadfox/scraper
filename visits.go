package scraper

import (
	cmap "github.com/streamrail/concurrent-map"
)

type Visiter interface {
	Visit(string) bool
	ResetVisit(string) error
	Drop() error
}

type memoryVisits struct {
	m cmap.ConcurrentMap
}

func (v *memoryVisits) Visit(u string) bool {
	return !v.m.SetIfAbsent(u, 1)
}

func (v *memoryVisits) ResetVisit(u string) error {
	v.m.Remove(u)
	return nil
}

func (v *memoryVisits) Drop() error {
	return nil
}

func NewMemoryVisits() Visiter {
	return &memoryVisits{
		m: cmap.New(),
	}
}
