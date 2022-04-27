module github.com/monasuite/monad

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd
	github.com/btcsuite/goleveldb v1.0.0
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792
	github.com/btcsuite/winsvc v1.0.0
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/lru v1.0.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jrick/logrotate v1.0.0
	github.com/monasuite/monad/monautil v1.0.0
	github.com/monasuite/monad/btcec/v2 v2.1.0
	github.com/shopspring/decimal v1.3.1
	github.com/wakiyamap/lyra2rev2 v0.1.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

require (
	github.com/aead/siphash v1.0.1 // indirect
	github.com/aead/skein v0.0.0-20160722084837-9365ae6e95d2 // indirect
	github.com/btcsuite/golangcrypto v0.0.0-20150304025918-53f62d9b43e8 // indirect
	github.com/btcsuite/snappy-go v1.0.0 // indirect
	github.com/dchest/blake256 v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/kkdai/bstream v1.0.0 // indirect
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/onsi/gomega v1.10.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/monasuite/monad/monautil => ./monautil

replace github.com/monasuite/monad/btcec/v2 => ./btcec

go 1.17
