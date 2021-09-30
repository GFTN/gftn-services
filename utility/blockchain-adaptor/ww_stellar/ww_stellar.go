// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package ww_stellar

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	hClient "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/xdr"
	crypto_client "github.com/GFTN/gftn-services/crypto-service-client/crypto-client"
	gasserviceclient "github.com/GFTN/gftn-services/gas-service-client"
	"github.com/GFTN/gftn-services/gftn-models/model"
	ast "github.com/GFTN/gftn-services/utility/asset"
)

var (
	ERROR_GET_STELLAR_ASSET            = 1
	ERROR_INSUFFICIENT_BALANCE         = 2
	ERROR_MAP_STELLAR_ADDRESS          = 3
	ERROR_ACCOUNT_DOES_NOT_TRUST_ASSET = 4
	ERROR_ACCOUNT_TRUST_LIMIT_EXCEEDED = 5
	ERROR_TIMEOUT                      = 6
	ERROR_SIGNING_TRANSACTION          = 7
	ERROR_MAP_STELLAR_ASSET            = 8
	ERROR_SUBMIT_TRANSACTION           = 9
	ERROR_AMOUNT_LESSTHEN_EQUAL_ZERO   = 10
	ERROR_NO_ISSUER_INVOLVED           = 11
	ERROR_ASSET_DOES_NOT_EXIST         = 12
)

type StellarTransactionConstructor struct {
	txeB  b.TransactionEnvelopeBuilder
	mutex sync.Mutex
}
type Account struct {
	ParticipantID string
	Account       string
}

type Payment struct {
	SourceAccount Account
	TargetAccount Account
	Asset         model.Asset
	Amount        decimal.Decimal
}

type Balance struct {
	Amount     decimal.Decimal
	TrustLimit decimal.Decimal // Amount <= TrustLimit
}

// outstanding balance of an asset, ie the the account which holds the asset and its respective account
type OutstandingBalance struct {
	Account Account
	Balance Balance
}

type CryptoClientGlobal interface {
	RequestSigning(txeBase64 string, requestBase64 string, signedRequestBase64 string, accountName string, participant model.Participant) (string, error)
}

type PRClient interface {
	GetParticipantAccount(domain string, account string) (string, error) // return participant account address'
	GetAllParticipants() ([]model.Participant, error)                    // return all participants
}

type Error interface {
	error
	Code() int
	Msg() string
}

// wwError: implement Error
type wwError struct {
	error
	errCode int
	errMsg  string
}

// CheckResult : go channel result checking
type CheckResult struct {
	Error Error
}

func (self wwError) Code() int {
	return self.errCode
}
func (self wwError) Msg() string {
	return self.errMsg
}

// SendTransaction : send transaction to Gas service
func (self *StellarTransactionConstructor) SendTransaction(gsc gasserviceclient.GasServiceClient) (string, uint64, error) {
	// filter redundent signatures
	sigs := self.txeB.E.Signatures
	sigshm := make(map[string]xdr.DecoratedSignature)
	for _, sig := range sigs {
		sigshm[string(sig.Signature)] = sig
	}
	sigs = nil
	for _, val := range sigshm {
		sigs = append(sigs, val)
	}
	self.txeB.E.Signatures = sigs

	// set fees
	LOGGER.Info("Sending Transcation: ", self.Base64())
	hash, ledger, err := gsc.SubmitTxe(self.Base64())
	if err != nil {
		LOGGER.Error("Send Transaction to Gas Service Failed")
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_SUBMIT_TRANSACTION
		return "", 0, wwErr
	}
	LOGGER.Info("Send Transaction to Gas Service Success: " + "hash: " + hash + "ledger: " + fmt.Sprint(ledger))
	return hash, ledger, err
}

