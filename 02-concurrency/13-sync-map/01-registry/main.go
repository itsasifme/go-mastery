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

	var online sync.Map // user -> *session

	var wg sync.WaitGroup
	for _, v := range visits {
		wg.Go(func() { visit(&online, &session{user: v.user, device: v.device}, v.stay) })
	}

	time.Sleep(broadcastAt)
	broadcast(&online, "server restarts in 5 minutes")
	fmt.Println("rahim online?", isOnline(&online, "rahim"))
	fmt.Println("karim online?", isOnline(&online, "karim"))

	wg.Wait()
	fmt.Println("karim online at the end?", isOnline(&online, "karim"))
}

// visit adds s to online, stays for a while, then removes it.
func visit(online *sync.Map, s *session, stay time.Duration) {
	online.Store(s.user, s)
	fmt.Println(s.user, "connected")

	time.Sleep(stay)

	online.Delete(s.user)
	fmt.Println(s.user, "left")
}

// broadcast sends msg to every session that's online right now.
func broadcast(online *sync.Map, msg string) {
	fmt.Printf("broadcast %q:\n", msg)
	online.Range(func(_, v any) bool { return send(v.(*session), msg) })
}

// send pretends to send msg to one session. It returns true, so Range
// goes on to the next one.
func send(s *session, msg string) bool {
	fmt.Printf("  to %s (%s)\n", s.user, s.device)
	return true
}

// isOnline reports whether user has a session.
func isOnline(online *sync.Map, user string) bool {
	_, ok := online.Load(user)
	return ok
}
