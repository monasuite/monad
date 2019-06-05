// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/monasuite/monad/chaincfg/chainhash"
)

// TestBlock tests the MsgBlock API.
func TestBlock(t *testing.T) {
	pver := ProtocolVersion

	// Block 1 header.
	prevHash := &blockOne.Header.PrevBlock
	merkleHash := &blockOne.Header.MerkleRoot
	bits := blockOne.Header.Bits
	nonce := blockOne.Header.Nonce
	bh := NewBlockHeader(1, prevHash, merkleHash, bits, nonce)

	// Ensure the command is expected value.
	wantCmd := "block"
	msg := NewMsgBlock(bh)
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgBlock: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	// Num addresses (varInt) + max allowed addresses.
	wantPayload := uint32(4000000)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Ensure we get the same block header data back out.
	if !reflect.DeepEqual(&msg.Header, bh) {
		t.Errorf("NewMsgBlock: wrong block header - got %v, want %v",
			spew.Sdump(&msg.Header), spew.Sdump(bh))
	}

	// Ensure transactions are added properly.
	tx := blockOne.Transactions[0].Copy()
	msg.AddTransaction(tx)
	if !reflect.DeepEqual(msg.Transactions, blockOne.Transactions) {
		t.Errorf("AddTransaction: wrong transactions - got %v, want %v",
			spew.Sdump(msg.Transactions),
			spew.Sdump(blockOne.Transactions))
	}

	// Ensure transactions are properly cleared.
	msg.ClearTransactions()
	if len(msg.Transactions) != 0 {
		t.Errorf("ClearTransactions: wrong transactions - got %v, want %v",
			len(msg.Transactions), 0)
	}
}

// TestBlockTxHashes tests the ability to generate a slice of all transaction
// hashes from a block accurately.
func TestBlockTxHashes(t *testing.T) {
	// Block 1, transaction 1 hash.
	hashStr := "10067abeabcd96a1261bc542b16d686d083308304923d74cb8f3bab4209cc3b9"
	wantHash, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
		return
	}

	wantHashes := []chainhash.Hash{*wantHash}
	hashes, err := blockOne.TxHashes()
	if err != nil {
		t.Errorf("TxHashes: %v", err)
	}
	if !reflect.DeepEqual(hashes, wantHashes) {
		t.Errorf("TxHashes: wrong transaction hashes - got %v, want %v",
			spew.Sdump(hashes), spew.Sdump(wantHashes))
	}
}

// TestBlockHash tests the ability to generate the hash of a block accurately.
func TestBlockHash(t *testing.T) {
	// Block 1 hash.
	hashStr := "a378b7a601e00a877160357137c72befc6c06c22c5997eedb30ce9976710d28b"
	wantHash, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Ensure the hash produced is expected.
	blockHash := blockOne.BlockHash()
	if !blockHash.IsEqual(wantHash) {
		t.Errorf("BlockHash: wrong hash - got %v, want %v",
			spew.Sprint(blockHash), spew.Sprint(wantHash))
	}
}

// TestBlockWire tests the MsgBlock wire encode and decode for various numbers
// of transaction inputs and outputs and protocol versions.
func TestBlockWire(t *testing.T) {
	tests := []struct {
		in     *MsgBlock       // Message to encode
		out    *MsgBlock       // Expected decoded message
		buf    []byte          // Wire encoding
		txLocs []TxLoc         // Expected transaction locations
		pver   uint32          // Protocol version for wire encoding
		enc    MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
			ProtocolVersion,
			BaseEncoding,
		},

		// Protocol version BIP0035Version.
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
			BIP0035Version,
			BaseEncoding,
		},

		// Protocol version BIP0031Version.
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
			BIP0031Version,
			BaseEncoding,
		},

		// Protocol version NetAddressTimeVersion.
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
			NetAddressTimeVersion,
			BaseEncoding,
		},

		// Protocol version MultipleAddressVersion.
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
			MultipleAddressVersion,
			BaseEncoding,
		},
		// TODO(roasbeef): add case for witnessy block
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		err := test.in.BtcEncode(&buf, test.pver, test.enc)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from wire format.
		var msg MsgBlock
		rbuf := bytes.NewReader(test.buf)
		err = msg.BtcDecode(rbuf, test.pver, test.enc)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(&msg), spew.Sdump(test.out))
			continue
		}
	}
}

