package karman

import (
    "errors"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/garyburd/redigo/redis"
    "os"
    "time"
)

type OurBot struct {
    pool *redis.Pool
}

func New() *OurBot {
    return &OurBot{}
}

func (b *OurBot) Start() {
    dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
    if err != nil {
        fmt.Println("Error creating session!", err)
        return
    }

    dg.AddHandler(b.ready)
    dg.AddHandler(b.guildCreate)
    dg.AddHandler(b.reactionAdd)
    dg.AddHandler(b.reactionRemove)
    dg.AddHandler(b.handleCommand)

    b.pool = &redis.Pool{
        MaxIdle:   80,
        MaxActive: 5, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.DialURL(os.Getenv("REDIS_URL"))
            if err != nil {
                return nil, err
            }
            return c, nil
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            reply, err := redis.String(c.Do("PING"))

            if err != nil {
                return err
            }
            if reply != "PONG" {
                return errors.New("Response was not PONG")
            }

            return nil
        },
    }

    if err != nil {
        fmt.Println("Error connecting to redis:", err)
        return
    }

    err = dg.Open()
    if err != nil {
        fmt.Println("Error starting websocket:", err)
    }
}

func (b *OurBot) Close() {
    b.pool.Close()
}