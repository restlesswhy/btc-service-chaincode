package chaincode

import "github.com/restlwsswhy/btc-service-chaincode/proto"

type QuotaDetailCommon struct {
	Symbol string `json:"symbol"`
	Buy    int64  `json:"buy"`
	Sell   int64  `json:"sell"`
	Time   int64  `json:"time"`
}

func (q QuotaDetailCommon) ConvertProto() *proto.QuotaDetail {
	return &proto.QuotaDetail{
		Buy:  q.Buy,
		Sell: q.Sell,
		Time: q.Time,
	}
}

type CurrencyCommon struct {
	Symbol string `json:"symbol"`
	RUName string `json:"ru_name"`
	ENName string `json:"en_name"`
	Code   string `json:"code"`
}

func (c CurrencyCommon) ConvertProto() *proto.Currency {
	return &proto.Currency{
		RuName: c.RUName,
		EnName: c.ENName,
		Code:   c.Code,
	}
}
