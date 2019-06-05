// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/monasuite/monad/chaincfg/chainhash"
	"github.com/monasuite/monad/wire"
)

// These variables are the chain proof-of-work limit parameters for each default
// network.
var (
	// bigOne is 1 represented as a big.Int.  It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)

	// mainPowLimit is the highest proof of work value a Bitcoin block can
	// have for the main network.  It is the value 2^224 - 1.
	//mainPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	mainPowLimit, _ = new(big.Int).SetString("0x0fffff000000000000000000000000000000000000000000000000000000", 0)

	// regressionPowLimit is the highest proof of work value a Bitcoin block
	// can have for the regression test network.  It is the value 2^255 - 1.
	regressionPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)

	// testNet4PowLimit is the highest proof of work value a Bitcoin block
	// can have for the test network (version 3).  It is the value
	// 2^224 - 1.
	//testNet4PowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
	testNet4PowLimit, _ = new(big.Int).SetString("0x0ffffff00000000000000000000000000000000000000000000000000000", 0)

	// simNetPowLimit is the highest proof of work value a Bitcoin block
	// can have for the simulation test network.  It is the value 2^255 - 1.
	simNetPowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 255), bigOne)
)

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
//
// Each checkpoint is selected based upon several factors.  See the
// documentation for blockchain.IsCheckpointCandidate for details on the
// selection criteria.
type Checkpoint struct {
	Height int32
	Hash   *chainhash.Hash
}

// DNSSeed identifies a DNS seed.
type DNSSeed struct {
	// Host defines the hostname of the seed.
	Host string

	// HasFiltering defines whether the seed supports filtering
	// by service flags (wire.ServiceFlag).
	HasFiltering bool
}

// ConsensusDeployment defines details related to a specific consensus rule
// change that is voted in.  This is part of BIP0009.
type ConsensusDeployment struct {
	// BitNumber defines the specific bit number within the block version
	// this particular soft-fork deployment refers to.
	BitNumber uint8

	// StartTime is the median block time after which voting on the
	// deployment starts.
	StartTime uint64

	// ExpireTime is the median block time after which the attempted
	// deployment expires.
	ExpireTime uint64
}

// Constants that define the deployment offset in the deployments field of the
// parameters for each deployment.  This is useful to be able to get the details
// of a specific deployment by name.
const (
	// DeploymentTestDummy defines the rule change deployment ID for testing
	// purposes.
	DeploymentTestDummy = iota

	// DeploymentCSV defines the rule change deployment ID for the CSV
	// soft-fork package. The CSV package includes the deployment of BIPS
	// 68, 112, and 113.
	DeploymentCSV

	// DeploymentSegwit defines the rule change deployment ID for the
	// Segregated Witness (segwit) soft-fork package. The segwit package
	// includes the deployment of BIPS 141, 142, 144, 145, 147 and 173.
	DeploymentSegwit

	// NOTE: DefinedDeployments must always come last since it is used to
	// determine how many defined deployments there currently are.

	// DefinedDeployments is the number of currently defined deployments.
	DefinedDeployments
)

