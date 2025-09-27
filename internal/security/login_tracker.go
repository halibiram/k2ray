package security

import (
	"sync"
	"time"
)

const (
	maxLoginAttempts = 5
	lockoutDuration  = 15 * time.Minute
)

// loginAttempt holds the details of recent login attempts for a given key (username or IP).
type loginAttempt struct {
	count      int
	lastAttempt time.Time
}

// failedLoginTracker is a thread-safe in-memory store for tracking failed login attempts.
type failedLoginTracker struct {
	attempts map[string]*loginAttempt
	mu       sync.Mutex
}

// Global instance of the login tracker.
var tracker = &failedLoginTracker{
	attempts: make(map[string]*loginAttempt),
}

// RecordFailedAttempt logs a failed login attempt for a given key.
// The key can be a username or an IP address.
func RecordFailedAttempt(key string) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	attempt, exists := tracker.attempts[key]
	if !exists {
		tracker.attempts[key] = &loginAttempt{count: 1, lastAttempt: time.Now()}
		return
	}

	// If the last attempt was outside the lockout window, reset the count.
	if time.Since(attempt.lastAttempt) > lockoutDuration {
		attempt.count = 1
		attempt.lastAttempt = time.Now()
		return
	}

	attempt.count++
	attempt.lastAttempt = time.Now()
}

// IsLockedOut checks if a given key is currently locked out due to too many failed attempts.
func IsLockedOut(key string) bool {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	attempt, exists := tracker.attempts[key]
	if !exists {
		return false
	}

	// If the lockout period has expired, the user is no longer locked out.
	if time.Since(attempt.lastAttempt) > lockoutDuration {
		// We can also remove the entry to free up memory.
		delete(tracker.attempts, key)
		return false
	}

	return attempt.count >= maxLoginAttempts
}

// ResetAttempts clears the attempt count for a given key.
// This should be called after a successful login.
func ResetAttempts(key string) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	delete(tracker.attempts, key)
}