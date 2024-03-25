package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gabrielluizsf/go-web3/core"
	"github.com/gabrielluizsf/go-web3/crypto"
	"github.com/gabrielluizsf/go-web3/network"
	"github.com/gabrielluizsf/go-web3/types"
	"github.com/gabrielluizsf/go-web3/util"
)

func main() {
	validatorPrivKey := crypto.GeneratePrivateKey()
	localNode := makeServer("LOCAL_NODE", &validatorPrivKey, ":3000", []string{":4000"}, ":9000")
	go localNode.Start()

	remoteNode := makeServer("REMOTE_NODE", nil, ":4000", []string{":5000"}, "")
	go remoteNode.Start()

	remoteNodeB := makeServer("REMOTE_NODE_B", nil, ":5000", nil, "")
	go remoteNodeB.Start()

	go func() {
		time.Sleep(11 * time.Second)

		lateNode := makeServer("LATE_NODE", nil, ":6000", []string{":4000"}, "")
		go lateNode.Start()
	}()

	time.Sleep(1 * time.Second)

	select {}
}

func sendTransaction(privKey crypto.PrivateKey) error {
	toPrivKey := crypto.GeneratePrivateKey()

	transaction := core.NewTransaction(nil)
	transaction.To = toPrivKey.PublicKey()
	transaction.Value = 42

	if err := transaction.Sign(privKey); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTransactionEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/Transaction", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)

	return err
}

func makeServer(id string, pk *crypto.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *network.Server {
	opts := network.ServerOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes:     seedNodes,
		ListenAddr:    addr,
		PrivateKey:    pk,
		ID:            id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func createCollectionTransaction(privKey crypto.PrivateKey) types.Hash {
	transaction := core.NewTransaction(nil)
	transaction.TransactionInner = core.CollectionTransaction{
		Fee:      200,
		MetaData: []byte("chicken and egg collection!"),
	}
	transaction.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTransactionEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/Transaction", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	return transaction.Hash(core.TransactionHasher{})
}

func nftMinter(privKey crypto.PrivateKey, collection types.Hash) {
	metaData := map[string]any{
		"power":  8,
		"health": 100,
		"color":  "green",
		"rare":   "yes",
	}

	metaBuf := new(bytes.Buffer)
	if err := json.NewEncoder(metaBuf).Encode(metaData); err != nil {
		panic(err)
	}

	transaction := core.NewTransaction(nil)
	transaction.TransactionInner = core.MintTransaction{
		Fee:             200,
		NFT:             util.RandomHash(),
		MetaData:        metaBuf.Bytes(),
		Collection:      collection,
		CollectionOwner: privKey.PublicKey(),
	}
	transaction.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := transaction.Encode(core.NewGobTransactionEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/Transaction", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