// Params defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name string

	// Net defines the magic bytes used to identify the network.
	Net wire.BitcoinNet

	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []DNSSeed

	// GenesisBlock defines the first block of the chain.
	GenesisBlock *wire.MsgBlock

	// GenesisHash is the starting block hash.
	GenesisHash *chainhash.Hash

	// PowLimit defines the highest allowed proof of work value for a block
	// as a uint256.
	PowLimit *big.Int

	// PowLimitBits defines the highest allowed proof of work value for a
	// block in compact form.
	PowLimitBits uint32

	// These fields define the block heights at which the specified softfork
	// BIP became active.
	BIP0034Height int32
	BIP0065Height int32
	BIP0066Height int32

	// CoinbaseMaturity is the number of blocks required before newly mined
	// coins (coinbase transactions) can be spent.
	CoinbaseMaturity uint16

	// SubsidyReductionInterval is the interval of blocks before the subsidy
	// is reduced.
	SubsidyReductionInterval int32

	// TargetTimespan is the desired amount of time that should elapse
	// before the block difficulty requirement is examined to determine how
	// it should be changed in order to maintain the desired block
	// generation rate.
	TargetTimespan time.Duration

	// TargetTimePerBlock is the desired amount of time to generate each
	// block.
	TargetTimePerBlock time.Duration

	// RetargetAdjustmentFactor is the adjustment factor used to limit
	// the minimum and maximum amount of adjustment that can occur between
	// difficulty retargets.
	RetargetAdjustmentFactor int64

	// ReduceMinDifficulty defines whether the network should reduce the
	// minimum required difficulty after a long enough period of time has
	// passed without finding a block.  This is really only useful for test
	// networks and should not be set on a main network.
	ReduceMinDifficulty bool

	// MinDiffReductionTime is the amount of time after which the minimum
	// required difficulty should be reduced when a block hasn't been found.
	//
	// NOTE: This only applies if ReduceMinDifficulty is true.
	MinDiffReductionTime time.Duration

	// GenerateSupported specifies whether or not CPU mining is allowed.
	GenerateSupported bool

	// Checkpoints ordered from oldest to newest.
	Checkpoints []Checkpoint

	// These fields are related to voting on consensus rule changes as
	// defined by BIP0009.
	//
	// RuleChangeActivationThreshold is the number of blocks in a threshold
	// state retarget window for which a positive vote for a rule change
	// must be cast in order to lock in a rule change. It should typically
	// be 95% for the main network and 75% for test networks.
	//
	// MinerConfirmationWindow is the number of blocks in each threshold
	// state retarget window.
	//
	// Deployments define the specific consensus rule changes to be voted
	// on.
	RuleChangeActivationThreshold uint32
	MinerConfirmationWindow       uint32
	Deployments                   [DefinedDeployments]ConsensusDeployment

	// Mempool parameters
	RelayNonStdTxs bool

	// Human-readable part for Bech32 encoded segwit addresses, as defined
	// in BIP 173.
	Bech32HRPSegwit string

	// Address encoding magics
	PubKeyHashAddrID        byte // First byte of a P2PKH address
	ScriptHashAddrID        byte // First byte of a P2SH address
	PrivateKeyID            byte // First byte of a WIF private key
	WitnessPubKeyHashAddrID byte // First byte of a P2WPKH address
	WitnessScriptHashAddrID byte // First byte of a P2WSH address

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID [4]byte
	HDPublicKeyID  [4]byte

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType uint32

	// vAlertPubKey is used checkpoint's deliverty
	AlertPubMainKey []byte
	AlertPubSubKey  []byte

	// DGWv3's hardfork height
	DGWv3Height int32

	// Lyra2re2's hardfork height
	Lyra2re2Height int32
}

