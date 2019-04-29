package payment

import (
	"context"

	transactionSDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/entity/transaction"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
)

// Payer a payment contract should be complied
type Payer interface {
	Pay(ctx context.Context, request transactionSDK.Payload) (result transactionSDK.Response)
}

// Payment is an inquiry schema which can be used for inquiry transaction
type Payment struct {
	payers map[string]Payer
}

// New will instantiate a new instance of Inquiry Type
func New(payers map[string]Payer) *Payment {
	if len(payers) == 0 {
		log.Error("incomplete dependencies")
		return nil
	}

	return &Payment{
		payers: payers,
	}
}

// Handle is a method that keep the whole inquiry flow,
// all transaction status will be handled in this blocks
func (p *Payment) Handle(ctx context.Context, payload transactionSDK.Payload) (result transactionSDK.Response) {
	productPayer, ok := p.payers[payload.BillerProductCode]
	if !ok {
		log.Warn("no payment service available for that billerProductCode", koollog.String("biller_product_code", payload.BillerProductCode))
		return
	}

	result = productPayer.Pay(ctx, payload)
	return
}
