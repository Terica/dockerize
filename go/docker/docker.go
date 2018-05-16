// +build !mocker

// Package docker is a convenience wrapper for fsouza's go-dockerclient
package docker

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/fayep/dockerize/go/progress"
	"io"
	"io/ioutil"
)

// Docker interface.
type Docker interface {
	PStat(map[string][]string) []APIContainers
	Pull(string, string, *progress.Progress)
	Run(string, string, []string, []string, []string) (string, error)
	Exec(string, []string, string, []string) (int, error)
}

// Docker interface.
type Mocker interface {
	Docker
	AddContainer(APIContainers)
	RemoveContainer(string)
	AddImage(APIContainers)
	RemoveImage(string)
	ClearFlags()
	SetFlag(string, bool)
	GetFlag(string) bool
}

// ProgressDetail used by status updates of at least "Pull"
type ProgressDetail struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

// StatusUpdate used for docker status updates
type StatusUpdate struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Detail *ProgressDetail `json:"progressDetail"`
}

// APIPort is a type that represents a port mapping returned by the Docker API
type APIPort struct {
	PrivatePort int64  `json:"PrivatePort,omitempty" yaml:"PrivatePort,omitempty" toml:"PrivatePort,omitempty"`
	PublicPort  int64  `json:"PublicPort,omitempty" yaml:"PublicPort,omitempty" toml:"PublicPort,omitempty"`
	Type        string `json:"Type,omitempty" yaml:"Type,omitempty" toml:"Type,omitempty"`
	IP          string `json:"IP,omitempty" yaml:"IP,omitempty" toml:"IP,omitempty"`
}

// APIMount represents a mount point for a container.
type APIMount struct {
	Name        string `json:"Name,omitempty" yaml:"Name,omitempty" toml:"Name,omitempty"`
	Source      string `json:"Source,omitempty" yaml:"Source,omitempty" toml:"Source,omitempty"`
	Destination string `json:"Destination,omitempty" yaml:"Destination,omitempty" toml:"Destination,omitempty"`
	Driver      string `json:"Driver,omitempty" yaml:"Driver,omitempty" toml:"Driver,omitempty"`
	Mode        string `json:"Mode,omitempty" yaml:"Mode,omitempty" toml:"Mode,omitempty"`
	RW          bool   `json:"RW,omitempty" yaml:"RW,omitempty" toml:"RW,omitempty"`
	Propogation string `json:"Propogation,omitempty" yaml:"Propogation,omitempty" toml:"Propogation,omitempty"`
}

// ContainerNetwork represents the networking settings of a container per network.
type ContainerNetwork struct {
	Aliases             []string `json:"Aliases,omitempty" yaml:"Aliases,omitempty" toml:"Aliases,omitempty"`
	MacAddress          string   `json:"MacAddress,omitempty" yaml:"MacAddress,omitempty" toml:"MacAddress,omitempty"`
	GlobalIPv6PrefixLen int      `json:"GlobalIPv6PrefixLen,omitempty" yaml:"GlobalIPv6PrefixLen,omitempty" toml:"GlobalIPv6PrefixLen,omitempty"`
	GlobalIPv6Address   string   `json:"GlobalIPv6Address,omitempty" yaml:"GlobalIPv6Address,omitempty" toml:"GlobalIPv6Address,omitempty"`
	IPv6Gateway         string   `json:"IPv6Gateway,omitempty" yaml:"IPv6Gateway,omitempty" toml:"IPv6Gateway,omitempty"`
	IPPrefixLen         int      `json:"IPPrefixLen,omitempty" yaml:"IPPrefixLen,omitempty" toml:"IPPrefixLen,omitempty"`
	IPAddress           string   `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty" toml:"IPAddress,omitempty"`
	Gateway             string   `json:"Gateway,omitempty" yaml:"Gateway,omitempty" toml:"Gateway,omitempty"`
	EndpointID          string   `json:"EndpointID,omitempty" yaml:"EndpointID,omitempty" toml:"EndpointID,omitempty"`
	NetworkID           string   `json:"NetworkID,omitempty" yaml:"NetworkID,omitempty" toml:"NetworkID,omitempty"`
}

// NetworkList encapsulates a map of networks, as returned by the Docker API in
// ListContainers.
type NetworkList struct {
	Networks map[string]ContainerNetwork `json:"Networks" yaml:"Networks,omitempty" toml:"Networks,omitempty"`
}

// APIContainers show API Container state

type APIContainers struct {
	ID         string            `json:"Id" yaml:"Id" toml:"Id"`
	Image      string            `json:"Image,omitempty" yaml:"Image,omitempty" toml:"Image,omitempty"`
	Command    string            `json:"Command,omitempty" yaml:"Command,omitempty" toml:"Command,omitempty"`
	Created    int64             `json:"Created,omitempty" yaml:"Created,omitempty" toml:"Created,omitempty"`
	State      string            `json:"State,omitempty" yaml:"State,omitempty" toml:"State,omitempty"`
	Status     string            `json:"Status,omitempty" yaml:"Status,omitempty" toml:"Status,omitempty"`
	Ports      []APIPort         `json:"Ports,omitempty" yaml:"Ports,omitempty" toml:"Ports,omitempty"`
	SizeRw     int64             `json:"SizeRw,omitempty" yaml:"SizeRw,omitempty" toml:"SizeRw,omitempty"`
	SizeRootFs int64             `json:"SizeRootFs,omitempty" yaml:"SizeRootFs,omitempty" toml:"SizeRootFs,omitempty"`
	Names      []string          `json:"Names,omitempty" yaml:"Names,omitempty" toml:"Names,omitempty"`
	Labels     map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty" toml:"Labels,omitempty"`
	Networks   NetworkList       `json:"NetworkSettings,omitempty" yaml:"NetworkSettings,omitempty" toml:"NetworkSettings,omitempty"`
	Mounts     []APIMount        `json:"Mounts,omitempty" yaml:"Mounts,omitempty" toml:"Mounts,omitempty"`
}

// DeepCopy uses reflection to map a source struct into a destination one.
// gob is used as the encoding.
func DeepCopy(from interface{}, to interface{}) error {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	if err := enc.Encode(from); err != nil {
		return err
	}
	err := dec.Decode(to)
	return err
}

func manageProgress(status io.Reader, pb *progress.Progress) {
	if pb == nil {
		// Junk the status feed.  It has to be read until its end.
		io.Copy(ioutil.Discard, status)
	} else {
		decoder := json.NewDecoder(status)
		// We use Number to ensure that large numbers work ok.
		decoder.UseNumber()
		for decoder.More() {
			var (
				m StatusUpdate
			)
			decoder.Decode(&m)
			if m.Status == "Downloading" || m.Status == "Extracting" {
				if m.Status == "Downloading" {
					pb.OnlyAdd("Extracting "+m.ID, 0, m.Detail.Total)
				}
				pb.Add(m.Status+" "+m.ID, m.Detail.Current, m.Detail.Total)
				pb.Display()
			}
		}
		pb.Done()
	}
}
