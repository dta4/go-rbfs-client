/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package rbfs

import (
	"context"
	"net/url"
	"testing"

	"github.com/dtA4/go-rbfs-client/pkg/rbfs/state"
	"github.com/stretchr/testify/require"
)

// ensure, that rbfsContext does implement RbfsContext.
var _ RbfsContext = &rbfsContext{}

func TestNewRbfsContext(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		endpointURL *url.URL
		elementName string
		options     []RbfsContextOption
		wantErr     string
		validate    func(t *testing.T, r *rbfsContext)
	}{
		{
			name:        "no option",
			ctx:         context.Background(),
			endpointURL: mustParse(t, "http://192.168.0.1"),
			elementName: "rtbrick",
			validate: func(t *testing.T, r *rbfsContext) {
				u, err := r.GetServiceEndpoint("test")
				require.NoError(t, err)
				want := mustParse(t, "http://192.168.0.1/api/v1/rbfs/elements/rtbrick/services/test/proxy")
				require.Equal(t, want, u)
			},
		}, {
			name:        "empty service name is not supported",
			ctx:         context.Background(),
			endpointURL: mustParse(t, "http://192.168.0.1"),
			elementName: "rtbrick",
			validate: func(t *testing.T, r *rbfsContext) {
				_, err := r.GetServiceEndpoint("")
				require.EqualError(t, err, "empty service name is not supported")
			},
		}, {
			name:        "Rbfs Token",
			ctx:         context.Background(),
			endpointURL: mustParse(t, "http://192.168.0.1"),
			elementName: "rtbrick",
			options:     []RbfsContextOption{RbfsAccessToken("token")},
			validate: func(t *testing.T, r *rbfsContext) {
				token := r.Value(state.ContextAccessToken)
				require.Equal(t, "token", token)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := NewRbfsContext(tt.ctx, tt.endpointURL, tt.elementName, tt.options...)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			tt.validate(t, ctx)
		})
	}
}

func TestMustRbfsContext(t *testing.T) {
	ctx := context.Background()
	require.PanicsWithValue(t, "ctrldEndpoint not set", func() { MustRbfsContext(ctx) })
	ctx = context.WithValue(ctx, ctrldURLKey, mustParse(t, "http://192.168.0.1"))
	require.PanicsWithValue(t, "elementName not set", func() { MustRbfsContext(ctx) })
	ctx = context.WithValue(ctx, elementNameKey, "rtbrick")
	require.NotPanics(t, func() { MustRbfsContext(ctx) })
}

func mustParse(t *testing.T, v string) *url.URL {
	t.Helper()
	u, err := url.Parse(v)
	require.NoError(t, err)
	return u
}
