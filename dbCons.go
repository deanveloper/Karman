package karman

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/garyburd/redigo/redis"
)

func (b *OurBot) plusOne(userId string) error {
    c := b.pool.Get()

    defer func() {
        err := c.Close()
        if err != nil {
            fmt.Println("Error closing connection for plusOne(" + userId + ")")
            fmt.Println(err)
        }
    }()

    _, err := c.Do("INCR", userId)
    if err != nil {
        fmt.Println("Error incrementing karma:", err)
    }
    return err
}

func (b *OurBot) minusOne(userId string) error {
    c := b.pool.Get()

    defer func() {
        err := c.Close()
        if err != nil {
            fmt.Println("Error closing connection for getKarma(" + userId + ")")
            fmt.Println(err)
        }
    }()

    _, err := c.Do("DECR", userId)
    if err != nil {
        fmt.Println("Error decrementing karma:", err)
    }
    return err
}

func (b *OurBot) getKarma(user *discordgo.User) (int, error) {
    c := b.pool.Get()

    defer func() {
        err := c.Close()
        if err != nil {
            fmt.Println("Error closing connection for getKarma(" + user.Username + ")")
            fmt.Println(err)
        }
    }()

    reply, err := redis.Int(c.Do("GET", user.ID))

    if err == redis.ErrNil {
        return 0, nil
    }
    return reply, err
}

func (b *OurBot) getKarmaMulti(users ... *discordgo.User) (map[*discordgo.User]int, error) {
    c := b.pool.Get()

    defer func() {
        err := c.Close()
        if err != nil {
            fmt.Printf("Error closing connection for getKarmaMulti(%q)", users)
            fmt.Println(err)
        }
    }()

    ids := make([]interface{}, len(users))
    for i, user := range users {
        ids[i] = user.ID
    }
    reply, err := redis.Ints(c.Do("MGET", ids...))
    karmas := make(map[*discordgo.User]int)
    if err == redis.ErrNil {
        return karmas, nil
    }
    if err != nil {
        return nil, err
    }

    for index, user := range users {
        karmas[user] = reply[index]
    }
    return karmas, nil
}