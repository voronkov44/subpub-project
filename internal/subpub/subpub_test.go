package subpub

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMultipleIndependentSubscribers(t *testing.T) {
	pubsub := NewSubPub()
	defer pubsub.Close(context.Background())

	var wg sync.WaitGroup
	messagesToSend := []string{"msg1", "msg2", "msg3"}

	type result struct {
		mu       sync.Mutex
		messages []string
	}

	subscriberResults := make([]*result, 3)

	for i := 0; i < 3; i++ {
		subscriberResults[i] = &result{}
		wg.Add(1)
		subIndex := i

		_, err := pubsub.Subscribe("test-topic", func(msg interface{}) {
			r := subscriberResults[subIndex]
			r.mu.Lock()
			r.messages = append(r.messages, msg.(string))
			r.mu.Unlock()

			// Симулируем "медленного" подписчика на втором подписчике
			if subIndex == 1 {
				time.Sleep(200 * time.Millisecond)
			}

			// Если получили все сообщения — завершаем горутину
			if len(r.messages) == len(messagesToSend) {
				wg.Done()
			}
		})
		require.NoError(t, err)
	}

	// Публикуем все сообщения
	for _, m := range messagesToSend {
		err := pubsub.Publish("test-topic", m)
		require.NoError(t, err)
	}

	// Ждём завершения всех подписчиков
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("test timed out waiting for subscribers")
	}

	// Проверяем, что все подписчики получили все сообщения в правильном порядке
	for i, res := range subscriberResults {
		res.mu.Lock()
		require.Equal(t, messagesToSend, res.messages, "subscriber %d did not receive correct messages", i)
		res.mu.Unlock()
	}
}

func TestPublishSubscribe(t *testing.T) {
	sp := NewSubPub()
	defer sp.Close(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	sub, err := sp.Subscribe("test-topic", func(msg interface{}) {
		defer wg.Done()
		if msg != "hello" {
			t.Errorf("expected message 'hello', got '%v'", msg)
		}
	})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	err = sp.Publish("test-topic", "hello")
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	wg.Wait()
	sub.Unsubscribe()
}

func TestUnsubscribe(t *testing.T) {
	sp := NewSubPub()
	defer sp.Close(context.Background())

	called := false

	sub, _ := sp.Subscribe("topic", func(msg interface{}) {
		called = true
	})
	sub.Unsubscribe()

	sp.Publish("topic", "hello")
	time.Sleep(100 * time.Millisecond)

	if called {
		t.Errorf("handler was called after unsubscribe")
	}
}

func TestCloseWithContextCancel(t *testing.T) {
	sp := NewSubPub()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := sp.Close(ctx)
	if err != context.Canceled {
		t.Errorf("Close should return context.Canceled if context is canceled, got: %v", err)
	}
}
