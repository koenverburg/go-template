package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Task struct {
  ticker *time.Ticker
}

func (t *Task) Run() {
  for {
    select {
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
	fmt.Println("Hello World!")

  task := &Task{
    ticker: time.NewTicker(time.Second * 2),
  }

  task.Run()

  go gracefulShutdown()
  forever := make(chan int)
  <-forever

}

func gracefulShutdown() {
	// create a "returnCode" channel which will be the return code of the application
  returnCode := make(chan int)

	// finishUP channel signals the application to finish up
  finishUP := make(chan struct{})

	// done channel signals the signal handler that the application has completed
  done := make(chan struct{})
  listener := make(chan os.Signal, 1)

  signal.Notify(listener, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

  go func() {
    // wait for our os signal to stop the app
    // on the graceful stop channel
    // this goroutine will block until we get an OS signal
    sig := <- listener
    fmt.Printf("caught sig: %+v", sig)

    // send message on "finish up" channel to tell the app to
    // gracefully shutdown
    finishUP<-struct{}{}

    // wait for word back if we finished or not
    select {
    case <-time.After(30*time.Second):
      // timeout after 30 seconds waiting for app to finish,
      // our application should Exit(1)
      returnCode<-1
    case <-done:
      // if we got a message on done, we finished, so end app
      // our application should Exit(0)
      returnCode<-0
    }
  }()

  // ... Do business Logic in goroutines

  fmt.Println("waiting for finish")
  // wait for finishUP channel write to close the app down
  <-finishUP
  fmt.Println("stopping things, might take 2 seconds")

  fmt.Println("Do business Logic for shutdown simulated by Sleep 2 seconds")

  // ... Do business Logic for shutdown simulated by Sleep 2 seconds
  time.Sleep(2*time.Second)

  // write to the done channel to signal we are done.
  done <-struct{}{}
  os.Exit(<-returnCode)
}
