package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"
	"database/sql"
	"errors"

	"github.com/zenmuharom/zenlogger"
)

type DefaultUserRepo struct {
	logger zenlogger.Zenlogger
}

type UserRepo interface {
	Find(username string) (domain.User, error)
}

func NewUserRepo(logger zenlogger.Zenlogger) UserRepo {
	return &DefaultUserRepo{
		logger: logger,
	}
}

func (repo *DefaultUserRepo) Find(username string) (domain.User, error) {
	var user domain.User

	sqlStmt := `SELECT username, password, is_deposit, neva_username, neva_password, account_number, mitraco FROM user WHERE deleted_at IS NULL AND username = ? LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "username", Value: username})

	err := util.GetDB().Get(&user, sqlStmt, username)
	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("username not found")
		}

		repo.logger.Error(err.Error())
		return user, err
	}

	return user, nil
}
