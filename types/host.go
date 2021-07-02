package types

import "context"

type Host struct {
	Requests           chan int64    // Philosopher can notify host that they want to start eating
	Notify             chan struct{} // Philosopher can notify host that they are done eating
	MembersEating      int
	PermissionChannels map[int64]chan struct{}
}

func NewHost(philos []*Philosopher, requests chan int64, notify chan struct{}) *Host {
	chans := make(map[int64]chan struct{}, 5)
	for _, p := range philos {
		chans[p.ID] = p.PermissionGranted
	}

	return &Host{
		Requests:           requests,
		Notify:             notify,
		PermissionChannels: chans,
	}
}

func (h *Host) ManageForever(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		select {
		case r := <-h.Requests:
			// Consume the request and ignore if more than 2 are already eating
			if h.MembersEating >= 2 {
				// Try cycling the requests to keep it going until a Philosopher has stopped eating and notified
				h.Requests <- r
				continue
			}

			// Add the Philosopher's id to the slice of those allowed to eat and emit the event that they can eat
			h.MembersEating++
			// Notify the specific philosopher that they have permission to eat
			h.PermissionChannels[r] <- struct{}{}
			continue
		case <-h.Notify:
			h.MembersEating--
		default:
			continue
		}
	}
}
