package xgrpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	pool := New()
	if pool == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNewPoolWithContext(t *testing.T) {
	ctx := context.Background()
	pool := NewWithContext(ctx)
	if pool == nil {
		t.Fatal("NewWithContext() returned nil")
	}

	// Test with nil context
	pool = NewWithContext(nil)
	if pool == nil {
		t.Fatal("NewWithContext(nil) returned nil")
	}
}

func TestGoroutinePool_Add(t *testing.T) {
	pool := New()
	counter := atomic.Int32{}

	// Add 10 goroutines
	for i := 0; i < 10; i++ {
		pool.Add(func(ctx context.Context) error {
			counter.Add(1)
			return nil
		})
	}

	pool.Wait()
	if counter.Load() != 10 {
		t.Errorf("Expected counter to be 10, got %d", counter.Load())
	}
}

func TestGoroutinePool_Stop(t *testing.T) {
	pool := New()
	counter := atomic.Int32{}

	// Add a long-running goroutine
	pool.Add(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Hour):
			counter.Add(1)
			return nil
		}
	})

	// Stop the pool immediately
	pool.Stop()

	if counter.Load() != 0 {
		t.Error("Goroutine was not cancelled")
	}
}

func TestGoroutinePool_AddWithCancel(t *testing.T) {
	pool := New()
	cancelCalled := atomic.Bool{}

	// Add a task with cancel function
	pool.AddWithCancel(
		func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
		func(ctx context.Context) error {
			cancelCalled.Store(true)
			return nil
		},
	)()

	pool.Stop()

	if !cancelCalled.Load() {
		t.Error("Cancel function was not called")
	}
}

func TestGoroutinePool_AddWithRecover(t *testing.T) {
	pool := New()
	panicCaught := atomic.Bool{}

	// Add a task that will panic
	pool.AddWithRecover(
		func(ctx context.Context) error {
			panic("test panic")
		},
		func(ctx context.Context, err interface{}) {
			if err != nil {
				panicCaught.Store(true)
			}
		},
	)()

	pool.Wait()

	if !panicCaught.Load() {
		t.Error("Panic was not caught")
	}
}

func TestGoroutinePool_ErrorHandling(t *testing.T) {
	pool := New()
	errorHandled := atomic.Bool{}

	// Set custom error handler
	pool.UncaughtErrorHandler(func(err any) {
		if err != nil {
			errorHandled.Store(true)
		}
	})

	// Add a task that returns an error
	pool.Add(func(ctx context.Context) error {
		return errors.New("test error")
	})

	pool.Wait()

	if !errorHandled.Load() {
		t.Error("Error was not handled")
	}
}

func BenchmarkGoroutinePool_Add(b *testing.B) {
	pool := New()
	counter := atomic.Int32{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Add(func(ctx context.Context) error {
			counter.Add(1)
			return nil
		})
	}
	pool.Wait()
}
