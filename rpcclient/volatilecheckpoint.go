// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"encoding/json"

	"github.com/monasuite/monad/btcjson"
)

// VolatileCheckpointCommand enumerates the available commands that the VolatileCheckpoint function
// accepts.
type VolatileCheckpointCommand string

// Constants used to indicate the command for the VolatileCheckpoint function.
const (
	// VCSet indicates the specified VolatileCheckpoint should be added as a persistent
	// volatilecheckpoint.
	VCSet VolatileCheckpointCommand = "set"

	// VCClear indicates the specified VolatileCheckpoint should be cleared.
	VCClear VolatileCheckpointCommand = "clear"
)

// String returns the VolatileCheckpointCommand in human-readable form.
func (cmd VolatileCheckpointCommand) String() string {
	return string(cmd)
}

// FutureVolatileCheckpointResult is a future promise to deliver the result of an
// VolatileCheckpointAsync RPC invocation (or an applicable error).
type FutureVolatileCheckpointResult chan *response

// Receive waits for the response promised by the future and returns an error if
// any occurred when performing the specified command.
func (r FutureVolatileCheckpointResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// VolatileCheckpointAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See VolatileCheckpoint for the blocking version and more details.
func (c *Client) VolatileCheckpointAsync(command VolatileCheckpointCommand, hash *string, blockheight *int64) FutureVolatileCheckpointResult {
	cmd := btcjson.NewVolatileCheckpointCmd(btcjson.VolatileCheckpointSubCmd(command), hash, blockheight)
	return c.sendCmd(cmd)
}

// VolatileCheckpoint attempts to perform the passed command on the passed persistent peer.
// For example, it can be used to add or a delete a persistent peer, or to do
// a one time connection to a peer.
//
// It may not be used to delete non-persistent peers.
func (c *Client) VolatileCheckpoint(command VolatileCheckpointCommand, hash *string, blockheight *int64) error {
	return c.VolatileCheckpointAsync(command, hash, blockheight).Receive()
}

// FutureDumpVolatileCheckpointResult is a future promise to deliver the result of a
// DumpVolatileCheckpointAsync RPC invocation (or an applicable error).
type FutureDumpVolatileCheckpointResult chan *response

// Receive waits for the response promised by the future and returns the hash of
// the block in the best block chain at the given height.
func (r FutureDumpVolatileCheckpointResult) Receive() ([]btcjson.DumpVolatileCheckpointResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as a string-encoded sha.
	var volatilecheckpoints []btcjson.DumpVolatileCheckpointResult
	err = json.Unmarshal(res, &volatilecheckpoints)
	if err != nil {
		return nil, err
	}
	return volatilecheckpoints, nil
}

// DumpVolatileCheckpointAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See DumpVolatileCheckpoint for the blocking version and more details.
func (c *Client) DumpVolatileCheckpointAsync() FutureDumpVolatileCheckpointResult {
	cmd := btcjson.NewDumpVolatileCheckpointCmd()
	return c.sendCmd(cmd)
}

// DumpVolatileCheckpoint returns the hash of the block in the best block chain at the
// given height.
func (c *Client) DumpVolatileCheckpoint() ([]btcjson.DumpVolatileCheckpointResult, error) {
	return c.DumpVolatileCheckpointAsync().Receive()
}
