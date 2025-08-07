/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package ping

import (
	"context"
	"net/http"

	"github.com/dtA4/go-rbfs-client/pkg/rbfs/state"
	"github.com/stretchr/testify/mock"
)

// ensure, that mockActionsAPI does implement ActionsAPI.
var _ ActionsAPI = &mockActionsAPI{}

type mockActionsAPI struct {
	mock.Mock
}

func (m *mockActionsAPI) Ping(ctx context.Context, localVarOptionals *state.ActionsApiPingOpts) (state.PingStatus, *http.Response, error) {
	args := m.Called(ctx, localVarOptionals)
	status, ok := args.Get(0).(state.PingStatus)
	if !ok {
		status = state.PingStatus{}
	}
	return status, nil, args.Error(1)
}
