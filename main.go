package main

import (
    "fmt"
)

func main() {
    Start()

    fmt.Println("When done, press CTRL-C to exit.")
    <-make(chan struct{})
    return
}
