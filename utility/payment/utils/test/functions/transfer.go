// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package functions

import (
	"fmt"
	"net/http"
	"time"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func Transfer() {
	from := ""

	seed := ""
	to := ""

	client := &horizon.Client{
		URL: "http://35.197.35.7:1234",
		HTTP: &http.Client{
			Timeout: time.Second * time.Duration(120),
		},
	}

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: from},
		b.Network{"Standalone World Wire Network ; Mar 2019"},
		//b.Sequence{uint64(seq) + 1},
		b.AutoSequence{client},
		b.Payment(
			b.Destination{to},
			b.CreditAmount{
				Code:   "TWD",
				Issuer: from,
				Amount: "100",
			},
		),
	)
	if err != nil {
		panic(err)
	}

	txe, err := tx.Sign(seed)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()

	if err != nil {
		panic(err)
	}

	fmt.Printf("tx base64: %s", txeB64)

	// And finally, send it off to Stellar!
	resp, err := client.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
}
