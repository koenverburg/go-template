package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Task struct {
	ticker *time.Ticker
}

func (t *Task) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("task.run stopping")
			return
		case <-t.ticker.C:
			handle()
		}
	}
}

func handle() {
	for i := 0; i < 5; i++ {
		fmt.Print("#")
		time.Sleep(time.Millisecond * 200)
	}
	fmt.Println()
}

func main() {
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		fmt.Println("defer stopping...")

		cancel()
	}()

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got '%s' signal. Aborting...\n", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	task := &Task{
		ticker: time.NewTicker(time.Second * 2),
	}
	task.Run(ctx)
}
