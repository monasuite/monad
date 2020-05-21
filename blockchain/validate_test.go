// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/monasuite/monad/chaincfg"
	"github.com/monasuite/monad/chaincfg/chainhash"
	"github.com/monasuite/monad/wire"
	"github.com/monasuite/monautil"
)

// TestSequenceLocksActive tests the SequenceLockActive function to ensure it
// works as expected in all possible combinations/scenarios.
func TestSequenceLocksActive(t *testing.T) {
	seqLock := func(h int32, s int64) *SequenceLock {
		return &SequenceLock{
			Seconds:     s,
			BlockHeight: h,
		}
	}

	tests := []struct {
		seqLock     *SequenceLock
		blockHeight int32
		mtp         time.Time

		want bool
	}{
		// Block based sequence lock with equal block height.
		{seqLock: seqLock(1000, -1), blockHeight: 1001, mtp: time.Unix(9, 0), want: true},

		// Time based sequence lock with mtp past the absolute time.
		{seqLock: seqLock(-1, 30), blockHeight: 2, mtp: time.Unix(31, 0), want: true},

		// Block based sequence lock with current height below seq lock block height.
		{seqLock: seqLock(1000, -1), blockHeight: 90, mtp: time.Unix(9, 0), want: false},

		// Time based sequence lock with current time before lock time.
		{seqLock: seqLock(-1, 30), blockHeight: 2, mtp: time.Unix(29, 0), want: false},

		// Block based sequence lock at the same height, so shouldn't yet be active.
		{seqLock: seqLock(1000, -1), blockHeight: 1000, mtp: time.Unix(9, 0), want: false},

		// Time based sequence lock with current time equal to lock time, so shouldn't yet be active.
		{seqLock: seqLock(-1, 30), blockHeight: 2, mtp: time.Unix(30, 0), want: false},
	}

	t.Logf("Running %d sequence locks tests", len(tests))
	for i, test := range tests {
		got := SequenceLockActive(test.seqLock,
			test.blockHeight, test.mtp)
		if got != test.want {
			t.Fatalf("SequenceLockActive #%d got %v want %v", i,
				got, test.want)
		}
	}
}

// TestCheckConnectBlockTemplate tests the CheckConnectBlockTemplate function to
// ensure it fails.
func TestCheckConnectBlockTemplate(t *testing.T) {
	// Create a new database and chain instance to run tests against.
	chain, teardownFunc, err := chainSetup("checkconnectblocktemplate",
		&chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("Failed to setup chain instance: %v", err)
		return
	}
	defer teardownFunc()

	// Since we're not dealing with the real block chain, set the coinbase
	// maturity to 1.
	chain.TstSetCoinbaseMaturity(1)

	// Load up blocks such that there is a side chain.
	// (genesis block) -> 1 -> 2 -> 3 -> 4
	//                          \-> 3a
	testFiles := []string{
		"blk_0_to_4.dat.bz2",
		"blk_3A.dat.bz2",
	}

	var blocks []*monautil.Block
	for _, file := range testFiles {
		blockTmp, err := loadBlocks(file)
		if err != nil {
			t.Fatalf("Error loading file: %v\n", err)
		}
		blocks = append(blocks, blockTmp...)
	}

	for i := 1; i <= 3; i++ {
		isMainChain, _, err := chain.ProcessBlock(blocks[i], BFNone)
		if err != nil {
			t.Fatalf("CheckConnectBlockTemplate: Received unexpected error "+
				"processing block %d: %v", i, err)
		}
		if !isMainChain {
			t.Fatalf("CheckConnectBlockTemplate: Expected block %d to connect "+
				"to main chain", i)
		}
	}

	// Block 3 should fail to connect since it's already inserted.
	err = chain.CheckConnectBlockTemplate(blocks[3])
	if err == nil {
		t.Fatal("CheckConnectBlockTemplate: Did not received expected error " +
			"on block 3")
	}

	// Block 4 should connect successfully to tip of chain.
	err = chain.CheckConnectBlockTemplate(blocks[4])
	if err != nil {
		t.Fatalf("CheckConnectBlockTemplate: Received unexpected error on "+
			"block 4: %v", err)
	}

	// Block 3a should fail to connect since does not build on chain tip.
	err = chain.CheckConnectBlockTemplate(blocks[5])
	if err == nil {
		t.Fatal("CheckConnectBlockTemplate: Did not received expected error " +
			"on block 3a")
	}

	// Block 4 should connect even if proof of work is invalid.
	invalidPowBlock := *blocks[4].MsgBlock()
	invalidPowBlock.Header.Nonce++
	err = chain.CheckConnectBlockTemplate(monautil.NewBlock(&invalidPowBlock))
	if err != nil {
		t.Fatalf("CheckConnectBlockTemplate: Received unexpected error on "+
			"block 4 with bad nonce: %v", err)
	}

	// Invalid block building on chain tip should fail to connect.
	invalidBlock := *blocks[4].MsgBlock()
	invalidBlock.Header.Bits--
	err = chain.CheckConnectBlockTemplate(monautil.NewBlock(&invalidBlock))
	if err == nil {
		t.Fatal("CheckConnectBlockTemplate: Did not received expected error " +
			"on block 4 with invalid difficulty bits")
	}
}

