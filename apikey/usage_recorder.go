package apikey

import (
	"context"
	"log"
	"sync"
	"time"
)

type UsageEvent struct {
	KeyID     int64
	At        time.Time
	IP        *string
	UserAgent *string
}

type UsageRecorder interface {
	Record(event UsageEvent)
	Close(ctx context.Context) error
}

type NoopUsageRecorder struct{}

func (NoopUsageRecorder) Record(event UsageEvent) {}

func (NoopUsageRecorder) Close(ctx context.Context) error { return nil }

type MemoryUsageRecorder struct {
	mu     sync.Mutex
	Events []UsageEvent
}

func (r *MemoryUsageRecorder) Record(event UsageEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Events = append(r.Events, event)
}

func (r *MemoryUsageRecorder) Close(ctx context.Context) error { return nil }

type AsyncUsageRecorder struct {
	store   Store
	timeout time.Duration
	events  chan UsageEvent
	done    chan struct{}
	once    sync.Once
}

func NewAsyncUsageRecorder(store Store, queueSize int, timeout time.Duration) *AsyncUsageRecorder {
	if queueSize < 1 {
		queueSize = 1
	}
	if timeout <= 0 {
		timeout = 500 * time.Millisecond
	}
	r := &AsyncUsageRecorder{
		store:   store,
		timeout: timeout,
		events:  make(chan UsageEvent, queueSize),
		done:    make(chan struct{}),
	}
	go r.run()
	return r
}

func (r *AsyncUsageRecorder) Record(event UsageEvent) {
	select {
	case r.events <- event:
	default:
	}
}

func (r *AsyncUsageRecorder) Close(ctx context.Context) error {
	r.once.Do(func() {
		close(r.events)
	})
	select {
	case <-r.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *AsyncUsageRecorder) run() {
	defer close(r.done)
	for event := range r.events {
		ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
		err := r.store.MarkUsed(ctx, event.KeyID, event.At, event.IP, event.UserAgent)
		cancel()
		if err != nil {
			log.Printf("apikey usage update failed key_id=%d: %v", event.KeyID, err)
		}
	}
}
