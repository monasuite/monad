// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/monasuite/monad/chaincfg/chainhash"
	"github.com/monasuite/monad/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xFF, 0xFF, 0x00, 0x1D, 0x01, 0x04, 0x4C, /* |.....LVD| */
				0x56, 0x44, 0x65, 0x63, 0x2E, 0x20, 0x33, 0x31, /* |ec. 31th| */
				0x74, 0x68, 0x20, 0x32, 0x30, 0x31, 0x33, 0x20, /* | 2013 Ja| */
				0x4A, 0x61, 0x70, 0x61, 0x6E, 0x2C, 0x20, 0x54, /* |pan, The| */
				0x68, 0x65, 0x20, 0x77, 0x69, 0x6E, 0x6E, 0x69, /* |winning | */
				0x6E, 0x67, 0x20, 0x6E, 0x75, 0x6D, 0x62, 0x65, /* |numbers | */
				0x72, 0x73, 0x20, 0x6F, 0x66, 0x20, 0x74, 0x68, /* |of the 2| */
				0x65, 0x20, 0x32, 0x30, 0x31, 0x33, 0x20, 0x59, /* |013 Year| */
				0x65, 0x61, 0x72, 0x2D, 0x45, 0x6E, 0x64, 0x20, /* |-End Jum| */
				0x4A, 0x75, 0x6D, 0x62, 0x6F, 0x20, 0x4C, 0x6F, /* |bo Lotte| */
				0x74, 0x74, 0x65, 0x72, 0x79, 0x3A, 0x32, 0x33, /* |ry:23-13| */
				0x2D, 0x31, 0x33, 0x30, 0x39, 0x31, 0x36, /* |130916| */
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x41, 0x04, 0x01, 0x84, 0x71, 0x0f, 0xa6, 0x89,
				0xad, 0x50, 0x23, 0x69, 0x0c, 0x80, 0xf3, 0xa4,
				0x9c, 0x8f, 0x13, 0xf8, 0xd4, 0x5b, 0x8c, 0x85,
				0x7f, 0xbc, 0xbc, 0x8b, 0xc4, 0xa8, 0xe4, 0xd3,
				0xeb, 0x4b, 0x10, 0xf4, 0xd4, 0x60, 0x4f, 0xa0,
				0x8d, 0xce, 0x60, 0x1a, 0xaf, 0x0f, 0x47, 0x02,
				0x16, 0xfe, 0x1b, 0x51, 0x85, 0x0b, 0x4a, 0xcf,
				0x21, 0xb1, 0x79, 0xc4, 0x50, 0x70, 0xac, 0x7b,
				0x03, 0xa9, 0xac,
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xb6, 0x8b, 0x8c, 0x41, 0x0d, 0x2e, 0xa4, 0xaf,
	0xd7, 0x4f, 0xb5, 0x6e, 0x37, 0x0b, 0xfc, 0x1b,
	0xed, 0xf9, 0x29, 0xe1, 0x45, 0x38, 0x96, 0xc9,
	0xe7, 0x9d, 0xd1, 0x16, 0x01, 0x1c, 0x9f, 0xff,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xa6, 0x4b, 0xac, 0x07, 0xfe, 0x31, 0x87, 0x7f,
	0x31, 0xd0, 0x32, 0x52, 0x95, 0x3b, 0x3c, 0x32,
	0x39, 0x89, 0x33, 0xaf, 0x7a, 0x72, 0x41, 0x19,
	0xbc, 0x4d, 0x6f, 0xa4, 0xa8, 0x05, 0xe4, 0x35,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // 35e405a8a46f4dbc1941727aaf338939323c3b955232d0317f8731fe07ac4ba6
		Timestamp:  time.Unix(0x52c283f0, 0), // 2013-12-31 08:44:32 +0000 UTC
		Bits:       0x1e0ffff0,               // 504365040
		Nonce:      0x0012d666,               // 1234534
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x9c, 0xa2, 0x90, 0xa5, 0x3a, 0x47, 0x49, 0xa7,
	0xc8, 0x72, 0x30, 0x6f, 0x3b, 0x25, 0x35, 0x4a,
	0x07, 0x4c, 0x06, 0xfc, 0xc2, 0xbe, 0x5e, 0x9a,
	0xb2, 0xcd, 0x2f, 0x7c, 0x9d, 0xa6, 0x43, 0x75,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 35e405a8a46f4dbc1941727aaf338939323c3b955232d0317f8731fe07ac4ba6
		Timestamp:  time.Unix(0x4d49e5da, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      1,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet4GenesisHash is the hash of the first block in the block chain for the
// test network (version 4).
var testNet4GenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xb2, 0xe0, 0x61, 0x10, 0x32, 0x9c, 0x44, 0x8f,
	0x15, 0x78, 0xe4, 0x8a, 0x25, 0xa8, 0x8b, 0x63,
	0x9d, 0xcf, 0xaa, 0xa6, 0xa6, 0xb2, 0x97, 0xd0,
	0xc6, 0xe0, 0x3b, 0xba, 0xce, 0x06, 0xb1, 0xa2,
})

// testNet4GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 4).  It is the same as the merkle root
// for the main network.
var testNet4GenesisMerkleRoot = genesisMerkleRoot

// testNet4GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNet4GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet4GenesisMerkleRoot, // 35e405a8a46f4dbc1941727aaf338939323c3b955232d0317f8731fe07ac4ba6
		Timestamp:  time.Unix(0x58bf2dec, 0),  // 2017-03-07 22:02:20 +0000 UTC
		Bits:       0x1e0ffff0,                // 504365040
		Nonce:      0x0020646c,                // 2122860
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xf6, 0x7a, 0xd7, 0x69, 0x5d, 0x9b, 0x66, 0x2a,
	0x72, 0xff, 0x3d, 0x8e, 0xdb, 0xbb, 0x2d, 0xe0,
	0xbf, 0xa6, 0x7b, 0x13, 0x97, 0x4b, 0xb9, 0x91,
	0x0d, 0x11, 0x6d, 0x5c, 0xbd, 0x86, 0x3e, 0x68,
})

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
	0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
	0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
	0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a,
})

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 35e405a8a46f4dbc1941727aaf338939323c3b955232d0317f8731fe07ac4ba6
		Timestamp:  time.Unix(1401292357, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// sigNetGenesisHash is the hash of the first block in the block chain for the
// signet test network.
var sigNetGenesisHash = chainhash.Hash{
	0x42, 0x4e, 0x20, 0x12, 0xf3, 0x1b, 0x58, 0x41,
	0xe8, 0x5b, 0x07, 0x54, 0x63, 0x14, 0x59, 0xa8,
	0x99, 0xc8, 0xe6, 0x15, 0x9e, 0x5e, 0x00, 0xb7,
	0xd2, 0x34, 0xd5, 0xfb, 0xc4, 0x6a, 0xc5, 0xc1,
}

// sigNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the signet test network. It is the same as the merkle root for
// the main network.
var sigNetGenesisMerkleRoot = genesisMerkleRoot

// sigNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the signet test network.
var sigNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: sigNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1598918400, 0), // 2020-09-01 00:00:00 +0000 UTC
		Bits:       0x1e0377ae,               // 503543726 [00000377ae000000000000000000000000000000000000000000000000000000]
		Nonce:      52613770,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