// TestCheckBlockSanity tests the CheckBlockSanity function to ensure it works
// as expected.
func TestCheckBlockSanity(t *testing.T) {
	powLimit := chaincfg.MainNetParams.PowLimit
	block := monautil.NewBlock(&Block1000000)
	timeSource := NewMedianTime()
	err := CheckBlockSanity(block, powLimit, timeSource)
	if err != nil {
		t.Errorf("CheckBlockSanity: %v", err)
	}

	// Ensure a block that has a timestamp with a precision higher than one
	// second fails.
	timestamp := block.MsgBlock().Header.Timestamp
	block.MsgBlock().Header.Timestamp = timestamp.Add(time.Nanosecond)
	err = CheckBlockSanity(block, powLimit, timeSource)
	if err == nil {
		t.Errorf("CheckBlockSanity: error is nil when it shouldn't be")
	}
}

// TestCheckSerializedHeight tests the checkSerializedHeight function with
// various serialized heights and also does negative tests to ensure errors
// and handled properly.
func TestCheckSerializedHeight(t *testing.T) {
	// Create an empty coinbase template to be used in the tests below.
	coinbaseOutpoint := wire.NewOutPoint(&chainhash.Hash{}, math.MaxUint32)
	coinbaseTx := wire.NewMsgTx(1)
	coinbaseTx.AddTxIn(wire.NewTxIn(coinbaseOutpoint, nil, nil))

	// Expected rule errors.
	missingHeightError := RuleError{
		ErrorCode: ErrMissingCoinbaseHeight,
	}
	badHeightError := RuleError{
		ErrorCode: ErrBadCoinbaseHeight,
	}

	tests := []struct {
		sigScript  []byte // Serialized data
		wantHeight int32  // Expected height
		err        error  // Expected error type
	}{
		// No serialized height length.
		{[]byte{}, 0, missingHeightError},
		// Serialized height length with no height bytes.
		{[]byte{0x02}, 0, missingHeightError},
		// Serialized height length with too few height bytes.
		{[]byte{0x02, 0x4a}, 0, missingHeightError},
		// Serialized height that needs 2 bytes to encode.
		{[]byte{0x02, 0x4a, 0x52}, 21066, nil},
		// Serialized height that needs 2 bytes to encode, but backwards
		// endianness.
		{[]byte{0x02, 0x4a, 0x52}, 19026, badHeightError},
		// Serialized height that needs 3 bytes to encode.
		{[]byte{0x03, 0x40, 0x0d, 0x03}, 200000, nil},
		// Serialized height that needs 3 bytes to encode, but backwards
		// endianness.
		{[]byte{0x03, 0x40, 0x0d, 0x03}, 1074594560, badHeightError},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		msgTx := coinbaseTx.Copy()
		msgTx.TxIn[0].SignatureScript = test.sigScript
		tx := monautil.NewTx(msgTx)

		err := checkSerializedHeight(tx, test.wantHeight)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("checkSerializedHeight #%d wrong error type "+
				"got: %v <%T>, want: %T", i, err, err, test.err)
			continue
		}

		if rerr, ok := err.(RuleError); ok {
			trerr := test.err.(RuleError)
			if rerr.ErrorCode != trerr.ErrorCode {
				t.Errorf("checkSerializedHeight #%d wrong "+
					"error code got: %v, want: %v", i,
					rerr.ErrorCode, trerr.ErrorCode)
				continue
			}
		}
	}
}

