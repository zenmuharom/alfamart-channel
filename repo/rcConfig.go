package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultRcConfigRepo struct {
	logger zenlogger.Zenlogger
}

type RcConfigRepo interface {
	FindRc(productCode string, GoCode int64) (*domain.RcConfig, error)
}

func NewRcConfigRepo(logger zenlogger.Zenlogger) RcConfigRepo {
	return &DefaultRcConfigRepo{
		logger: logger,
	}
}

func (repo *DefaultRcConfigRepo) FindRc(productCode string, GoCode int64) (*domain.RcConfig, error) {
	var rcConfig domain.RcConfig

	sqlStmt := `SELECT go_code, code, product_code, httpstatus, desc_eng, desc_id, created_at, updated_at, deleted_at FROM rc_config WHERE (product_code = ? OR product_code = '000000') AND go_code = ? ORDER BY CASE WHEN product_code <> '000000' THEN 0 ELSE 1 END, id LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "product_code", Value: productCode}, zenlogger.ZenField{Key: "go_code", Value: GoCode})

	err := util.GetDB().Get(&rcConfig, sqlStmt, productCode, GoCode)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	repo.logger.Debug("found RC", zenlogger.ZenField{Key: "data", Value: rcConfig})

	return &rcConfig, nil
}
