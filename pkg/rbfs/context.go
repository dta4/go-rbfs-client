/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package rbfs

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dta4/go-rbfs-client/pkg/rbfs/state"
)

const (
	ctrldURLKey           = contextKey("CtrldURL")
	elementNameKey        = contextKey("ElementName")
	OpsdServiceName       = ServiceName("opsd")
	RestconfdServiceName  = ServiceName("restconfd")
	PrometheusServiceName = ServiceName("prometheus")
)

type (
	contextKey string

	ServiceName string

	// RbfsContext provides access to all request-specific settings to invoke the RBFS REST API.
	// It is also a go context that allows cancelling a request.
	RbfsContext interface {
		context.Context

		// GetServiceEndpoint computes the REST API endpoint for the given service.
		GetServiceEndpoint(ServiceName) (*url.URL, error)
		// GetCtrldElementsEndpoint computes the REST API for the elements collection.
		GetCtrldElementsEndpoint() (*url.URL, error)
		// GetCtrldElementEndpoint computes the REST API endpoint for the element resource accessible via the given path segments.
		GetCtrldElementEndpoint(...string) (*url.URL, error)
	}

	rbfsContext struct {
		context.Context
	}

	// RbfsContextOption allows applying an optional RBFS context setting.
	RbfsContextOption func(ctx context.Context) (context.Context, error)
)

// RbfsAccessToken adds an access token the a RBFS context.
func RbfsAccessToken(token string) RbfsContextOption {
	return func(ctx context.Context) (context.Context, error) {
		if token != "" {
			// Populate access token to all APIs
			ctx = context.WithValue(ctx, state.ContextAccessToken, token)
		}
		return ctx, nil
	}
}

// NewRbfsContext creates a new RBFS context from the given context to access an RBFS instance available under the
// given endpointURL with the specified elementName and the request options.
func NewRbfsContext(ctx context.Context, endpointURL *url.URL, elementName string, options ...RbfsContextOption) (*rbfsContext, error) {
	ctx = context.WithValue(ctx, ctrldURLKey, endpointURL)
	ctx = context.WithValue(ctx, elementNameKey, elementName)
	var err error
	for _, option := range options {
		// Apply optional settings to the diagnostic context
		ctx, err = option(ctx)
		if err != nil {
			return nil, err
		}
	}
	return &rbfsContext{Context: ctx}, nil
}

// MustRbfsContext creates a new RBFS context from the given context.
func MustRbfsContext(ctx context.Context) *rbfsContext {
	_, ok := ctx.Value(ctrldURLKey).(*url.URL)
	if !ok {
		panic("ctrldEndpoint not set")
	}
	_, ok = ctx.Value(elementNameKey).(string)
	if !ok {
		panic("elementName not set")
	}
	return &rbfsContext{Context: ctx}
}

func (r *rbfsContext) GetServiceEndpoint(serviceName ServiceName) (*url.URL, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("empty service name is not supported")
	}
	ctrldEndpoint := r.Value(ctrldURLKey).(*url.URL)
	elementName := r.Value(elementNameKey).(string)
	serviceEndpoint := fmt.Sprintf("%v/api/v1/rbfs/elements/%v/services/%v/proxy", ctrldEndpoint, elementName, serviceName)

	return url.Parse(serviceEndpoint)
}

func (r *rbfsContext) GetCtrldElementsEndpoint() (*url.URL, error) {
	ctrldEndpoint := r.Value(ctrldURLKey).(*url.URL)
	endpoint := fmt.Sprintf("%v/api/v1/ctrld/elements", ctrldEndpoint)
	return url.Parse(endpoint)
}

func (r *rbfsContext) GetCtrldElementEndpoint(pathSegments ...string) (*url.URL, error) {
	ctrldEndpoint := r.Value(ctrldURLKey).(*url.URL)
	elementName := r.Value(elementNameKey).(string)
	endpoint := ""
	if len(pathSegments) > 0 {
		endpoint = fmt.Sprintf("%v/api/v1/ctrld/elements/%v/%v", ctrldEndpoint, elementName, strings.Join(pathSegments, "/"))
	} else {
		endpoint = fmt.Sprintf("%v/api/v1/ctrld/elements/%v", ctrldEndpoint, elementName)
	}
	return url.Parse(endpoint)
}
