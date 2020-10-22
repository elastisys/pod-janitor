package main

import (
	"context"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var podNamespace = os.Getenv("POD_NAMESPACE")

func main() {
	if podNamespace == "" {
		log.Fatalf("init failed: POD_NAMESPACE environment variable not set")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get cluster config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialise client")
	}

	pods, err := client.CoreV1().Pods(podNamespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Succeeded"})
	if err != nil {
		log.Fatalf("Failed to get list of pods")
	}

	log.Printf("There are %d pods in the cluster with the status of Succeeded\n", len(pods.Items))

	for _, pod := range pods.Items {
		err := client.CoreV1().Pods(podNamespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			log.Printf("Failed to delete pod: %s %v\n", pod.Name, err)
			continue
		}
		
		log.Printf("Cleaned up pod %s\n", pod.Name)
	}
}