package main

import (
    "fmt"
)

func main() {
    Start()

    fmt.Println("Press CTRL-C to exit.")
    <-make(chan struct{})
    return
}
