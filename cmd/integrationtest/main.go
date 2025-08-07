package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dtA4/go-rbfs-client/pkg/rbfs"
	"github.com/dtA4/go-rbfs-client/pkg/rbfs/alerts"
	"github.com/dtA4/go-rbfs-client/pkg/rbfs/elements"
	"github.com/dtA4/go-rbfs-client/pkg/rbfs/metrics"
	"github.com/dtA4/go-rbfs-client/pkg/rbfs/services"
)

var client http.Client

const (
	endpointURL string = "http://10.200.134.26:19091"
	elementName string = "ufi09.q2c.u19.r4.nbg.rtbrick.net"
	accessToken string = ""
)

//nolint:forbidigo  // this is an integration test
func main() {
	m := metrics.NewClient(&client)
	endpoint, _ := url.Parse(endpointURL)
	ctx, _ := rbfs.NewRbfsContext(context.Background(), endpoint, elementName, rbfs.RbfsAccessToken(accessToken))

	metric, err := m.QueryMetric(ctx, "chassis_temperature_millicelsius")
	fmt.Println(err)
	b, _ := json.MarshalIndent(metric, " ", " ")
	fmt.Println(string(b))

	a := alerts.NewClient(&client)
	alerts, err := a.QueryAlerts(ctx)
	fmt.Println(err)
	b, _ = json.MarshalIndent(alerts, " ", " ")
	fmt.Println(string(b))

	s := services.NewClient(&client)
	services, err := s.ListServices(ctx)
	fmt.Println(err)
	b, _ = json.MarshalIndent(services, " ", " ")
	fmt.Println(string(b))

	e := elements.NewClient(&client)
	elements, err := e.ListElements(ctx)
	fmt.Println(err)
	b, _ = json.MarshalIndent(elements, " ", " ")
	fmt.Println(string(b))

	element, err := e.GetElement(ctx, "ufi09.q2c.u19.r4.nbg.rtbrick.net")
	fmt.Println(err)
	b, _ = json.MarshalIndent(element, " ", " ")
	fmt.Println(string(b))
}
