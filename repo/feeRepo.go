package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultFeeRepo struct {
	logger zenlogger.Zenlogger
}

type FeeRepo interface {
	Find(username, productCode, amount string) (fee *domain.Fee, err error)
}

func NewFeeRepo(logger zenlogger.Zenlogger) FeeRepo {
	return &DefaultFeeRepo{
		logger: logger,
	}
}

func (repo *DefaultFeeRepo) Find(username, productCode, amount string) (fee *domain.Fee, err error) {
	var data domain.Fee
	sqlStmt := `SELECT id, username, product_code, fee_min, fee_max, fee_ca_fix, fee_ca_percentage, fee_biller_fix, fee_biller_percentage, fee_finnet_fix, fee_finnet_percentage, is_surcharge FROM fee WHERE deleted_at IS NULL AND username = ? AND product_code = ? AND ? BETWEEN fee_min AND fee_max ORDER BY id DESC LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "username", Value: username}, zenlogger.ZenField{Key: "product_code", Value: productCode}, zenlogger.ZenField{Key: "amount", Value: amount})

	err = util.GetDB().Get(&data, sqlStmt, username, productCode, amount)
	if err != nil {
		repo.logger.Error(err.Error())
	} else {
		fee = &data
	}
	return
}
