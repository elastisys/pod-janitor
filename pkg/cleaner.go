package cleaner

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type CleanerArgs struct {
	PodNamespace            string
	Client                  *kubernetes.Clientset
	DeleteSuccessfulAfter   time.Duration
	DeleteFailedAfter       time.Duration
}

func NewCleanerArgs(podNamespace string, deleteSuccessfulAfter, deleteFailedAfter time.Duration)(*CleanerArgs, error){
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cleanerArgs := &CleanerArgs {
		PodNamespace:            podNamespace,
		Client:                  client,
		DeleteSuccessfulAfter:   deleteSuccessfulAfter,
		DeleteFailedAfter:       deleteFailedAfter,
	}

	return cleanerArgs, nil
}

func (ca CleanerArgs) RunCleaner() error {
	succeededPods, err := ca.Client.CoreV1().Pods(ca.PodNamespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Succeeded"})
	if err != nil {
		err = fmt.Errorf("Failed to get list of successful pods: %v", err)
		return err
	}

	failedPods, err := ca.Client.CoreV1().Pods(ca.PodNamespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Failed"})
	if err != nil {
		err = fmt.Errorf("Failed to get list of failed pods: %v", err)
		return err
	}

	pods := succeededPods.Items[:0]
	pods = append(pods, succeededPods.Items[:]...)
	pods = append(pods, failedPods.Items[:]...)

	for _, pod := range pods {
		
		if shouldDeletePod(&pod, ca.DeleteSuccessfulAfter, ca.DeleteFailedAfter){
			err := ca.Client.CoreV1().Pods(ca.PodNamespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})

			if err != nil {
				log.Printf("Failed to delete pod: %s %v", pod.Name, err)
				continue
			}
			
			log.Printf("Cleaned up pod %s", pod.Name)
		}
	}

	return nil
}

func shouldDeletePod(pod *corev1.Pod, successful, failed time.Duration) bool {
	podFinishTime := podFinishTime(pod)

	if !podFinishTime.IsZero(){
		age := time.Since(podFinishTime)

		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			if (successful > 0 && age >= successful){
				return true
			}
		case corev1.PodFailed:
			if (failed > 0 && age >= failed) {
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