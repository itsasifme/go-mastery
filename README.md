# Go Mastery

> A hands-on journey from Go fundamentals to advanced concepts through practical examples, experiments, and coding exercises.

I'm a backend engineer, and I use this repo to learn Go step by step. It starts with the basics, then moves to concurrency and the standard library. Every example is small, runs on its own, and shows one idea.

## What I want to learn

- Write clean and simple Go code.
- Understand pointers, types, methods, and interfaces.
- Get comfortable with goroutines, channels, and the `sync` tools.
- Know the standard library well.

## Topics

✅ = has examples · 📝 = planned

### 01 · Fundamentals

| # | Topic | Status |
| --- | --- | --- |
| 01 | [Variables](01-fundamentals/01-variables/) | 📝 |
| 02 | [Data types](01-fundamentals/02-data-types/) | 📝 |
| 03 | [Functions](01-fundamentals/03-functions/) | 📝 |
| 04 | [Pointers](01-fundamentals/04-pointers/) | ✅ |
| 05 | [Arrays and slices](01-fundamentals/05-arrays-slices/) | 📝 |
| 06 | [Maps](01-fundamentals/06-maps/) | 📝 |
| 07 | [Structs](01-fundamentals/07-structs/) | ✅ |
| 08 | [Methods](01-fundamentals/08-methods/) | ✅ |
| 09 | [Interfaces](01-fundamentals/09-interfaces/) | ✅ |
| 10 | [Error handling](01-fundamentals/10-error-handling/) | 📝 |
| 11 | [Generics](01-fundamentals/11-generics/) | 📝 |
| 12 | [Defer, panic, and recover](01-fundamentals/12-defer-panic-recover/) | 📝 |

### 02 · Concurrency

| # | Topic | Status |
| --- | --- | --- |
| 01 | [Goroutines](02-concurrency/01-goroutines/) | ✅ |
| 02 | [Channels](02-concurrency/02-channels/) | ✅ |
| 03 | [Select](02-concurrency/03-select/) | 📝 |
| 04 | [Context](02-concurrency/04-context/) | 📝 |
| 05 | [WaitGroup](02-concurrency/05-waitgroup/) | ✅ |
| 06 | [errgroup](02-concurrency/06-errgroup/) | 📝 |
| 07 | [Mutex](02-concurrency/07-mutex/) | 📝 |
| 08 | [RWMutex](02-concurrency/08-rwmutex/) | 📝 |
| 09 | [Atomic](02-concurrency/09-atomic/) | 📝 |
| 10 | [Once](02-concurrency/10-once/) | 📝 |
| 11 | [Semaphore](02-concurrency/11-semaphore/) | 📝 |
| 12 | [singleflight](02-concurrency/12-singleflight/) | 📝 |
| 13 | [Patterns](02-concurrency/13-patterns/) | 📝 |

### 03 · Standard library

| # | Topic | Status |
| --- | --- | --- |
| 01 | [strings](03-standard-library/01-strings/) | 📝 |
| 02 | [strconv](03-standard-library/02-strconv/) | 📝 |
| 03 | [time](03-standard-library/03-time/) | 📝 |
| 04 | [encoding/json](03-standard-library/04-json/) | 📝 |
| 05 | [io](03-standard-library/05-io/) | 📝 |
| 06 | [Files and folders](03-standard-library/06-filesystem/) | 📝 |
| 07 | [Logging](03-standard-library/07-logging/) | 📝 |
| 08 | [regexp](03-standard-library/08-regexp/) | 📝 |
| 09 | [net/http](03-standard-library/09-net-http/) | 📝 |
| 10 | [os/exec](03-standard-library/10-os-exec/) | 📝 |

## How the folders are numbered

Every folder starts with a number, so they show up in learning order. Numbers start at `01` inside every folder. Just go through them in order.

Every planned topic already has its folder, so adding code later won't change any numbers.

## Folder structure

```text
01-fundamentals/
├── 01-variables/
├── 02-data-types/
├── 03-functions/
├── 04-pointers/
├── 05-arrays-slices/
├── 06-maps/
├── 07-structs/
├── 08-methods/
├── 09-interfaces/
├── 10-error-handling/
├── 11-generics/
└── 12-defer-panic-recover/
02-concurrency/
├── 01-goroutines/
├── 02-channels/
│   ├── 01-unbuffered/
│   ├── 02-buffered/
│   ├── 03-directional/
│   └── 04-closing/
├── 03-select/
├── 04-context/
│   ├── 01-cancellation/
│   ├── 02-timeout/
│   └── 03-deadline/
├── 05-waitgroup/
├── 06-errgroup/
├── 07-mutex/
├── 08-rwmutex/
├── 09-atomic/
├── 10-once/
├── 11-semaphore/
├── 12-singleflight/
└── 13-patterns/
    ├── 01-worker-pool/
    ├── 02-fan-in/
    ├── 03-fan-out/
    ├── 04-pipeline/
    └── 05-backpressure/
03-standard-library/
├── 01-strings/
├── 02-strconv/
├── 03-time/
├── 04-json/
├── 05-io/
├── 06-filesystem/
├── 07-logging/
├── 08-regexp/
├── 09-net-http/
└── 10-os-exec/
```

## Run an example

Each example has its own folder with a `main.go`. Run it from the repo root:

```bash
go run ./01-fundamentals/09-interfaces/01-shapes
go run ./02-concurrency/02-channels/02-buffered/01-basic
```

## Checks

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

There are no tests yet. These commands still check that all the code builds and is formatted. GitHub Actions runs them on every push.

`go test -race` needs cgo and a C compiler. On Windows without `gcc` it fails with `-race requires cgo`. On Linux (and in GitHub Actions) it works.

Some examples print things in a different order each run, or leave a goroutine stuck on purpose. None of them crash or hang, so the checks still pass.

## Go version

This repo uses Go 1.23 (set in `go.mod`). Older versions may not work. When a topic uses a newer Go feature, its README says which version added it.

There are no outside packages yet, so there's no `go.sum`. The errgroup, singleflight, and semaphore examples will add `golang.org/x/sync` later.

## Notes

- One example per folder. Never put two `main` functions in the same folder.
- `_playground/` folders are for my own experiments. Git ignores them, and Go tools skip them because the name starts with `_`.
