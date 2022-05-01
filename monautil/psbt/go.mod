module github.com/monasuite/monad/monautil/psbt

go 1.13

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/monasuite/monad v0.22.3-beta.0.20220423152854-630e468daa19
	github.com/monasuite/monad/btcec/v2 v2.1.2
	github.com/monasuite/monad/chaincfg/chainhash v1.0.1
	github.com/monasuite/monad/monautil v1.1.0
)

replace github.com/monasuite/monad/monautil => ../

replace github.com/monasuite/monad => ../..
