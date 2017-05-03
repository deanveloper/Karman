package karman

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "strings"
)

func (b *Karman) ready(s *discordgo.Session, ev *discordgo.Ready) {
    err := s.UpdateStatus(0, "Karma Counter")
    if err != nil {
        b.log.Println("Error while readying:", err)
    } else {
        b.log.Println("I'm ready to count some karma!")
    }
}

func (b *Karman) guildCreate(s *discordgo.Session, ev *discordgo.GuildCreate) {
    _, err := s.Request("PATCH", discordgo.EndpointGuildMembers(ev.ID)+"/@me/nick", struct{ nick string }{"Karman"})
    if err != nil {
        b.log.Println("Error while joining guild "+ev.Name+":", err)
    }
}

func (b *Karman) handleCommand(s *discordgo.Session, ev *discordgo.MessageCreate) {
    if strings.HasPrefix(strings.ToLower(ev.Content), "!karma") {
        if ev.MentionEveryone {
            s.ChannelMessageSend(ev.ChannelID, "Getting everyone's karma is not allowed.")
            return
        }

        mentions := ev.Mentions

        if len(mentions) < 2 {
            if len(mentions) == 0 { // if someone was mentioned
                karma, err := b.getKarma(ev.Author)
                if err != nil {
                    b.log.Println("Error getting karma:", err)
                    s.ChannelMessageSend(ev.ChannelID, "Error getting karma: `"+err.Error()+"`")
                    return
                }

                s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("You have **%d** karma", karma))

            } else {
                user := mentions[0]
                karma, err := b.getKarma(mentions[0])
                if err != nil {
                    b.log.Println("Error getting karma:", err)
                    s.ChannelMessageSend(ev.ChannelID, "Error getting karma: `"+err.Error()+"`")
                    return
                }

                s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("**%s** has **%d** karma", user.Username, karma))
            }

        } else { // if multiple people were mentioned
            for _, user := range mentions {
                // get each one asynchronously
                go func(user *discordgo.User) {
                    karma, err := b.getKarma(user)

                    if err != nil {
                        b.log.Println("Error getting karma for", user.Username, ":", err)
                        s.ChannelMessageSend(ev.ChannelID, "Error getting karma for "+user.Username+": `"+err.Error()+"`")
                        return
                    }

                    s.ChannelMessageSend(ev.ChannelID, fmt.Sprintf("**%s** has **%d** karma", user.Username, karma))
                }(user)
            }
        }
    }
}

func (b *Karman) reactionAdd(s *discordgo.Session, ev *discordgo.MessageReactionAdd) {
    if ev.Emoji.APIName() == "⬆" || ev.Emoji.APIName() == "⬇" { // up or down
        msg, err := s.ChannelMessage(ev.ChannelID, ev.MessageID)
        if err != nil {
            fmt.Println("Error getting message", ev.MessageID, "for channel", ev.ChannelID, err)
            return
        }

        if ev.Emoji.Name == "⬆" { // up
            err = b.plusOne(msg.Author.ID)
        } else if ev.Emoji.Name == "⬇" { // down
            err = b.minusOne(msg.Author.ID)
        }
        if err != nil {
            fmt.Println("Error changing karma for", msg.Author.Username, ":", err)
            return
        }
    }
}

func (b *Karman) reactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove) {
    if ev.Emoji.APIName() == "⬆" || ev.Emoji.APIName() == "⬇" { // up or down
        msg, err := s.ChannelMessage(ev.ChannelID, ev.MessageID)
        if err != nil {
            fmt.Println("Error getting message", ev.MessageID, "for channel", ev.ChannelID, err)
            return
        }

        if ev.Emoji.Name == "⬇" { // down
            err = b.plusOne(msg.Author.ID)
        } else if ev.Emoji.Name == "⬆" { // up
            err = b.minusOne(msg.Author.ID)
        }
        if err != nil {
            fmt.Println("Error changing karma for", msg.Author.Username, ":", err)
            return
        }
    }
}
