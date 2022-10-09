package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	p "github.com/restlwsswhy/btc-service-chaincode/proto"
	"google.golang.org/protobuf/proto"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) PutCurrency(ctx contractapi.TransactionContextInterface, data string) error {
	if data == "" {
		return fmt.Errorf("empty input data error")
	}

	curr := &CurrencyCommon{}
	if err := json.Unmarshal([]byte(data), curr); err != nil {
		return fmt.Errorf("unmarshal input data error: %w", err)
	}

	b, err := proto.Marshal(curr.ConvertProto())
	if err != nil {
		return fmt.Errorf("marshal to proto erro: %v", err)
	}

	if err := ctx.GetStub().PutState(curr.Symbol, b); err != nil {
		return fmt.Errorf("put currency to state error: %w", err)
	}

	return nil
}

func (s *SmartContract) GetCurrency(ctx contractapi.TransactionContextInterface, symbol string) (*CurrencyCommon, error) {
	if symbol == "" {
		return nil, fmt.Errorf("empty input data error")
	}

	b, err := ctx.GetStub().GetState(symbol)
	if err != nil {
		return nil, fmt.Errorf("get currency from state error: %w", err)
	}
	if b == nil {
		return nil, fmt.Errorf("currency is not exist")
	}

	currProto := &p.Currency{}
	if err := proto.Unmarshal(b, currProto); err != nil {
		return nil, fmt.Errorf("unmarshal state data from proto error: %w", err)
	}

	return &CurrencyCommon{
		Symbol: symbol,
		RUName: currProto.GetRuName(),
		ENName: currProto.GetEnName(),
		Code:   currProto.GetCode(),
	}, nil
}

func (s *SmartContract) PutCurrencyPrice(ctx contractapi.TransactionContextInterface, data string) error {
	if data == "" {
		return fmt.Errorf("empty input data error")
	}

	quota := &QuotaDetailCommon{}
	if err := json.Unmarshal([]byte(data), quota); err != nil {
		return fmt.Errorf("unmarshal input data error: %w", err)
	}

	b, err := proto.Marshal(quota.ConvertProto())
	if err != nil {
		return fmt.Errorf("marshal to proto erro: %v", err)
	}

	quotaKey := fmt.Sprintf("%s.%d", quota.Symbol, quota.Time)
	if err := ctx.GetStub().PutState(quotaKey, b); err != nil {
		return fmt.Errorf("put currency to state error: %w", err)
	}

	return nil
}

func (s *SmartContract) GetCurrencyPrice(ctx contractapi.TransactionContextInterface, symbol string) (*QuotaDetailCommon, error) {
	if symbol == "" {
		return nil, fmt.Errorf("empty input data error")
	}

	t, _ := ctx.GetStub().GetTxTimestamp()
	quotaKeya := fmt.Sprintf("%s.%d", symbol, t.Seconds-3000)
	quotaKeyb := fmt.Sprintf("%s.%d", symbol, t.Seconds)

	resultsIterator, err := ctx.GetStub().GetStateByRange(quotaKeya, quotaKeyb)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var quotas [][]byte
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		quotas = append(quotas, queryResult.Value)
	}

	if len(quotas) == 0 {
		return nil, fmt.Errorf("quotas not found")
	}

	lastQuoataByte := quotas[len(quotas)-1]
	quota := &p.QuotaDetail{}
	if err := proto.Unmarshal(lastQuoataByte, quota); err != nil {
		return nil, fmt.Errorf("unmarshal state data from proto error: %w", err)
	}

	return &QuotaDetailCommon{
		Symbol: symbol,
		Buy:    quota.GetBuy(),
		Sell:   quota.GetSell(),
		Time:   quota.GetTime(),
	}, nil
}

func (s *SmartContract) GetAllCurrentPrices(ctx contractapi.TransactionContextInterface, currencies string) ([]*QuotaDetailCommon, error) {
	curr := make([]string, 0)
	if err := json.Unmarshal([]byte(currencies), &curr); err != nil {
		return nil, fmt.Errorf("unmarshal body error: %v", err)
	}

	if len(curr) == 0 {
		return nil, fmt.Errorf("empty input data")
	}

	res := make([]*QuotaDetailCommon, 0, len(curr))
	for _, v := range curr {
		q, err := s.GetCurrencyPrice(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("get currency price error: %v", err)
		}

		res = append(res, q)
	}

	return res, nil
}

func (s *SmartContract) GetCurrencyPriceFromHistory(ctx contractapi.TransactionContextInterface, symbol string) ([]*QuotaDetailCommon, error) {
	if symbol == "" {
		return nil, fmt.Errorf("empty input data error")
	}

	t, _ := ctx.GetStub().GetTxTimestamp()
	quotaKeya := fmt.Sprintf("%s.%d", symbol, t.Seconds-3600)
	quotaKeyb := fmt.Sprintf("%s.%d", symbol, t.Seconds)

	resultsIterator, err := ctx.GetStub().GetStateByRange(quotaKeya, quotaKeyb)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var quotas []*QuotaDetailCommon
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		quota := &p.QuotaDetail{}
		if err := proto.Unmarshal(queryResult.Value, quota); err != nil {
			return nil, fmt.Errorf("unmarshal state data from proto error: %w", err)
		}

		quotas = append(quotas, &QuotaDetailCommon{
			Symbol: symbol,
			Buy: quota.GetBuy(),
			Sell: quota.GetSell(),
			Time: quota.GetTime(),
		})
	}

	return quotas, nil
}
