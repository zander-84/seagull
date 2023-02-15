package registry

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/contrib/lb"
	"time"
)

type ServiceDiscovery struct {
	d        contract.Discovery
	lb       contract.Balancer
	listener lb.Listener
	watcher  contract.Watcher
}

func NewServiceDiscovery(serviceName string, d contract.Discovery, p lb.Policy) (*ServiceDiscovery, error) {
	var err error
	sd := new(ServiceDiscovery)
	sd.d = d
	sd.listener = lb.NewListener(serviceName)
	sd.lb = lb.NewBalancer(sd.listener, p, false)
	sd.watcher, err = sd.d.Watch(context.Background(), serviceName)
	if err != nil {
		return nil, err
	}
	_ = sd.set()

	go func() {
		for {
			if err := sd.set(); err != nil {
				if sd.watcher.IsStop() {
					break
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return sd, nil
}

func (sd *ServiceDiscovery) set() error {
	serviceInstances, err := sd.watcher.Next()
	if err != nil {
		sd.listener.SetErr(err)
		return err
	}
	in := make(map[any]int, 0)
	for _, serviceInstance := range serviceInstances {
		in[serviceInstance] = serviceInstance.Weight
	}
	sd.listener.Set(in)
	sd.listener.SetErr(nil)
	return nil
}

func (sd *ServiceDiscovery) GetServiceInstance() (contract.ServiceInstance, error) {
	outAny, err := sd.lb.Next()
	if err != nil {
		return contract.ServiceInstance{}, err
	}
	out := outAny.(*contract.ServiceInstance)
	return *out, nil
}
