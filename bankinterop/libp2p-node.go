package bankinterop

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	. "github.com/san-lab/immudb-tests/datastructs"
)

var NET string
var LIBP2P_NODE *Node

// MPCNode represents a subscription to a single PubSub topic. Messages
// can be published to the topic with MPCNode.Publish, and received
// messages are pushed to the Messages channel.
type Node struct {
	// Messages is a channel of messages received from other peers in the chat room
	ctx   context.Context
	ps    *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription
	self  peer.ID
}

type GenericMessage struct {
	MessageType string
	Recipient   string
	Data        []byte
}

// readLoop pulls messages from the pubsub topic and pushes them onto the Messages channel.
func (node *Node) readLoop() {
	for {
		msg, err := node.sub.Next(node.ctx)
		if err != nil {
			log.Fatalln(err)
			return
		}
		// only forward messages delivered by others
		if msg.ReceivedFrom == node.self {
			continue
		}
		node.ProcessMessage(msg)
	}
}

func (node *Node) ProcessMessage(msg *pubsub.Message) {
	genmsg := new(GenericMessage)
	err := json.Unmarshal(msg.Data, genmsg)
	if err != nil {
		fmt.Println("bad frame:", err)
		return
	}
	if genmsg.Recipient != "" && genmsg.Recipient != THIS_BANK.Name {
		// Not addressed to us
		return
	}
	HandleMessage(genmsg.MessageType, genmsg.Data)
}

func (node *Node) SendMessage(dataType, recipient string, data []byte) {
	genmsg := GenericMessage{MessageType: dataType, Recipient: recipient, Data: data}
	b, _ := json.Marshal(genmsg)
	node.topic.Publish(context.Background(), b)
}

func (node *Node) PeerCount() int {
	mutex.Lock()
	defer mutex.Unlock()
	return len(node.topic.ListPeers())
}

// tries to subscribe to the PubSub topic for the room name, returning
// an MPCNode on success.
func JoinNet(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, roomName string) (*Node, error) {
	// join the pubsub topic
	mutex.Lock()
	defer mutex.Unlock()
	topic, err := ps.Join(NET)
	if err != nil {
		return nil, err
	}

	// and subscribe to it
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	node := &Node{
		ctx:   ctx,
		ps:    ps,
		topic: topic,
		sub:   sub,
		self:  selfID,
	}

	// start reading messages from the subscription in a loop
	go node.readLoop()
	return node, nil
}

func (node *Node) GetNodeID() string {
	return string(node.self) // not working??
}
