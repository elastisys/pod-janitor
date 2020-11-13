package main

import (
	"flag"
	"log"
	"time"

	cleaner "github.com/filetrust/pod-janitor/pkg"
)

var podNamespace = flag.String("pod-namespace", "", "The kubernetes namespace to run in")
var deleteSuccessfulAfter = flag.Duration("delete-successful-after", 0*time.Minute, "Delete pods in succeeded state after X duration (golang duration format, e.g 5m), 0 - never delete")
var deleteFailedAfter = flag.Duration("delete-failed-after", 0*time.Minute, "Delete pods in failed state after X duration (golang duration format, e.g 5m), 0 - never delete")

func main() {
	flag.Parse()

	if *podNamespace == "" {
		log.Fatalf("init failed: pod-namespace argument not set")
	}

	cleanerArgs, err := cleaner.NewCleanerArgs(*podNamespace, *deleteSuccessfulAfter, *deleteFailedAfter)
	if err != nil {
		log.Fatalf("Failed to initialise Cleaner: %v", err)
	}

	cleanerArgs.RunCleaner()
}
