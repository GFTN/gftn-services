// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package common

import (
	"net/http"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	hClient "github.com/stellar/go/clients/horizonclient"
	n "github.com/stellar/go/network"
)

// We should use this method only to point to horizon client that we want this way we can use our own standalone horizon client
func GetHorizonClient(urlString string) horizon.Client {
	horizonClient := horizon.Client{
		URL:  urlString,
		HTTP: http.DefaultClient,
	}
	return horizonClient
}

func GetStellarNetwork(passphrase string) b.Network {
	if passphrase == n.PublicNetworkPassphrase {
		return b.PublicNetwork
	} else if passphrase == n.TestNetworkPassphrase {
		return b.TestNetwork
	}
	return b.Network{Passphrase: passphrase}
}

func GetNewHorizonClient(urlString string) hClient.Client {
	horizonClient := hClient.Client{
		HorizonURL: urlString,
		HTTP:       http.DefaultClient,
	}
	return horizonClient
}
