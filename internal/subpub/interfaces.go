package subpub

import "context"

type MessageHandler func(msg interface{})

type Subscriber interface {
	Unsubscribe()
}

type PubSub interface {
	Subscribe(subject string, cb MessageHandler) (Subscriber, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

func NewSubPub() PubSub {
	return newSubPub()
}
