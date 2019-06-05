// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"encoding/json"

	"github.com/monasuite/monad/btcjson"
)

// CheckpointCommand enumerates the available commands that the Checkpoint function
// accepts.
type CheckpointCommand string

// Constants used to indicate the command for the Checkpoint function.
const (
	// CAdd indicates the specified Checkpoint should be added as a persistent
	// checkpoint.
	CAdd CheckpointCommand = "add"

	// CDelete indicates the specified Checkpoint should be deleted.
	CDelete CheckpointCommand = "delete"
)

// String returns the CheckpointCommand in human-readable form.
func (cmd CheckpointCommand) String() string {
	return string(cmd)
}

// FutureCheckpointResult is a future promise to deliver the result of an
// CheckpointAsync RPC invocation (or an applicable error).
type FutureCheckpointResult chan *response

// Receive waits for the response promised by the future and returns an error if
// any occurred when performing the specified command.
func (r FutureCheckpointResult) Receive() error {
	_, err := receiveFuture(r)
	return err
}

// CheckpointAsync returns an instance of a type that can be used to get the result
// of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See Checkpoint for the blocking version and more details.
func (c *Client) CheckpointAsync(command CheckpointCommand, hash *string, blockheight int64) FutureCheckpointResult {
	cmd := btcjson.NewCheckpointCmd(btcjson.CheckpointSubCmd(command), hash, blockheight)
	return c.sendCmd(cmd)
}

// Checkpoint attempts to perform the passed command on the passed persistent peer.
// For example, it can be used to add or a delete a persistent peer, or to do
// a one time connection to a peer.
//
// It may not be used to delete non-persistent peers.
func (c *Client) Checkpoint(command CheckpointCommand, hash *string, blockheight int64) error {
	return c.CheckpointAsync(command, hash, blockheight).Receive()
}

// FutureDumpCheckpointResult is a future promise to deliver the result of a
// DumpCheckpointAsync RPC invocation (or an applicable error).
type FutureDumpCheckpointResult chan *response

// Receive waits for the response promised by the future and returns the hash of
// the block in the best block chain at the given height.
func (r FutureDumpCheckpointResult) Receive() ([]btcjson.DumpCheckpointResult, error) {
	res, err := receiveFuture(r)
	if err != nil {
		return nil, err
	}

	// Unmarshal the result as a string-encoded sha.
	var checkpoints []btcjson.DumpCheckpointResult
	err = json.Unmarshal(res, &checkpoints)
	if err != nil {
		return nil, err
	}
	return checkpoints, nil
}

// DumpCheckpointAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on the
// returned instance.
//
// See DumpCheckpoint for the blocking version and more details.
func (c *Client) DumpCheckpointAsync(maxnum *int32) FutureDumpCheckpointResult {
	cmd := btcjson.NewDumpCheckpointCmd(maxnum)
	return c.sendCmd(cmd)
}

// DumpCheckpoint returns the hash of the block in the best block chain at the
// given height.
func (c *Client) DumpCheckpoint(maxnum *int32) ([]btcjson.DumpCheckpointResult, error) {
	return c.DumpCheckpointAsync(maxnum).Receive()
}
