package main

import (
	"fmt"
	"sync"

	pubnub "github.com/pubnub/go"
)

var pn *pubnub.PubNub

func init() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = pubnub.NewPubNub(config)
}

func pubnubCopy() *pubnub.PubNub {
	_pn := new(pubnub.PubNub)
	*_pn = *pn
	return _pn
}

func gettingStarted() {
	listener := pubnub.NewListener()
	doneConnect := make(chan bool)
	donePublish := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNDisconnectedCategory:
				case pubnub.PNConnectedCategory:
					doneConnect <- true
				case pubnub.PNReconnectedCategory:
				}
			case msg := <-listener.Message:
				fmt.Println(msg)
				donePublish <- true
			case <-listener.Presence:
				// handle presence
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"hello_world"},
	})

	<-doneConnect

	response, status, err := pn.Publish().
		Channel("hello_world").Message("Hello!").Execute()

	fmt.Println(response, status, err)

	<-donePublish
}

func listeners() {
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneSubscribe <- true
					return
				case pubnub.PNDisconnectedCategory:
					//
				case pubnub.PNReconnectedCategory:
					//
				case pubnub.PNAccessDeniedCategory:
					//
				case pubnub.PNUnknownCategory:
					//
				}
			case <-listener.Message:
				//
			case <-listener.Presence:
				//
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	<-doneSubscribe
}

func time() {
	res, status, err := pn.Time().Execute()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func publish() {
	res, status, err := pn.Publish().
		Channel("ch").
		Message("hey").
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func hereNow() {
	res, status, err := pn.HereNow().
		Channels([]string{"ch"}).
		IncludeUuids(true).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func presence() {
	// await both connected event on emitter and join presence event received
	var wg sync.WaitGroup
	wg.Add(2)

	donePresenceConnect := make(chan bool)
	doneJoin := make(chan bool)
	doneLeave := make(chan bool)
	errChan := make(chan string)
	ch := "my-channel"

	configPresenceListener := pubnub.NewConfig()
	configPresenceListener.SubscribeKey = "demo"
	configPresenceListener.PublishKey = "demo"

	pnPresenceListener := pubnub.NewPubNub(configPresenceListener)

	pn.Config.Uuid = "my-emitter"
	pnPresenceListener.Config.Uuid = "my-listener"

	listenerEmitter := pubnub.NewListener()
	listenerPresenceListener := pubnub.NewListener()

	// emitter
	go func() {
		for {
			select {
			case status := <-listenerEmitter.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					wg.Done()
					return
				}
			case <-listenerEmitter.Message:
				errChan <- "Got message while awaiting for a status event"
				return
			case <-listenerEmitter.Presence:
				errChan <- "Got presence while awaiting for a status event"
				return
			}
		}
	}()

	// listener
	go func() {
		for {
			select {
			case status := <-listenerPresenceListener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					donePresenceConnect <- true
				}
			case message := <-listenerPresenceListener.Message:
				errChan <- fmt.Sprintf("Unexpected message: %s",
					message.Message)
			case presence := <-listenerPresenceListener.Presence:
				fmt.Println(presence, "\n", configPresenceListener)
				// ignore join event of presence listener
				if presence.Uuid == configPresenceListener.Uuid {
					continue
				}

				if presence.Event == "leave" {
					doneLeave <- true
					return
				} else {
					wg.Done()
				}
			}
		}
	}()

	pn.AddListener(listenerEmitter)
	pnPresenceListener.AddListener(listenerPresenceListener)

	pnPresenceListener.Subscribe(&pubnub.SubscribeOperation{
		Channels:        []string{ch},
		PresenceEnabled: true,
	})

	select {
	case <-donePresenceConnect:
	case err := <-errChan:
		panic(err)
		return
	}

	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{ch},
	})

	go func() {
		wg.Wait()
		doneJoin <- true
	}()

	select {
	case <-doneJoin:
	case err := <-errChan:
		panic(err)
		return
	}

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{ch},
	})

	select {
	case <-doneLeave:
	case err := <-errChan:
		panic(err)
		return
	}
}

func history() {
	res, status, err := pn.History().
		Channel("ch").
		Count(2).
		IncludeTimetoken(true).
		Reverse(true).
		Start(int64(1)).
		End(int64(2)).
		Execute()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)
	fmt.Println(res)
}

func unsubscribe() {
	pn.Subscribe(&pubnub.SubscribeOperation{
		Channels: []string{"ch"},
	})

	// t.Sleep(3 * t.Second)

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		Channels: []string{"ch"},
	})
}

func main() {
	// gettingStarted()
	// listeners()
	// time()
	// publish()
	// hereNow()
	presence()
	// history()
	// unsubscribe()
}
