package pod

import "testing"

func TestIsOverloaded(t *testing.T) {
	p := Pod{
		Name:       "nginx",
		CPURequest: 0.5,
		CPUUsage:   0.8,
	}

	if !p.IsOverloaded() {
		t.Errorf("expected pod to be overloaded")
	}
}

func TestPodNotOverloaded(t *testing.T) {
	p := Pod{
		Name:       "redis",
		CPURequest: 1.0,
		CPUUsage:   0.4,
	}

	if p.IsOverloaded() {
		t.Errorf("expected pod not to be overloaded")
	}
}
