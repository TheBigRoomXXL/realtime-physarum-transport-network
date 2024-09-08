package main

import (
	"fmt"
	"time"
)

type Timing struct {
	name    string
	labels  []string
	timings []time.Time
}

func (t *Timing) Step(label string) {
	t.labels = append(t.labels, label)
	t.timings = append(t.timings, time.Now())
}

func (t *Timing) Print() {
	total := t.timings[len(t.timings)-1].Sub(t.timings[0])
	fmt.Printf("%s: took %v \n", t.name, total)

	for i := 0; i < len(t.timings)-1; i++ {
		duration := t.timings[i+1].Sub(t.timings[i])
		percentage := float64(duration.Microseconds()) / float64(total.Microseconds()) * 100
		fmt.Printf("  ↳ %s: %.0f%% - %v\n", t.labels[i], percentage, duration)
	}
}
