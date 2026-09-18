package fanoutfanin_solution

import (
	"sync"
)

// FanOutFanIn distributes inputs across workerCount workers and consolidates results.
func FanOutFanIn(inputs []int, workerCount int, transform func(int) int) ([]int, error) {
	if len(inputs) == 0 {
		return []int{}, nil
	}

	if workerCount <= 0 {
		workerCount = 1
	}

	jobs := make(chan int, len(inputs))
	results := make(chan int, len(inputs))

	var wg sync.WaitGroup

	// Spawn workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for val := range jobs {
				results <- transform(val)
			}
		}()
	}

	// Feed jobs
	for _, val := range inputs {
		jobs <- val
	}
	close(jobs)

	// Close results channel when all workers exit
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect outputs
	out := make([]int, 0, len(inputs))
	for res := range results {
		out = append(out, res)
	}

	return out, nil
}
