# CLAUDE.md

Rules for working on this repo. The goal: learn Go from basic to advanced, with code that looks like real production Go.

## Repo layout

- One Go module (`github.com/itsasifme/go-mastery`, Go 1.23) at the root. No third-party packages unless there's no other way.
- Three parts: `01-fundamentals/`, `02-concurrency/`, `03-standard-library/`.
- Every folder starts with a 2-digit number in learning order. Numbers start at `01` in every folder and have no gaps. Planned topics already have a folder, so numbers never change.
- Each example is its own folder with one `main.go`, run with `go run ./<path>`.
- `_playground/` is private and gitignored. Never touch it.

## Writing style

- Plain, easy English everywhere: comments, READMEs, commit messages. Short sentences, common words. It should sound like a person wrote it.
- Don't use words like "leverage", "robust", "deliberately", "seamlessly", or long sentences with many parts.

## Go code

- Run `gofmt` (or `goimports`). `go vet ./...` and `go build ./...` must pass.
- Keep the same order in every `main.go`, top to bottom:
  1. header comment, `package main`, imports
  2. package-level `type` (with its methods right after it), `const`, `var`
  3. `func main`
  4. the other functions, in the order `main` uses them

  Never put a function or a `const` block above or below `main` out of this order. Every function has a doc comment.
- Goroutines are **named functions** that take their channels as parameters, with a direction where it fits:

  ```go
  go produce(ch)                     // good
  func produce(out chan<- int) {...}

  go func() { ch <- 1 }()            // avoid: hides how the channel is passed
  ```

  Anonymous functions are only allowed in these three cases:
  - the closure is the lesson (like the loop variable example)
  - `defer func() { recover() }()`
  - `defer func() { done <- struct{}{} }()` when several goroutines share one `done`

- A function name says exactly what it does. `printBoth` prints two channels. `merge` would sound like fan-in, which returns a channel. Other good names: `readUntilDone`, `stopAfter`, `handleUrgentFirst`.
- Examples under `02-concurrency/` run their work in a goroutine, even when the lesson could be shown in `main` alone. That's what the topic is about.
- A goroutine that gets stopped (by `stop`, `ctx`, a timeout) must be **doing work until then**, like real code. Don't let it just sit on `<-ctx.Done()`. Loop over small steps and watch for the stop signal in the same `select`, and show how far it got when it stops:

  ```go
  for {
  	select {
  	case <-time.After(50 * time.Millisecond): // simulate one step of work
  		steps++
  	case <-ctx.Done():
  		fmt.Println("stopped after", steps, "steps:", ctx.Err())
  		return
  	}
  }
  ```

- How to name the goroutine function:
  - **Simple demo work** (sleep, send a number, print a step): call it `worker`.
  - **A specific job the example is about** (a database query, an export, a heartbeat, a fetch): name it for that job, like `query`, `export`, `heartbeat`, `fetch`.
  - Don't call it `send` or `run`. Those say nothing about the job.

  ```go
  go worker(ch, done)          // good: simple demo work, main does the receiving
  go query(ctx, done)          // good: the example is about a slow database call
  go send(ch, done)            // avoid
  go run(ctx, "big", d, done)  // avoid
  ```

  Keep send/receive names only when both sides are in the same file (`sendData` and `readData` in 01-bidirectional), or when the direction itself is the lesson (the directional channel examples).

- Loop over a number with `range`, not a C-style loop:

  ```go
  for i := range 3 { ... }        // not: for i := 0; i < 3; i++
  for range workerCount { <-done }
  ```

  Values start at 0. If you need 1 to n, write `i + 1`, and make the comments match.

- **No magic numbers, anywhere.** Every count, buffer size, duration, loop limit and rounding step gets a named constant (`const workerCount = 3`, `stepTime = 100 * time.Millisecond`, `bufferSize = 2`), or comes from the code (`cap(ch)`, `bufferSize + 1`). This includes `time.Sleep`, `time.After`, `WithTimeout`, `Round` and `make(chan T, n)`. Better still, don't count at all when you can stop on a signal, like reading until `ok` is false.
  - Add a short comment when the name alone doesn't say why, like `// longer than timeout, so the timeout wins` or "main and pong must use the same number, or one of them waits forever".
  - Not magic: sample data (values sent on a channel like `ch <- 5`, IDs like `"req-42"`, struct field values in fundamentals), counters starting at `0` or `1`, and a buffer of `1` for "room for one result" when the header says why.
- Use `time.Duration` for any amount of time, and always write the unit (`200 * time.Millisecond`).

### Signal channels

- Signal channels are `chan struct{}` and use the same names everywhere:
  - `stop`: main tells a goroutine to stop. Close it to reach everyone.
  - `done`: a goroutine tells main it finished.
- Label each one with its direction where it's made:

  ```go
  stop := make(chan struct{}) // main -> workers: "please stop"
  done := make(chan struct{}) // workers -> main: "I'm finished"
  ```

