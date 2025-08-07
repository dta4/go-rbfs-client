/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package ping

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/dtA4/go-rbfs-client/pkg/rbfs"
	"github.com/dtA4/go-rbfs-client/pkg/rbfs/state"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func expectedPing(hostname string) *state.ActionsApiPingOpts {
	return &state.ActionsApiPingOpts{
		DestinationIp:   rbfs.OptionalIP(nil),
		DestinationAaaa: rbfs.OptionalString(""),
		DestinationA:    rbfs.OptionalString(hostname),
		SourceIp:        rbfs.OptionalIP(nil),
		SourceIfl:       rbfs.OptionalString(""),
		Count:           rbfs.OptionalInt(5),
		Interval:        rbfs.OptionalFloat32(1.0),
		InstanceName:    rbfs.OptionalString("default"),
		Size:            rbfs.OptionalInt(56),
		Ttl:             rbfs.OptionalInt(64),
	}
}

func Test_defaultService_Run(t *testing.T) {
	endpoint, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)
	ctx, err := rbfs.NewRbfsContext(context.Background(), endpoint, "test")
	require.NoError(t, err)

	tests := []struct {
		name string

		ping    func(t *testing.T) *Ping
		setup   func(t *testing.T, actionApi *mockActionsAPI)
		want    state.PingStatus
		wantErr bool
	}{
		{
			ping: func(t *testing.T) *Ping {
				p, err := NewPing(DestinationHostNameA("www.rtbrick.com"))
				require.NoError(t, err)
				return p
			},
			setup: func(t *testing.T, actionApi *mockActionsAPI) {
				actionApi.On("Ping", mock.Anything, expectedPing("www.rtbrick.com")).
					Return(state.PingStatus{Command: "www.rtbrick.com"}, nil)
			},
			want: state.PingStatus{Command: "www.rtbrick.com"},
		},
		{
			ping: func(t *testing.T) *Ping {
				p, err := NewPing(DestinationHostNameA("www.rtbrick.com"))
				require.NoError(t, err)
				return p
			},
			setup: func(t *testing.T, actionApi *mockActionsAPI) {
				actionApi.On("Ping", mock.Anything, expectedPing("www.rtbrick.com")).
					Return(state.PingStatus{}, fmt.Errorf("test"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actionAPI := &mockActionsAPI{}
			getActionsAPIFuncOrig := getActionsAPIFunc
			defer func() { getActionsAPIFunc = getActionsAPIFuncOrig }()
			getActionsAPIFunc = func(c *http.Client, endpoint *url.URL) (ActionsAPI, error) {
				return actionAPI, nil
			}
			actionAPI.Test(t)
			tt.setup(t, actionAPI)
			defer actionAPI.AssertExpectations(t)

			s := NewPingService(nil)
			got, err := s.Run(ctx, tt.ping(t))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_defaultService_RunAll(t *testing.T) {
	endpoint, err := url.Parse("http://localhost:8080")
	require.NoError(t, err)
	ctx, err := rbfs.NewRbfsContext(context.Background(), endpoint, "test")
	require.NoError(t, err)

	tests := []struct {
		name string

		ping    func(t *testing.T) []*Ping
		setup   func(t *testing.T, actionApi *mockActionsAPI)
		want    []state.PingStatus
		wantErr bool
	}{
		{
			ping: func(t *testing.T) []*Ping {
				var result []*Ping
				p, err := NewPing(DestinationHostNameA("www.rtbrick.com"))
				require.NoError(t, err)
				result = append(result, p)
				p, err = NewPing(DestinationHostNameA("www.google.de"))
				require.NoError(t, err)
				result = append(result, p)
				return result
			},
			setup: func(t *testing.T, actionApi *mockActionsAPI) {
				actionApi.On("Ping", mock.Anything, expectedPing("www.rtbrick.com")).
					Return(state.PingStatus{Command: "www.rtbrick.com"}, nil).
					Run(func(args mock.Arguments) { time.Sleep(time.Second * 3) })
				actionApi.On("Ping", mock.Anything, expectedPing("www.google.de")).
					Return(state.PingStatus{Command: "www.google.de"}, nil).
					Run(func(args mock.Arguments) { time.Sleep(time.Second * 2) })
			},
			want: []state.PingStatus{{Command: "www.rtbrick.com"}, {Command: "www.google.de"}},
		},
		{
			ping: func(t *testing.T) []*Ping {
				var result []*Ping
				p, err := NewPing(DestinationHostNameA("www.rtbrick.com"))
				require.NoError(t, err)
				result = append(result, p)
				return result
			},
			setup: func(t *testing.T, actionApi *mockActionsAPI) {
				actionApi.On("Ping", mock.Anything, expectedPing("www.rtbrick.com")).
					Return(state.PingStatus{}, fmt.Errorf("test"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actionAPI := &mockActionsAPI{}
			getActionsAPIFuncOrig := getActionsAPIFunc
			defer func() { getActionsAPIFunc = getActionsAPIFuncOrig }()
			getActionsAPIFunc = func(c *http.Client, endpoint *url.URL) (ActionsAPI, error) {
				return actionAPI, nil
			}
			actionAPI.Test(t)
			tt.setup(t, actionAPI)
			defer actionAPI.AssertExpectations(t)

			s := NewPingService(nil)
			got, err := s.RunAll(ctx, tt.ping(t)...)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, got)
		})
	}
}