// TestBlockWireErrors performs negative tests against wire encode and decode
// of MsgBlock to confirm error paths work correctly.
func TestBlockWireErrors(t *testing.T) {
	// Use protocol version 60002 specifically here instead of the latest
	// because the test data is using bytes encoded with that protocol
	// version.
	pver := uint32(60002)

	tests := []struct {
		in       *MsgBlock       // Value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for wire encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Force error in version.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 0, io.ErrShortWrite, io.EOF},
		// Force error in prev block hash.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 4, io.ErrShortWrite, io.EOF},
		// Force error in merkle root.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 36, io.ErrShortWrite, io.EOF},
		// Force error in timestamp.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 68, io.ErrShortWrite, io.EOF},
		// Force error in difficulty bits.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 72, io.ErrShortWrite, io.EOF},
		// Force error in header nonce.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 76, io.ErrShortWrite, io.EOF},
		// Force error in transaction count.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 80, io.ErrShortWrite, io.EOF},
		// Force error in transactions.
		{&blockOne, blockOneBytes, pver, BaseEncoding, 81, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		w := newFixedWriter(test.max)
		err := test.in.BtcEncode(w, test.pver, test.enc)
		if err != test.writeErr {
			t.Errorf("BtcEncode #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from wire format.
		var msg MsgBlock
		r := newFixedReader(test.max, test.buf)
		err = msg.BtcDecode(r, test.pver, test.enc)
		if err != test.readErr {
			t.Errorf("BtcDecode #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestBlockSerialize tests MsgBlock serialize and deserialize.
func TestBlockSerialize(t *testing.T) {
	tests := []struct {
		in     *MsgBlock // Message to encode
		out    *MsgBlock // Expected decoded message
		buf    []byte    // Serialized data
		txLocs []TxLoc   // Expected transaction locations
	}{
		{
			&blockOne,
			&blockOne,
			blockOneBytes,
			blockOneTxLocs,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Serialize the block.
		var buf bytes.Buffer
		err := test.in.Serialize(&buf)
		if err != nil {
			t.Errorf("Serialize #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("Serialize #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Deserialize the block.
		var block MsgBlock
		rbuf := bytes.NewReader(test.buf)
		err = block.Deserialize(rbuf)
		if err != nil {
			t.Errorf("Deserialize #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&block, test.out) {
			t.Errorf("Deserialize #%d\n got: %s want: %s", i,
				spew.Sdump(&block), spew.Sdump(test.out))
			continue
		}

		// Deserialize the block while gathering transaction location
		// information.
		var txLocBlock MsgBlock
		br := bytes.NewBuffer(test.buf)
		txLocs, err := txLocBlock.DeserializeTxLoc(br)
		if err != nil {
			t.Errorf("DeserializeTxLoc #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&txLocBlock, test.out) {
			t.Errorf("DeserializeTxLoc #%d\n got: %s want: %s", i,
				spew.Sdump(&txLocBlock), spew.Sdump(test.out))
			continue
		}
		if !reflect.DeepEqual(txLocs, test.txLocs) {
			t.Errorf("DeserializeTxLoc #%d\n got: %s want: %s", i,
				spew.Sdump(txLocs), spew.Sdump(test.txLocs))
			continue
		}
	}
}

// TestBlockSerializeErrors performs negative tests against wire encode and
// decode of MsgBlock to confirm error paths work correctly.
func TestBlockSerializeErrors(t *testing.T) {
	tests := []struct {
		in       *MsgBlock // Value to encode
		buf      []byte    // Serialized data
		max      int       // Max size of fixed buffer to induce errors
		writeErr error     // Expected write error
		readErr  error     // Expected read error
	}{
		// Force error in version.
		{&blockOne, blockOneBytes, 0, io.ErrShortWrite, io.EOF},
		// Force error in prev block hash.
		{&blockOne, blockOneBytes, 4, io.ErrShortWrite, io.EOF},
		// Force error in merkle root.
		{&blockOne, blockOneBytes, 36, io.ErrShortWrite, io.EOF},
		// Force error in timestamp.
		{&blockOne, blockOneBytes, 68, io.ErrShortWrite, io.EOF},
		// Force error in difficulty bits.
		{&blockOne, blockOneBytes, 72, io.ErrShortWrite, io.EOF},
		// Force error in header nonce.
		{&blockOne, blockOneBytes, 76, io.ErrShortWrite, io.EOF},
		// Force error in transaction count.
		{&blockOne, blockOneBytes, 80, io.ErrShortWrite, io.EOF},
		// Force error in transactions.
		{&blockOne, blockOneBytes, 81, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Serialize the block.
		w := newFixedWriter(test.max)
		err := test.in.Serialize(w)
		if err != test.writeErr {
			t.Errorf("Serialize #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Deserialize the block.
		var block MsgBlock
		r := newFixedReader(test.max, test.buf)
		err = block.Deserialize(r)
		if err != test.readErr {
			t.Errorf("Deserialize #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}

		var txLocBlock MsgBlock
		br := bytes.NewBuffer(test.buf[0:test.max])
		_, err = txLocBlock.DeserializeTxLoc(br)
		if err != test.readErr {
			t.Errorf("DeserializeTxLoc #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestBlockOverflowErrors  performs tests to ensure deserializing blocks which
// are intentionally crafted to use large values for the number of transactions
// are handled properly.  This could otherwise potentially be used as an attack
// vector.
func TestBlockOverflowErrors(t *testing.T) {
	// Use protocol version 70001 specifically here instead of the latest
	// protocol version because the test data is using bytes encoded with
	// that version.
	pver := uint32(70001)

	tests := []struct {
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for wire encoding
		enc  MessageEncoding // Message encoding format
		err  error           // Expected error
	}{
		// Block that claims to have ~uint64(0) transactions.
		{
			[]byte{
				0x01, 0x00, 0x00, 0x00, // Version 1
				0xb6, 0x8b, 0x8c, 0x41, 0x0d, 0x2e, 0xa4, 0xaf,
				0xd7, 0x4f, 0xb5, 0x6e, 0x37, 0x0b, 0xfc, 0x1b,
				0xed, 0xf9, 0x29, 0xe1, 0x45, 0x38, 0x96, 0xc9,
				0xe7, 0x9d, 0xd1, 0x16, 0x01, 0x1c, 0x9f, 0xff, // PrevBlock
				0xb9, 0xc3, 0x9c, 0x20, 0xb4, 0xba, 0xf3, 0xb8,
				0x4c, 0xd7, 0x23, 0x49, 0x30, 0x08, 0x33, 0x08,
				0x6d, 0x68, 0x6d, 0xb1, 0x42, 0xc5, 0x1b, 0x26,
				0xa1, 0x96, 0xcd, 0xab, 0xbe, 0x7a, 0x06, 0x10, // MerkleRoot
				0xd4, 0x5a, 0xc3, 0x52, // Timestamp
				0xf0, 0xff, 0x0f, 0x1e, // Bits
				0x00, 0x10, 0xbb, 0x75, // Nonce
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, // TxnCount
			}, pver, BaseEncoding, &MessageError{},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from wire format.
		var msg MsgBlock
		r := bytes.NewReader(test.buf)
		err := msg.BtcDecode(r, test.pver, test.enc)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("BtcDecode #%d wrong error got: %v, want: %v",
				i, err, reflect.TypeOf(test.err))
			continue
		}

		// Deserialize from wire format.
		r = bytes.NewReader(test.buf)
		err = msg.Deserialize(r)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Deserialize #%d wrong error got: %v, want: %v",
				i, err, reflect.TypeOf(test.err))
			continue
		}

		// Deserialize with transaction location info from wire format.
		br := bytes.NewBuffer(test.buf)
		_, err = msg.DeserializeTxLoc(br)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("DeserializeTxLoc #%d wrong error got: %v, "+
				"want: %v", i, err, reflect.TypeOf(test.err))
			continue
		}
	}
}

// TestBlockSerializeSize performs tests to ensure the serialize size for
// various blocks is accurate.
func TestBlockSerializeSize(t *testing.T) {
	// Block with no transactions.
	noTxBlock := NewMsgBlock(&blockOne.Header)

	tests := []struct {
		in   *MsgBlock // Block to encode
		size int       // Expected serialized size
	}{
		// Block with no transactions.
		{noTxBlock, 81},

		// First block in the mainnet block chain.
		{&blockOne, len(blockOneBytes)},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		serializedSize := test.in.SerializeSize()
		if serializedSize != test.size {
			t.Errorf("MsgBlock.SerializeSize: #%d got: %d, want: "+
				"%d", i, serializedSize, test.size)
			continue
		}
	}
}

// blockOne is the first block in the mainnet block chain.
var blockOne = MsgBlock{
	Header: BlockHeader{
		Version: 2,
		PrevBlock: chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
			0xb6, 0x8b, 0x8c, 0x41, 0x0d, 0x2e, 0xa4, 0xaf,
			0xd7, 0x4f, 0xb5, 0x6e, 0x37, 0x0b, 0xfc, 0x1b,
			0xed, 0xf9, 0x29, 0xe1, 0x45, 0x38, 0x96, 0xc9,
			0xe7, 0x9d, 0xd1, 0x16, 0x01, 0x1c, 0x9f, 0xff,
		}),
		MerkleRoot: chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
			0xb9, 0xc3, 0x9c, 0x20, 0xb4, 0xba, 0xf3, 0xb8,
			0x4c, 0xd7, 0x23, 0x49, 0x30, 0x08, 0x33, 0x08,
			0x6d, 0x68, 0x6d, 0xb1, 0x42, 0xc5, 0x1b, 0x26,
			0xa1, 0x96, 0xcd, 0xab, 0xbe, 0x7a, 0x06, 0x10,
		}),

		Timestamp: time.Unix(0x52c35ad4, 0), // 2014-01-01 09:01:24 +0900
		Bits:      0x1e0ffff0,
		Nonce:     0x75bb1000, // 1975193600
	},
	Transactions: []*MsgTx{
		{
			Version: 1,
			TxIn: []*TxIn{
				{
					PreviousOutPoint: OutPoint{
						Hash:  chainhash.Hash{},
						Index: 0xffffffff,
					},
					SignatureScript: []byte{
						0x51, 0x01, 0x01, 0x06, 0x2f,
						0x50, 0x32, 0x53, 0x48, 0x2f,
					},
					Sequence: 0xffffffff,
				},
			},
			TxOut: []*TxOut{
				{
					Value: 0x12a05f200,
					PkScript: []byte{
						0x21, // OP_DATA_33
						0x03, 0x60, 0xbc, 0x5a, 0xd9, 0x66, 0xa9, 0xa8,
						0xe6, 0xb0, 0x06, 0x1b, 0x25, 0xd8, 0xab, 0x0c,
						0x7a, 0x90, 0xd9, 0x39, 0xb5, 0x21, 0x0b, 0xf0,
						0x85, 0x5f, 0x43, 0x1a, 0x32, 0x3c, 0x1d, 0xaa,
						0x11, // 33-byte signature
						0xac, // OP_CHECKSIG
					},
				},
			},
			LockTime: 0,
		},
	},
}

// Block one serialized bytes.
var blockOneBytes = []byte{
	0x02, 0x00, 0x00, 0x00, // Version 1
	0xb6, 0x8b, 0x8c, 0x41, 0x0d, 0x2e, 0xa4, 0xaf,
	0xd7, 0x4f, 0xb5, 0x6e, 0x37, 0x0b, 0xfc, 0x1b,
	0xed, 0xf9, 0x29, 0xe1, 0x45, 0x38, 0x96, 0xc9,
	0xe7, 0x9d, 0xd1, 0x16, 0x01, 0x1c, 0x9f, 0xff, // PrevBlock
	0xb9, 0xc3, 0x9c, 0x20, 0xb4, 0xba, 0xf3, 0xb8,
	0x4c, 0xd7, 0x23, 0x49, 0x30, 0x08, 0x33, 0x08,
	0x6d, 0x68, 0x6d, 0xb1, 0x42, 0xc5, 0x1b, 0x26,
	0xa1, 0x96, 0xcd, 0xab, 0xbe, 0x7a, 0x06, 0x10, // MerkleRoot
	0xd4, 0x5a, 0xc3, 0x52, // Timestamp
	0xf0, 0xff, 0x0f, 0x1e, // Bits
	0x00, 0x10, 0xbb, 0x75, // Nonce
	0x01,                   // TxnCount
	0x01, 0x00, 0x00, 0x00, // Version
	0x01, // Varint for number of transaction inputs
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Previous output hash
	0xff, 0xff, 0xff, 0xff, // Prevous output index
	0x0a,                                                       // Varint for length of signature script
	0x51, 0x01, 0x01, 0x06, 0x2f, 0x50, 0x32, 0x53, 0x48, 0x2f, // Signature script (coinbase)
	0xff, 0xff, 0xff, 0xff, // Sequence
	0x01,                                           // Varint for number of transaction outputs
	0x00, 0xf2, 0x05, 0x2a, 0x01, 0x00, 0x00, 0x00, // Transaction amount
	0x23, // Varint for length of pk script
	0x21, // OP_DATA_33
	0x03, 0x60, 0xbc, 0x5a, 0xd9, 0x66, 0xa9, 0xa8,
	0xe6, 0xb0, 0x06, 0x1b, 0x25, 0xd8, 0xab, 0x0c,
	0x7a, 0x90, 0xd9, 0x39, 0xb5, 0x21, 0x0b, 0xf0,
	0x85, 0x5f, 0x43, 0x1a, 0x32, 0x3c, 0x1d, 0xaa,
	0x11,                   // 33-byte uncompressed public key
	0xac,                   // OP_CHECKSIG
	0x00, 0x00, 0x00, 0x00, // Lock time
}

// Transaction location information for block one transactions.
var blockOneTxLocs = []TxLoc{
	{TxStart: 81, TxLen: 105},
}