// MainNetParams defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "9401",
	DNSSeeds: []DNSSeed{
		{"monacoin.org", false},
		{"dnsseed.monacoin.org", false},
		{"dnsseed.tamami-foundation.org", true},
	},

	// Chain parameters
	GenesisBlock:             &genesisBlock,
	GenesisHash:              &genesisHash,
	PowLimit:                 mainPowLimit,
	PowLimitBits:             0x1e0ffff0,
	BIP0034Height:            0,
	BIP0065Height:            977759,
	BIP0066Height:            977759,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 1051200,
	TargetTimespan:           time.Second * 95040, // 1.1 weeks
	TargetTimePerBlock:       time.Second * 90,    // 90 seconds
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{1500, newHashFromStr("9f42d51d18d0a8914a00664c433a0ca4be3eed02f9374d790bffbd3d3053d41d")},
		{4000, newHashFromStr("2c60edac7d9f44d90d1e218af2a8085e78b735185c5bf42f9fe9dbd0e604c97b")},
		{8000, newHashFromStr("61d4d053b1a4c6deb4c7e806cedd876f25b51da6c51b209109579c7b9892e5c2")},
		{16000, newHashFromStr("3c4a8887bb3ae0599abfefe765f7c911fbfe98b3f23d7f70b05bf49cf62ebdaf")},
		{32000, newHashFromStr("c0703986c1c6a9052478db5e52432e5a1e55d6b6362b85f0ffdbb61ce3311b77")},
		{58700, newHashFromStr("a9c5d9878864b77ba52b068787b83ce2fcf526c5899f40af51c9d441eeb4c84d")},
		{80000, newHashFromStr("c99b83da7328b58251d16f4646da222b0280f180bd208efa5e3256c9eb6ea2be")},
		{115000, newHashFromStr("75e642c003e5bd748b679472e981b7b2f81f344b3f197029f84470256cef33e4")},
		{189250, newHashFromStr("1bea3d5c25a8097eef2e70ece4beb6c502b895fe00056552948309beb3497c99")},
		{300000, newHashFromStr("11095515590421444ba29396d9122c234baced79be8b32604acc37cf094558ab")},
		{444000, newHashFromStr("3ed05516cdce4db93b135189592c7e2b37d768f99a1819a1d2ea3a8e5b8439a8")},
		{904000, newHashFromStr("353f5b7f9440e1d830bd1c265c69fb0e7c7988e343b2202a704406d04a8cd02e")},
		{1045000, newHashFromStr("562593372b8cee7d2eac325ef8874532d93cc031a75a00cec0086d2d943f40dc")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 7560,  // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       10080, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1488931200, // March 8, 2017
			ExpireTime: 1520467200, // March 8, 2018
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1488931200, // March 8, 2017
			ExpireTime: 1520467200, // March 8, 2018
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "mona", // always mona for main net

	// Address encoding magics
	PubKeyHashAddrID:        0x32, // starts with M
	ScriptHashAddrID:        0x37, // starts with P
	PrivateKeyID:            0xB0,
	WitnessPubKeyHashAddrID: 0x06, // starts with p2
	WitnessScriptHashAddrID: 0x0A, // starts with 7Xh

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 22,

	// vAlertPubKey is used checkpoint's deliverty
	AlertPubMainKey: []byte{
		0x04, 0xfc, 0x55, 0xd9, 0x19, 0xcd, 0xea, 0x21,
		0xe8, 0x17, 0x1f, 0xea, 0x24, 0xc3, 0xcf, 0x23,
		0x86, 0x4f, 0xf6, 0x4a, 0x53, 0x40, 0x48, 0x29,
		0xad, 0x50, 0xaf, 0x86, 0xe1, 0xbe, 0x1b, 0x82,
		0x17, 0x11, 0x57, 0x01, 0xb3, 0x48, 0xd5, 0x0c,
		0x6a, 0xab, 0xa6, 0x98, 0x3b, 0xc1, 0x48, 0xd3,
		0xc9, 0xe6, 0xfa, 0x8c, 0x11, 0x36, 0x58, 0x89,
		0x77, 0x4f, 0xc1, 0xdb, 0x6d, 0xa6, 0x84, 0x0c,
		0x06,
	},
	AlertPubSubKey: []byte{
		0x04, 0x23, 0x67, 0xfa, 0x06, 0xa8, 0x24, 0x5f,
		0x17, 0xc2, 0x13, 0x7f, 0x69, 0x46, 0xe6, 0x08,
		0x4a, 0x22, 0x51, 0xb5, 0xa6, 0xa7, 0x12, 0x33,
		0x80, 0xcd, 0x0f, 0xf0, 0xd5, 0x82, 0xd8, 0x2e,
		0xc3, 0x2d, 0x93, 0x9e, 0xd1, 0x1c, 0x77, 0x74,
		0xe9, 0x49, 0x63, 0x7d, 0xc4, 0x79, 0x50, 0x81,
		0x7a, 0x80, 0x06, 0x0a, 0x99, 0xfe, 0x7e, 0x51,
		0x6c, 0x0a, 0xf4, 0x3f, 0xa5, 0xa1, 0xa1, 0xb6,
		0x82,
	},

	// DGWv3's HF height. Used to calculate target.
	DGWv3Height: 450025,

	// Lyra2re2's HF height. Used to calculate Pow.
	Lyra2re2Height: 450000,
}

