package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dta4/go-rbfs-client/pkg/rbfs"
	"github.com/dta4/go-rbfs-client/pkg/rbfs/state"
)

type (

	// Service describes a running brick daemon or services.
	Service struct {
		// ServiceName holds the daemon name.
		ServiceName string `json:"service_name"`
		// AdministrativeState holds the administrative daemon state.
		AdministrativeState string `json:"administrative_state"`
		// OperationalState holds the operational daemon state.
		OperationalState string `json:"operational_state"`
	}

	// Client provides access to the switch metrics.
	Client interface {
		// ListServices returns all daemons and their respective state.
		ListServices(ctx rbfs.RbfsContext) ([]Service, error)
	}

	client struct {
		rbfs *http.Client
	}
)

// NewClient creates a new client to query running RBFS services and daemons.
func NewClient(c *http.Client) Client {
	return &client{c}
}

func (c *client) ListServices(ctx rbfs.RbfsContext) ([]Service, error) {
	endpoint, err := ctx.GetCtrldElementEndpoint("services")
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

	var services []Service

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	err = decoder.Decode(&services)
	if err != nil {
		return nil, fmt.Errorf("cannot read service list: %v", err)
	}
	return services, nil
}
