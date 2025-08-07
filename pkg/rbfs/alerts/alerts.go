// Package alerts contains a client to query RBFS alerts
package alerts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dta4/go-rbfs-client/pkg/rbfs"
	"github.com/dta4/go-rbfs-client/pkg/rbfs/state"
)

type (

	// Alert describes a single switch alert.
	Alert struct {
		// AlertName holds the alert rule name
		AlertName string `json:"name"`
		// Summary holds the alert summary.
		Summary string `json:"summary"`
		// Level holds the alert level.
		Level int `json:"level"`
		// DateCreated holds the alert creation time.
		DateCreated time.Time `json:"date_created"`
	}

	// Client providess access to the switch metrics.
	Client interface {
		// QueryAlerts returns a list of firing alerts.
		QueryAlerts(ctx rbfs.RbfsContext) ([]Alert, error)
	}

	client struct {
		rbfs *http.Client
	}
)

// NewClient creates a new client to query switch alerts.
func NewClient(c *http.Client) Client {
	return &client{c}
}

func (c *client) QueryAlerts(ctx rbfs.RbfsContext) ([]Alert, error) {
	endpoint, err := ctx.GetServiceEndpoint(rbfs.PrometheusServiceName)
	if err != nil {
		return nil, err
	}

	// Compose the metric query
	queryURL := fmt.Sprintf("%s/api/v1/alerts", endpoint)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
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

	var responseJSON map[string]interface{}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(&responseJSON); err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		var aa []Alert
		if data, ok := responseJSON["data"].(map[string]interface{}); ok {
			if alerts, ok := data["alerts"].([]interface{}); ok {
				for _, i := range alerts {
					item := i.(map[string]interface{})
					state, _ := item["state"].(string)
					if state == "firing" {
						annotations := item["annotations"].(map[string]interface{})
						if level, ok := annotations["level"].(string); ok {
							labels := item["labels"].(map[string]interface{})
							name, _ := labels["alertname"].(string)
							summary, _ := annotations["summary"].(string)
							activeAt, _ := item["activeAt"].(string)

							var dateCreated time.Time
							if activeAt != "" {
								dateCreated, _ = time.Parse(time.RFC3339Nano, activeAt)
							}
							l, _ := strconv.Atoi(level)
							a := Alert{
								AlertName:   name,
								Summary:     summary,
								Level:       l,
								DateCreated: dateCreated,
							}
							aa = append(aa, a)
						}
					}

				}
			}
		}
		return aa, nil
	}
	return nil, fmt.Errorf("cannot read switch alerts. Status: %d", response.StatusCode)
}
