package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bluenviron/gomavlib/v2"
)

const ConnectionInactiveDuration = 1 * time.Second
const channelAddrA = "127.0.0.1:14557"
const channelAddrB = "127.0.0.1:14558"
const downLinkAddr = "127.0.0.1:14560"

func main() {

	endpointA := gomavlib.EndpointUDPServer{Address: channelAddrA}
	endpointB := gomavlib.EndpointUDPServer{Address: channelAddrB}
	endpointDown := gomavlib.EndpointUDPClient{Address: downLinkAddr}

	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			endpointA,
			endpointB,
			endpointDown,
		},
		Dialect:     nil,
		OutVersion:  gomavlib.V2,
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	var channelA *gomavlib.Channel
	var channelB *gomavlib.Channel
	var downChannel *gomavlib.Channel

	// wait till number of channels is 3
	for {
		if len(node.Channels) == 3 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	// set channels
	for channel := range node.Channels {
		// address without protocol tag (assumed UDP)
		splitTmpAddr := strings.Split(channel.String(), ":")
		fmtTmpAddr := splitTmpAddr[1] + ":" + splitTmpAddr[2]

		if fmtTmpAddr != downLinkAddr {
			if channelA == nil {
				channelA = channel
			} else if channelB == nil {
				channelB = channel
			}
		} else {
			downChannel = channel
		}
	}

	fmt.Printf("channels set!")

	// Initialize the active channel
	activeChannel := channelA
	lastActiveChannelMessageTime := time.Now()
	fmt.Printf("Active channel is %+v\n", activeChannel)

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			switch frm.Channel {

			case activeChannel:
				node.WriteFrameTo(downChannel, frm.Frame)
				lastActiveChannelMessageTime = time.Now()

			case downChannel:
				node.WriteFrameTo(activeChannel, frm.Frame)
			}

			// Monitor the active channel status
			activeChannelMsgTimeDiff := time.Since(lastActiveChannelMessageTime)

			if activeChannelMsgTimeDiff >= ConnectionInactiveDuration {
				fmt.Printf("Active channel down for > %v seconds\n", activeChannelMsgTimeDiff.Seconds())
				// Switch to the other channel
				if activeChannel == channelA {
					activeChannel = channelB
				} else {
					activeChannel = channelA
				}

				fmt.Printf("Switching to other channel\n")
				fmt.Printf("Active channel is %+v\n", activeChannel)

				lastActiveChannelMessageTime = time.Now()
			}
		}
	}
}
