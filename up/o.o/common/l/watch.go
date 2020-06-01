package l

import (
	"fmt"
	"sync"

	"go.uber.org/atomic"
)

// TODO(vu): implement a watcher system instead

var initialized atomic.Bool
var channels sync.Map

type Messenger interface {
	SendMessage(msg string)
}

func RegisterChannels(chans map[string]Messenger) {
	if initialized.Load() {
		panic("already init")
	}
	initialized.Store(true)
	prevChans := make([]*emptyMessenger, 0, 16)
	for name, messenger := range chans {
		ch := getChannel(name)
		m := (*ch).(*emptyMessenger)
		if len(m.msgs) > 0 {
			prevChans = append(prevChans, m)
		}
		*ch = messenger
	}
	channels.Range(func(key, _ interface{}) bool {
		name := key.(string)
		if chans[name] == nil {
			panic(fmt.Sprintf("channel %v was not registered", name))
		}
		return true
	})

	// now send the old messages
	for _, m := range prevChans {
		ch := getChannel(m.name)
		for _, msg := range m.msgs {
			(*ch).SendMessage(msg)
		}
	}
}

// getChannel allocates a new slot for registering messenger later
func getChannel(channel string) *Messenger {
	item, _ := channels.Load(channel)
	ch, _ := item.(*Messenger)
	if ch == nil {
		ch = new(Messenger)
		*ch = &emptyMessenger{name: channel}
		channels.Store(channel, ch)
	}
	return ch
}

type emptyMessenger struct {
	name string
	msgs []string
}

func (m *emptyMessenger) SendMessage(msg string) {
	if initialized.Load() {
		errMsg := fmt.Sprintf("%v (message sent on channel %v without registering a messenger)", msg, m.name)
		sendMessageOnDefaultChannel(errMsg)
		return
	}
	if len(m.msgs) < 16 {
		m.msgs = append(m.msgs, msg)
		return
	}
	errMsg := fmt.Sprintf("%v (too many messages sent on channel %v without registering a messenger)", msg, m.name)
	sendMessageOnDefaultChannel(errMsg)
}

func sendMessageOnDefaultChannel(msg string) {
	defaultCh := getChannel("")
	if _, ok := (*defaultCh).(*emptyMessenger); !ok {
		(*defaultCh).SendMessage(msg)
	}
}

type MockMessenger struct {
	Name string
}

func (m MockMessenger) SendMessage(msg string) {
	ll.Info("[" + m.Name + "] " + msg)
}
