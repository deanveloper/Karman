package karman

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/bwmarrin/discordgo"
    "github.com/guregu/dynamo"
    "io/ioutil"
    "log"
    "os"
    "path"
)

type User struct {
    User  string     `dynamo:"user"`
    Karma int       `dynamo:"karma"`
}

type Karman struct {
    con   *dynamo.DB
    table *dynamo.Table
    dg    *discordgo.Session
    log   *log.Logger
}

func New(log *log.Logger) *Karman {
    return &Karman{log: log}
}

func (b *Karman) Start() {
    b.log.Println("Starting Karman...")

    // start DynamoDB session
    sess, err := session.NewSession()
    if err != nil {
        b.log.Println("Error connecting to DB:", err)
        return
    }
    b.con = dynamo.New(sess, aws.NewConfig().WithRegion("us-west-2"))
    temp := b.con.Table("Karma")
    b.table = &temp

    test := User{}
    err = b.table.Get("user", "test").One(&test)

    if err != nil {
        b.log.Println("Error getting test value from DB:", err)
        return
    }
    if test.Karma != 1337 {
        b.log.Println("Test test's Karma was not 1337! Instead was", test.Karma)
        return
    }
    b.log.Println("Successfully connected to DynamoDB!")

    dat, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), "botfiles", "karman", "SECRET"))
    if err != nil {
        b.log.Println("Error reading secret key!", err)
        return
    }
    dg, err := discordgo.New("Bot " + string(dat))
    if err != nil {
        b.log.Println("Error creating session!", err)
        return
    }

    // start discord stuff
    dg.AddHandler(b.ready)
    dg.AddHandler(b.reactionAdd)
    dg.AddHandler(b.reactionRemove)
    dg.AddHandler(b.handleCommand)

    err = dg.Open()
    if err != nil {
        b.log.Println("Error starting websocket:", err)
        return
    }
    b.dg = dg
    b.log.Println("Successfully connected to Discord!")
}

func (b *Karman) Stop() {
    // DB Connection automatically closes?
    b.dg.Close()
}
