package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Noder interface {
	Node() *Node
}

type noder struct {
	getter kv.Getter
	once   comfig.Once
}

func NewNoder(getter kv.Getter) Noder {
	return &noder{
		getter: getter,
	}
}

type Node struct {
	NodeUrl string
	ApiKey  string
}

func (n *noder) Node() *Node {
	result := n.once.Do(func() interface{} {
		var populatedNode Node
		err := figure.Out(&populatedNode).From(kv.MustGetStringMap(n.getter, "node")).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out node"))
		}
		return populatedNode
	})

	node := result.(Node)
	return &node
}

func (n *Node) GetNodeUrl() string {
	if n == nil {
		return ""
	}
	return n.NodeUrl
}

func (n *Node) GetApiKey() string {
	if n == nil {
		return ""
	}
	return n.ApiKey
}
