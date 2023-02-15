package etcd

import (
	"encoding/json"
	"github.com/zander-84/seagull/contract"
)

func marshal(si *contract.ServiceInstance) (string, error) {
	data, err := json.Marshal(si)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshal(data []byte) (si *contract.ServiceInstance, err error) {
	err = json.Unmarshal(data, &si)
	return
}
