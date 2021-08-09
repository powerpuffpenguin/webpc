package forward

import (
	"context"
)

type State struct {
	ID    int64
	Ready bool
}

type Subscription struct {
	forward *Forward
	ch      chan State
	ctx     context.Context
	keys    map[int64]bool
}

func (s *Subscription) put(state State) {
	done := s.ctx.Done()

	if _, exists := s.keys[state.ID]; exists {
		select {
		case s.ch <- state:
			return
		case <-done:
			return
		default:
		}

		select {
		case s.ch <- state:
			return
		case <-s.ch:
		case <-done:
			return
		default:
		}

		select {
		case s.ch <- state:
			return
		case <-done:
			return
		}
	}
}
func (s *Subscription) Change(targets []int64) {
	s.forward.change(s, targets)
}
func (s *Subscription) Get() (items []State, e error) {
	done := s.ctx.Done()
	select {
	case v := <-s.ch:
		items = append(items, v)
	case <-done:
		e = s.ctx.Err()
		return
	}
	for len(items) < 100 {
		select {
		case v := <-s.ch:
			items = append(items, v)
		case <-done:
			e = s.ctx.Err()
			return
		default:
			return
		}
	}
	return
}
