package forward

import "errors"

var errClosed = errors.New(`subscription closed`)

type State struct {
	ID    int64
	Ready bool
}

type Subscription struct {
	forward *Forward
	ch      chan State
	close   chan struct{}
	keys    map[int64]bool
}

func (s *Subscription) put(state State) {
	if _, exists := s.keys[state.ID]; exists {
		select {
		case s.ch <- state:
			return
		case <-s.close:
			return
		default:
		}

		select {
		case <-s.ch:
		case <-s.close:
			return
		default:
		}

		select {
		case s.ch <- state:
			return
		case <-s.close:
			return
		}
	}
}
func (s *Subscription) Get() (items []State, e error) {
	select {
	case v := <-s.ch:
		items = append(items, v)
	case <-s.close:
		e = errClosed
	}
	return
}
func (s *Subscription) Unsubscribe() bool {
	return s.forward.Unsubscribe(s)
}
