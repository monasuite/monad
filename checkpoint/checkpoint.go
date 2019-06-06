package checkpoint

import (
	"fmt"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/monasuite/monad/chaincfg"
	"github.com/monasuite/monad/wire"
	"github.com/monasuite/monautil"
)

const (
	// userCheckpointDbNamePrefix is the prefix for the monad usercheckpoint database.
	userCheckpointDbNamePrefix = "usercheckpoints"

	// volatileCheckpointDbNamePrefix is the prefix for the monad volatilecheckpoint database.
	volatileCheckpointDbNamePrefix = "volatilecheckpoints"

	// AlertKeyDbNamePrefix is the prefix for the monad volatilecheckpoint database.
	alertKeyDbNamePrefix = "alertkey"

	defaultDbType = "leveldb"
)

var (
	monadHomeDir   = monautil.AppDataDir("monad", false)
	defaultDataDir = filepath.Join(monadHomeDir, "data")
)

type UserCheckpoint struct {
	Ucdb *leveldb.DB
}

var uinstance *UserCheckpoint
var uonce sync.Once

// netName returns the name used when referring to a bitcoin network.  At the
// time of writing, monad currently places blocks for testnet version 3 in the
// data and log directory "testnet", which does not match the Name field of the
// chaincfg parameters.  This function can be used to override this directory name
// as "testnet" when the passed active network matches wire.TestNet4.
//
// A proper upgrade to move the data and log directories for this network to
// "testnet4" is planned for the future, at which point this function can be
// removed and the network parameter's name used instead.
func netName(chainParams *chaincfg.Params) string {
	switch chainParams.Net {
	case wire.TestNet4:
		return "testnet"
	default:
		return chainParams.Name
	}
}

// open usercheckpointDB. Basically it is called only at startup.
func (uc *UserCheckpoint) OpenDB(dataDir string, chainParams *chaincfg.Params) (*leveldb.DB, error) {
	if uc.Ucdb != nil {
		return uc.Ucdb, nil
	}

	var err error
	dbpath := GetUserCheckpointDbPath(dataDir, chainParams)
	uc.Ucdb, err = leveldb.OpenFile(dbpath, nil)
	return uc.Ucdb, err
}

// Basically it is called only at the end.
func (uc *UserCheckpoint) CloseDB() {
	if uc.Ucdb == nil {
		return
	}
	uc.Ucdb.Close()
	uc.Ucdb = nil
}

func (uc *UserCheckpoint) Add(height int64, hash string) {
	_ = uc.Ucdb.Put([]byte(fmt.Sprintf("%020d", height)), []byte(hash), nil)
}

func (uc *UserCheckpoint) Delete(height int64) {
	_ = uc.Ucdb.Delete([]byte(fmt.Sprintf("%020d", height)), nil)
}

func (uc *UserCheckpoint) GetMaxCheckpointHeight() (height int64) {
	height = 0
	iter := uc.Ucdb.NewIterator(nil, nil)
	iter.Last()

	if !iter.Valid() {
		return height
	}

	height, _ = strconv.ParseInt(string(iter.Key()), 10, 64)
	iter.Release()
	return height
}

func GetUserCheckpointDbInstance() *UserCheckpoint {
	uonce.Do(func() {
		time.Sleep(1 * time.Second)
		uinstance = &UserCheckpoint{nil}
	})
	return uinstance
}

func GetUserCheckpointDbPath(dataDir string, chainParams *chaincfg.Params) (dbPath string) {
	dbName := userCheckpointDbNamePrefix + "_" + defaultDbType
	dbPath = filepath.Join(dataDir, dbName)

	return dbPath
}

type VolatileCheckpoint struct {
	Vcdb *leveldb.DB
}

var vinstance *VolatileCheckpoint
var vonce sync.Once

// open volatilecheckpointDB. Basically it is called only at startup.
func (vc *VolatileCheckpoint) OpenDB(dataDir string, chainParams *chaincfg.Params) (*leveldb.DB, error) {
	if vc.Vcdb != nil {
		return vc.Vcdb, nil
	}

	var err error
	dbpath := GetVolatileCheckpointDbPath(dataDir, chainParams)
	vc.Vcdb, err = leveldb.OpenFile(dbpath, nil)
	return vc.Vcdb, err
}

// close volatilecheckpointDB. Basically it is called only at the end.
func (vc *VolatileCheckpoint) CloseDB() {
	if vc.Vcdb == nil {
		return
	}
	vc.Vcdb.Close()
	vc.Vcdb = nil
}

func (vc *VolatileCheckpoint) Set(height int64, hash string) {
	_ = vc.Vcdb.Put([]byte(fmt.Sprintf("%020d", height)), []byte(hash), nil)
}

func (vc *VolatileCheckpoint) ClearDB() {
	if vc.Vcdb == nil {
		return
	}
	iter := vc.Vcdb.NewIterator(nil, nil)
	for iter.Next() {
		err := vc.Vcdb.Delete([]byte(string(iter.Key())), nil)
		if err != nil {
			break
		}
	}
	iter.Release()
}

func GetVolatileCheckpointDbInstance() *VolatileCheckpoint {
	vonce.Do(func() {
		time.Sleep(1 * time.Second)
		vinstance = &VolatileCheckpoint{nil}
	})
	return vinstance
}

func GetVolatileCheckpointDbPath(dataDir string, chainParams *chaincfg.Params) (dbPath string) {
	dbName := volatileCheckpointDbNamePrefix + "_" + defaultDbType
	dbPath = filepath.Join(dataDir, dbName)

	return dbPath
}

type AlertKey struct {
	Akdb *leveldb.DB
}

var ainstance *AlertKey
var aonce sync.Once

func (ak *AlertKey) OpenDB(dataDir string, chainParams *chaincfg.Params) (*leveldb.DB, error) {
	if ak.Akdb != nil {
		return ak.Akdb, nil
	}
	var err error
	dbpath := GetAlertKeyDbPath(dataDir, chainParams)
	ak.Akdb, err = leveldb.OpenFile(dbpath, nil)
	return ak.Akdb, err
}

func (ak *AlertKey) CloseDB() {
	if ak.Akdb == nil {
		return
	}
	ak.Akdb.Close()
	ak.Akdb = nil
}

// Alertkey is disabled when you came here.irreversible.
func (ak *AlertKey) Set(key string) {
	_ = ak.Akdb.Put([]byte(key), []byte("true"), nil)
}

// Add alertkey if it is not in database.
// Returns true if both of the public keys alertkey are OK.
func (ak *AlertKey) IsValid(chainParams *chaincfg.Params) bool {
	if ak.Akdb == nil {
		return true
	}
	d1, err := ak.Akdb.Get(chainParams.AlertPubMainKey, nil)
	if err != nil {
		_ = ak.Akdb.Put(chainParams.AlertPubMainKey, []byte("false"), nil)
	}

	d2, err := ak.Akdb.Get(chainParams.AlertPubSubKey, nil)
	if err != nil {
		_ = ak.Akdb.Put(chainParams.AlertPubSubKey, []byte("false"), nil)
	}

	if string(d1) == "false" && string(d2) == "false" {
		return true
	}
	return false
}

func GetAlertKeyDbInstance() *AlertKey {
	aonce.Do(func() {
		time.Sleep(1 * time.Second)
		ainstance = &AlertKey{nil}
	})
	return ainstance
}

func GetAlertKeyDbPath(dataDir string, chainParams *chaincfg.Params) (dbPath string) {
	dbName := alertKeyDbNamePrefix + "_" + defaultDbType
	dbPath = filepath.Join(dataDir, dbName)
	return dbPath
}
