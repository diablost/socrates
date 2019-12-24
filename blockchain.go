package main

import (
	"os"
	"os/exec"
)

var pluginCmd *exec.Cmd


func startParity() {
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
		Path:   "parity --chain genesis-spec.json -d /tmp/parity0 --port 30300 --jsonrpc-port 8540 --jsonrpc-apis web3,eth,net,personal,parity,parity_set,traces,rpc,parity_accounts ",
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