//InitTransaction : Init Transaction
func (self *StellarTransactionConstructor) InitTransaction(gsc gasserviceclient.GasServiceClient) error {
	//init mutex
	self.mutex = sync.Mutex{}
	//init txeB
	//obtain IBM account and squence number from gas service
	LOGGER.Info("Obtaining IBM account and squence number from gas service")
	ibmAccount, squenceNum, err := gsc.GetAccountAndSequence()
	if err != nil {
		LOGGER.Error("Failed to get IBMAccount and Sequence Number from Gas Service")
		LOGGER.Error(err)
		return err
	}
	LOGGER.Info("IBM account and squence number obtained:", ibmAccount, squenceNum)
	//construct atomic transaction
	LOGGER.Info("Constructing Transcation...")
	txeB := b.TransactionEnvelopeBuilder{}
	txeB.MutateTX(b.SourceAccount{AddressOrSeed: ibmAccount})
	txeB.MutateTX(b.Sequence{Sequence: squenceNum})
	var fee uint64
	fee = 100
	txeB.MutateTX(b.BaseFee{Amount: fee})
	LOGGER.Info("Constructed Transcation...")
	self.txeB = txeB
	// self.txeB.Init()
	txeB64, _ := self.txeB.Base64()
	LOGGER.Debug("initiated transaction:", txeB64)
	return nil
}

// AddPayment : Add Payment operation to transaction envelop
func (self *StellarTransactionConstructor) AddPayment(payment Payment, prc PRClient) error {

	scAsset, err := getCreditAsset(payment.Asset, payment.Amount, prc)
	LOGGER.Debug(scAsset)
	if err != nil {
		LOGGER.Error("Error getting stellar asset")
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_MAP_STELLAR_ASSET
		return err
	}
	type GetAddressResult struct {
		Address string
		Error   error
	}
	sourceAddress := ""
	targetAddress := ""
	ch1 := make(chan GetAddressResult)
	ch2 := make(chan GetAddressResult)
	go func() {
		_sourceAddress, err := prc.GetParticipantAccount(payment.SourceAccount.ParticipantID, payment.SourceAccount.Account)
		if err != nil {
			msg := "Error getting stellar source address, account:" + payment.SourceAccount.ParticipantID + ":" + payment.SourceAccount.Account
			wwErr := wwError{}
			wwErr.error = err
			wwErr.errCode = ERROR_MAP_STELLAR_ADDRESS
			wwErr.errMsg = msg
			ch1 <- GetAddressResult{_sourceAddress, wwErr}
			return
		}
		ch1 <- GetAddressResult{_sourceAddress, err}
		return
	}()
	go func() {
		_targetAddress, err := prc.GetParticipantAccount(payment.TargetAccount.ParticipantID, payment.TargetAccount.Account)
		if err != nil {
			msg := "Error getting stellar target address, account:" + payment.SourceAccount.ParticipantID + ":" + payment.SourceAccount.Account
			wwErr := wwError{}
			wwErr.error = err
			wwErr.errCode = ERROR_MAP_STELLAR_ADDRESS
			wwErr.errMsg = msg
			ch2 <- GetAddressResult{_targetAddress, wwErr}
			return
		}
		ch2 <- GetAddressResult{_targetAddress, err}
		return
	}()
	for i := 0; i < 2; i++ {
		select {
		case sourceAddressResult := <-ch1:
			err = sourceAddressResult.Error
			if err != nil {
				return err
			} else {
				sourceAddress = sourceAddressResult.Address
			}
		case targetAddressResult := <-ch2:
			err = targetAddressResult.Error
			if err != nil {
				return err
			} else {
				targetAddress = targetAddressResult.Address
			}
		case <-time.After(time.Second * 20):
			msg := "Error: Timeout when getting address from participant registry"
			LOGGER.Error(msg)
			return errors.New(msg)
		}
	}
	close(ch1)
	close(ch2)
	paymentStellar := b.Payment(
		b.SourceAccount{AddressOrSeed: sourceAddress},
		b.Destination{AddressOrSeed: targetAddress},
		scAsset,
	)
	LOGGER.Debug(paymentStellar)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err = self.AddOperation(paymentStellar)
	if err != nil {
		LOGGER.Error(err)
		return err
	}
	return nil
}

