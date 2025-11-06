/*
 * Copyright (C) 2021, RtBrick, Inc.
 * SPDX-License-Identifier: BSD-3-Clause
 */

package rbfs

import (
	"context"
	"net"
	"net/http"
	"syscall"
)

// createDialerWithTOS creates a DialContext function with TOS (Type of Service) marking capability
func createDialerWithTOS(tosValue int) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{}

		// Apply TOS before connect to server
		dialer.Control = func(network, address string, c syscall.RawConn) error {
			var controlErr error
			err := c.Control(func(fd uintptr) {
				controlErr = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_TOS, tosValue)
			})
			if err != nil {
				return err
			}
			return controlErr
		}

		conn, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// CreateNewHTTPClient creates and returns an *http.Client that applies the specified TOS value
// to all outgoing network packets.
func CreateNewHTTPClient(tosValue int) *http.Client {
	transport := &http.Transport{
		DialContext: createDialerWithTOS(tosValue),
	}
	return &http.Client{
		Transport: transport,
	}
}
