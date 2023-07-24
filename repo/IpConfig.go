package repo

import (
	"alfamart-channel/util"
	"database/sql"

	"github.com/zenmuharom/zenlogger"
)

type DefaultIpConfigRepo struct {
	logger zenlogger.Zenlogger
}

type IpConfigRepo interface {
	FindIp(ClientIP string) (bool, error)
	FindIpByIpUsernameAndProductCode(ClientIP string, userName string, productCode string) (bool, error)
}

func NewIpConfigRepo(logger zenlogger.Zenlogger) IpConfigRepo {
	return &DefaultIpConfigRepo{
		logger: logger,
	}
}

func (repo *DefaultIpConfigRepo) FindIp(ClientIP string) (bool, error) {
	var ip string

	sqlStmt := `SELECT ip FROM ip_config WHERE deleted_at IS NULL AND (ip = ? OR ip = '0.0.0.0') AND status = 1 LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "ClientIP", Value: ClientIP})

	err := util.GetDB().Get(&ip, sqlStmt, ClientIP)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false, nil
		}

		repo.logger.Error(err.Error())
		return false, err
	}

	return true, nil
}

func (repo *DefaultIpConfigRepo) FindIpByIpUsernameAndProductCode(ClientIP string, userName string, productCode string) (bool, error) {
	var ip string

	sqlStmt := `SELECT ip FROM ip_config WHERE deleted_at IS NULL AND (ip = ? OR ip = '0.0.0.0') AND (username IS NULL OR username = ?) AND status = 1 AND (product_code = '000000' OR product_code = ?) LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "ClientIP", Value: ClientIP}, zenlogger.ZenField{Key: "userName", Value: userName}, zenlogger.ZenField{Key: "ProductCode", Value: productCode})

	err := util.GetDB().Get(&ip, sqlStmt, ClientIP, userName, productCode)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return false, nil
		}

		repo.logger.Error(err.Error())
		return false, err
	}

	return true, nil
}
