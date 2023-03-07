package tunnel

import (
	"errors"
	"sync"
)

type VhostStorage interface {
	// AddHost adds the given host and identifier to the storage
	AddHost(host VirtualHost, identifier string)

	// DeleteHost deletes the given host
	DeleteHost(host string)

	// GetHost returns the host name for the given identifier
	GetHost(identifier string) (VirtualHost, bool)

	// GetIdentifier returns the identifier for the given host
	GetIdentifier(host string) (string, bool)

	// GetNextHost returns next host on mapping
	GetNextHost() (VirtualHost, error)
}

type VirtualHost struct {
	Identifier    string
	Port          string
	RemoteAddress string
}

// virtualHosts is used for mapping host to users example: host
// "fs-1-fatih.kd.io" belongs to user "arslan"
type VirtualHosts struct {
	Mapping map[string]*VirtualHost
	Last    int
	sync.Mutex
}

// newVirtualHosts provides an in memory virtual host storage for mapping
// virtual hosts to identifiers.
func NewVirtualHosts() *VirtualHosts {
	return &VirtualHosts{
		Mapping: make(map[string]*VirtualHost),
		Last:    0,
	}
}

func (v *VirtualHosts) AddHost(vHost VirtualHost, identifier string) {
	v.Lock()
	v.Mapping[identifier] = &VirtualHost{
		Identifier:    vHost.Identifier,
		Port:          vHost.Port,
		RemoteAddress: vHost.RemoteAddress,
	}
	v.Unlock()
}

func (v *VirtualHosts) DeleteHost(identifier string) {
	v.Lock()
	delete(v.Mapping, identifier)
	v.Unlock()
}

// GetIdentifier returns the identifier associated with the given host
func (v *VirtualHosts) GetIdentifier(hostID string) (string, bool) {
	v.Lock()
	ht, ok := v.Mapping[hostID]
	v.Unlock()

	if !ok {
		return "", false
	}

	return ht.Identifier, true
}

// GetNextHost returns the next host in the map
func (v *VirtualHosts) GetNextHost() (VirtualHost, error) {

	v.Lock()
	defer v.Unlock()

	if v.Last == 0 && len(v.Mapping) == 0 {
		return VirtualHost{}, errors.New("there is no open connections")
	}

	index := 0
	for _, value := range v.Mapping {
		if index == v.Last {
			if v.Last == (len(v.Mapping) - 1) {
				v.Last = 0
			}
			v.Last += 1
			return *value, nil
		}
		index += 1
	}
	return VirtualHost{}, errors.New("could not find next host")

}

// GetHost returns the host associated with the given identifier
func (v *VirtualHosts) GetHost(identifier string) (VirtualHost, bool) {
	v.Lock()
	ht, ok := v.Mapping[identifier]
	v.Unlock()
	if !ok {
		return VirtualHost{}, false
	}

	return *ht, true
}
