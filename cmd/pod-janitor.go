package main

import (
	"context"
	"flag"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var podNamespace = flag.String("pod-namespace", "", "The kubernetes namespace to run in")
var deleteSuccessfulAfter = flag.Duration("delete-successful-after", 0*time.Minute, "The kubernetes namespace to run in")
var deleteFailedAfter = flag.Duration("delete-failed-after", 0*time.Minute, "The kubernetes namespace to run in")

func main() {
	flag.Parse()
	
	if *podNamespace == "" {
		log.Fatalf("init failed: pod-namespace argument not set")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get cluster config: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialise client: %v", err)
	}

	succeededPods, err := client.CoreV1().Pods(*podNamespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Succeeded"})
	if err != nil {
		log.Fatalf("Failed to get list of successful pods: %v", err)
	}

	failedPods, err := client.CoreV1().Pods(*podNamespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Failed"})
	if err != nil {
		log.Fatalf("Failed to get list of failed pods: %v", err)
	}

	pods := succeededPods.Items[:0]
	pods = append(pods, succeededPods.Items[:]...)
	pods = append(pods, failedPods.Items[:]...)

	for _, pod := range pods {
		
		if shouldDeletePod(&pod){
			err := client.CoreV1().Pods(*podNamespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})

			if err != nil {
				log.Printf("Failed to delete pod: %s %v", pod.Name, err)
				continue
			}
			
			log.Printf("Cleaned up pod %s", pod.Name)
		}
	}
}

func shouldDeletePod(pod *corev1.Pod) bool {
	podFinishTime := podFinishTime(pod)

	if !podFinishTime.IsZero(){
		age := time.Since(podFinishTime)

		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			if (*deleteSuccessfulAfter > 0 && age >= *deleteSuccessfulAfter){
				return true
			}
		case corev1.PodFailed:
			if (*deleteFailedAfter > 0 && age >= *deleteFailedAfter) {
				return true
			}
		default:
			return false
		}
	}

	return false
}

func podFinishTime(podObj *corev1.Pod) time.Time {
	for _, pc := range podObj.Status.Conditions {
		// Looking for the time when pod's condition "Ready" became "false" (equals end of execution)
		if pc.Type == corev1.PodReady && pc.Status == corev1.ConditionFalse {
			return pc.LastTransitionTime.Time
		}
	}

	return time.Time{}
}