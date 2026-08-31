// Package retry provides a small, dependency-free exponential-backoff retry
// helper shared by all sdup adapters. It is deliberately generic: callers pass
// a function that performs one attempt and classify their own errors as
// transient (return them as-is) or fatal (wrap with Permanent).
package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// Config controls the backoff schedule used by Do.
type Config struct {
	// MaxAttempts is the total number of attempts, including the first.
	// Values <= 0 mean "keep retrying until the context is cancelled".
	MaxAttempts int
	// BaseDelay is the wait before the second attempt. Each subsequent wait is
	// the previous one multiplied by Multiplier, capped at MaxDelay.
	BaseDelay time.Duration
	// MaxDelay caps any single backoff wait.
	MaxDelay time.Duration
	// Multiplier is the growth factor between attempts. Coerced to >= 1.
	Multiplier float64
	// Jitter is the fraction (0..1) of each wait that is randomised, so a fleet
	// of adapters restarting together does not retry in lockstep.
	Jitter float64
}

// DefaultConfig returns a schedule of 8 attempts spanning a little over a
// minute: roughly 0.5s, 1s, 2s, 4s, 8s, 16s, 30s (capped) between tries,
// each +/- 20% jitter.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 8,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    30 * time.Second,
		Multiplier:  2,
		Jitter:      0.2,
	}
}

func (c Config) withDefaults() Config {
	if c.BaseDelay <= 0 {
		c.BaseDelay = 500 * time.Millisecond
	}
	if c.MaxDelay <= 0 {
		c.MaxDelay = 30 * time.Second
	}
	if c.Multiplier < 1 {
		c.Multiplier = 2
	}
	if c.Jitter < 0 {
		c.Jitter = 0
	}
	if c.Jitter > 1 {
		c.Jitter = 1
	}
	return c
}

// permanent wraps an error that Do must not retry.
type permanent struct{ err error }

func (p permanent) Error() string { return p.err.Error() }
func (p permanent) Unwrap() error { return p.err }

// Permanent marks err as non-retryable. Returning it (or anything that wraps
// it) from the function passed to Do stops the loop immediately; Do then
// returns the underlying error, not the wrapper.
func Permanent(err error) error {
	if err == nil {
		return nil
	}
	return permanent{err: err}
}

// IsPermanent reports whether err was marked with Permanent anywhere in its chain.
func IsPermanent(err error) bool {
	var p permanent
	return errors.As(err, &p)
}

// Do calls fn until it returns nil, returns a Permanent error, exhausts the
// attempt budget, or ctx is cancelled. attempt is 1-based. On give-up it
// returns the last error fn produced; if ctx ended first, that error wraps
// ctx.Err() with the last error for context.
func Do(ctx context.Context, cfg Config, fn func(ctx context.Context, attempt int) error) error {
	cfg = cfg.withDefaults()
	var lastErr error
	for attempt := 1; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return ctxError(err, lastErr)
		}

		lastErr = fn(ctx, attempt)
		if lastErr == nil {
			return nil
		}
		if IsPermanent(lastErr) {
			var p permanent
			errors.As(lastErr, &p)
			return p.err
		}
		if cfg.MaxAttempts > 0 && attempt >= cfg.MaxAttempts {
			return lastErr
		}

		timer := time.NewTimer(cfg.delay(attempt))
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctxError(ctx.Err(), lastErr)
		}
	}
}

func ctxError(ctxErr, lastErr error) error {
	if lastErr == nil {
		return ctxErr
	}
	return fmt.Errorf("%w (last error: %v)", ctxErr, lastErr)
}

// delay returns the backoff wait that follows the given 1-based attempt.
func (c Config) delay(attempt int) time.Duration {
	d := float64(c.BaseDelay) * math.Pow(c.Multiplier, float64(attempt-1))
	if d > float64(c.MaxDelay) {
		d = float64(c.MaxDelay)
	}
	if c.Jitter > 0 {
		delta := c.Jitter * d
		d = d - delta + rand.Float64()*(2*delta)
	}
	if d < 0 {
		d = 0
	}
	return time.Duration(d)
}
