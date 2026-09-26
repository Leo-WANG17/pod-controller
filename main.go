package main

import (
	"context"
	"fmt"
	"github.com/Leo-WANG17/pod-controller/pod"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// fmt.Println("hello go")

	// result := add(1, 3)
	// fmt.Println(result)

	nginx := pod.Pod{
		Name:       "nginx",
		CPURequest: 0.5,
		CPUUsage:   0.8,
	}
	fmt.Println(nginx.IsOverloaded())

	pods := []Pod{
		{
			Name:       "nginx",
			CPURequest: 0.5,
			CPUUsage:   0.8,
		},
		{
			Name:       "redis",
			CPURequest: 1.0,
			CPUUsage:   0.4,
		},
		{
			Name:       "mysql",
			CPURequest: 2.0,
			CPUUsage:   2.5,
		},
	}

	// ch := make(chan string, 2)

	// ch <- "one"
	// ch <- "twe"
	// fmt.Println("sent 2 messages")

	// fmt.Println(<-ch)
	// ch <- "three"
	// fmt.Println(<-ch)
	// fmt.Println(<-ch)

	// ch := make(chan string)

	// go pod.ProducePods(ch)

	// for value := range ch {
	// 	fmt.Println(value)
	// }

	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	ch := make(chan string)

	go pod.CheckPod(ctx, "nginx", ch)

	msg := <-ch
	fmt.Println(msg)

	// for i := range pods {
	// 	// if pods[i].isOverloaded() {
	// 	// 	pods[i].CPURequest = pods[i].CPUUsage
	// 	// 	fmt.Println(pods[i].CPURequest)
	// 	// } else {
	// 	// 	fmt.Println(pods[i].CPURequest)
	// 	// }
	// 	p := &pods[i]

	// 	if p.isOverloaded() {
	// 		p.setCPURequest(p.CPUUsage)
	// 		fmt.Println(p.CPURequest)
	// 	}
	// }

	checkWorkload(pods[2])

	pod := Pod{
		Name:       "nginx",
		CPURequest: 0.5,
		CPUUsage:   0.8,
	}

	checkWorkload(&pod)
	checkWorkload(pod)

	// pods = append(pods, Pod{
	// 	Name:       "api-server",
	// 	CPURequest: 1.5,
	// 	CPUUsage:   0.9,
	// })

	// for _, pod := range pods {
	// 	if pod.isOverloaded() {
	// 		fmt.Println(pod.Name)
	// 	}
	// }

	// podMap := make(map[string]float64)

	// podMap["nginx"] = 0.8
	// podMap["redis"] = 0.4
	// podMap["mysql"] = 2.5

	// fmt.Println(podMap)

	// podMap["nginx"] = 1.2
	// podMap["api-server"] = 0.6
	// fmt.Println(podMap["nginx"])
	// fmt.Println(podMap["api-server"])

	// if _, ok := podMap["redis"]; ok {
	// 	fmt.Println("redis exists")
	// }

	// if _, ok := podMap["mongodb"]; !ok {
	// 	fmt.Println("mongodb not found")
	// }

	// for _, value := range podMap {
	// 	fmt.Println(value)
	// }

	podsMap := make(map[string]Pod)

	podsMap["nginx"] = Pod{
		Name:       "nginx",
		CPURequest: 0.5,
		CPUUsage:   0.8,
	}

	podsMap["redis"] = Pod{
		Name:       "redis",
		CPURequest: 1.0,
		CPUUsage:   0.4,
	}

	podsMap["mysql"] = Pod{
		Name:       "mysql",
		CPURequest: 2.0,
		CPUUsage:   2.5,
	}

	// pod, err := findPod("api-server", podsMap)

	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return
	// }

	// fmt.Println(pod.Name)

	printPod("mongodb", podsMap)

	// for _, pod := range podsMap {
	// 	if pod.isOverloaded() {
	// 		fmt.Println(pod.Name)
	// 	}
	// }

	// fmt.Println(podsMap)

	// println(podMap["nginx"])
	// println(podMap["redis"])
	// println(podMap["mysql"])

	// fmt.Println(pod_1)

	// updateCPU(&pod_1)
	// fmt.Println(pod_1)

	// pod_1.PrintUsage()

	// pod_2 := Pod{
	// 	Name:       "nginx",
	// 	CPURequest: 0.5,
	// 	CPUUsage:   0.8,
	// }

	// fmt.Println(isOverloaded(&pod_2))

	// fmt.Println(pod_1.isOverloaded())
	// pod_1.setCPURequest(1.0)
	// fmt.Println(pod_1.CPURequest)
	// fmt.Println(pod_1.isOverloaded())

	// pod_test := corev1.Pod{
	// 	Spec: corev1.PodSpec{
	// 		Containers: []corev1.Container{
	// 			{
	// 				Name: "nginx",
	// 			},
	// 		},
	// 	},
	// }

	// pod_test.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
	// 	corev1.ResourceCPU: resource.MustParse("500m"),
	// }

	// cpu := pod_test.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU]
	// fmt.Println(cpu.String())

	// pod_test.Name = "nginx"
	// pod_test.Namespace = "default"

	// fmt.Println(pod_test.Name)
	// fmt.Println(pod_test.Namespace)

	//pod_redis
	pod_redis := corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "redis",
					Image: "redis:latest",

					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},

						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},

			NodeSelector: map[string]string{
				"gpu": "nvidia",
			},

			Affinity: &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
						NodeSelectorTerms: []corev1.NodeSelectorTerm{
							{
								MatchExpressions: []corev1.NodeSelectorRequirement{
									{
										Key:      "gpu",
										Operator: corev1.NodeSelectorOpIn,
										Values:   []string{"nvidia", "amd"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	cpuRequest := pod_redis.Spec.Containers[0].
		Resources.Requests[corev1.ResourceCPU]

	cpuLimit := pod_redis.Spec.Containers[0].
		Resources.Limits[corev1.ResourceCPU]

	memoryRequest := pod_redis.Spec.Containers[0].
		Resources.Requests[corev1.ResourceMemory]

	memoryLimit := pod_redis.Spec.Containers[0].
		Resources.Limits[corev1.ResourceMemory]

	fmt.Println(cpuRequest.String())
	fmt.Println(cpuLimit.String())
	fmt.Println(memoryRequest.String())
	fmt.Println(memoryLimit.String())
}

var cpu = 4
var usage = 0.75
var name = "nginx"
var ready = true

func add(a int, b int) int {
	return a + b
}

type Pod struct {
	Name       string
	CPURequest float64
	CPUUsage   float64
}

func updateCPU(p *Pod) {
	p.CPURequest = 1.0
}

func (p *Pod) PrintUsage() {
	fmt.Println(p.Name, p.CPUUsage)
}

func isOverloaded(p *Pod) bool {
	if p.CPUUsage > p.CPURequest {
		return true
	}

	return false
}

func (p Pod) isOverloaded() bool {
	return p.CPUUsage > p.CPURequest
}

func (p *Pod) setCPURequest(value float64) {
	p.CPURequest = value
}

type Workload interface {
	isOverloaded() bool
}

func checkWorkload(w Workload) {
	fmt.Println(w.isOverloaded())
}

// error handling
func findPod(name string, pods map[string]Pod) (Pod, error) {
	pod, ok := pods[name]

	if !ok {
		return pod, fmt.Errorf("pod %q dose not exist", name)
	}

	return pod, nil
}

// error handling - 2
func printPod(name string, pods map[string]Pod) error {
	result, err := findPod(name, pods)

	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	fmt.Println(result.Name)
	return nil
}

func work(ch chan string) {
	fmt.Println("working")

	ch <- "done"
}
