// Upgrade: you can't turn a read lock into a write lock.
//
// A user cache has a "get or load" method. It takes RLock to look in the
// map. On a miss, it loads the user and needs Lock to save it. But Lock
// waits until every read lock is gone, including our own. Go has no way
// to "upgrade" RLock to Lock.
//
// Problem: getOrLoadBroken still holds RLock when it calls Lock, so it
// waits for itself forever. Go stops with "fatal error: all goroutines
// are asleep - deadlock!".
//
// Fix: GetOrLoad drops RLock first, then takes Lock and looks again,
// because another goroutine may have loaded the user in between. Comment
// out the line marked "problem", then uncomment the "fix:" line.
//
// Note: same cause as mutex/04-nested-lock. A lock doesn't know who holds
// it, so a goroutine that locks again waits for itself.
//
// Note: GetOrLoad holds Lock during the database call, so every reader
// waits for it. For slow loads, use singleflight (see singleflight/02-cache).
//
// Real use:
//   - a lazy cache in front of a database: db.QueryRowContext on a miss,
//     then save the row. The "look again after Lock" step is called
//     double-checked locking.
//
// Run: go run ./02-concurrency/08-rwmutex/03-upgrade
// Race check: go run -race ./02-concurrency/08-rwmutex/03-upgrade
package main

import (
	"fmt"
	"sync"
	"time"
)

// userCache maps user IDs to names, and loads missing users from the
// database.
type userCache struct {
	mu    sync.RWMutex // guards users
	users map[string]string
}

// getOrLoadBroken tries to take Lock while it still holds RLock.
func (c *userCache) getOrLoadBroken(id string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if name, ok := c.users[id]; ok {
		return name
	}

	name := loadUser(id)
	c.mu.Lock() // waits for every RLock to be released, including ours above
	c.users[id] = name
	c.mu.Unlock()
	return name
}

// GetOrLoad returns the user's name, loading it from the database on a miss.
func (c *userCache) GetOrLoad(id string) string {
	c.mu.RLock()
	name, ok := c.users[id]
	c.mu.RUnlock()
	if ok {
		return name
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if name, ok := c.users[id]; ok { // look again: someone may have loaded it after our RUnlock
		return name
	}
	name = loadUser(id)
	c.users[id] = name
	return name
}

// loadTime is how long the fake database takes.
const loadTime = 50 * time.Millisecond

func main() {
	cache := &userCache{users: map[string]string{"u1": "Rahim"}}

	done := make(chan struct{}) // lookup -> main: "I'm finished"
	go lookup(cache, "u2", done)
	<-done // a stuck lookup shows up here as a deadlock
}

// lookup gets one user through the cache and prints the name.
func lookup(c *userCache, id string, done chan<- struct{}) {
	defer close(done)

	name := c.getOrLoadBroken(id) // problem: takes Lock while holding RLock
	// name := c.GetOrLoad(id) // fix: drop RLock, take Lock, then look again
	fmt.Println(id, "is", name)
}

// loadUser pretends to read a user from the database.
func loadUser(id string) string {
	time.Sleep(loadTime)
	fmt.Println("loaded", id, "from the database")
	return "Karim"
}
