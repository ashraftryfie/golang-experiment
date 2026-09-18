package main

import (
	"fmt"
	"time"
)

func generateStream(label string, interval time.Duration, count int) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 1; i <= count; i++ {
			time.Sleep(interval)
			ch <- fmt.Sprintf("[%s] event #%d", label, i)
		}
	}()
	return ch
}

func main() {
	streamA := generateStream("Sensor-A", 30*time.Millisecond, 4)
	streamB := generateStream("Sensor-B", 50*time.Millisecond, 3)

	timer := time.NewTimer(300 * time.Millisecond)
	defer timer.Stop()

	doneA, doneB := false, false

	for !doneA || !doneB {
		select {
		case msg, ok := <-streamA:
			if !ok {
				doneA = true
				streamA = nil // Nil channels never select!
			} else {
				fmt.Println("Received:", msg)
			}

		case msg, ok := <-streamB:
			if !ok {
				doneB = true
				streamB = nil // Nil channels never select!
			} else {
				fmt.Println("Received:", msg)
			}

		case <-timer.C:
			fmt.Println("Overall timeout reached!")
			return
		}
	}

	fmt.Println("All streams consumed cleanly.")
}