// Block1000000 defines block 1,000,000 of the block chain.  It is used to
// test Block operations.
var Block1000000 = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version: 536870912,
		PrevBlock: chainhash.Hash([32]byte{ // Make go vet happy.
			0x09, 0x11, 0xdf, 0x29, 0x27, 0x4b, 0x6c, 0x4d,
			0x9f, 0x16, 0xc9, 0x2b, 0xeb, 0x05, 0x66, 0xa1,
			0x98, 0xb8, 0x7a, 0x51, 0x66, 0x5c, 0x97, 0xc5,
			0xec, 0xfc, 0xde, 0x38, 0x61, 0x9e, 0xf4, 0xbb,
		}), // bbf49e6138defcecc5975c66517ab898a16605eb2bc9169f4d6c4b2729df1109
		MerkleRoot: chainhash.Hash([32]byte{ // Make go vet happy.
			0x8e, 0x16, 0x9c, 0x69, 0x99, 0x36, 0xe8, 0x3b,
			0xcc, 0x35, 0xaa, 0xab, 0xda, 0x7b, 0x2a, 0x96,
			0xa5, 0xc9, 0x14, 0x81, 0x4a, 0x92, 0x60, 0x6a,
			0x99, 0x1e, 0x1d, 0x5c, 0x72, 0xcb, 0x86, 0x66,
		}), // 6686cb725c1d1e996a60924a8114c9a5962a7bdaabaa35cc3be83699699c168e
		Timestamp: time.Unix(1495380733, 0), // 2017-05-22 00:32:13 JST
		Bits:      0x1b18a61e,               // 454600222
		Nonce:     0x0594c0af,               // 93634735
	},
	Transactions: []*wire.MsgTx{
		{
			Version: 1,
			TxIn: []*wire.TxIn{
				{
					PreviousOutPoint: wire.OutPoint{
						Hash:  chainhash.Hash{},
						Index: 0xffffffff,
					},
					SignatureScript: []byte{
						0x03, 0x40, 0x42, 0x0f, 0x04, 0xfd, 0xb2, 0x21,
						0x59, 0x08, 0x68, 0x00, 0x01, 0x37, 0x6c, 0x0a,
						0x00, 0x00, 0x0d, 0x2f, 0x6e, 0x6f, 0x64, 0x65,
						0x53, 0x74, 0x72, 0x61, 0x74, 0x75, 0x6d, 0x2f,
					},
					Sequence: 0x00,
				},
			},
			TxOut: []*wire.TxOut{
				{
					Value: 0x12a05f200, // 5000000000
					PkScript: []byte{
						0x76, // OP_DATA_118
						0xa9, 0x14, 0x90, 0x7e, 0xf7, 0x62, 0x54, 0x8d,
						0xd3, 0xf4, 0xe4, 0xad, 0xf6, 0x4e, 0x3a, 0x62,
						0x58, 0xc1, 0x58, 0x41, 0x71, 0x62,
						0x88, // 118-byte signature
						0xac, // OP_CHECKSIG
					},
				},
			},
			LockTime: 0,
		},
	},
}
