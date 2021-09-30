// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package listeners

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/stellar/go/clients/horizon"
	"github.com/GFTN/gftn-services/gftn-models/model"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
	"github.com/GFTN/gftn-services/utility/kafka"
)

// StellarAccountPaymentListener is listening for a new payments received by ReceivingAccount
type StellarAccountPaymentListener struct {
	horizon horizon.ClientInterface
	account model.Account
	cursor  horizon.Cursor
	//wg        sync.WaitGroup
	channel   chan string
	initFetch bool
}

// NewStellarAccountStellarAccountPaymentListener creates a new StellarAccountPaymentListener
func NewStellarAccountStellarAccountPaymentListener(
	horizon horizon.ClientInterface,
) (pl StellarAccountPaymentListener, err error) {
	pl.horizon = horizon
	pl.account = model.Account{}
	ch := make(chan string)
	pl.channel = ch
	pl.initFetch = true

	return pl, nil
}

//SetDistAccount - Save account ID in memory
func (pl *StellarAccountPaymentListener) SetDistAccount(account model.Account) {
	pl.account.Address = account.Address
	pl.account.Name = account.Name
	return
}

// Listen starts listening for new payments
func (pl *StellarAccountPaymentListener) Listen(ch chan string) (err error) {

	if *pl.account.Address == "" {
		err = errors.New("account ID is empty to listen")
		return
	}
	_, err = pl.horizon.LoadAccount(*pl.account.Address)
	if err != nil {
		LOGGER.Errorf("Account %v is invalid on Stellar NW", *pl.account.Address)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)
	//defer cancel()

	//routine to listen on payments stream
	go pl.ListenLoop(ctx, ch, *pl.account.Address, cancel, &wg)

	//routine to receive signal from parent listener
	go func(ch chan string, cancelFunc context.CancelFunc, loopContext context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case p := <-ch:
				LOGGER.Infof("Listener on account %v got signal: %v \n", pl.account.Address, p)
				cancelFunc()
				time.Sleep(1 * time.Second)

				return
			case <-loopContext.Done():
				LOGGER.Infof("Done: %v ", pl.account.Address)
				return
			}
		}
	}(ch, cancel, ctx, &wg)

	//routine to wait for wait group to end
	go func(ch chan string, wg *sync.WaitGroup) {
		wg.Wait()
		LOGGER.Debugf("*************** Restarting Listener %v ***************", pl.account.Address)
		time.Sleep(1 * time.Second)
		cursor, err := dynamoClient.GetAccountCursor(pl.account.Name)
		if err != nil {
			LOGGER.Errorf("Encounter error while updating new cursor for account: %v", pl.account.Name)
			return
		}
		pl.cursor = horizon.Cursor(cursor.Cursor)
		err = pl.Listen(ch)
	}(ch, &wg)
	return
}

func (pl *StellarAccountPaymentListener) ListenLoop(ctx context.Context, ch chan string, accountID string, cancelFunc context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		LOGGER.Debugf("Started listening for new payments: account: %v, cursor: %v\n", accountID, pl.cursor)
		err := pl.horizon.StreamPayments(ctx, accountID, &pl.cursor, pl.onPayment)
		if err != nil {
			LOGGER.Error("Error while streamining ", err.Error())
			go func() { cancelFunc() }()
		}
	}()

	select {
	case p := <-ch:
		LOGGER.Debugf("Got signal %v, account Name: %v", p, pl.account.Name)
	case <-ctx.Done():
		LOGGER.Debugf("ListenLoop Done: %v", pl.account.Name)
		return
	}
}

//onPayment - Callback onPayment method
func (pl *StellarAccountPaymentListener) onPayment(payment horizon.Payment) {

	LOGGER.Infof("onPayment called")
	transaction, err := pl.horizon.LoadTransaction(payment.TransactionHash)
	cursor := transaction.PagingToken()
	if err != nil {
		LOGGER.Errorf("Error while retrieving transaction %s: %s, cursor: %s, %s", payment.TransactionHash, err.Error(),
			cursor, transaction.Memo)
		return
	}

	LOGGER.Infof("id %s New received payment %s, Cursor:%s, %s ", payment.ID, payment.TransactionHash,
		cursor, transaction.Memo)

	//paging token is the last read cursor, save it
	pl.cursor = horizon.Cursor(cursor)

	//send notification to gateway
	pNotification := model.Receive{}
	pNotification.TransactionID = &payment.TransactionHash
	pNotification.TransactionMemo = payment.Memo.Value
	pNotification.TransactionReference = &cursor
	pNotification.AccountName = &pl.account.Name

	if !pl.initFetch {

		/*
			sending message to Kafka
		*/

		ofiId := os.Getenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME)

		b, _ := json.Marshal(pNotification)
		err = kafkaOperations.Produce(ofiId+"_"+kafka.PAYMENT_TOPIC, b)
		if err != nil {
			LOGGER.Errorf("Kafka producer failed")
			return
		}
		/*
			write to dynamo
		*/

		err = dynamoClient.UpdateAccountCursor(pl.account.Name, string(cursor))
		if err != nil {
			LOGGER.Error("Error updating cursor to Dynamo DB: %v", err.Error())
		}
	}
	pl.initFetch = false
}
