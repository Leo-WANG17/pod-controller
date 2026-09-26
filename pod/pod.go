package pod

import (
	"context"
	"fmt"
	"time"
)

type Pod struct {
	Name       string
	CPURequest float64
	CPUUsage   float64
}

type PodManager struct {
	Pods map[string]*Pod
}

func NewPodManager() *PodManager {
	return &PodManager{
		Pods: make(map[string]*Pod),
	}
}

func (p Pod) IsOverloaded() bool {
	return p.CPUUsage > p.CPURequest
}

func (pm *PodManager) AddPod(p *Pod) {
	pm.Pods[p.Name] = p
}

func (pm *PodManager) FindPod(name string) (*Pod, error) {
	pod, ok := pm.Pods[name]

	if !ok {
		return nil, fmt.Errorf("pod %q does not exist", name)
	}

	return pod, nil
}

func (pm *PodManager) PrintOverloadedPods() {
	for _, pod := range pm.Pods {
		if pod.IsOverloaded() {
			fmt.Println(pod.Name)
		}
	}
}

func (pm *PodManager) UpdateCPURequest(name string, value float64) error {
	pod, err := pm.FindPod(name)

	if err != nil {
		return err
	}

	pod.CPURequest = value
	return nil
}

// context
func WaitForPod(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("pod is ready")

	case <-ctx.Done():
		fmt.Println("operation cancelled")
	}
}

func ProducePods(ch chan string) {
	ch <- "nginx"
	ch <- "redis"
	ch <- "mysql"

	close(ch)
}

func CheckPod(ctx context.Context, podName string, ch chan string) {
	select {
	case <-time.After(3 * time.Second):
		ch <- podName + "ready"

	case <-ctx.Done():
		ch <- podName + " " + "cancelled"
	}
}
