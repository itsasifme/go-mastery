// Registry: track who's online with sync.Map.
//
// A chat server keeps one entry per connected user. Each user connects
// and leaves in their own goroutine, and only touches their own key. A
// broadcast walks over everyone who's online with Range. This code needs
// no mutex: sync.Map is safe to use from many goroutines.
//
// This is one of the two cases sync.Map is built for: goroutines that
// read and write different keys. The other is a key that's written once
// and read many times (see 02-load-or-store). For anything else, a plain
// map with a Mutex or RWMutex is simpler (see rwmutex/01-cache).
//
// The sync.Map is kept inside a registry type, and only its methods touch
// it. Callers can only add a *session, so the type assertion in Broadcast
// can't fail. This is how production code usually uses sync.Map, because
// its keys and values are any.
//
// Range visits the keys in no fixed order, so the broadcast order can
// change between runs.
//
// Real use:
//   - a WebSocket hub that maps user IDs to their connections, and sends
//     a message to everyone who's online.
//   - tracking running jobs by ID, so an admin page can list them.
//
// Run: go run ./02-concurrency/13-sync-map/01-registry
// Race check: go run -race ./02-concurrency/13-sync-map/01-registry
package main

import (
	"fmt"
	"sync"
	"time"
)

// session is one connected user.
type session struct {
	user   string
	device string
}

// Send pretends to send msg to s.
func (s *session) Send(msg string) {
	fmt.Printf("  to %s (%s)\n", s.user, s.device)
}

// registry tracks who's online. Use it through a pointer, because a
// sync.Map must not be copied after first use.
type registry struct {
	sessions sync.Map // user -> *session. Only the methods below touch it
}

// Add marks s as online.
func (r *registry) Add(s *session) {
	r.sessions.Store(s.user, s)
}

// Remove marks user as offline.
func (r *registry) Remove(user string) {
	r.sessions.Delete(user)
}

// IsOnline reports whether user has a session.
func (r *registry) IsOnline(user string) bool {
	_, ok := r.sessions.Load(user)
	return ok
}

// Broadcast sends msg to every session that's online right now.
func (r *registry) Broadcast(msg string) {
	fmt.Printf("broadcast %q:\n", msg)
	r.sessions.Range(func(_, v any) bool {
		v.(*session).Send(msg) // v is any. Add only stores *session, so this is safe
		return true            // keep going to the next session
	})
}

// broadcastAt is when main sends the broadcast. One user has left by then.
const broadcastAt = 150 * time.Millisecond

func main() {
	visits := []struct {
		user, device string
		stay         time.Duration // how long the user stays online
	}{
		{"rahim", "phone", 100 * time.Millisecond}, // leaves before the broadcast
		{"karim", "laptop", 300 * time.Millisecond},
		{"nadia", "phone", 300 * time.Millisecond},
	}

	online := &registry{}

	var wg sync.WaitGroup
	for _, v := range visits {
		wg.Go(func() { visit(online, &session{user: v.user, device: v.device}, v.stay) })
	}

	time.Sleep(broadcastAt)
	online.Broadcast("server restarts in 5 minutes")
	fmt.Println("rahim online?", online.IsOnline("rahim"))
	fmt.Println("karim online?", online.IsOnline("karim"))

	wg.Wait()
	fmt.Println("karim online at the end?", online.IsOnline("karim"))
}

// visit adds s to online, stays for a while, then removes it.
func visit(online *registry, s *session, stay time.Duration) {
	online.Add(s)
	fmt.Println(s.user, "connected")

	time.Sleep(stay)

	online.Remove(s.user)
	fmt.Println(s.user, "left")
}
