module github.com/monasuite/monad/monautil

go 1.16

require (
	github.com/aead/siphash v1.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1
	github.com/kkdai/bstream v1.0.0
	github.com/monasuite/monad v0.22.3-beta.0.20220423152854-630e468daa19
	github.com/monasuite/monad/btcec/v2 v2.1.3
	github.com/monasuite/monad/chaincfg/chainhash v1.0.2
	github.com/shopspring/decimal v1.3.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
)

replace github.com/monasuite/monad => ../
