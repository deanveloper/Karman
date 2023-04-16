package karman

import (
    "github.com/guregu/dynamo"
)

func (b *Karman) plusOne(userId string) error {
    old := User{}
    err := b.table.Put(User{userId, 1}).OldValue(&old)

    // if no old value, we are done
    if err == dynamo.ErrNotFound {
        return nil
    }
    if err != nil {
        return err
    }

    err = b.table.Put(User{userId, old.Karma + 1}).Run()
    return err
}

func (b *Karman) minusOne(userId string) error {
    old := User{}
    err := b.table.Put(User{userId, -1}).OldValue(&old)

    // if no old value, we are done
    if err == dynamo.ErrNotFound {
        return nil
    }
    if err != nil {
        return err
    }

    err = b.table.Put(User{userId, old.Karma - 1}).Run()
    return err
}

func (b *Karman) getKarma(userId string) (int, error) {
    resp := User{}
    err := b.table.Get("user", userId).One(&resp)

    if err == dynamo.ErrNotFound {
        return 0, nil
    }
    if err != nil {
        return 0, err
    }
    return resp.Karma, err
}
