package types

import "sync"

type Chopstick struct {
	sync.Mutex
}