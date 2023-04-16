package karman

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "strconv"
    "strings"
)

func (b *Karman) Command(args []string) {
    if len(args) == 0 {
        b.log.Println("No arguments provided")
    } else {
        switch strings.ToLower(args[0]) {
        case "get":
            if len(args) == 1 {
                b.log.Println("Please provide a user id!")
            } else {
                _, err := strconv.ParseUint(args[1], 10, 64)
                if err != nil {
                    b.log.Println("Invalid user id: must be uint64")
                    return
                }
                karma, err := b.getKarma(args[1])
                if err != nil {
                    b.log.Println("Error:", err)
                    return
                }
                b.log.Println("Karma:", karma)
            }
        case "announce":
            msg := fmt.Sprintf("%s", args[1:])
            msg = "**Developer announcement:** " + msg[1:len(msg)-1] // add prefix, trim brackets
            for _, g := range b.dg.State.Guilds {
                go func(g *discordgo.Guild) {
                    for _, ch := range g.Channels {
                        _, err := b.dg.ChannelMessageSend(ch.ID, msg)
                        if err == nil {
                            break
                        }
                    }
                }(g)
            }
        case "list":
            b.log.Println("Guilds: ")
            for _, g := range b.dg.State.Guilds {
                b.log.Println(g.Name)
            }
        default:
            b.log.Println("Command unknown:", args[0])
        }
    }
}