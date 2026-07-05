package state

import (
	"fan2go-tui/internal/client"
	"sync"
)

type Poller struct {
	apiClient client.Fan2goApiClient
	store     *Store
	onUpdate  func(error)
}

func NewPoller(apiClient client.Fan2goApiClient, store *Store, onUpdate func(error)) *Poller {
	return &Poller{
		apiClient: apiClient,
		store:     store,
		onUpdate:  onUpdate,
	}
}

func (p *Poller) FetchAndUpdate() {
	var wg sync.WaitGroup
	wg.Add(3)

	var fanErr, curveErr, sensorErr error

	go func() {
		defer wg.Done()
		fans, err := p.apiClient.GetFans()
		if err == nil && fans != nil {
			p.store.UpdateFans(*fans)
		} else if err != nil {
			fanErr = err
		}
	}()

	go func() {
		defer wg.Done()
		curves, err := p.apiClient.GetCurves()
		if err == nil && curves != nil {
			p.store.UpdateCurves(*curves)
		} else if err != nil {
			curveErr = err
		}
	}()

	go func() {
		defer wg.Done()
		sensors, err := p.apiClient.GetSensors()
		if err == nil && sensors != nil {
			p.store.UpdateSensors(*sensors)
		} else if err != nil {
			sensorErr = err
		}
	}()

	wg.Wait()
	if p.onUpdate != nil {
		var combinedErr error
		if fanErr != nil {
			combinedErr = fanErr
		} else if curveErr != nil {
			combinedErr = curveErr
		} else if sensorErr != nil {
			combinedErr = sensorErr
		}
		p.onUpdate(combinedErr)
	}
}
