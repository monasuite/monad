module github.com/monasuite/monad/monautil/psbt

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/monasuite/monad v0.22.3-beta.0.20220423152854-630e468daa19
	github.com/monasuite/monad/btcec/v2 v2.1.0
	github.com/monasuite/monad/monautil v1.0.0 // indirect
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
)

replace github.com/monasuite/monad/btcec/v2 => ../../btcec

replace github.com/monasuite/monad/monautil => ../

replace github.com/monasuite/monad => ../..
