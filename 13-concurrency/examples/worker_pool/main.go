package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID    int
	Input int
}

type Result struct {
	Job    Job
	Output int
	Worker int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// Simulate processing workload
		time.Sleep(10 * time.Millisecond)
		results <- Result{
			Job:    j,
			Output: j.Input * j.Input,
			Worker: id,
		}
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	// Start worker pool
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Input: j}
	}
	close(jobs) // Signals to workers that no more jobs will arrive

	// Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for res := range results {
		fmt.Printf("Worker %d completed Job %d: input=%d -> output=%d\n",
			res.Worker, res.Job.ID, res.Job.Input, res.Output)
	}
	fmt.Println("All jobs processed successfully.")
}
