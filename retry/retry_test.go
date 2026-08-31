package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fastConfig(maxAttempts int) Config {
	return Config{
		MaxAttempts: maxAttempts,
		BaseDelay:   time.Millisecond,
		MaxDelay:    2 * time.Millisecond,
		Multiplier:  2,
		Jitter:      0,
	}
}

func TestDoSucceedsAfterTransientErrors(t *testing.T) {
	calls := 0
	err := Do(context.Background(), fastConfig(5), func(_ context.Context, attempt int) error {
		calls = attempt
		if attempt < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}

func TestDoExhaustsAttemptsAndReturnsLastError(t *testing.T) {
	last := errors.New("still failing")
	calls := 0
	err := Do(context.Background(), fastConfig(4), func(_ context.Context, _ int) error {
		calls++
		return last
	})
	if !errors.Is(err, last) {
		t.Fatalf("expected last error %v, got %v", last, err)
	}
	if calls != 4 {
		t.Fatalf("expected 4 attempts, got %d", calls)
	}
}

func TestDoStopsOnPermanent(t *testing.T) {
	sentinel := errors.New("bad request")
	calls := 0
	err := Do(context.Background(), fastConfig(10), func(_ context.Context, _ int) error {
		calls++
		return Permanent(sentinel)
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if IsPermanent(err) {
		t.Fatalf("Do should unwrap Permanent before returning")
	}
	if calls != 1 {
		t.Fatalf("expected 1 attempt, got %d", calls)
	}
}

func TestDoRespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()
	err := Do(ctx, Config{MaxAttempts: 0, BaseDelay: 2 * time.Millisecond, MaxDelay: 2 * time.Millisecond, Multiplier: 2}, func(_ context.Context, _ int) error {
		calls++
		return errors.New("keep going")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls == 0 {
		t.Fatalf("expected at least one attempt before cancellation")
	}
}

func TestDelayCapsAtMaxDelay(t *testing.T) {
	c := Config{BaseDelay: time.Second, MaxDelay: 10 * time.Second, Multiplier: 2, Jitter: 0}.withDefaults()
	if got := c.delay(1); got != time.Second {
		t.Fatalf("attempt 1 delay = %v, want 1s", got)
	}
	if got := c.delay(20); got != 10*time.Second {
		t.Fatalf("attempt 20 delay = %v, want capped 10s", got)
	}
}
