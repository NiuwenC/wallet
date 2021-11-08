package client

import (
	"context"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	sol "wallet/abi"
	"wallet/hdwallet"
)

//处理命令行调用，需要连接到节点
type CmdClient struct {
	network string //区块链地址
	dataDir string //数据路径
}

func NewCmdClient(network, datadir string) *CmdClient {
	return &CmdClient{
		network: network,
		dataDir: datadir,
	}
}

//cmdClient创建钱包
func (cmd CmdClient) createWallet(pass string) error {
	w, err := hdwallet.NewWallet(cmd.dataDir)
	if err != nil {
		log.Panic("Failed to create Wallet", err)
	}
	//存储密钥
	return w.StoreKey(pass)
}

//Run 方法
func (c CmdClient) Help() {
	fmt.Println("./walletdevlop createwallet  -pass PASSWORD --for create new wallet")
	fmt.Println("./walletdevlop transfer -prefix prefix -from FROM -toaddr TOADDR -value VALUE --for transfer from acct to toaddr")
	fmt.Println("./walletdevelop balance  -from FROM --for get balance")

	fmt.Println("./walletdevelop sendtoken -from FROM -toaddr TOADDR -value VALUE --for sendtoken")
}

func (c CmdClient) Run() {
	if len(os.Args) < 2 {
		c.Help()
		os.Exit(-1)
	}

	//1. 立Flag
	cw_cmd := flag.NewFlagSet("Createwallet", flag.ExitOnError)
	transfer_cmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	balance_cmd := flag.NewFlagSet("balance", flag.ExitOnError)

	//sendtoken解析
	sendtoken_cmd := flag.NewFlagSet("sendtpken", flag.ExitOnError)

	//2. 立flag参数
	cw_cmd_pass := cw_cmd.String("pass", "", "PASSWORD")
	transfer_cmd_from := transfer_cmd.String("from", "", "FROM")
	transfer_cmd_toaddr := transfer_cmd.String("toaddr", "", "TOADDR")
	transfer_cmd_value := transfer_cmd.Int64("value", 0, "VALUE")
	balance_cmd_from := balance_cmd.String("from", "", "FROM")

	//参数
	sendtoken_cmd_from := sendtoken_cmd.String("from", "", "FROM")
	sendtoken_cmd_toaddr := sendtoken_cmd.String("toaddr", "", "TOADDR")
	sendtoken_cmd_value := sendtoken_cmd.Int64("value", 0, "VALUE")

	//3 解析
	switch os.Args[1] {
	case "createwallet":
		err := cw_cmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to Parse cw_cmd", err)
			return
		}
	case "transfer":
		err := transfer_cmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to Parse transfer_cmd", err)
			return
		}
	case "balance":
		err := balance_cmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to Parse transfer_cmd", err)
			return
		}
	case "sendtoken":
		err := sendtoken_cmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Failed to parse sendtoken_cmd", err)
			return
		}

	}

	//4 确认flag参数
	if cw_cmd.Parsed() {
		fmt.Println("params is ", *cw_cmd_pass)
		c.createWallet(*cw_cmd_pass)
	}
	//处理Coin转移
	if transfer_cmd.Parsed() {
		fmt.Println(*transfer_cmd_from, *transfer_cmd_toaddr, *transfer_cmd_value)
		c.Transfer(*transfer_cmd_from, *transfer_cmd_toaddr, *transfer_cmd_value)
	}

	if balance_cmd.Parsed() {
		c.balance(*balance_cmd_from)
	}

	if sendtoken_cmd.Parsed() {
		c.sendtoken(*sendtoken_cmd_from, *sendtoken_cmd_toaddr, *sendtoken_cmd_value)
	}
}

func (cmd CmdClient) Transfer(from, toaddr string, value int64) error {
	//1 钱包加载
	w, _ := hdwallet.LoadWallet(from, cmd.dataDir)
	//2 连接到以太坊,拨号连接
	cli, _ := ethclient.Dial(cmd.network)
	defer cli.Close()

	//3 获取nonce 官方文件中提供
	nonce, _ := cli.NonceAt(context.Background(), common.HexToAddress(from), nil)

	// 4 创建交易
	gaslimit := uint64(300000)
	gasprice := big.NewInt(21000000000)
	amount := big.NewInt(value)
	tx := types.NewTransaction(nonce, common.HexToAddress(toaddr), amount, gaslimit, gasprice, []byte("JackNiu"))

	//5 签名
	stx, err := w.HdKeyStore.SignTx(common.HexToAddress(from), tx, nil)
	if err != nil {
		log.Panic("Failed to Signtx", err)

	}
	//6 发送交易
	return cli.SendTransaction(context.Background(), stx)
}

//查询余额
func (cmd CmdClient) balance(from string) (int64, error) {
	//1 连接到以太坊
	cli, err := ethclient.Dial(cmd.network)
	if err != nil {
		log.Panic("Failed to ethclient dial. ", err)
	}
	defer cli.Close()

	//2 查询余额
	addr := common.HexToAddress(from)
	value, err := cli.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Panic("Failed to BalanceAt ", err)
	}
	fmt.Printf("%s's balance is %d\n", from, value)
	return value.Int64(), nil
}

const TokenCOntractAddr = ""

func (cmd CmdClient) sendtoken(from, toaddr string, value int64) error {
	//1 连接到以太坊
	cli, _ := ethclient.Dial(cmd.network)
	defer cli.Close()

	//2 创建token合约实例，需要合约地址
	token, _ := sol.NewToken(common.HexToAddress(TokenCOntractAddr), cli)

	//3 设置调用身份
	//3.1 钱包加载
	w, _ := hdwallet.LoadWallet(from, cmd.dataDir)
	//3.2 利用钱包私钥创建身份
	auth := w.HdKeyStore.NewTransactOpts()
	//4 调用转移
	_, err := token.Transfer(auth, common.HexToAddress(toaddr), big.NewInt(value))
	return err

}
