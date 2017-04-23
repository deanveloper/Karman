package main

func main() {
    StartBotService()
    StartWebService()

    // c from bot.go
    defer pool.Close()

    // just keep it running until force closed
    <-make(chan int)
}