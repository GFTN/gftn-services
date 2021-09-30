// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package ww_stellar

import (
	"bytes"
	"encoding/base64"
	"log"
	"strings"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/stellar/go/clients/horizonclient"
	hClient "github.com/stellar/go/clients/horizonclient"

	. "github.com/smartystreets/goconvey/convey"
	b "github.com/stellar/go/build"
	hProtocol "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/mock"
	"github.com/GFTN/gftn-services/gftn-models/model"
)

// type HTTPClientMock struct {
// 	GetFunc  func(url string) (*http.Response, error)
// 	PostFunc func(url, contentType string, body io.Reader) (*http.Response, error)
// }

// func (m *HTTPClientMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
// 	return m.PostFunc(url, contentType, body)
// }
// func (m *HTTPClientMock) Get(url string) (*http.Response, error) {
// 	return m.GetFunc(url)
// }

type MockPRClient struct {
	mock.Mock
}

func (m *MockPRClient) GetParticipantAccount(domain string, account string) (string, error) {
	a := m.Called(domain, account)
	return a.Get(0).(string), a.Error(1)
}
func (m *MockPRClient) GetAllParticipants() ([]model.Participant, error) {
	a := m.Called()
	return a.Get(0).([]model.Participant), a.Error(1)
}

func xdrToTransactionE(txeBase64 string) xdr.TransactionEnvelope {
	LOGGER.Debug(txeBase64)
	txeBase64R := strings.NewReader(txeBase64)
	txeByteR := base64.NewDecoder(base64.StdEncoding, txeBase64R)
	var txe xdr.TransactionEnvelope
	xdr.Unmarshal(txeByteR, &txe)
	return txe
}

func TestAppendSignature(t *testing.T) {
	txe1 := xdrToTransactionE("AAAAAMmq4ZF1kf71Z9lcoxr+UR6pVXegHA6i/jI+Py0WY3IvAAAAZAABhp8AAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAYafAAAABgAAAAAAAAAA")
	txeB1 := &b.TransactionEnvelopeBuilder{E: &txe1}
	txeB1.Init()
	txeB1.MutateTX(b.TestNetwork)
	nodeseed1 := "(seed value)"
	sig1 := b.Sign{Seed: nodeseed1}
	err := txeB1.Mutate(sig1)
	if err != nil {
		LOGGER.Error(err)
		log.Fatal()
	}
	txe2 := xdrToTransactionE("AAAAAMmq4ZF1kf71Z9lcoxr+UR6pVXegHA6i/jI+Py0WY3IvAAAAZAABhp8AAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAYafAAAABgAAAAAAAAAA")
	txeB2 := &b.TransactionEnvelopeBuilder{E: &txe2}
	txeB2.Init()
	txeB2.MutateTX(b.TestNetwork)
	nodeseed2 := "(seed value)"
	sig2 := b.Sign{Seed: nodeseed2}
	err = txeB2.Mutate(sig2)
	if err != nil {
		LOGGER.Error(err)
		log.Fatal()
	}
	Convey("Successful get caller identity", t, func() {
		txe2Byte, _ := txeB2.Bytes()
		// txe1Byte, _ := txeB1.Bytes()
		var txeSigned xdr.TransactionEnvelope
		xdr.Unmarshal(bytes.NewReader(txe2Byte), &txeSigned)
		AppendSignature(txeSigned, txeB1)
		txe1B64ss, _ := txeB1.Base64()
		txe2B64s, _ := txeB2.Base64()
		LOGGER.Debug(txe1B64ss)
		LOGGER.Debug(txe2B64s)
		So(err, ShouldBeNil)
		So(txe1B64ss, ShouldEqual, "AAAAAMmq4ZF1kf71Z9lcoxr+UR6pVXegHA6i/jI+Py0WY3IvAAAAZAABhp8AAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAYafAAAABgAAAAAAAAACpofgegAAAECzWpEZMybCND9rmUAW4Fn7MR+aO1bC89CTp3x+83Zt+5/tT6eNwNurTr/LVQwhWmdp9TUb00x6F+jHI35lWokNJMdI2QAAAEAd6V+I1cAIMWj+ry1cJXpRV7CpYACOcNrk1kMxG9yuuj4opFZTL7ZSDYX3DSGEwj5UMJRX6Mh0E0sKl9/kKJAG")
		So(txe2B64s, ShouldNotEqual, "AAAAAMmq4ZF1kf71Z9lcoxr+UR6pVXegHA6i/jI+Py0WY3IvAAAAZAABhp8AAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAYafAAAABgAAAAAAAAACpofgegAAAECzWpEZMybCND9rmUAW4Fn7MR+aO1bC89CTp3x+83Zt+5/tT6eNwNurTr/LVQwhWmdp9TUb00x6F+jHI35lWokNJMdI2QAAAEAd6V+I1cAIMWj+ry1cJXpRV7CpYACOcNrk1kMxG9yuuj4opFZTL7ZSDYX3DSGEwj5UMJRX6Mh0E0sKl9/kKJAG")
		So(txe2B64s, ShouldEqual, "AAAAAMmq4ZF1kf71Z9lcoxr+UR6pVXegHA6i/jI+Py0WY3IvAAAAZAABhp8AAAABAAAAAAAAAAAAAAABAAAAAAAAAAsAAYafAAAABgAAAAAAAAABJMdI2QAAAEAd6V+I1cAIMWj+ry1cJXpRV7CpYACOcNrk1kMxG9yuuj4opFZTL7ZSDYX3DSGEwj5UMJRX6Mh0E0sKl9/kKJAG")
	})
}

