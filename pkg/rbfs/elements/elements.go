// Package elements provides access to the element informations. A single CTRLD can maintain multiple RBFS instances, which are called elements.
package elements

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/dta4/go-rbfs-client/pkg/rbfs"
	"github.com/dta4/go-rbfs-client/pkg/rbfs/state"
)

const (
	ContainerStateInitializing = ContainerState("INITIALIZING")
	ContainerStateStopper      = ContainerState("STOPPED")
	ContainerStateStarting     = ContainerState("STARTING")
	ContainerStateRunning      = ContainerState("RUNNING")
	ContainerStateStopping     = ContainerState("STOPPING")
	ContainerStateAborting     = ContainerState("ABORTING")
	ContainerStateFreezing     = ContainerState("FREEZING")
	ContainerStateFrozen       = ContainerState("FROZEN")
	ConainerStateThawed        = ContainerState("THAWED")

	OperationalStateUp   = OperationalState("UP")
	OperationalStateDown = OperationalState("DOWN")
)

type (
	// ContainerState describes the operational state of the RBFS linux container.
	ContainerState string

	// OperationalState describes the state of the RBFS instance (running inside the RBFS linux container)
	OperationalState string

	// Element describes a single element, which is managed by CTRLD.
	Element struct {
		// ContainerName holds the name of the RBFS linux container.
		ContainerName string `json:"container_name"`
		// Container state holds the operational state of the RBFS linux container.
		ContainerState ContainerState `json:"container_state"`
		// IPAddresses holds the IP addresses associated with the container.
		IPAddresses []net.IP `json:"ip_addresses"`
		// ElementName holds the element name assigned to the RBFS container.
		ElementName string `json:"element_name"`
		// PodName hold the name of the pod the element belongs to.
		PodName string `json:"pod_name"`
		// OperationalState holds the operational state of the RBFS instance.
		OperationalState OperationalState `json:"operational_state"`
		// ZTPEnabled indicates whether RBFS pulls the startup configuration from the ZTP server after each reboot or not.
		ZTPEnabled bool `json:"ztp_enabled"`
	}

	// Client provides access to the available elements.
	Client interface {
		// ListElements returns all elements managed by the CTRLD instance.
		ListElements(ctx rbfs.RbfsContext) ([]Element, error)
		// GetElement returns the element with the given name managed by the CTRLD instance.
		GetElement(ctx rbfs.RbfsContext, elementName string) (*Element, error)
	}

	client struct {
		rbfs *http.Client
	}
)

// NewClient creates a new client to query managed elements.
func NewClient(c *http.Client) Client {
	return &client{c}
}

func (c *client) ListElements(ctx rbfs.RbfsContext) ([]Element, error) {
	endpoint, err := ctx.GetCtrldElementsEndpoint()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	if accessToken, ok := ctx.Value(state.ContextAccessToken).(string); ok {
		request.Header.Add("Authorization", "Bearer "+accessToken)
	}

	response, err := c.rbfs.Do(request)
	if err != nil {
		return nil, err
	}

	var elements []Element

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&elements)
	if err != nil {
		return nil, fmt.Errorf("cannot read element list: %v", err)
	}
	return elements, nil
}

func (c *client) GetElement(ctx rbfs.RbfsContext, elementName string) (*Element, error) {
	endpoint, err := ctx.GetCtrldElementEndpoint()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	if accessToken, ok := ctx.Value(state.ContextAccessToken).(string); ok {
		request.Header.Add("Authorization", "Bearer "+accessToken)
	}

	response, err := c.rbfs.Do(request)
	if err != nil {
		return nil, err
	}

	var element Element

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&element)
	if err != nil {
		return nil, fmt.Errorf("cannot read element list: %v", err)
	}
	return &element, nil
}
