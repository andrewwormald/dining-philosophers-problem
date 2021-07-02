package main

import (
	"context"
	"dinnertable/types"
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Process started")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var chopsticks []*types.Chopstick
	for i := 0; i < 5; i++ {
		chopsticks = append(chopsticks, &types.Chopstick{})
	}

	requests := make(chan int64, 5)
	notifies := make(chan struct{}, 5)

	var philos []*types.Philosopher
	for i := int64(0); i < 5; i++ {
		philos = append(philos, types.NewPhilosopher(i+1, chopsticks[i], chopsticks[(i+1) %5], requests, notifies))
	}

	host := types.NewHost(philos, requests, notifies)
	go host.ManageForever(ctx)

	var wg sync.WaitGroup
	wg.Add(len(philos))
	for _, philo := range philos {
		go philo.EatForever(ctx, &wg)
	}
	wg.Wait()
	fmt.Println("Process completed")
}