// AddPaymentWithSanityCheck : AddPaymentWithSanityCheck
func (self *StellarTransactionConstructor) AddPaymentWithSanityCheck(payment Payment, prc PRClient, hc *horizon.Client) error {

	scAsset, err := getCreditAsset(payment.Asset, payment.Amount, prc)
	LOGGER.Debug(scAsset)
	if err != nil {
		LOGGER.Error("Error getting stellar asset")
		return err
	}
	type GetAddressResult struct {
		Address string
		Error   error
	}
	sourceAddress := ""
	targetAddress := ""
	ch1 := make(chan GetAddressResult)
	ch2 := make(chan GetAddressResult)
	go func() {
		_sourceAddress, err := prc.GetParticipantAccount(payment.SourceAccount.ParticipantID, payment.SourceAccount.Account)
		if err != nil {
			LOGGER.Error("Error getting stellar source address, account:" + payment.SourceAccount.ParticipantID + ":" + payment.SourceAccount.Account)
			ch1 <- GetAddressResult{_sourceAddress, err}
			return
		}
		// checktrustline
		switch scAsset.(type) {
		case b.NativeAmount:
			// break
		case b.CreditAmount:
			account, _ := hc.LoadAccount(_sourceAddress)
			isTrust := false
			for _, balance := range account.Balances {
				if balance.Asset.Code == scAsset.(b.CreditAmount).Code &&
					balance.Asset.Issuer == scAsset.(b.CreditAmount).Issuer {
					isTrust = true
					break
				}
			}
			if isTrust == false {
				err = errors.New("Particiapant " + payment.SourceAccount.ParticipantID + " account: " + payment.SourceAccount.Account + " does not trust asset " + *payment.Asset.AssetCode)
			}
		}
		ch1 <- GetAddressResult{_sourceAddress, err}
		return
	}()
	go func() {
		_targetAddress, err := prc.GetParticipantAccount(payment.TargetAccount.ParticipantID, payment.TargetAccount.Account)
		if err != nil {
			LOGGER.Error("Error getting stellar target address, account:" + payment.SourceAccount.ParticipantID + ":" + payment.SourceAccount.Account)
			ch2 <- GetAddressResult{_targetAddress, err}
			return
		}
		// checktrustline
		switch scAsset.(type) {
		case b.NativeAmount:
			// break
		case b.CreditAmount:
			account, _ := hc.LoadAccount(_targetAddress)
			isTrust := false
			for _, balance := range account.Balances {
				if balance.Asset.Code == scAsset.(b.CreditAmount).Code &&
					balance.Asset.Issuer == scAsset.(b.CreditAmount).Issuer {
					isTrust = true
					break
				}
			}
			if isTrust == false {
				err = errors.New("Particiapant " + payment.TargetAccount.ParticipantID + " account: " + payment.TargetAccount.Account + " does not trust asset " + *payment.Asset.AssetCode)
			}
		}
		ch2 <- GetAddressResult{_targetAddress, err}
		return
	}()
	for i := 0; i < 2; i++ {
		select {
		case sourceAddressResult := <-ch1:
			err = sourceAddressResult.Error
			if err != nil {
				return err
			} else {
				sourceAddress = sourceAddressResult.Address
			}
		case targetAddressResult := <-ch2:
			err = targetAddressResult.Error
			if err != nil {
				return err
			} else {
				targetAddress = targetAddressResult.Address
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when signing transaction"
			LOGGER.Error(msg)
			return errors.New(msg)
		}
	}
	close(ch1)
	close(ch2)
	paymentStellar := b.Payment(
		b.SourceAccount{AddressOrSeed: sourceAddress},
		b.Destination{AddressOrSeed: targetAddress},
		scAsset,
	)
	LOGGER.Debug(paymentStellar)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err = self.AddOperation(paymentStellar)
	if err != nil {
		LOGGER.Error(err)
		return err
	}
	return nil
}

// AddOperation : add operation to transaction
func (self *StellarTransactionConstructor) AddOperation(operation b.TransactionMutator) error {
	LOGGER.Debug("before adding operation:", self.Base64())
	err := self.txeB.MutateTX(operation)
	LOGGER.Debug("after adding operation:", self.Base64())
	if err != nil {
		LOGGER.Error("Error Mutating transaction")
		LOGGER.Error(err)
		return err
	}
	return nil
}

// SignTransactionAndAppend : sign the transaction, should be called after the transaction is finalized
func (self *StellarTransactionConstructor) SignTransactionAndAppend(account string, cc crypto_client.CryptoServiceClient) error {
	LOGGER.Info("Signing Transcation wih account:", account)
	txebyte := self.Bytes()
	transactionSigned, err, _, _ := cc.ParticipantSignXdr(account, txebyte)
	if err != nil {
		LOGGER.Error(err)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_SIGNING_TRANSACTION
		wwErr.errMsg = "Error signing transaction"
		return err
	}
	var txe xdr.TransactionEnvelope
	xdr.Unmarshal(bytes.NewReader(transactionSigned), &txe)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err = AppendSignature(txe, &self.txeB)
	if err != nil {
		LOGGER.Error(err)
		return err
	}
	LOGGER.Info("Signature added with account:", account)
	return nil
}

// SignTransactionGlobalAndAppend : sign the transaction, should be called after the transaction is finalized
// for global micro service like quote service.
func (self *StellarTransactionConstructor) SignTransactionGlobalAndAppend(account string, cc CryptoClientGlobal, participant model.Participant, requestBase64 string, signedRequestBase64 string) error {
	LOGGER.Info("Signing Transcation wih account:", account)
	txeB64 := self.Base64()
	transactionSigned, err := cc.RequestSigning(txeB64, requestBase64, signedRequestBase64, account, participant)
	if err != nil {
		LOGGER.Error(err)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_SIGNING_TRANSACTION
		wwErr.errMsg = "Error signing transaction"
		return err
	}
	var txe xdr.TransactionEnvelope
	LOGGER.Debug(transactionSigned)
	raw := strings.NewReader(transactionSigned)
	b64r := base64.NewDecoder(base64.StdEncoding, raw)
	xdr.Unmarshal(b64r, &txe)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err = AppendSignature(txe, &self.txeB)
	if err != nil {
		LOGGER.Error(err)
		return err
	}
	LOGGER.Info("Signature added with account:", account)
	return nil
}

// AppendSignature : Append signature from signed xdr to another transactionEnvelopBuilder
func AppendSignature(transactionSigned xdr.TransactionEnvelope, txeB *b.TransactionEnvelopeBuilder) error {
	if len(transactionSigned.Signatures) == 0 {
		return errors.New("transactionEnvelop unsigned")
	}
	txeB.E.Signatures = append(txeB.E.Signatures, transactionSigned.Signatures[len(transactionSigned.Signatures)-1])
	return nil
}

func (self *StellarTransactionConstructor) Base64() string {
	txeBase64, _ := self.txeB.Base64()
	return txeBase64
}

func (self *StellarTransactionConstructor) Bytes() []byte {
	txeByte, _ := self.txeB.Bytes()
	return txeByte
}

func (self *StellarTransactionConstructor) SetFee() error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	err := self.txeB.MutateTX(b.Defaults{})
	return err
}

