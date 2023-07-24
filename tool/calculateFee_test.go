package tool

import (
	"alfamart-channel/domain"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zenmuharom/zenlogger"
)

func TestCalculateFee(t *testing.T) {

	logger := zenlogger.NewZenlogger()
	config := domain.Fee{
		Id:                  2,
		Username:            sql.NullString{String: "finnetDev"},
		ProductCode:         sql.NullString{String: "001001"},
		FeeMin:              sql.NullFloat64{Float64: 10000},
		FeeMax:              sql.NullFloat64{Float64: 20000},
		FeeCaFix:            sql.NullInt64{Int64: 300},
		FeeCaPercentage:     sql.NullFloat64{Float64: 0.03},
		FeeBillerFix:        sql.NullInt64{Int64: 200},
		FeeBillerPercentage: sql.NullFloat64{Float64: 0.02},
		FeeFinnetFix:        sql.NullInt64{Int64: 100},
		FeeFinnetPercentage: sql.NullFloat64{Float64: 0.01},
		IsSurcharge:         sql.NullBool{Bool: true},
	}
	calculatedFee := CalculateFee(logger, config, 10000)
	require.Equal(t, "500", calculatedFee)
	fmt.Println(calculatedFee)

}
