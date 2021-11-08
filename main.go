package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"wallet/client"
	hdwallet "wallet/hdwallet"
)

func main1() {
	w, err := hdwallet.NewWallet("./keystore")
	if err != nil {
		fmt.Println("Failed to NewWallet")
		return
	}
	w.StoreKey("123")

}

func main12() {
	filename := "6a659d835181320e5ff388f5ceae791bfca32749"
	datadir := "E:\\project\\dapp_project\\local_blockchain\\data\\keystore"
	transferTest(filename, datadir)
	w, _ := hdwallet.LoadWallet("6a659d835181320e5ff388f5ceae791bfca32749", "E:\\project\\dapp_project\\local_blockchain\\data\\keystore")
	fmt.Println(w.Address)

	//keyjson, err := ioutil.ReadFile("E:\\project\\dapp_project\\local_blockchain\\data\\keystore\\UTC--2021-10-20T07-46-37.222655700Z--7152af9ee51eba1f6045a0af7c0311aae4770866")
	//keyjson,err := ioutil.ReadFile("D:\\software\\GoPath\\src\\Go_blockchain\\wallet\\keystore\\0x1a6aBb5A3C1fEAa3D1521453971487A834703114")
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(keyjson)
}

func main() {
	c := client.NewCmdClient("http://localhost:8545", "E:\\project\\dapp_project\\local_blockchain\\data\\keystore")
	c.Run()
	//c.Transfer("7152af9ee51eba1f6045a0af7c0311aae4770866", "a71a7f5ab9e0b0bb848f5ed6e0e37a2e75e001e2", 2001)
}

//go run main.go transfer -prefix UTC--2021-10-20T07-46-37.222655700Z--  -from 0x7152af9ee51eba1f6045a0af7c0311aae4770866
//-toaddr 0xa71a7f5ab9e0b0bb848f5ed6e0e37a2e75e001e2  -value 2001
func transferTest(filename, datadir string) {
	//r := regexp.MustCompile(`{[^{}]*}`)
	//matches := r.FindAllString("{city}, {state} {zip}", -1)
	//for _, v := range matches {
	//	fmt.Println(v[1])
	//}
	tmpfile := ""
	fiels, _ := ioutil.ReadDir(datadir)
	for _, tmpf := range fiels {
		tmpfile = tmpf.Name()
		r, _ := regexp.MatchString(filename, tmpfile)
		if r {
			fmt.Println(tmpfile)
		}

	}
	fmt.Println(tmpfile)
}