// Set time-bounds using quote expiry date
func (self *StellarTransactionConstructor) SetTimeBound(lowerbound, upperbound int64) error {
	timeBounds := xdr.TimeBounds{}
	timeBounds.MinTime = xdr.TimePoint(lowerbound)
	timeBounds.MaxTime = xdr.TimePoint(upperbound)
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.txeB.E.Tx.TimeBounds = &timeBounds
	return nil
}

//stellar specific
func (self *StellarTransactionConstructor) ImportTransaction(txe xdr.TransactionEnvelope) {
	self.txeB = b.TransactionEnvelopeBuilder{E: &txe}
	self.txeB.Init()
	return
}

func (self *StellarTransactionConstructor) GetTxeB() b.TransactionEnvelopeBuilder {
	return self.txeB
}
func (self *StellarTransactionConstructor) SetTxeB(txeB b.TransactionEnvelopeBuilder) {
	self.txeB = txeB
	return
}

func GetBalance(account Account, asset model.Asset, prc PRClient, hc *horizon.Client) (Balance, error) {

	_address, err := prc.GetParticipantAccount(account.ParticipantID, account.Account)
	if err != nil {
		LOGGER.Error("Error getting stellar source address, account:" + account.ParticipantID + ":" + account.Account)
	}
	_account, err := hc.LoadAccount(_address)
	if err != nil {
		LOGGER.Error("Error loading account from horizon")
		return Balance{}, err
	}
	if ast.IsNative(*asset.AssetCode) {
		balanceStr, _ := _account.GetNativeBalance()
		balanceDecimal, _ := decimal.NewFromString(balanceStr)
		LOGGER.Info("pariticipant", account.ParticipantID, "asset assetcode:", *asset.AssetCode, "balance....:", balanceStr)
		return Balance{balanceDecimal, decimal.Decimal{}}, nil
	}
	issuerAddress, err := prc.GetParticipantAccount(asset.IssuerID, "issuing")
	if err != nil {
		LOGGER.Error("Error getting asset issuing address from participant registry, issuerID: " + asset.IssuerID)
		return Balance{}, err
	}
	balance := getCreditBalance(_account, *asset.AssetCode, issuerAddress)
	LOGGER.Info("pariticipant", account.ParticipantID, "asset assetcode:", *asset.AssetCode, "balance....:", balance)
	return balance, nil
}

