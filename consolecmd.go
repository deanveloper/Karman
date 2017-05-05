package karman

import (
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
        }
    }
}