- Report `done` in a `defer` on the first line of the goroutine, so it always runs:
  - one goroutine reports: `defer close(done)`
  - several share `done`: `defer func() { done <- struct{}{} }() // send, not close: others share done`
- main waits for every goroutine it starts (`<-done`, one per goroutine). Then a stuck goroutine shows up as a deadlock, instead of a silent leak. Say so in a comment, like `<-done // proves fetch stopped too`.

## Comments

- A comment doesn't repeat a number that can change. Put the number in the code once and describe it in words:

  ```go
  // ❌ The other 2 workers wait forever.    // 3 = main + 2 stuck
  // ✅ The other workers wait forever.      // main + (workerCount - 1) stuck
  ```

  The exception is when the numbers are the lesson and won't change, like a timeline of the output ("0ms sent 1 ... 200ms got 1").

- Check every comment against the real output. If the code changes (like `range 3` starting at 0), update the comments too.

## Header comment

Every `main.go` starts with this comment. `Real use` and `Run` are always there. Leave out other parts that don't apply.

```go
// <Title>: <what it shows in one line>.
//
// <A few short lines on what happens and why.>
//
// Problem: <what goes wrong>.          (problem examples only)
//
// Fix: <how to fix it>. Uncomment the "fix:" line.
//
// Try: <a small change to see something new>.
//
// Note: <a related tip or link to another example>.
//
// Real use:
//   - <where this shows up in real programs>
//
// Run: go run ./<path>
package main
```

- Problem examples ship **showing the problem**. The fix stays commented out with a `// fix:` marker on the line. If lines must be removed, mark them `// problem`.
- **Real use** gives concrete production cases that name real APIs: `r.Context()`, `db.QueryContext`, `signal.NotifyContext`, `json.Unmarshal(data, &user)`. Don't write vague bullets like "useful in servers".
- When an example builds on an earlier one, say what's different. For example: "Compare with select/04-timeout: there, main stopped waiting, but the work kept running." A small side-by-side table is fine too.
- Links to other examples:
  - another topic: the topic name without its number, then the example folder with its number, like `(see channels/02-closing/02-broadcast)` or `(see waitgroup/01-basic)`. Never `(see 05-waitgroup)` or `(see 02-closing)`.
  - the same topic: just the folder, like `(see 05-leak)`.
  - always point to a specific example, not a whole topic.

## Topic README

Each topic README has these sections, in this order:

1. Title and one or two lines on what it is. A sub-topic README adds `See the [full API overview](../README.md#api).` under it.
2. `## Why`: **Problem:** and **Solution:**, each in its own paragraph (with a blank line between them).
3. `## API`: a table of the syntax or functions, with a link to the example that shows each one (`none` if no example shows it).
   - Syntax tables (channels, select) have three columns: Syntax, What it does, Example.
   - Package function tables add `Returns` and `Go` (the version it was added in). Show newer APIs too, even ones this repo can't use yet: `wg.Go(f) | 1.25 | ... | none (this repo uses Go 1.23)`.
4. `## Examples`: a table with Example, Run, and What it shows.
5. `## Things to remember`: short bullets. Bold the most important part of each.
6. `## Go version notes` if something changed in a Go version.

When a topic gets code:

- remove the `> 📝 Planned` line
- replace "Examples to add" with the real Examples table
- change its row in the root `README.md` from 📝 to ✅

## Before you say it's done

- `gofmt -l .` prints nothing, and `go vet ./...` and `go build ./...` pass.
- Run every example you touched. Problem examples should fail the way the header says (panic, deadlock, leak).
- Try the `fix:` and `Try:` lines on a copy and check they do what the comment says.
- Check that the README links work.

## Git

- Commits are authored only by the repo owner. **Never add a `Co-Authored-By` line or any Claude attribution.**
- One commit per file or change.
- Don't push. The owner pushes.
- Format: `<emoji> <Type>: <description>`. Use an imperative, lowercase verb after the colon, no period at the end, and 72 characters at most. Prefer the subject only.
- Types: 🎉 Init, ✏️ Add, ✨ Feature, 🆙 Update, 🐛 Fix, 🔨 Refactor, 🔥 Remove, 📝 Docs, 🎨 Style, ⚡ Perf, 🧪 Test, 🔒 Security, 📦 Build, 🚀 CI, 🧹 Chore, ⏪ Revert.
- The subject must make sense on its own. Name the topic in words, and say what the change teaches or why.
- **Never use folder numbers** in a commit message. The same numbers repeat in every topic.

  ```text
  ❌ ✏️ Add: write drain after close example
  ❌ 🔨 Refactor: move range example to 02-closing
  ✅ ✏️ Add: buffered channel example showing values survive close
  ✅ ✏️ Add: select example making urgent jobs go before normal ones
  ```
