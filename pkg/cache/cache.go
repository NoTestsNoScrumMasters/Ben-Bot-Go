package cache

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	mutex        = sync.RWMutex{}
	objects      = make(map[string]interface{})
	objectMeta   = make(map[string]int64)
	session      *discordgo.Session
	sessionMutex sync.RWMutex
	timeout      int64 = 15
)

type ItemGetter func(id string) (interface{}, error)

func GetOrRequest(id string, cb ItemGetter) (item interface{}, e error) {
	dirty := false

	// Retrieve if not existant yet
	mutex.RLock()
	ok := false
	item, ok = objects[id]
	mutex.RUnlock()
	if !ok {
		dirty = true
		item, e = cb(id)
	}

	// Check if there is a timeout
	mutex.RLock()
	meta := objectMeta[id]
	mutex.RUnlock()
	if time.Now().Unix()-meta > timeout {
		dirty = true
		item, e = cb(id)
	}

	// Update the entry
	if dirty {
		mutex.Lock()
		objects[id] = item
		mutex.Unlock()
	}

	// Return data
	return
}

func Guild(id string) (*discordgo.Guild, error) {
	ch, err := GetOrRequest(id, func(id string) (interface{}, error) {
		return GetSession().Guild(id)
	})

	return ch.(*discordgo.Guild), err
}

func Channel(id string) (*discordgo.Channel, error) {
	ch, err := GetOrRequest(id, func(id string) (interface{}, error) {
		return GetSession().Channel(id)
	})

	return ch.(*discordgo.Channel), err
}

func SetSession(s *discordgo.Session) {
	sessionMutex.Lock()
	session = s
	sessionMutex.Unlock()
}

func GetSession() *discordgo.Session {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()

	if session == nil {
		panic(errors.New("Tried to get discord session before cache#setSession() was called"))
	}

	return session
}
