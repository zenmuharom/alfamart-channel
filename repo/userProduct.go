package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultUserProductRepo struct {
	logger zenlogger.Zenlogger
}

type UserProductRepo interface {
	Find(find domain.UserProduct) (domain.UserProduct, error)
}

func NewUserProductRepo(logger zenlogger.Zenlogger) UserProductRepo {
	return &DefaultUserProductRepo{
		logger: logger,
	}
}

func (repo *DefaultUserProductRepo) Find(find domain.UserProduct) (domain.UserProduct, error) {
	var UserProduct domain.UserProduct

	sqlStmt := `SELECT t.id, t.username, t.product_code, t.bit18, t.bit32, t.bit33, t.bit62, t.product_code_mapped, product.account_number AS account_number, t.rc_success, t.rc_error_nondebit, t.timeout_biller, t.switching_url FROM alfamart_product t JOIN product ON product.CODE = t.product_code WHERE t.deleted_at IS NULL AND t.username = ? AND t.product_code = ? LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "Username", Value: find.Username.String}, zenlogger.ZenField{Key: "ProductCode", Value: find.ProductCode.String})

	err := util.GetDB().Get(&UserProduct, sqlStmt, find.Username.String, find.ProductCode.String)
	if err != nil {
		repo.logger.Error(err.Error())
		return UserProduct, err
	}

	return UserProduct, nil
}
