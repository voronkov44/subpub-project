package subpub

import (
	"context"
	"sync"
)

type subscription struct {
	subject string
	handler MessageHandler
	ch      chan interface{}
	closeCh chan struct{}
	pub     *subPub
}

func (s *subscription) Unsubscribe() {
	s.pub.unsubscribe(s.subject, s)
	close(s.closeCh)
}

type subPub struct {
	mu          sync.RWMutex
	subscribers map[string][]*subscription
	closed      bool
}

func newSubPub() *subPub {
	return &subPub{
		subscribers: make(map[string][]*subscription),
	}
}

func (sp *subPub) Subscribe(subject string, handler MessageHandler) (Subscriber, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		return nil, context.Canceled
	}

	sub := &subscription{
		subject: subject,
		handler: handler,
		ch:      make(chan interface{}, 100),
		closeCh: make(chan struct{}),
		pub:     sp,
	}

	sp.subscribers[subject] = append(sp.subscribers[subject], sub)

	go func() {
		for {
			select {
			case msg := <-sub.ch:
				handler(msg)
			case <-sub.closeCh:
				return
			}
		}
	}()

	return sub, nil
}

func (sp *subPub) Publish(subject string, msg interface{}) error {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.closed {
		return context.Canceled
	}

	for _, sub := range sp.subscribers[subject] {
		select {
		case sub.ch <- msg:
		default:
			// если очередь подписчика забита — пропускаем, не блокируем
		}
	}

	return nil
}

func (sp *subPub) unsubscribe(subject string, sub *subscription) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	subs := sp.subscribers[subject]
	for i, s := range subs {
		if s == sub {
			sp.subscribers[subject] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

func (sp *subPub) Close(ctx context.Context) error {
	sp.mu.Lock()
	sp.closed = true
	sp.mu.Unlock()

	done := make(chan struct{})

	go func() {
		sp.mu.RLock()
		defer sp.mu.RUnlock()

		for _, subs := range sp.subscribers {
			for _, sub := range subs {
				close(sub.closeCh)
			}
		}

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
