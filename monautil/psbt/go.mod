module github.com/monasuite/monad/monautil/psbt

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/monasuite/monad v0.22.3-beta.0.20220528054306-a5fb7bce8b48
	github.com/monasuite/monad/btcec/v2 v2.1.3
	github.com/monasuite/monad/chaincfg/chainhash v1.0.2
	github.com/monasuite/monad/monautil v1.1.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/monasuite/monad/monautil => ../

replace github.com/monasuite/monad => ../..
