package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/bwmarrin/discordgo"
    "github.com/guregu/dynamo"
    "io/ioutil"
    "os"
    "path"
)

var table *dynamo.Table

type User struct {
    User  string     `dynamo:"user"`
    Karma int       `dynamo:"karma"`
}

func Start() {
    fmt.Println("Starting Karman...")

    // start DynamoDB session
    sess, err := session.NewSession()
    if err != nil {
        fmt.Println("Error connecting to DB:", err)
        return
    }
    temp := dynamo.New(sess, aws.NewConfig().WithRegion("us-west-2")).Table("Karma")
    table = &temp

    test := User{}
    err = table.Get("user", "test").One(&test)

    if err != nil {
        fmt.Println("Error getting test value from DB:", err)
        return
    }
    if test.Karma != 1337 {
        fmt.Println("Test test's Karma was not 1337! Instead was", test.Karma)
        return
    }
    fmt.Println("Successfully connected to DynamoDB!")

    dat, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), "KARMAN_SECRET"))
    if err != nil {
        fmt.Println("Error reading secret key!", err)
        return
    }
    dg, err := discordgo.New("Bot " + string(dat))
    if err != nil {
        fmt.Println("Error creating session!", err)
        return
    }

    // start discord stuff
    dg.AddHandler(ready)
    dg.AddHandler(guildCreate)
    dg.AddHandler(reactionAdd)
    dg.AddHandler(reactionRemove)
    dg.AddHandler(handleCommand)

    err = dg.Open()
    if err != nil {
        fmt.Println("Error starting websocket:", err)
        return
    }
    fmt.Println("Successfully connected to Discord!")
}
