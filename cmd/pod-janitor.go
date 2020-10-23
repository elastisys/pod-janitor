package main

import (
	"flag"
	"log"
	"time"

	cleaner "github.com/filetrust/pod-janitor/pkg"
)

var podNamespace = flag.String("pod-namespace", "", "The kubernetes namespace to run in")
var deleteSuccessfulAfter = flag.Duration("delete-successful-after", 0*time.Minute, "The kubernetes namespace to run in")
var deleteFailedAfter = flag.Duration("delete-failed-after", 0*time.Minute, "The kubernetes namespace to run in")

func main() {
	flag.Parse()

	if *podNamespace == "" {
		log.Fatalf("init failed: pod-namespace argument not set")
	}

	cleanerArgs, err := cleaner.NewCleanerArgs(*podNamespace, *deleteSuccessfulAfter, *deleteFailedAfter)
	if err != nil {
		log.Fatalf("Failed to initialise Cleaner: %v", err)
	}

	err = cleanerArgs.RunCleaner()
	if err != nil {
		log.Fatalf("Failed to run Cleaner: %v", err)
	}
}