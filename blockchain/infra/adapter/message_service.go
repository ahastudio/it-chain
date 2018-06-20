package adapter

import (
	"errors"

	"time"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyNodeId = errors.New("empty nodeid proposed")

// ToDo: 구현.(gitId:junk-sound)
type Publish func(exchange string, topic string, data interface{}) (err error)

type MessageService struct {
	publish Publish // midgard.client.Publish
}

func NewMessageService(publish Publish) *MessageService {
	return &MessageService{
		publish: publish,
	}
}

func (md *MessageService) RequestBlock(nodeId p2p.NodeId) error {

	if nodeId.Id == "" {
		return ErrEmptyNodeId
	}

	body := blockchain.BlockRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := createMessageDeliverCommand(event.BlockRequestProtocol, body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, nodeId.ToString())

	return md.publish("Command", "message.deliver", deliverCommand)
}

//func (md *MessageService) ResponseBlock(nodeId p2p.NodeId) error {
//
//	return md.publish()
//}

func createMessageDeliverCommand(protocol string, body interface{}) (blockchain.MessageDeliverCommand, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return blockchain.MessageDeliverCommand{}, err
	}

	return blockchain.MessageDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
