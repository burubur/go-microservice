package checkstatus

import (
	"context"

	trxSDK "bitbucket.org/kudoindonesia/frontier_biller_sdk/entity/transaction"
	"bitbucket.org/kudoindonesia/frontier_biller_sdk/log"
	"bitbucket.org/kudoindonesia/koolkit/koollog"
)

const csLabel = "CheckStatus"

// StatusChecker an interface for CheckStatus functionality
type StatusChecker interface {
	Check(ctx context.Context, request trxSDK.Payload) (result trxSDK.Response)
}

// CheckStatus a instance that responsible for all checkstatus status request
type CheckStatus struct {
	statusChecker map[string]StatusChecker
}

// New will instantiate a new instance of checkstatus.CheckStatus type
func New(statusChecker map[string]StatusChecker) *CheckStatus {
	if len(statusChecker) == 0 {
		log.Warn("incomplete checkstatus dependencies")
		return nil
	}

	return &CheckStatus{
		statusChecker: statusChecker,
	}
}

// Handle is a method that keep the whole inquiry flow,
// all transaction status will be handled in this blocks
func (cs *CheckStatus) Handle(ctx context.Context, payload trxSDK.Payload) (result trxSDK.Response) {
	productStatusChecker, ok := cs.statusChecker[payload.BillerProductCode]
	if !ok {
		log.Warn("no checkstatus service available for that billerProductCode", koollog.String("biller_product_code", payload.BillerProductCode))
		return
	}

	result = productStatusChecker.Check(ctx, payload)
	return
}