// RegressionNetParams defines the network parameters for the regression test
// Bitcoin network.  Not to be confused with the test Bitcoin network (version
// 3), this network is sometimes simply called "testnet".
var RegressionNetParams = Params{
	Name:        "regtest",
	Net:         wire.TestNet,
	DefaultPort: "18444",
	DNSSeeds:    []DNSSeed{},

	// Chain parameters
	GenesisBlock:             &regTestGenesisBlock,
	GenesisHash:              &regTestGenesisHash,
	PowLimit:                 regressionPowLimit,
	PowLimitBits:             0x207fffff,
	CoinbaseMaturity:         100,
	BIP0034Height:            100000000, // Not active - Permit ver 1 blocks
	BIP0065Height:            1351,      // Used by regression tests
	BIP0066Height:            1251,      // Used by regression tests
	SubsidyReductionInterval: 150,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 108, // 75%  of MinerConfirmationWindow
	MinerConfirmationWindow:       144,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "rmona", // always tmona for test net

	// Address encoding magics
	PubKeyHashAddrID: 0x6f, // starts with m or n
	ScriptHashAddrID: 0xc4, // starts with 2
	PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,

	// vAlertPubKey is used checkpoint's deliverty
	AlertPubMainKey: []byte{},
	AlertPubSubKey:  []byte{},

	// DGWv3's HF height. Used to calculate target.
	// regtest only works up to the height specified in DGWv3Height.
	// Please increase if you don't have enough value.
	DGWv3Height: 1000085,

	// Lyra2re2's HF height. Used to calculate PoW.
	Lyra2re2Height: 60,
}

