module github.com/monasuite/monad/monautil/psbt

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/monasuite/monad v0.22.3-beta.0.20220529213328-9a42bb8aabf7
	github.com/monasuite/monad/btcec/v2 v2.1.3
	github.com/monasuite/monad/chaincfg/chainhash v1.0.2
	github.com/monasuite/monad/monautil v1.1.2-0.20220529200742-de38c6468c62
)

replace github.com/monasuite/monad/monautil => ../

replace github.com/monasuite/monad => ../..

replace github.com/monasuite/monad/chaincfg/chainhash => ../../chaincfg/chainhash

replace github.com/monasuite/monad/btcec/v2 => ../../btcec
