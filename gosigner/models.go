package gosigner

import (
	"appengine"
	"appengine/datastore"
	"time"
	"os"
	"log"
	"appengine/user"
)

type Key struct {
	Alias   string
	Value   string
	Owner   string
	Created datastore.Time
}

func (self *Key) Save(c appengine.Context) (*datastore.Key, os.Error) {
	self.Created = datastore.SecondsToTime(time.Seconds())
	self.Owner = user.Current(c).Email
	if self.Owner == "" {
		self.Owner = "paulosuzart@gmail.com"
	}
	key, err := datastore.Put(c, datastore.NewIncompleteKey("key"), self)
	if err != nil {
		log.Printf("Error while saving key for %s", self.Owner)
		return nil, err
	}
	return key, nil
}

func iterateKeys(it *datastore.Iterator) []Key {

	var keys = []Key{}
	for i := it; ; {
		var key Key
		_, err := i.Next(&key)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return []Key{}
		}
		keys = append(keys, key)
	}
	return keys

}
func AllKeys(c appengine.Context, owner string) []Key {
	q := datastore.NewQuery("key")
	q.Filter("Owner=", owner).Order("Alias")
	return iterateKeys(q.Run(c))
}
