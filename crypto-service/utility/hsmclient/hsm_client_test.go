// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package hsmclient

import (
	"fmt"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	//from = "GA3U34DAANSCZ6TCNGF43JKLNB3YQKSUQVRZHW35ZYL3EXVYOUYR7AXC"
	to = "GBJWMHOQHXUE4CCIJ4GXRFF5VMOFRXDHFL7GU4GCFXTQAADRXX7BGBPF"
	//privateKeyLabel = "oVHjlfrwcbMdTzq"
	//publicKeyLabel = "kshhaijSSwhMQKm"
	c, slot, session = InitiateHSM()
	_, privateKeyObject, _ = FindHSMObject(*c, session, "test1-private-key")
	data, publicKeyObject, _ = FindHSMObject(*c, session, "test1-public-key")
	//_, _, publicKeyObject, privateKeyObject, data = GenerateKeyPair(*c, slot, "test1-public-key", "test1-private-key")
	tx build.TransactionEnvelopeBuilder
	txHash [32]byte
	publicKey string
	)

// Build a TransactionEnvelopeBuilder for the payment transaction.
// A transaction need to be hash before submit.
func initiateTransaction(from, to string) (build.TransactionEnvelopeBuilder, [32]byte){
	tx, err := build.Transaction(
		build.SourceAccount{AddressOrSeed: from},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: to},
			build.NativeAmount{Amount: "10"},
		),
	)
	if err != nil {
		LOGGER.Errorf("Error while creating the transaction: %v\n", err)
	}

	txHash, _ := tx.Hash()

	var teb build.TransactionEnvelopeBuilder

	err = teb.Mutate(tx)
	if err != nil {
		LOGGER.Errorf("Error while mutating the transaction: %v\n", err)
	}
	LOGGER.Infof("Successfully initiate transaction.\n")

	return teb, txHash
}

func submitTransactionToStellar(txeB64 string) (bool){
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		LOGGER.Errorf("Error while submitting a transaction to Stellar test-net: %v\n", err)
		return false
	}
	LOGGER.Infof("Transaction posted in ledger: %v\n", resp.Ledger)

	return true
}

func TestGenerateStellarAccount(t *testing.T){
	//data := []uint8{4, 32, 55, 77, 240, 96, 3, 100, 44, 250, 98, 105, 139, 205, 165, 75, 104, 119, 136, 42, 84, 133, 99, 147, 219, 125, 206, 23, 178, 94, 184, 117, 49, 31}
	publicKey = GenerateStellarAccount(data)
	Convey("Correct Stellar account is expected.", t, func(){
		addr := "GCE7DV4X5G5HD257D32PGWPRWSQ4CHU2FZX5RG5QQWWVQ6VTLI4EQYKK"
		So(publicKey, ShouldEqual, addr)
	})
	Convey("Wrong Stellar account is expected.", t, func(){
		addr := "GA3U34DAANSCZ6TCNGF44JKLNB3YQKSUQVRZHW35ZYL3EXVYOUYR7AXC"
		So(publicKey, ShouldNotEqual, addr)
	})
	tx, txHash = initiateTransaction(publicKey, to)
}

func TestInitTransaction(t *testing.T){
	txeB64, _:= GetSignatureAndAddToTransaction(*c, privateKeyObject, slot, txHash, tx, publicKey)
	Convey("Correct XDR is expected.", t, func(){
		So(txeB64, ShouldNotBeNil)
	})
}

func TestSubmitTransaction(t *testing.T) {
	Convey("Successfully submit transaction to Stellar test-net.", t, func(){
		txe := "AAAAADdN8GADZCz6YmmLzaVLaHeIKlSFY5Pbfc4Xsl64dTEfAAAAZAAFIIMAAAABAAAAAAAAAAAAAAABAAAAAAAAAAEAAAAAU2Yd0D3oTghITw14lL2rHFjcZyr+anDCLecAAHG9/hMAAAAAAAAAAAX14QAAAAAAAAAAAbh1MR8AAABAvYRP8pNW7VKTY1WBMsJK3+NX7ZQnwyKN4OkF8a2pvIyg0vo6LmJzYaQkT8Q8sjQVuTSW3fghTRW+pF2LHQuLBQ=="
		result := submitTransactionToStellar(txe)
		So(result, ShouldBeTrue)
	})
	Convey("Failed to submit transaction to Stellar test-net.", t, func(){
		txe := "abc"
		result := submitTransactionToStellar(txe)
		So(result, ShouldBeFalse)
	})
}

func TestVerifyDataWithPublicKey(t *testing.T) {
	sig, _ := SignDataWithPrivateKey(*c, privateKeyObject, 3, txHash)
	Convey("Successfully verify data with public key.", t, func(){
		err := VerifyDataWithPublicKey(*c, publicKeyObject, slot, txHash, sig)
		So(err, ShouldBeNil)
	})
	Convey("Failed to verify data with public key.", t, func(){
		err := VerifyDataWithPublicKey(*c, publicKeyObject, slot, txHash, []byte{'1'})
		So(err, ShouldNotBeNil)
	})
}

func BenchmarkAddSignatureAsync(b *testing.B) {
	rec := make(chan build.TransactionEnvelopeBuilder, 100)
	for i := 0 ; i < 10 ; i++ {
		go func() {
			txe,_:= GetSignatureAndAddToTransaction(*c, privateKeyObject, slot, txHash, tx, publicKey)
			rec <- txe
		}()
	}
	for t := 0 ; t < 10 ; t++ {
		<-rec
	}
}

func BenchmarkAddSignatureAsync10(b *testing.B) {
	for n := 1 ; n <= 10 ; n++ {
		fmt.Printf("\nTime %d.\n", n)
		BenchmarkAddSignatureAsync(b)
	}
}