// SanityCheck : Sanity check for payment in the following context:
// 1. trustline exist
// 2. trustlimit does not exceed
// 3. sufficient balance
func SanityCheck(payments []Payment, prc PRClient, hc *horizon.Client) Error {

	chBalance := make(chan CheckResult)
	chTrustline := make(chan CheckResult)

	for idx, _ := range payments {
		go func(ch chan CheckResult, idx int) {
			scAsset, err := getCreditAsset(payments[idx].Asset, payments[idx].Amount, prc)
			if err != nil {
				msg := "Error getting stellar asset"
				LOGGER.Error(msg)
				wwErr := wwError{}
				wwErr.error = err
				wwErr.errCode = ERROR_GET_STELLAR_ASSET
				wwErr.errMsg = msg
				ch <- CheckResult{wwErr}
				return
			}
			// check trustline
			go func(ch chan CheckResult, idx int) {
				// checktrustline
				switch scAsset.(type) {
				case b.NativeAmount:
					ch <- CheckResult{nil}
					ch <- CheckResult{nil}
				case b.CreditAmount:
					// 1. check if trustline is established if asset is DO
					if *payments[idx].Asset.AssetType == model.AssetAssetTypeDO {
						//  1. if asset is DO, DO issuer has to be either reciever or sender
						if payments[idx].SourceAccount.ParticipantID != payments[idx].Asset.IssuerID &&
							payments[idx].TargetAccount.ParticipantID != payments[idx].Asset.IssuerID {
							msg := "Failed Sanity check: DO is neither asset issued by " + payments[idx].TargetAccount.ParticipantID + " or " + payments[idx].SourceAccount.ParticipantID
							LOGGER.Error(msg)
							wwErr := wwError{}
							wwErr.error = errors.New(msg)
							wwErr.errCode = ERROR_NO_ISSUER_INVOLVED
							wwErr.errMsg = msg
							ch <- CheckResult{wwErr}
							return
						}
						// 2. check if trustline is established if asset is DO
						if payments[idx].SourceAccount.ParticipantID == payments[idx].Asset.IssuerID {
							// check trust line of target account
							go func(ch chan CheckResult, idx int) {
								ch <- CheckTrustline(payments[idx].TargetAccount, payments[idx].Asset, payments[idx].Amount, prc, hc)
							}(chTrustline, idx)
							// no checking is needed for source account
							ch <- CheckResult{nil}
						} else {
							// checktrustline of source account
							go func(ch chan CheckResult, idx int) {
								ch <- CheckTrustline(payments[idx].SourceAccount, payments[idx].Asset, payments[idx].Amount, prc, hc)
							}(chTrustline, idx)
							// no checking is needed for target account
							ch <- CheckResult{nil}
						}
						return
					}

					// 2. check if trustline is established if asset is DA
					if *payments[idx].Asset.AssetType == model.AssetAssetTypeDA {
						// checktrustline of source account
						go func(ch chan CheckResult, idx int) {
							ch <- CheckTrustline(payments[idx].SourceAccount, payments[idx].Asset, payments[idx].Amount, prc, hc)
						}(chTrustline, idx)
						// checktrustline of target account
						go func(ch chan CheckResult, idx int) {
							ch <- CheckTrustline(payments[idx].TargetAccount, payments[idx].Asset, payments[idx].Amount, prc, hc)
						}(chTrustline, idx)
					}
					return
				}
			}(chTrustline, idx)
			// check balance
			// a. balance checking for source account: sufficient balance
			go func(ch chan CheckResult, idx int) {
				// check amount compare to balance, if balance < amount || amount <=0, return err.
				if payments[idx].Amount.LessThanOrEqual(decimal.NewFromFloat(0)) {
					err := errors.New("Failed Sanity check: Amount less than or equal to 0 is not allowed, account: " + payments[idx].SourceAccount.Account)
					LOGGER.Info(err)
					wwErr := wwError{}
					wwErr.error = err
					wwErr.errCode = ERROR_AMOUNT_LESSTHEN_EQUAL_ZERO
					ch <- CheckResult{wwErr}
					return
				}
				// doesn't check the balance limit if the asset is issued by source account
				if payments[idx].SourceAccount.ParticipantID != payments[idx].Asset.IssuerID {
					balance, err := GetBalance(payments[idx].SourceAccount, payments[idx].Asset, prc, hc)
					if err != nil {
						msg := "Error getting balance account: " + payments[idx].SourceAccount.Account
						LOGGER.Error(msg)
						wwErr := wwError{}
						wwErr.error = err
						wwErr.errMsg = msg
						ch <- CheckResult{wwErr}
						return
					}
					if balance.Amount.LessThan(payments[idx].Amount) {
						msg := "Failed Sanity check: insufficent balance Account: " + payments[idx].SourceAccount.Account
						LOGGER.Error(msg)
						wwErr := wwError{}
						wwErr.error = errors.New(msg)
						wwErr.errCode = ERROR_INSUFFICIENT_BALANCE
						wwErr.errMsg = msg
						ch <- CheckResult{wwErr}
						return
					}
				}
				ch <- CheckResult{}
				return
			}(chBalance, idx)
			// b. balance checking for target account : check trust limit
			go func(ch chan CheckResult, idx int) {
				// doesn't check the trust limit if the asset is issued by target account
				if payments[idx].TargetAccount.ParticipantID != payments[idx].Asset.IssuerID {
					balance, err := GetBalance(payments[idx].TargetAccount, payments[idx].Asset, prc, hc)
					if err != nil {
						msg := "Error getting balance account: " + payments[idx].TargetAccount.Account
						LOGGER.Error(msg)
						wwErr := wwError{}
						wwErr.error = err
						wwErr.errMsg = msg
						ch <- CheckResult{wwErr}
						return
					}
					// check if amount exceed the rest of the trust limit
					limit := balance.TrustLimit
					existingBalance := balance.Amount
					if payments[idx].Amount.GreaterThan(limit.Sub(existingBalance)) {
						msg := "Failed Sanity check: Participant " + payments[idx].TargetAccount.ParticipantID + " account: " + payments[idx].TargetAccount.Account + " trust limit exceeded with asset " + *payments[idx].Asset.AssetCode
						LOGGER.Error(msg)
						wwErr := wwError{}
						wwErr.error = errors.New(msg)
						wwErr.errCode = ERROR_ACCOUNT_TRUST_LIMIT_EXCEEDED
						wwErr.errMsg = msg
						ch <- CheckResult{wwErr}
						return
					}
				}
				ch <- CheckResult{}
				return

			}(chBalance, idx)

		}(chBalance, idx)
	}
	// 1. check trusline channel results
	for i := 0; i < len(payments)*2; i++ {
		select {
		case checkResult := <-chTrustline:
			err := checkResult.Error
			if err != nil {
				return err
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when processing sanity check"
			LOGGER.Error(msg)
			wwErr := wwError{}
			wwErr.error = errors.New(msg)
			wwErr.errCode = ERROR_TIMEOUT
			wwErr.errMsg = msg
			return wwErr
		}
	}
	// 2. check balance channel results
	for i := 0; i < len(payments)*2; i++ {
		select {
		case checkResult := <-chBalance:
			err := checkResult.Error
			if err != nil {
				return err
			}
		case <-time.After(time.Second * 10):
			msg := "Error: Timeout when processing sanity check"
			LOGGER.Error(msg)
			wwErr := wwError{}
			wwErr.error = errors.New(msg)
			wwErr.errCode = ERROR_TIMEOUT
			wwErr.errMsg = msg
			return wwErr
		}
	}
	close(chBalance)
	close(chTrustline)
	return nil
}

// CheckTrustline : Check if trustline is established
func CheckTrustline(paymentAccount Account, paymentAsset model.Asset, paymentAmount decimal.Decimal, prc PRClient, hc *horizon.Client) CheckResult {
	scAsset, err := getCreditAsset(paymentAsset, paymentAmount, prc)
	if err != nil {
		msg := "Error getting stellar asset"
		LOGGER.Error(msg)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_GET_STELLAR_ASSET
		wwErr.errMsg = msg
		return CheckResult{wwErr}
	}
	_accountAddress, err := prc.GetParticipantAccount(paymentAccount.ParticipantID, paymentAccount.Account)
	if err != nil {
		msg := "Error getting stellar target address, account:" + paymentAccount.ParticipantID + ":" + paymentAccount.Account
		LOGGER.Error(msg)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_MAP_STELLAR_ADDRESS
		wwErr.errMsg = msg
		return CheckResult{wwErr}
	}
	// checktrustline
	switch scAsset.(type) {
	case b.NativeAmount:
		return CheckResult{nil}
	case b.CreditAmount:
		account, _ := hc.LoadAccount(_accountAddress)
		// check if trustline is established
		isTrust := false
		for _, balance := range account.Balances {
			if balance.Asset.Code == scAsset.(b.CreditAmount).Code &&
				balance.Asset.Issuer == scAsset.(b.CreditAmount).Issuer {
				isTrust = true
				break
			}
		}
		if isTrust == false {
			msg := "Failed Sanity check: Participant " + paymentAccount.ParticipantID + " account: " + paymentAccount.Account + " does not trust asset " + *paymentAsset.AssetCode
			LOGGER.Error(msg)
			wwErr := wwError{}
			wwErr.error = errors.New(msg)
			wwErr.errCode = ERROR_ACCOUNT_DOES_NOT_TRUST_ASSET
			wwErr.errMsg = msg
			return CheckResult{wwErr}
		}
	}
	return CheckResult{nil}
}

func getCreditBalance(a horizon.Account, code string, issuer string) Balance {
	for _, balance := range a.Balances {
		if balance.Asset.Code == code && balance.Asset.Issuer == issuer {
			b, _ := decimal.NewFromString(balance.Balance)
			l, _ := decimal.NewFromString(balance.Limit)
			return Balance{b, l}
		}
	}
	return Balance{
		decimal.NewFromFloat(0),
		decimal.Decimal{},
	}
}

func getCreditAsset(ast model.Asset, amount decimal.Decimal, prclient PRClient) (creditAsset interface{}, err error) {
	creditAsset = b.CreditAmount{}
	if *ast.AssetCode == "xlm" || *ast.AssetCode == "XLM" {
		creditAsset = b.NativeAmount{Amount: amount.Round(7).String()}
	} else {
		astIA, err := prclient.GetParticipantAccount(ast.IssuerID, "issuing")
		if err != nil {
			LOGGER.Error("Error getting asset issuing address from participant registry, issuerID: " + ast.IssuerID)
			return creditAsset, err
		}
		creditAsset = b.CreditAmount{Code: *ast.AssetCode, Issuer: astIA,
			Amount: amount.Round(7).String()}
	}
	return creditAsset, nil
}

//GetOutstandingBalances : GetsOutstanding Balances for specific asset
func GetOutstandingBalances(asset model.Asset, prc PRClient, nhc hClient.ClientInterface) ([]OutstandingBalance, error) {

	participants, err := prc.GetAllParticipants()
	issuerAddress := ""
	if err != nil {
		msg := "Error getting particiapants from participant registry"
		LOGGER.Error(msg)
		LOGGER.Debug(err)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_MAP_STELLAR_ADDRESS
		wwErr.errMsg = msg
		return nil, wwErr
	}

	// find issuer address
	for _, participant := range participants {
		if *participant.ID == asset.IssuerID {
			issuerAddress = participant.IssuingAccount
		}
	}

	type GetBalanceResult struct {
		OutstandingBalance OutstandingBalance
		Error              Error
	}

	//func to get asset balance of a specific account
	getBalance := func(address string, assetCode string, assetIssuerAddress string) (Balance, error) {
		accountRequest := hClient.AccountRequest{AccountID: address}

		// Load the account detail from the network
		account, err := nhc.AccountDetail(accountRequest)
		if err != nil {
			LOGGER.Error("get outstanding balance: ", accountRequest, err)
			return Balance{}, err
		}
		for _, sbalance := range account.Balances {
			if assetCode == sbalance.Code && assetIssuerAddress == sbalance.Issuer {
				// this account trust the asset
				LOGGER.Debug("get outstanding balance: ", sbalance.Code, sbalance.Issuer)
				balance := Balance{}
				balance.Amount, _ = decimal.NewFromString(sbalance.Balance)
				balance.TrustLimit, _ = decimal.NewFromString(sbalance.Limit)
				return balance, nil
			}
		}
		return Balance{}, nil
	}

	accountCount := 0
	resultOutStandingBalance := []OutstandingBalance{}
	chGetBalanceResult := make(chan GetBalanceResult)
	for i, participant := range participants {
		//skim through issuing account
		if participant.IssuingAccount != "" {
			accountCount++
			// goroutine
			go func(i int) {
				balance, _ := getBalance(participants[i].IssuingAccount, *asset.AssetCode, issuerAddress)
				outstandingBalance := OutstandingBalance{}
				outstandingBalance.Account = Account{*participants[i].ID, "issuing"}
				outstandingBalance.Balance = balance
				chGetBalanceResult <- GetBalanceResult{outstandingBalance, nil}
				return
			}(i)
		}
		//skim through operating accounts
		for j, opAccount := range participants[i].OperatingAccounts {
			if *opAccount.Address != "" {
				accountCount++
				// goroutine
				go func(i, j int) {
					opbalance, _ := getBalance(*participants[i].OperatingAccounts[j].Address, *asset.AssetCode, issuerAddress)

					outstandingBalance := OutstandingBalance{}
					outstandingBalance.Account = Account{*participants[i].ID, participants[i].OperatingAccounts[j].Name}
					outstandingBalance.Balance = opbalance
					chGetBalanceResult <- GetBalanceResult{outstandingBalance, nil}
				}(i, j)
			}
		}
	}
	LOGGER.Debug("total accounts count = ", accountCount)
	// collect results outstanding balance
	for k := 0; k < accountCount; k++ {
		select {
		case checkResult := <-chGetBalanceResult:
			err := checkResult.Error
			if err != nil {
				return nil, err
			}
			if checkResult.OutstandingBalance.Balance.TrustLimit.GreaterThan(decimal.NewFromFloat(0)) {
				resultOutStandingBalance = append(resultOutStandingBalance, checkResult.OutstandingBalance)
			}
		case <-time.After(time.Second * 10):
			LOGGER.Debug("Timeout when fetching outstanding balance")
			msg := "Error: Timeout when fetching outstanding balance"
			LOGGER.Error(msg)
			wwErr := wwError{}
			wwErr.error = errors.New(msg)
			wwErr.errCode = ERROR_TIMEOUT
			wwErr.errMsg = msg
			return nil, wwErr
		}
	}
	close(chGetBalanceResult)
	// check if the asset exits, this is put until the end of funciton for performance reason
	if len(resultOutStandingBalance) == 0 {
		isExist, err := checkAssetExist(issuerAddress, *asset.AssetCode, nhc)
		if err != nil {
			LOGGER.Error(err)
			return nil, err
		}
		if !isExist {
			msg := "Asset does not exist"
			wwErr := wwError{}
			wwErr.error = errors.New(msg)
			wwErr.errCode = ERROR_ASSET_DOES_NOT_EXIST
			wwErr.errMsg = msg
			return nil, wwErr
		}
	}
	return resultOutStandingBalance, nil
}

func checkAssetExist(issuerAddress, assetCode string, nhc hClient.ClientInterface) (bool, error) {
	assetList, err := nhc.Assets(hClient.AssetRequest{ForAssetCode: assetCode, ForAssetIssuer: issuerAddress})
	if err != nil {
		msg := "Error getting asset list from horizon"
		LOGGER.Debug(msg, err)
		wwErr := wwError{}
		wwErr.error = err
		wwErr.errCode = ERROR_GET_STELLAR_ASSET
		wwErr.errMsg = msg
		return false, wwErr
	}
	if len(assetList.Embedded.Records) == 0 {
		return false, nil
	}
	return true, nil
}
