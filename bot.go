package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/bwmarrin/discordgo"
    "github.com/guregu/dynamo"
    "os"
)

type OurBot struct {
    table *dynamo.Table
}

type User struct {
    User string     `dynamo:"user"`
    Karma int       `dynamo:"karma"`
}

func New() *OurBot {
    return &OurBot{}
}

func (b *OurBot) Start() {

    // start DynamoDB session
    sess, err := session.NewSession()
    if err != nil {
        fmt.Println("Error connecting to DB:", err)
        return
    }
    temp := dynamo.New(sess).Table("Karma")
    b.table = &temp

    test := User{}
    err = b.table.Get("user", "test").One(&test)

    if err != nil {
        fmt.Println("Error getting test value from DB:", err)
        return
    }
    if test.Karma != 1337 {
        fmt.Println("Test test's Karma was not 1337! Instead was", test.Karma)
        return
    }

    fmt.Println("Successfully connected to DynamoDB!")

    dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
    if err != nil {
        fmt.Println("Error creating session!", err)
        return
    }

    // start discord stuff
    dg.AddHandler(b.ready)
    dg.AddHandler(b.guildCreate)
    dg.AddHandler(b.reactionAdd)
    dg.AddHandler(b.reactionRemove)
    dg.AddHandler(b.handleCommand)

    err = dg.Open()
    if err != nil {
        fmt.Println("Error starting websocket:", err)
    }
    fmt.Println("Successfully connected to Discord!")
}

func (b *OurBot) Close() {

}