package chaincode_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/restlwsswhy/btc-service-chaincode/chaincode"
	"github.com/restlwsswhy/btc-service-chaincode/chaincode/mocks"
	"github.com/stretchr/testify/require"
)

func prepMocks() (*mocks.TransactionContext, *shimtest.MockStub, *chaincode.SmartContract) {
	stubMock := shimtest.NewMockStub(`btc`, nil)
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(stubMock)

	contract := &chaincode.SmartContract{}

	return transactionContext, stubMock, contract
}

func TestIntegrationCurrency(t *testing.T) {
	ctx, stub, contract := prepMocks()

	// PUT CURRENCY
	curr1 := &chaincode.CurrencyCommon{
		Symbol: "USD",
		RUName: "фыв",
		ENName: "asd",
		Code:   "098",
	}
	curr1b, err := json.Marshal(curr1)
	if err != nil {
		log.Fatal(err)
	}

	curr2 := &chaincode.CurrencyCommon{
		Symbol: "ASD",
		RUName: "павп",
		ENName: "sdgsg",
		Code:   "093",
	}
	curr2b, err := json.Marshal(curr2)
	if err != nil {
		log.Fatal(err)
	}

	stub.MockTransactionStart("1")
	err = contract.PutCurrency(ctx, string(curr1b))
	stub.MockTransactionEnd("1")
	require.NoError(t, err)

	stub.MockTransactionStart("2")
	err = contract.PutCurrency(ctx, string(curr2b))
	stub.MockTransactionEnd("2")
	require.NoError(t, err)

	// GET CURRENCY
	stub.MockTransactionStart("3")
	res, err := contract.GetCurrency(ctx, curr1.Symbol)
	stub.MockTransactionEnd("3")
	require.NoError(t, err)
	require.Equal(t, curr1, res)

	stub.MockTransactionStart("4")
	res, err = contract.GetCurrency(ctx, curr2.Symbol)
	stub.MockTransactionEnd("4")
	require.NoError(t, err)
	require.Equal(t, curr2, res)
}

func TestIntegrationPrices(t *testing.T) {
	ctx, stub, contract := prepMocks()

	// PUT PRICE
	price1 := &chaincode.QuotaDetailCommon{
		Symbol: "USD",
		Buy:    123,
		Sell:   123,
		Time:   time.Now().Unix() - 100,
	}
	price1b, err := json.Marshal(price1)
	if err != nil {
		log.Fatal(err)
	}

	price2 := &chaincode.QuotaDetailCommon{
		Symbol: "USD",
		Buy:    123,
		Sell:   123,
		Time:   time.Now().Unix() - 10,
	}
	price2b, err := json.Marshal(price2)
	if err != nil {
		log.Fatal(err)
	}

	price3 := &chaincode.QuotaDetailCommon{
		Symbol: "ASD",
		Buy:    123,
		Sell:   123,
		Time:   time.Now().Unix() - 100,
	}
	price3b, err := json.Marshal(price3)
	if err != nil {
		log.Fatal(err)
	}

	price4 := &chaincode.QuotaDetailCommon{
		Symbol: "ASD",
		Buy:    123,
		Sell:   123,
		Time:   time.Now().Unix() - 10,
	}
	price4b, err := json.Marshal(price4)
	if err != nil {
		log.Fatal(err)
	}

	stub.MockTransactionStart("1")
	err = contract.PutCurrencyPrice(ctx, string(price1b))
	stub.MockTransactionEnd("1")
	require.NoError(t, err)
	stub.MockTransactionStart("2")
	err = contract.PutCurrencyPrice(ctx, string(price2b))
	stub.MockTransactionEnd("2")
	require.NoError(t, err)
	stub.MockTransactionStart("3")
	err = contract.PutCurrencyPrice(ctx, string(price3b))
	stub.MockTransactionEnd("3")
	require.NoError(t, err)
	stub.MockTransactionStart("4")
	err = contract.PutCurrencyPrice(ctx, string(price4b))
	stub.MockTransactionEnd("4")
	require.NoError(t, err)

	// GET LAST PRICE
	stub.MockTransactionStart("5")
	res, err := contract.GetCurrencyPrice(ctx, price1.Symbol)
	stub.MockTransactionEnd("5")
	require.NoError(t, err)
	require.Equal(t, price2, res)

	stub.MockTransactionStart("6")
	res, err = contract.GetCurrencyPrice(ctx, price3.Symbol)
	stub.MockTransactionEnd("6")
	require.NoError(t, err)
	require.Equal(t, price4, res)

	// GET HISTORY PRICE
	stub.MockTransactionStart("7")
	history, err := contract.GetCurrencyPriceFromHistory(ctx, price1.Symbol)
	stub.MockTransactionEnd("7")
	require.NoError(t, err)
	require.Equal(t, []*chaincode.QuotaDetailCommon{price1, price2}, history)

	// GET ALL PRICES
	symbols := []string{price1.Symbol, price3.Symbol}
	symbolsb, err := json.Marshal(symbols)
	if err != nil {
		log.Fatal(err)
	}

	stub.MockTransactionStart("7")
	all, err := contract.GetAllCurrentPrices(ctx, string(symbolsb))
	stub.MockTransactionEnd("7")
	require.NoError(t, err)
	require.Equal(t, []*chaincode.QuotaDetailCommon{price2, price4}, all)
}
