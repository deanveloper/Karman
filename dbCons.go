package main

import (
    "github.com/bwmarrin/discordgo"
    "github.com/guregu/dynamo"
)

func plusOne(userId string) error {
    old := User{}
    err := table.Put(User{userId, 1}).OldValue(&old)

    // if no old value, we are done
    if err == dynamo.ErrNotFound {
        return nil
    }
    if err != nil {
        return err
    }

    err = table.Put(User{userId, old.Karma + 1}).Run()
    return err
}

func minusOne(userId string) error {
    old := User{}
    err := table.Put(User{userId, -1}).OldValue(&old)

    // if no old value, we are done
    if err == dynamo.ErrNotFound {
        return nil
    }
    if err != nil {
        return err
    }

    err = table.Put(User{userId, old.Karma - 1}).Run()
    return err
}

func getKarma(user *discordgo.User) (int, error) {
    resp := User{}
    err := table.Get("user", user.ID).One(&resp)

    if err == dynamo.ErrNotFound {
        return 0, nil
    }
    if err != nil {
        return 0, err
    }
    return resp.Karma, err
}
