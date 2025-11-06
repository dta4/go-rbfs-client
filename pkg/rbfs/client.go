/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package rbfs

import (
	"net/http"
	"net/url"

	"github.com/dtA4/go-rbfs-client/pkg/rbfs/state"
)

type Option func(*state.Configuration)

// DefaultHeader returns an option to add a default HTTP header
func DefaultHeader(name, value string) Option {
	return func(c *state.Configuration) {
		c.DefaultHeader[name] = value
	}
}

// UserAgent returns an option to set the user agent to the given value.
func UserAgent(value string) Option {
	return func(c *state.Configuration) {
		c.UserAgent = value
	}
}

// GetAPIClient creates a new API client for the specified endpoint, using TOS 192 by default.
// All RBFS network management communications use TOS 192 to ensure the highest transmission priority.
func GetAPIClient(client *http.Client, endpoint *url.URL, options ...Option) *state.APIClient {
	config := state.NewConfiguration()
	config.BasePath = endpoint.String()
	config.Host = endpoint.Host
	config.HTTPClient = CreateNewHTTPClient(192)
	config.UserAgent = "go-client"
	for _, option := range options {
		option(config)
	}
	return state.NewAPIClient(config)
}