// TestNet4Params defines the network parameters for the test Bitcoin network
// (version 4).  Not to be confused with the regression test network, this
// network is sometimes simply called "testnet".
var TestNet4Params = Params{
	Name:        "testnet4",
	Net:         wire.TestNet4,
	DefaultPort: "19403",
	DNSSeeds: []DNSSeed{
		{"testnet-dnsseed.monacoin.org", true},
		{"electrumx1.testnet.monacoin.nl", true},
		{"electrumx1.testnet.monacoin.ninja", true},
	},

	// Chain parameters
	GenesisBlock:             &testNet4GenesisBlock,
	GenesisHash:              &testNet4GenesisHash,
	PowLimit:                 testNet4PowLimit,
	PowLimitBits:             0x1e0fffff, // TODO
	BIP0034Height:            0,
	BIP0065Height:            -1,
	BIP0066Height:            -1,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 1051200,
	TargetTimespan:           time.Second * 95040, // 1.1 weeks
	TargetTimePerBlock:       time.Second * 90,    // 90 seconds
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     0,
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []Checkpoint{
		{1500, newHashFromStr("fa09400a2b4b62e852383d7abd67023e0aaad0e484e190cde9b6d6c079cf2a19")},
		{4000, newHashFromStr("6f554b86ba4379e181f890fb7b7eb14a47a04661b15cc42b3ab8dba48184bd13")},
		{8000, newHashFromStr("32c1f92d0e4968e26577859c31643a60650464953d85dde92f4b434e1af35258")},
		{16000, newHashFromStr("df4b1eec11c3ae392923278c0ec70beb162c5c80e2563208b9ea2b40e72d9128")},
		{32000, newHashFromStr("0e96364a1a7813a3bb7015ea6ee99fcbd76fbc198b13d4eb33915a0d1340f357")},
		{58700, newHashFromStr("668f0cfaca30134792208fd77ffb9f7405aec52b25d75ce5df032342d7b5c216")},
		{80000, newHashFromStr("4666cafb5b29e23acc991a9850ad4f21969695e52794f782fcb18eddc6cc0033")},
		{115000, newHashFromStr("8bd40f1c82cd0572b9e5e5eab29262acbb237597c487903fc53a3d5c07edda1a")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 75,  // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       100, //
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  1199145601, // January 1, 2008 UTC
			ExpireTime: 1230767999, // December 31, 2008 UTC
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  1488931200, // March 8, 2017
			ExpireTime: 1520467200, // March 8, 2018
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  1488931200, // March 1, 2017
			ExpireTime: 1520467200, // March 8, 2018
		},
	},

	// Mempool parameters
	RelayNonStdTxs: false,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "tmona", // always tmona for test net

	// Address encoding magics
	PubKeyHashAddrID:        0x6f, // starts with m or n
	ScriptHashAddrID:        0x75, // starts with p
	WitnessPubKeyHashAddrID: 0x03, // starts with QW
	WitnessScriptHashAddrID: 0x28, // starts with T7n
	PrivateKeyID:            0xef,

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,

	// vAlertPubKey is used checkpoint's deliverty
	AlertPubMainKey: []byte{
		0x04, 0x88, 0x76, 0x65, 0x07, 0x0e, 0x79, 0xd2,
		0x0f, 0x72, 0x28, 0x57, 0xe5, 0x8e, 0xc8, 0xf4,
		0x02, 0x73, 0x3f, 0x71, 0x01, 0x35, 0x52, 0x1a,
		0x0b, 0x63, 0x44, 0x14, 0x19, 0xbf, 0x56, 0x65,
		0xba, 0x46, 0x23, 0xbe, 0xd1, 0x3f, 0xca, 0x0c,
		0xb2, 0x33, 0x86, 0x82, 0xab, 0x2a, 0x54, 0xad,
		0x13, 0xce, 0x07, 0xfb, 0xc8, 0x1c, 0x3c, 0x2f,
		0x09, 0x12, 0xa4, 0xeb, 0x85, 0x21, 0xdd, 0x3c,
		0xfb,
	},
	AlertPubSubKey: []byte{
		0x04, 0x33, 0x68, 0xff, 0xcc, 0x74, 0xf2, 0x9c,
		0x9e, 0x67, 0x7f, 0x1b, 0x4a, 0x5d, 0x58, 0x34,
		0x5e, 0x67, 0x2a, 0x19, 0x78, 0x25, 0x14, 0xc8,
		0xe1, 0xaf, 0x11, 0x5f, 0x17, 0xd1, 0x10, 0x76,
		0x97, 0xdd, 0x70, 0x17, 0x59, 0x3c, 0x7c, 0xb3,
		0x91, 0x01, 0x01, 0x1f, 0x2d, 0xaf, 0xe5, 0xb0,
		0xb1, 0x41, 0x67, 0x07, 0x2d, 0xe7, 0x76, 0xc2,
		0x79, 0x63, 0xa9, 0x5a, 0x6a, 0xa9, 0x4d, 0x10,
		0x96,
	},

	// DGWv3's HF height. Used to calculate target.
	DGWv3Height: 55,

	// Lyra2re2's HF height. Used to calculate Pow.
	Lyra2re2Height: 30,
}

