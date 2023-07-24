package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultPrefixSuffixRepo struct {
	logger zenlogger.Zenlogger
}

type PrefixSuffixRepo interface {
	FindAllByProductCode(productCode string) ([]domain.PrefixSuffixConfig, error)
}

func NewPrefixSuffixRepo(logger zenlogger.Zenlogger) PrefixSuffixRepo {
	return &DefaultPrefixSuffixRepo{logger: logger}
}

func (repo *DefaultPrefixSuffixRepo) FindAllByProductCode(productCode string) ([]domain.PrefixSuffixConfig, error) {
	prefixSuffixConfigs := make([]domain.PrefixSuffixConfig, 0)

	sqlStmt := `SELECT id, presuf_type, fix_number_start, fix_number_end, fix_character, prefix, suffix_length, length, product_code FROM prefix_suffix_config WHERE product_code = ?`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "product_code", Value: productCode})

	err := util.GetDB().Select(&prefixSuffixConfigs, sqlStmt, productCode)
	if err != nil {
		repo.logger.Error(err.Error())
		return prefixSuffixConfigs, err
	}

	return prefixSuffixConfigs, nil
}
