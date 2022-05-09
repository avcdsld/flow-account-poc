package main

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
)

func main() {
	ctx := context.Background()
	flowClient := examples.NewFlowClient()
	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	// Create account

	privateKey1, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, []byte("seedseedseedseedseedseedseedseed"))
	examples.Handle(err)
	privateKey2, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, []byte("seedseedseedseedseedseedseedseed"))
	examples.Handle(err)
	privateKey3, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, []byte("seedseedseedseedseedseedseedseed"))
	examples.Handle(err)

	tx := flow.NewTransaction().
		SetScript(readFile("../../transactions/1_create_account.cdc")).
		SetProposalKey(serviceAcctAddr, serviceAcctKey.Index, serviceAcctKey.SequenceNumber).
		SetPayer(serviceAcctAddr).
		AddAuthorizer(serviceAcctAddr).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))
	tx.AddArgument(cadenceHexString(privateKey1.PublicKey().String()))
	tx.AddArgument(cadenceHexString(privateKey2.PublicKey().String()))
	tx.AddArgument(cadenceHexString(privateKey3.PublicKey().String()))

	err = tx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	res := examples.WaitForSeal(ctx, flowClient, tx.ID())

	var address sdk.Address
	for _, event := range res.Events {
		if event.Type == sdk.EventAccountCreated {
			address = sdk.Address(event.Value.Fields[0].(cadence.Address))
			break
		}
	}
	examples.FundAccountInEmulator(flowClient, address, 1.0)

	// Revoke keys

	account, _ := flowClient.GetAccount(ctx, address)
	signer := crypto.NewInMemorySigner(privateKey2, account.Keys[1].HashAlgo)

	tx = flow.NewTransaction().
		SetScript(readFile("../../transactions/2_revoke_keys.cdc")).
		SetProposalKey(account.Address, account.Keys[1].Index, account.Keys[1].SequenceNumber).
		SetPayer(account.Address).
		AddAuthorizer(account.Address).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))
	tx.AddArgument(cadence.NewInt(0))
	tx.AddArgument(cadence.NewInt(2))

	err = tx.SignEnvelope(account.Address, account.Keys[1].Index, signer)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	res = examples.WaitForSeal(ctx, flowClient, tx.ID())
	if res.Error != nil {
		panic(res.Error)
	}

	// Add new key

	account, _ = flowClient.GetAccount(ctx, address)

	newPrivateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, []byte("seedseedseedseedseedseedseedseed"))
	examples.Handle(err)

	tx = flow.NewTransaction().
		SetScript(readFile("../../transactions/3_add_key.cdc")).
		SetProposalKey(account.Address, account.Keys[1].Index, account.Keys[1].SequenceNumber).
		SetPayer(account.Address).
		AddAuthorizer(account.Address).
		SetReferenceBlockID(examples.GetReferenceBlockId(flowClient))
	tx.AddArgument(cadenceHexString(newPrivateKey.PublicKey().String()))

	err = tx.SignEnvelope(account.Address, account.Keys[1].Index, signer)
	examples.Handle(err)

	err = flowClient.SendTransaction(ctx, *tx)
	examples.Handle(err)

	res = examples.WaitForSeal(ctx, flowClient, tx.ID())
	if res.Error != nil {
		panic(res.Error)
	}
}

func cadenceHexString(value string) cadence.Value {
	return cadence.String(strings.Replace(value, "0x", "", 1))
}

func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return contents
}
