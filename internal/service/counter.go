package service

import (
	"context"
	"sync"
)

// CounterService реализует взаимодействие с сервисом счётчиков. Возможно надо было для него сделать репозиторий,
// но не хочу слишком усложнять код
type CounterService struct {
	sync.RWMutex
	counters map[string]uint64
}

// Init готовит внутренную мапу сервиса к работе
func (cs *CounterService) Init() {
	cs.counters = make(map[string]uint64, 0)
}

// Get получает значение счётчика, если он не найден, то возвращается 0, false
func (cs *CounterService) Get(ctx context.Context, key string) (value uint64, found bool) {
	cs.RLock()
	defer cs.RUnlock()
	value, found = cs.counters[key]
	if !found {
		return 0, false
	}
	return value, true
}

// Increment увеличивает значение счётчика
func (cs *CounterService) Increment(ctx context.Context, key string, delta uint64) uint64 {
	if delta == 0 {
		return 0
	}
	cs.RLock()
	value, found := cs.counters[key]
	cs.RUnlock()
	if !found {
		cs.Lock()
		cs.counters[key] = delta
		cs.Unlock()
		return delta
	}
	cs.Lock()
	cs.counters[key] = value + delta
	cs.Unlock()
	return value + delta
}