func TestAddOperation(t *testing.T) {
	txe1 := xdrToTransactionE("AAAAAJL43xZ6EVH/JWpJrBDayD1ukInRdU2lWHfhG4jtH3MlAAAAAAAKB3oAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAA=")
	tc := StellarTransactionConstructor{}
	tc.ImportTransaction(txe1)
	tex1Base64before := tc.Base64()
	// txeB1.MutateTX(b.TestNetwork)
	scAsset := b.CreditAmount{"test", "SC5HQEXEYBYFQ5ZRJ6IVPWW24EBKZBAYM565JBSXQDI4RGDSYBR72EZE", "123456"}
	paymentStellar := b.Payment(
		b.SourceAccount{AddressOrSeed: "(seed value)"},
		b.Destination{AddressOrSeed: "(seed value)"},
		scAsset,
	)
	err := tc.AddOperation(paymentStellar)
	tex1Base64after := tc.Base64()

	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
		So(tex1Base64before, ShouldNotEqual, tex1Base64after)
		So(tex1Base64after, ShouldEqual, "AAAAAJL43xZ6EVH/JWpJrBDayD1ukInRdU2lWHfhG4jtH3MlAAAAAAAKB3oAAAABAAAAAAAAAAAAAAABAAAAAQAAAABePcCYp0nnRM7HeWloMEFkkLNmpCDiPfgPB50HpofgegAAAAEAAAAAN56OXjAeFHiGiWIiJUAocnJK3tU6wm3JxUfakiTHSNkAAAABdGVzdAAAAAA3no5eMB4UeIaJYiIlQChyckre1TrCbcnFR9qSJMdI2QAAAR9xgqAAAAAAAAAAAAA=")
	})
}

func strPtr(i string) *string {
	return &i
}
func TestGetOutstandingBalance(t *testing.T) {

	// create instances of mock object
	nhcMock := horizonclient.MockClient{}
	prcMock := MockPRClient{}

	// construct expected outputs
	balance := hProtocol.Balance{}
	balance.Code = "USDDO"
	balance.Issuer = "issuer1address"
	balance.Balance = "100"
	balance.Limit = "1000"
	balances := []hProtocol.Balance{balance}
	// account detail of issuer1 address
	accountDetail := hProtocol.Account{}
	accountDetail.Balances = balances

	participantIssuer1 := model.Participant{
		ID:             strPtr("issuer1"),
		IssuingAccount: "issuer1address",
	}
	participants := []model.Participant{participantIssuer1}

	asset := model.Asset{
		AssetCode: &balance.Code,
		IssuerID:  "issuer1",
	}
	// setup expectations
	nhcMock.On("AccountDetail", hClient.AccountRequest{AccountID: "issuer1address"}).Return(accountDetail, nil)
	prcMock.On("GetAllParticipants").Return(participants, nil)

	// call the code we are testing
	outstandingBalances, _ := GetOutstandingBalances(asset, &prcMock, &nhcMock)
	LOGGER.Debug(outstandingBalances)
	// assert that the expectations were met
	Convey("Successful return balance", t, func() {
		So(outstandingBalances[0].Balance.Amount, ShouldEqual, decimal.NewFromFloat(100))
		So(outstandingBalances[0].Balance.TrustLimit, ShouldEqual, decimal.NewFromFloat(1000))
		So(outstandingBalances[0].Account.ParticipantID, ShouldEqual, "issuer1")
	})
}
