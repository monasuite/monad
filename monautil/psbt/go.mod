module github.com/monasuite/monad/monautil/psbt

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/monasuite/monad v0.23.0
	github.com/monasuite/monad/btcec/v2 v2.2.0
	github.com/monasuite/monad/chaincfg/chainhash v1.0.2
	github.com/monasuite/monad/monautil v1.1.2-0.20220529200742-de38c6468c62
	github.com/stretchr/testify v1.7.0
)

replace github.com/monasuite/monad/monautil => ../

replace github.com/monasuite/monad => ../..
