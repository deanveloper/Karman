package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/garyburd/redigo/redis"
    "os"
    "strings"
)

var pool *redis.Pool
var session *discordgo.Session

func StartBotService() {
    dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
    if err != nil {
        fmt.Println("Error creating session!", err)
        return
    }

    dg.AddHandler(ready)
    dg.AddHandler(guildCreate)
    dg.AddHandler(reactionAdd)
    dg.AddHandler(reactionRemove)
    dg.AddHandler(handleCommand)
    session = dg

    pool = &redis.Pool{
        MaxIdle: 80,
        MaxActive: 5, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.DialURL(os.Getenv("REDIS_URL"))
            if err != nil {
                return nil, err
            }
            return c, nil
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

func ready(s *discordgo.Session, ev *discordgo.Ready) {
    err := s.UpdateStatus(0, "Karma Counter")
    if err != nil {
        fmt.Println("Error while readying:", err)
    }
}

func guildCreate(s *discordgo.Session, ev *discordgo.GuildCreate) {
    _, err := s.Request("PATCH", discordgo.EndpointGuildMembers(ev.ID) + "/@me/nick", struct{nick string}{"Karman"})
    if err != nil {
        fmt.Println("Error while joining guild " + ev.Name + ":", err)
    }
}

func handleCommand(s *discordgo.Session, ev *discordgo.MessageCreate) {
    if strings.HasPrefix(strings.ToLower(ev.Content), "!karma") {
        mentions := ev.Mentions

        if len(mentions) < 2 {
            if len(mentions) == 0 {
                karma, err := getKarma(ev.Author)
                if err != nil {
                    fmt.Println("Error getting karma:", err)
                    s.ChannelMessageSend(ev.ChannelID, "Error getting karma: `" + err.Error() + "`")
                    return
                }

                s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("You have **%s** karma", karma))
            } else { // len is 1
                user := mentions[0]
                karma, err := getKarma(mentions[0])
                if err != nil {
                    fmt.Println("Error getting karma:", err)
                    s.ChannelMessageSend(ev.ChannelID, "Error getting karma: `" + err.Error() + "`")
                    return
                }

                s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("**%s** has **%s** karma", user.Username, karma))
            }

        } else {
            karmas, err := getKarmaMulti(mentions...)
            if err != nil {
                fmt.Println("Error getting karma:", err)
                s.ChannelMessageSend(ev.ChannelID, "Error getting karma: `" + err.Error() + "`")
                return
            }

            for user, karma := range karmas {
                s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("**%s** has **%s** karma", user.Username, karma))
            }
        }
    }
}

func getKarma(user *discordgo.User) (int, error) {
    c := pool.Get()
    defer c.Close()

    rawReply, err := pool.Get().Do("GET", user.ID)

    return redis.Int(rawReply, err)
}

func getKarmaMulti(users ... *discordgo.User) (map[*discordgo.User]int, error) {

    ids := make([]interface{}, len(users))
    for i, user := range users {
        ids[i] = user.ID
    }
    rawReply, err := pool.Get().Do("MGET", ids...)
    if err != nil {
        return nil, err
    }

    reply, err := redis.Ints(rawReply, err)

    karmas := make(map[*discordgo.User]int)
    for index, user := range users {
        karmas[user] = reply[index]
    }
    return karmas, nil
}

func reactionAdd(s *discordgo.Session, ev *discordgo.MessageReactionAdd) {
    if ev.Emoji.APIName() == "⬆" || ev.Emoji.APIName() == "⬇" { // up or down
        msg, err := s.ChannelMessage(ev.ChannelID, ev.MessageID)
        if err != nil {
            fmt.Println("Error getting message", ev.MessageID, "for channel", ev.ChannelID, err)
            return
        }

        if ev.Emoji.Name == "⬆" { // up
            plusOne(msg.Author.ID)
        } else if ev.Emoji.Name == "⬇" { // down
            minusOne(msg.Author.ID)
        }
    }
}

func reactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove) {
    if ev.Emoji.APIName() == "⬆" || ev.Emoji.APIName() == "⬇" { // up or down
        msg, err := s.ChannelMessage(ev.ChannelID, ev.MessageID)
        if err != nil {
            fmt.Println("Error getting message", ev.MessageID, "for channel", ev.ChannelID, err)
            return
        }

        if ev.Emoji.Name == "⬇" { // down
            plusOne(msg.Author.ID)
        } else if ev.Emoji.Name == "⬆" { // up
            minusOne(msg.Author.ID)
        }
    }
}

func plusOne(userId string) error {
    _, err := pool.Get().Do("INCR", userId)
    if err != nil {
        fmt.Println("Error incrementing karma:", err)
    }
    return err
}

func minusOne(userId string) error {
    _, err := pool.Get().Do("DECR", userId)
    if err != nil {
        fmt.Println("Error decrementing karma:", err)
    }
    return err
}

func onKill() {
    session.UserUpdateStatus(discordgo.StatusOffline)
}