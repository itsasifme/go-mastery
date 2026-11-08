# Channels

A channel lets goroutines send values to each other. It also makes them wait for each other.

1. [Directional channels](01-directional/): how to pass a channel to a function as send-only or receive-only
2. [Closing channels](02-closing/): how to say "no more values", and how to signal many goroutines at once
3. [Unbuffered channels](03-unbuffered/): a send waits for a receiver
4. [Buffered channels](04-buffered/): a send only waits when the buffer is full

Directional channels and closing come first, because every example after them passes channels as `chan<- int` / `<-chan int` and uses `close`.
