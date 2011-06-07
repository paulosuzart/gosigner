package gosigner

import (
        "appengine"
        "appengine/datastore"
        "time"
        "os"
        "log"
)

type Key struct {
        Alias string
        Value string
        Owner string
        Created datastore.Time
}

func (self *Key) Save(c appengine.Context) (*datastore.Key, os.Error) {
        self.Created = datastore.SecondsToTime(time.Seconds())
        key, err := datastore.Put(c, datastore.NewIncompleteKey("key"), self)
        if err != nil {
                log.Printf("Error while saving key for %s", self.Owner)
                return nil, err
        }
        return key, nil
}