// SimNetParams defines the network parameters for the simulation test Bitcoin
// network.  This network is similar to the normal test network except it is
// intended for private use within a group of individuals doing simulation
// testing.  The functionality is intended to differ in that the only nodes
// which are specifically specified are used to create the network rather than
// following normal discovery rules.  This is important as otherwise it would
// just turn into another public testnet.
var SimNetParams = Params{
	Name:        "simnet",
	Net:         wire.SimNet,
	DefaultPort: "18555",
	DNSSeeds:    []DNSSeed{}, // NOTE: There must NOT be any seeds.

	// Chain parameters
	GenesisBlock:             &simNetGenesisBlock,
	GenesisHash:              &simNetGenesisHash,
	PowLimit:                 simNetPowLimit,
	PowLimitBits:             0x207fffff,
	BIP0034Height:            0, // Always active on simnet
	BIP0065Height:            0, // Always active on simnet
	BIP0066Height:            0, // Always active on simnet
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        true,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: nil,

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 75, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       100,
	Deployments: [DefinedDeployments]ConsensusDeployment{
		DeploymentTestDummy: {
			BitNumber:  28,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentCSV: {
			BitNumber:  0,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires
		},
		DeploymentSegwit: {
			BitNumber:  1,
			StartTime:  0,             // Always available for vote
			ExpireTime: math.MaxInt64, // Never expires.
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "smona", // always sb for sim net

	// Address encoding magics
	PubKeyHashAddrID:        0x3f, // starts with S
	ScriptHashAddrID:        0x7b, // starts with s
	PrivateKeyID:            0x64, // starts with 4 (uncompressed) or F (compressed)
	WitnessPubKeyHashAddrID: 0x19, // starts with Gg
	WitnessScriptHashAddrID: 0x28, // starts with ?

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x20, 0xb9, 0x00}, // starts with sprv
	HDPublicKeyID:  [4]byte{0x04, 0x20, 0xbd, 0x3a}, // starts with spub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 115, // ASCII for s

	// vAlertPubKey is used checkpoint's deliverty
	AlertPubMainKey: []byte{},
	AlertPubSubKey:  []byte{},

	// DGWv3's HF height. Used to calculate target.
	// simnet only works up to the height specified in DGWv3Height.
	// Please increase if you don't have enough value.
	DGWv3Height: 1000025,

	// Lyra2re2's HF height. Used to calculate Pow.
	Lyra2re2Height: 0,
}

var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate Bitcoin network")

	// ErrUnknownHDKeyID describes an error where the provided id which
	// is intended to identify the network for a hierarchical deterministic
	// private extended key is not registered.
	ErrUnknownHDKeyID = errors.New("unknown hd private extended key bytes")
)

var (
	registeredNets       = make(map[wire.BitcoinNet]struct{})
	pubKeyHashAddrIDs    = make(map[byte]struct{})
	scriptHashAddrIDs    = make(map[byte]struct{})
	bech32SegwitPrefixes = make(map[string]struct{})
	hdPrivToPubKeyIDs    = make(map[[4]byte][]byte)
)

// String returns the hostname of the DNS seed in human-readable form.
func (d DNSSeed) String() string {
	return d.Host
}

// Register registers the network parameters for a Bitcoin network.  This may
// error with ErrDuplicateNet if the network is already registered (either
// due to a previous Register call, or the network being one of the default
// networks).
//
// Network parameters should be registered into this package by a main package
// as early as possible.  Then, library packages may lookup networks or network
// parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]

	// A valid Bech32 encoded segwit address always has as prefix the
	// human-readable part for the given net followed by '1'.
	bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}
	return nil
}

// mustRegister performs the same function as Register except it panics if there
// is an error.  This should only be called from package init functions.
func mustRegister(params *Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-pubkey-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsScriptHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.  This is
// used when decoding an address string into a specific address type.  It is up
// to the caller to check both this and IsPubKeyHashAddrID and decide whether an
// address is a pubkey hash address, script hash address, neither, or
// undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic
// extended key id and returns the associated public key id.  When the provided
// id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}

	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}

	return pubBytes, nil
}

// newHashFromStr converts the passed big-endian hex string into a
// chainhash.Hash.  It only differs from the one available in chainhash in that
// it panics on an error since it will only (and must only) be called with
// hard-coded, and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, err := chainhash.NewHashFromStr(hexStr)
	if err != nil {
		// Ordinarily I don't like panics in library code since it
		// can take applications down without them having a chance to
		// recover which is extremely annoying, however an exception is
		// being made in this case because the only way this can panic
		// is if there is an error in the hard-coded hashes.  Thus it
		// will only ever potentially panic on init and therefore is
		// 100% predictable.
		panic(err)
	}
	return hash
}

func init() {
	// Register all default networks when the package is initialized.
	mustRegister(&MainNetParams)
	mustRegister(&TestNet4Params)
	mustRegister(&RegressionNetParams)
	mustRegister(&SimNetParams)
}
