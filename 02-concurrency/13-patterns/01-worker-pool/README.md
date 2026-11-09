# Worker Pool

> 📝 **Planned.** No code yet. This is what I plan to add here.

A fixed number of workers take jobs from one channel and send results to another.

## Examples to add

- `01-basic`: workers, a jobs channel, and a results channel
- `02-with-context`: stop the pool early with a context

## Things to remember

- Close the jobs channel after the last job, so the workers can stop.
- **Close the results channel only after all workers are done.** Use a WaitGroup for that.
