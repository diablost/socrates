package socrates

import (
	"os"
	"os/exec"
	"log"
	"context"

	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var pluginCmd *exec.Cmd

func StartParity() {
	//logf("starting light parity with option (%s)....", pluginOpts)
	logf("starting light parity ....")

	localHost := "127.0.0.1"
	newAddr := localHost + ":" + "30030"

	_ = execPlugin(newAddr)
	return
}

func execPlugin(addr string) (err error) {

	logH := newLogHelper("[parity]: ")
	//parity  --chain genesis-spec.json -d /tmp/parity0 --port 30300 --jsonrpc-port 8540 --jsonrpc-apis web3,eth,net,personal,parity,parity_set,traces,rpc,parity_accounts
	cmd := &exec.Cmd{
		Path:   "parity --light --chain genesis-spec.json -d /tmp/parity0 --port 30300 --jsonrpc-port 8540 --jsonrpc-apis web3,eth,net,personal,parity,parity_set,traces,rpc,parity_accounts ",
		//Env:    env,
		Stdout: logH,
		Stderr: logH,
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	pluginCmd = cmd
	go func() {
		if err := cmd.Wait(); err != nil {
			logf("parity exited (%v)\n", err)
			os.Exit(2)
		}
		logf("parity exited\n")
		os.Exit(0)
	}()
	return nil
}

func getBlockNumber() {
	cli, err := ethclient.Dial("127.0.0.1:30030")
	if err != nil {
		log.Fatal(err)
	}
	cli.HeaderByNumber(context.Background(), nil);
}

func getProxyServices() {
	// dicover
	c := make(chan int, 1)
	services := Discovery(c)
	for _, proxy := range services {
		logf(proxy)
	}
}

func registerProxyService() {
	// push transaction, update proxy services with my IP

}

func clientStake () {

}

func proofOfProxy() {
	// push transaction

}

func proofOfAlive() {
	// push transaction
	// ping on chain

}