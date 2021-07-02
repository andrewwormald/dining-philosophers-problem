package types

import (
	"context"
	"fmt"
	"sync"
)

func NewPhilosopher(id int64, left, right *Chopstick, hostRequestChan chan int64, notifyChan chan struct{}) *Philosopher {
	return &Philosopher{
		ID:    id,
		Left:  left,
		Right: right,
		PermissionRequest: hostRequestChan,
		PermissionGranted: make(chan struct{}),
		FinishedEating: notifyChan,
	}
}

type Philosopher struct {
	ID              int64
	Left            *Chopstick
	Right           *Chopstick
	TimesEaten      int64
	PermissionRequest chan int64
	PermissionGranted chan struct{}
	FinishedEating chan struct{}
}

func (p *Philosopher) EatForever(ctx context.Context, wg *sync.WaitGroup) {
	for {
		if ctx.Err() != nil {
			wg.Done()
			fmt.Println("context cancelled, calling it a day...", p.ID)
			return
		}

		if p.TimesEaten >= 3 {
			wg.Done()
			return
		}

		// Request permission to eat
		p.PermissionRequest <- p.ID
		// Block until permission granted
		<- p.PermissionGranted

		pickupChopsticks(p.Left, p.Right, p.ID)

		fmt.Println(fmt.Sprintf("starting to eat %v", p.ID))
		p.TimesEaten++
		fmt.Println(fmt.Sprintf("finished eating %v", p.ID))
		p.Left.Unlock()
		p.Right.Unlock()
		p.FinishedEating <- struct{}{}
	}
}

// pickupChopsticks is a blocking call that finishes once the Philosopher has picked up both a left and right chopstick.
func pickupChopsticks(left, right *Chopstick, name int64) {
	var wg sync.WaitGroup
	wg.Add(2)
	go obtainChopstick(left, &wg, name, "left")
	go obtainChopstick(right, &wg, name, "right")
	wg.Wait()
}

func obtainChopstick(cs *Chopstick, wg *sync.WaitGroup, name int64, side string) {
	fmt.Println(fmt.Sprintf("%v is waiting for the %v chopstick", name, side))
	cs.Lock()
	fmt.Println(fmt.Sprintf("%v got the %v chopstick", name, side))
	wg.Done()
}
