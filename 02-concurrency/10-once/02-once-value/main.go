// Once value: OnceValues and OnceFunc, without a separate Once variable.
//
// sync.OnceValues wraps a function that returns a value and an error.
// The first call runs it. Every later call returns the same value and
// error, without running it again. Here several goroutines read the
// config, but the file is read only once.
//
// sync.OnceFunc wraps a function with no result. Here it prints a warning
// about an old setting only once, no matter how many requests hit it.
//
// Note: compare with 01-lazy-init. There, a Once and a package variable
// had to be kept in sync by hand. OnceValues keeps the value for you.
//
// Real use:
//   - loading a config file or reading secrets once, the first time
//     they're needed: var loadConfig = sync.OnceValues(readConfig).
//   - logging a deprecation warning once, instead of on every request.
//
// Run: go run ./02-concurrency/10-once/02-once-value
package main

import (
	"fmt"
	"sync"
	"time"
)

// appConfig is what readConfig returns.
type appConfig struct {
	port     int
	logLevel string
}

// readTime is how long reading the fake config file takes.
const readTime = 100 * time.Millisecond

var (
	loadConfig     = sync.OnceValues(readConfig)
	warnOldSetting = sync.OnceFunc(printOldSettingWarning)
)

func main() {
	const workerCount = 4

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(id+1, &wg)
	}
	wg.Wait()
}

// worker loads the config and uses it.
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("worker", id, "config error:", err)
		return
	}
	warnOldSetting()
	fmt.Printf("worker %d: port %d, log level %s\n", id, cfg.port, cfg.logLevel)
}

// readConfig pretends to read and parse a config file.
func readConfig() (appConfig, error) {
	fmt.Println("reading config.json...")
	time.Sleep(readTime)
	return appConfig{port: 8080, logLevel: "info"}, nil
}

// printOldSettingWarning tells the user about a setting that will go away.
func printOldSettingWarning() {
	fmt.Println(`warning: "log_file" is deprecated, use "log_dir" instead`)
}
