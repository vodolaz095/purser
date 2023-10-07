package config

import "testing"

func TestConfig(t *testing.T) {
	t.Logf("Address: %s", Address)
	t.Logf("Port: %s", Port)
}
