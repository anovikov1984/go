package main

import (
	"fmt"

	pubnub "github.com/pubnub/go"
)

var pn *pubnub.PubNub

func init() {
	config := pubnub.NewConfig()
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	pn = pubnub.NewPubNub(config)
}

func main() {
	channelGroup := "family"

	res, status, err := pn.ListChannelsInChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(res, status, err)

	resAdd, statusAdd, err := pn.AddChannelToChannelGroup().
		Channels([]string{"wife"}).
		Group(channelGroup).
		Execute()

	fmt.Println(resAdd, statusAdd, err)

	resAdd, statusAdd, err = pn.AddChannelToChannelGroup().
		Channels([]string{"son", "daughter"}).
		Group(channelGroup).
		Execute()

	fmt.Println(resAdd, statusAdd, err)

	listener := pubnub.NewListener()
	doneConnect := make(chan bool)

	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					doneConnect <- true
					return
				case pubnub.PNReconnectedCategory:
				}
			case <-listener.Message:
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{channelGroup},
		Timetoken:       int64(1337),
		PresenceEnabled: true,
	})

	<-doneConnect

	pn.Unsubscribe(&pubnub.UnsubscribeOperation{
		ChannelGroups: []string{channelGroup},
	})

	pn.Subscribe(&pubnub.SubscribeOperation{
		ChannelGroups:   []string{"cg1", "cg2"},
		Timetoken:       int64(1337),
		PresenceEnabled: true,
	})

	resRemove, statusRemove, err := pn.RemoveChannelFromChannelGroup().
		Channels([]string{"son"}).
		Group("family").
		Execute()

	fmt.Println(resRemove, statusRemove, err)

	res, status, err = pn.ListChannelsInChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(res, status, err)

	resDel, statusDel, err := pn.DeleteChannelGroup().
		ChannelGroup("family").
		Execute()

	fmt.Println(resDel, statusDel, err)
}
