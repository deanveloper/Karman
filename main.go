package main

import (
    "os"
    "os/signal"
    "syscall"
)

func main() {
    StartBotService()
    StartWebService()

    // c from bot.go
    defer pool.Close()

    // just keep it running until force closed
    sigChan := make(chan os.Signal)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // wait for SIGINT or SIGTERM
    <-sigChan

    // at this point, pool.Close() will be called (as it was deferred)
}
