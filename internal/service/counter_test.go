package service

import (
	"context"
	"testing"
)

var testCounterService CounterService

func TestCounterService(t *testing.T) {
	testCounterService.Init()
	notFoundValue, notFound := testCounterService.Get(context.Background(), "something")
	if notFound {
		t.Errorf("счётчик найден?")
	}
	if notFoundValue != 0 {
		t.Errorf("значение ненайденного счётчика - не нуль")
	}
	testCounterService.Increment(context.Background(), "something", 3)

	value, found := testCounterService.Get(context.Background(), "something")
	if !found {
		t.Errorf("счётчик не найден")
	}
	if value != 3 {
		t.Errorf("значение найденного счётчика - не 3")
	}
}
