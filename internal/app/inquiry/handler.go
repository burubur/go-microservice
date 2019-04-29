package inquiry

import (
	"context"
	"errors"

	inquirySDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/entity/inquiry"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
)

// Inquirer a defined Inquiry contract should be complied
type Inquirer interface {
	Inquire(ctx context.Context, request inquirySDK.Request) (result inquirySDK.Response)
}

// Inquiry is an inquiry schema which can be used for inquiry transaction
type Inquiry struct {
	inquirers map[string]Inquirer
}

// New will instantiate a new instance of Inquiry Type
func New(inquirers map[string]Inquirer) *Inquiry {
	if len(inquirers) == 0 {
		log.Warn("incomplete dependencies")
		return nil
	}

	return &Inquiry{
		inquirers: inquirers,
	}
}

// Handle is a method that keep the whole inquiry flow,
// all transaction status will be handled in this blocks
func (i *Inquiry) Handle(ctx context.Context, payload inquirySDK.Request) (result inquirySDK.Response, err error) {
	productInquirer, ok := i.inquirers[payload.BillerProductCode]
	if !ok {
		err = errors.New("no available inquiry service")
		log.Warn("no inquiry service available for that billerProductCode", koollog.String("biller_product_code", payload.BillerProductCode))
		return
	}

	result = productInquirer.Inquire(ctx, payload)
	return
}
