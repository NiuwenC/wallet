package hdkeystore

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
)

// 存储密钥 keystore 的模拟实现

type HDKeyStore struct {
	keysDirPath string
	scryptN     int
	scryptP     int
	Key         keystore.Key
}

//全局加密随机阅读器
var rander = rand.Reader

//生成UUID
func NewRandom() uuid.UUID {
	return uuid.New()
}

//3 编写HDKeyStore构造函数
//给出一个生成HDkeyStore对象的方法，通过privateKey生成
func NewHDkeyStore(path string, privateKey *ecdsa.PrivateKey) *HDKeyStore {
	//获得uuid
	uuid := NewRandom()
	key := keystore.Key{
		Id:         uuid,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	return &HDKeyStore{
		keysDirPath: path,
		scryptN:     keystore.LightScryptN,
		scryptP:     keystore.LightScryptP,
		Key:         key,
	}

}

// 写入文件实现,源码中的KeyStore实际是以太坊的一个接口，内部定义了三个方法，都需要实现

func (ks HDKeyStore) StoreKey(filename string, key *keystore.Key, auth string) error {
	//编码key为json
	keyjson, err := keystore.EncryptKey(key, auth, ks.scryptN, ks.scryptP)
	if err != nil {
		return err
	}
	return WriteKeyFile(filename, keyjson)
}

func WriteKeyFile(file string, content []byte) error {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()

	return os.Rename(f.Name(), file)
}

func (ks HDKeyStore) JoinPath(filename string) string {
	fmt.Println(ks.keysDirPath)
	//如果filename是绝对路径，则直接返回
	if filepath.IsAbs(filename) {
		return filename
	}
	//正则匹配获取完整路径
	fiels, _ := ioutil.ReadDir(ks.keysDirPath)
	tmpfile := ""
	for _, tmpf := range fiels {
		tmpfile = tmpf.Name()
		r, _ := regexp.MatchString(filename, tmpfile)
		if r {
			fmt.Println(tmpfile)
			break
		}
	}

	//将路径与文件拼接
	return filepath.Join(ks.keysDirPath, tmpfile)
}

//解析key，keystore文件解析。
func (ks *HDKeyStore) GetKey(addr common.Address, filename, auth string) (*keystore.Key, error) {
	//读取文件内容

	//TODO: 要在这里读取文件内容，但是geth生成的filename 并不只是address，还包括前缀
	keyjson, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//利用以太坊DecryptKey解码json文件
	key, err := keystore.DecryptKey(keyjson, auth)
	if err != nil {
		return nil, err
	}
	// 如果地址不同代表解析失败
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	ks.Key = *key
	return key, nil
}

func NewHDkeyStoreNoKey(path string) *HDKeyStore {
	return &HDKeyStore{
		keysDirPath: path,
		scryptN:     keystore.LightScryptN,
		scryptP:     keystore.LightScryptP,
		Key:         keystore.Key{},
	}
}

func (ks HDKeyStore) getRegexFile() {

}

//对交易进行签名，交易使用官方的结构
func (ks HDKeyStore) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// 交易签名
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, ks.Key.PrivateKey)
	if err != nil {
		return nil, err
	}
	//验证 签名
	msg, err := signedTx.AsMessage(types.HomesteadSigner{}, nil)
	if err != nil {
		return nil, err
	}
	sender := msg.From()
	if sender != address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

func (ks HDKeyStore) NewTransactOpts() *bind.TransactOpts {
	return bind.NewKeyedTransactor(ks.Key.PrivateKey)
}

func testMatch() {

}
