// Once value: OnceValues and OnceFunc, without a separate Once variable.
//
// You give sync.OnceValues a function that returns a value and an error.
// It gives you back a new function. The first call runs yours. Every later
// call gets the same value and error, and yours doesn't run again. Here
// several goroutines ask for the config, but the file is read only once.
//
// sync.OnceFunc does the same for a function that returns nothing. Here it
// prints a warning about an old setting only once, even when many workers
// call it.
//
// Note: compare with 01-lazy-init. There, you needed a Once and a separate
// variable to hold the pool. OnceValues keeps the value for you.
//
// Real use:
//   - loading a config file or reading secrets once, the first time
//     they're needed: var loadConfig = sync.OnceValues(readConfig).
//   - logging "this setting will go away soon" once, instead of on every
//     request.
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

// readConfig pretends to read a config file and turn it into an appConfig.
func readConfig() (appConfig, error) {
	fmt.Println("reading config.json...")
	time.Sleep(readTime)
	return appConfig{port: 8080, logLevel: "info"}, nil
}

// printOldSettingWarning tells the user about a setting that will go away.
func printOldSettingWarning() {
	fmt.Println(`warning: "log_file" is deprecated, use "log_dir" instead`)
}
