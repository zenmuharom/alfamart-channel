package tool

import (
	"alfamart-channel/domain"

	"github.com/zenmuharom/zenlogger"
)

func CalculateFee(logger zenlogger.Zenlogger, config domain.Fee, amount int) (amountToDebet int) {
	logger.Debug("CalculateFee", zenlogger.ZenField{Key: "config", Value: config}, zenlogger.ZenField{Key: "amount", Value: amount})
	var calculatedFee = 0
	amountToDebet = amount

	feeCaPercentage := float64(amount) * config.FeeCaPercentage.Float64
	feeCaFix := config.FeeCaFix.Int64
	feeBillerPercentage := float64(amount) * config.FeeBillerPercentage.Float64
	feeBillerFix := config.FeeBillerFix.Int64
	feeFinnetPercentage := float64(amount) * config.FeeFinnetPercentage.Float64
	feeFinnetFix := config.FeeFinnetFix.Int64

	calculatedFee += int(feeCaPercentage)
	calculatedFee += int(feeCaFix)
	calculatedFee += int(feeBillerPercentage)
	calculatedFee += int(feeBillerFix)
	calculatedFee += int(feeFinnetPercentage)
	calculatedFee += int(feeFinnetFix)

	if config.IsSurcharge.Bool {
		amountToDebet += calculatedFee
	} else {
		amountToDebet -= calculatedFee
	}

	params := []zenlogger.ZenField{
		{Key: "config", Value: config},
		{Key: "amount", Value: amount},
		{Key: "feeCaPercentage", Value: feeCaPercentage},
		{Key: "feeCaFix", Value: feeCaFix},
		{Key: "feeBillerPercentage", Value: feeBillerPercentage},
		{Key: "feeBillerFix", Value: feeBillerFix},
		{Key: "feeFinnetPercentage", Value: feeFinnetPercentage},
		{Key: "feeFinnetFix", Value: feeFinnetFix},
		{Key: "calculatedFee", Value: calculatedFee},
	}

	logger.Debug("CalculateFee", params...)
	return
}
