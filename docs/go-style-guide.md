<a id="top"></a>

<div align="center">

# 🐹 Go Style Guide

**How Go code is usually written, shown with small examples.**

[![Official Go](https://img.shields.io/badge/Official-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](#official)
[![Google](https://img.shields.io/badge/Google-Go_Style_Guide-4285F4?style=for-the-badge&logo=google&logoColor=white)](#google)
[![Uber](https://img.shields.io/badge/Uber-Go_Style_Guide-000000?style=for-the-badge&logo=uber&logoColor=white)](#uber)

**In a hurry? 👉 [Jump to the summary](#summary)**

</div>

---

## 📚 Contents

| Part | What it covers |
| --- | --- |
| [1. Official Go](#official) | `gofmt`, the compiler, and Go Code Review Comments. **Everyone follows these.** |
| [2. Google](#google) | What Google adds on top of Part 1 |
| [3. Uber](#uber) | What Uber adds on top of Part 1 |
| [4. Summary](#summary) | The most important rules from all three, in one table |
| [5. Tools](#tools) | Linters that check all of this |

---

<a id="official"></a>

## 1. Official Go

![Official Go](https://img.shields.io/badge/Official-Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

**Sources:** [`gofmt`](https://pkg.go.dev/cmd/gofmt) · [Effective Go](https://go.dev/doc/effective_go) · [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

These rules come from the Go team. Google and Uber both follow all of them.

### 🎨 Formatting

`gofmt` formats everything. It has no options, so all Go code looks the same.

```go
// ❌ Before gofmt
func add(a int,b int)int{
return a+b}

// ✅ After gofmt: tabs, spaces, braces all fixed
func add(a int, b int) int {
    return a + b
}
```

```bash
gofmt -w .        # format all files
go fmt ./...      # the same, for every package
goimports -w .    # gofmt + add missing imports + remove unused ones
```

The compiler enforces some style too:

```go
// ❌ Compile error: brace on the next line
if ok
{
}

// ❌ Compile error: imported and not used
import "os"
```

### 🏷️ Naming

```go
// MixedCaps, never snake_case. The first letter decides public or private.
var userName string        // ✅ private
const MaxSize = 10         // ✅ public
var user_name string       // ❌
const MAX_SIZE = 10        // ❌

// Acronyms keep one case
var userID int             // ✅
func ServeHTTP() {}        // ✅
var userId int             // ❌

// Getters have no "Get", setters use "Set"
func (u *User) Name() string      // ✅
func (u *User) SetName(n string)  // ✅
func (u *User) GetName() string   // ❌

// Constructors: New + type
func NewServer() *Server          // ✅

// Receivers: 1–2 letters, the same in every method
func (s *Server) Start()          // ✅
func (this *Server) Start()       // ❌

// Short names for small scopes
for i, u := range users {}        // ✅
```

**Packages:** short, lowercase, one word. Don't repeat the package name in what it exports.

```go
package user               // ✅
package user_utils         // ❌
package utils              // ❌ says nothing about what's inside

// the caller writes user.New(), user.Store
func New() *User           // ✅
type Store struct{}        // ✅

// the caller writes user.NewUser(), user.UserStore
func NewUser() *User       // ❌
type UserStore struct{}    // ❌
```

**Interfaces:** one method → method name + `er`.

```go
type Reader interface { Read(p []byte) (n int, err error) }  // ✅
type IReader interface { ... }                               // ❌
```

**Files:** lowercase, `_` allowed: `user_store.go`, `user_store_test.go`.

> [!WARNING]
> `_windows.go`, `_linux.go`, and `_amd64.go` are **only built on that OS or CPU**. A file named `notes_windows.go` silently disappears on Linux.

### 💬 Comments

A doc comment is a full sentence that **starts with the name** it describes.

```go
// Package user manages user accounts.
package user

// Server handles HTTP requests for the user API.
type Server struct{ ... }

// Parse reads a config file and returns the settings.
func Parse(path string) (*Config, error)
```

For a program, start with the program name:

```go
// Seedgen generates seed data for tests.
package main
```

### 🚨 Errors

```go
// ❌ error thrown away, and the normal code is nested in else
x, _ := load()
if err == nil {
    use(x)
} else {
    return err
}

// ✅ check the error first, keep the normal code at the left edge
x, err := load()
if err != nil {
    return err
}
use(x)
```

```go
// Error messages: lowercase, no period. They get joined together.
errors.New("user not found")   // ✅ prints: load: user not found
errors.New("User not found.")  // ❌ prints: load: User not found.

// Don't return special values to mean "error"
func Find(key string) string                   // ❌ "" when not found?
func Lookup(key string) (value string, ok bool) // ✅
```

Don't use `panic` for normal errors. Return an `error`.

### 🧭 Functions and methods

**Value or pointer receiver?**

```go
type Counter struct{ n int }

func (c Counter) Inc()  { c.n++ }  // ❌ changes a copy, the caller sees nothing
func (c *Counter) Inc() { c.n++ }  // ✅ changes the original
```

| Use a pointer `*T` when… | Use a value `T` when… |
| --- | --- |
| the method changes the receiver | it's a map, func, or channel type |
| the struct has a `sync.Mutex` | it's small and never changes, like `time.Time` |
| the struct is large, or you're not sure | |

Don't mix them. If one method needs a pointer, use pointers for all methods of that type.

**Return values:** `error` last. Name results only when it helps.

```go
func Load(path string) (*Config, error)   // ✅
func Load(path string) (error, *Config)   // ❌
func Location() (lat, long float64)       // ✅ names make it clear
```

**Interfaces:** define them where they're **used**, and return structs.

```go
// ✅ the package that needs a store defines the small interface it uses
type userGetter interface {
    Get(ctx context.Context, id string) (*User, error)
}

// ✅ constructors return the concrete type
func NewPostgresStore(db *sql.DB) *PostgresStore
```

### 🔀 Concurrency

```go
// ❌ context stored in a struct
type Worker struct{ ctx context.Context }

// ✅ context is the first parameter
func (w *Worker) Run(ctx context.Context) error
```

Make it clear when every goroutine stops. A goroutine that never stops is a leak.

### 📦 Imports, slices, tests, security

```go
// Imports: standard library, blank line, everything else
import (
    "fmt"

    "github.com/gin-gonic/gin"
)

// Empty slice: prefer the nil form
var s []string     // ✅
s := []string{}    // ❌ (except for JSON: nil → null, []string{} → [])

// Test failures: what was called, got, then want
t.Errorf("Count(%q) = %d, want %d", in, got, want)

// Keys and tokens: crypto/rand, never math/rand
```

**Line length:** no fixed limit. Break lines where it reads well.

<div align="right"><a href="#top">⬆ back to top</a></div>

---

<a id="google"></a>

## 2. Google Go Style Guide

![Google](https://img.shields.io/badge/Google-Go_Style_Guide-4285F4?style=for-the-badge&logo=google&logoColor=white)

**Source:** [google.github.io/styleguide/go](https://google.github.io/styleguide/go/)

Google follows all of Part 1. This part shows what it adds.

### 💡 Core ideas

**Order of importance:** 1. Clarity → 2. Simplicity → 3. Concision → 4. Maintainability → 5. Consistency

**Least mechanism:** use plain language features (slice, map, channel, loop) before a library, and the standard library before an outside package.

### 🏷️ Naming

```go
// Don't put the type in the name
var users []User           // ✅
var userSlice []User       // ❌

// Don't repeat the receiver or package
func (u *User) Save() error      // ✅ u.Save()
func (u *User) SaveUser() error  // ❌ u.SaveUser()

// "Must" means it panics instead of returning an error. Only use at startup.
var re = regexp.MustCompile(`^\d+$`)

// Underscores only in test names, generated code, and cgo
func TestParse_EmptyInput(t *testing.T)  // ✅

// Test helper packages end with "test"
package creditcardtest

// Protobuf imports end with "pb"
import foopb "example.com/foo/proto"
```

### 💬 Comments

Say what the caller must know: is it safe for goroutines, does the zero value work, what errors come back.

```go
// Server handles HTTP requests for the user API.
// A Server is safe for use by multiple goroutines.
// The zero value is not usable; call NewServer.
type Server struct{ ... }

// Get returns the user with the given id.
// It returns ErrNotFound if there is no such user.
func (s *Server) Get(ctx context.Context, id string) (*User, error)

// Deprecated: Use NewClient instead.
func Dial(addr string) *Client
```

### 🚨 Errors

Choose `%w` or `%v` on purpose, and put it at the end.

```go
// %w: callers can check the inner error with errors.Is / errors.As
return fmt.Errorf("load config: %w", err)

// %v: callers only get the text. Use it to hide internal details.
return fmt.Errorf("load config: %v", err)
```

### 🧭 Functions

```go
// Many parameters → an options struct
type ServerOptions struct {
    Port    int
    Timeout time.Duration
}

func NewServer(opts ServerOptions) *Server

// Give channels a direction
func produce(out chan<- int)   // send only
func consume(in <-chan int)    // receive only

// Modern Go
var v any           // ✅ instead of interface{}
// Use generics only when they remove real repetition
```

### 🧪 Tests

```go
// ❌ assertion library
assert.Equal(t, want, got)

// ✅ plain check with a clear message
if got != want {
    t.Errorf("Parse(%q) = %v, want %v", in, got, want)
}

// ✅ compare structs with go-cmp and print the difference
if diff := cmp.Diff(want, got); diff != "" {
    t.Errorf("Parse(%q) mismatch (-want +got):\n%s", in, diff)
}
```

Use `t.Error` to keep going. Use `t.Fatal` only when the test can't continue.

<div align="right"><a href="#top">⬆ back to top</a></div>

---

<a id="uber"></a>

## 3. Uber Go Style Guide

![Uber](https://img.shields.io/badge/Uber-Go_Style_Guide-000000?style=for-the-badge&logo=uber&logoColor=white)

**Source:** [github.com/uber-go/guide](https://github.com/uber-go/guide/blob/master/style.md)

Uber follows all of Part 1. This part shows what it adds.

### 🏷️ Naming

```go
// Error values: Err (public) or err (private)
var ErrNotFound = errors.New("not found")
var errClosed = errors.New("closed")

// Error types end with Error
type NotFoundError struct{ ID string }

// Private package-level variables start with _ (Google doesn't do this)
var _defaultPort = 8080

// Printf-style functions end with f
func Logf(format string, args ...any)
```

### 🚨 Errors

```go
// Short context, no "failed to"
return fmt.Errorf("get user %q: %w", id, err)          // ✅
return fmt.Errorf("failed to get user: %w", err)       // ❌

// Handle an error once: log it OR return it
if err != nil {
    log.Printf("get user: %v", err)  // ❌ logged here...
    return err                       // ...and again by the caller
}
```

Pick the right kind:

```go
// Caller doesn't check it → errors.New or fmt.Errorf
return fmt.Errorf("user %q has no email", id)

// Caller checks it, fixed message → sentinel value
if errors.Is(err, ErrNotFound) { ... }

// Caller checks it and needs details → custom type
var nf *NotFoundError
if errors.As(err, &nf) { fmt.Println(nf.ID) }
```

```go
// Type assertions: always comma-ok
s := i.(string)       // ❌ panics if it's not a string
s, ok := i.(string)   // ✅

// Exit only in main, and only once
func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}
```

### 🧱 Structs

```go
// Use field names, leave out zero values
u := User{Name: "Asif", Age: 30}  // ✅
u := User{"Asif", 30, ""}         // ❌

p := &Point{}       // ✅ instead of new(Point)
var cfg Config      // ✅ all zero

// JSON fields always get a tag
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Check at compile time that a type implements an interface
var _ http.Handler = (*Handler)(nil)
```

Don't embed types in exported structs. Their methods leak into your public API.

### 🍰 Slices and maps

```go
// nil is a valid slice
return nil                 // ✅ instead of []User{}
if len(s) == 0 { ... }     // ✅ instead of s == nil

// Give a capacity when you know the size
users := make([]User, 0, len(ids))
seen := make(map[string]bool, len(ids))

// Copy at the edges, so callers can't change your data later
func (d *Driver) SetTrips(trips []Trip) {
    d.trips = make([]Trip, len(trips))
    copy(d.trips, trips)
}
```

### 🔀 Concurrency

```go
// Mutex: a named field, never embedded, never a pointer
type Cache struct {
    mu   sync.Mutex  // ✅ zero value is ready to use
    data map[string]string
}

// Channels: size 0 or 1
ch := make(chan int)        // ✅
ch := make(chan int, 1)     // ✅
ch := make(chan int, 64)    // ❌ needs a good reason

// Typed atomics (Go 1.19+)
var count atomic.Int64
count.Add(1)
```

**No fire-and-forget goroutines.** Every goroutine needs a way to stop, and someone waits for it:

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    for {
        select {
        case <-ticker.C:
            flush()
        case <-ctx.Done():
            return       // ✅ stops on cancel
        }
    }
}()

cancel()   // on shutdown
wg.Wait()  // wait for it to finish
```

No goroutines in `init()`. Catch leaks in tests with [`goleak`](https://github.com/uber-go/goleak).

### 🧭 Functions and code shape

```go
// Functional options for optional settings
type Option func(*Server)

func WithPort(p int) Option {
    return func(s *Server) { s.port = p }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080}  // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}

s := NewServer(WithPort(9000))
```

```go
// No bare true/false arguments
printInfo("foo", true, true)                          // ❌
printInfo("foo", true /* isLocal */, true /* done */) // ✅

// Use time types, not ints
func Poll(delay int)            // ❌ seconds? milliseconds?
func Poll(delay time.Duration)  // ✅ Poll(10 * time.Second)

// Enums start at 1, so zero means "not set"
const (
    ColorRed Color = iota + 1
    ColorGreen
)
```

**Order in a file:** `const` → `var` → types → `NewX()` → methods → plain functions.

**Also:** avoid mutable globals and `init()`, keep variable scope small, use raw strings to avoid escaping.

**Line length:** soft limit of **99 characters**.

### 🚀 Performance

Only in hot code. Write clear code first, then measure.

```go
strconv.Itoa(n)        // ✅ faster than fmt.Sprint(n)

b := []byte("hello")   // ✅ convert once, before the loop
for i := 0; i < n; i++ {
    w.Write(b)
}
```

### 🧪 Tests

Use table-driven tests, and keep the tables simple:

```go
func TestCount(t *testing.T) {
    tests := []struct {
        name string
        give string
        want int
    }{
        {name: "empty", give: "", want: 0},
        {name: "one word", give: "go", want: 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Count(tt.give); got != tt.want {
                t.Errorf("Count(%q) = %d, want %d", tt.give, got, tt.want)
            }
        })
    }
}
```

<div align="right"><a href="#top">⬆ back to top</a></div>

---

<a id="summary"></a>

## ⭐ 4. Summary: the rules that matter most

The most important and most popular rules from all three, in one place.

| # | Rule | Example | From |
| :---: | --- | --- | --- |
| 1 | Run `gofmt` on save | `go fmt ./...` | Go |
| 2 | MixedCaps, never snake_case | `userName`, `MaxSize` | Go |
| 3 | Acronyms in one case | `userID`, `ServeHTTP` | Go |
| 4 | Short package names, no repeating | `user.New()`, not `user.NewUser()` | Go |
| 5 | No `Get` in getters | `u.Name()` | Go |
| 6 | Receivers: 1–2 letters | `func (s *Server)` | Go |
| 7 | `error` last, always checked first | `if err != nil { return err }` | Go |
| 8 | Wrap errors with short context | `fmt.Errorf("get user: %w", err)` | Go · Google · Uber |
| 9 | Log or return an error, not both | `return fmt.Errorf(...)` | Uber |
| 10 | `ctx` first, never in a struct | `Run(ctx context.Context)` | Go |
| 11 | Every goroutine can stop | `case <-ctx.Done(): return` | Go · Uber |
| 12 | Small interfaces, defined where used | `type Reader interface{ Read(...) }` | Go |
| 13 | Pointer receivers when changing state | `func (c *Counter) Inc()` | Go |
| 14 | Field names in struct literals | `User{Name: "Asif"}` | Uber |
| 15 | Mutex as a named field | `mu sync.Mutex` | Uber |
| 16 | Table-driven tests | `for _, tt := range tests` | Google · Uber |
| 17 | Doc comments start with the name | `// Parse reads ...` | Go |

<div align="right"><a href="#top">⬆ back to top</a></div>

---

<a id="tools"></a>

## 🛠️ 5. Tools

```bash
gofmt -l .           # formatting
go vet ./...         # common mistakes
staticcheck ./...    # deeper checks, naming, unused code
golangci-lint run    # many linters at once
```

| Linter | Checks |
| --- | --- |
| `govet` | the same as `go vet` |
| `staticcheck` | bugs, simplifications, style |
| `errcheck` | ignored errors |
| `revive` | naming and comments |
| `goimports` | import grouping |
| `lll` | line length (default 120) |
| `funlen` | function length (default 60 lines) |

<div align="right"><a href="#top">⬆ back to top</a></div>
