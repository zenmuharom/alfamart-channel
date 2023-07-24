package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultProductRepo struct {
	logger zenlogger.Zenlogger
}

type ProductRepo interface {
	Find(productCode string) (domain.Product, error)
}

func NewProductRepo(logger zenlogger.Zenlogger) ProductRepo {
	return &DefaultProductRepo{
		logger: logger,
	}
}

func (repo *DefaultProductRepo) Find(code string) (domain.Product, error) {
	var mitraConfig domain.Product

	sqlStmt := `SELECT code, name, bit103, bit104, mitra_id FROM product WHERE deleted_at IS NULL AND code = ? LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "code", Value: code})

	err := util.GetDB().Get(&mitraConfig, sqlStmt, code)
	if err != nil {
		repo.logger.Error(err.Error())
		return mitraConfig, err
	}

	return mitraConfig, nil
}
