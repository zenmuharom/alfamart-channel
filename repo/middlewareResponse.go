package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultMiddlewareResponseRepo struct {
	logger zenlogger.Zenlogger
}

type MiddlewareResponseRepo interface {
	FindAllByMiddlewareAndProcess(middleware, process string) ([]domain.MiddlewareResponse, error)
	FindByProcessAndFieldAs(middleware, process, fieldAs string) (*domain.MiddlewareResponse, error)
}

func NewMiddlewareResponseRepo(logger zenlogger.Zenlogger) MiddlewareResponseRepo {
	return &DefaultMiddlewareResponseRepo{
		logger: logger,
	}
}

func (repo *DefaultMiddlewareResponseRepo) FindAllByMiddlewareAndProcess(middleware, process string) ([]domain.MiddlewareResponse, error) {
	tsAdapterResponseConfig := make([]domain.MiddlewareResponse, 0)

	sqlStmt := `SELECT id, field, type FROM middleware_response WHERE deleted_at IS NULL AND middleware = ? AND process = ?`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "process", Value: process})

	err := util.GetDB().Select(&tsAdapterResponseConfig, sqlStmt, middleware, process)
	if err != nil {
		repo.logger.Error(err.Error())
		return tsAdapterResponseConfig, err
	}

	return tsAdapterResponseConfig, nil
}

func (repo *DefaultMiddlewareResponseRepo) FindByProcessAndFieldAs(middleware, process, fieldAs string) (*domain.MiddlewareResponse, error) {

	var data domain.MiddlewareResponse

	sqlStmt := `SELECT id, field, type FROM middleware_response WHERE deleted_at IS NULL AND middleware = ? AND process = ? AND field_as = ? LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "process", Value: process}, zenlogger.ZenField{Key: "field_as", Value: fieldAs})

	err := util.GetDB().Get(&data, sqlStmt, middleware, process, fieldAs)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return &data, nil
}
