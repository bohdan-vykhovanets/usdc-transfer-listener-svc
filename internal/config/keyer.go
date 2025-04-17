package config

import "gitlab.com/distributed_lab/kit/kv"

type Keyer interface {
	InfuraApiKey() string
}

type keyer struct {
	getter kv.Getter
}

func NewKeyer(getter kv.Getter) Keyer {
	return &keyer{
		getter: getter,
	}
}

func (k *keyer) InfuraApiKey() string {
	key := kv.MustGetStringMap(k.getter, "keys")
	return key["infura_api"].(string)
}
