package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultRouteRepo struct {
	logger zenlogger.Zenlogger
}

type RouteRepo interface {
	FindAll() (routes []domain.Route, err error)
}

func NewRouteRepo(logger zenlogger.Zenlogger) RouteRepo {
	return &DefaultRouteRepo{
		logger: logger,
	}
}

func (repo *DefaultRouteRepo) FindAll() (routes []domain.Route, err error) {
	datas := make([]domain.Route, 0)
	sqlStmt := `SELECT id, method, path, created_at, updated_at, deleted_at FROM route WHERE deleted_at IS NULL`
	repo.logger.Debug(sqlStmt)

	err = util.GetDB().Select(&datas, sqlStmt)
	if err != nil {
		repo.logger.Error(err.Error())
		return
	}
	routes = datas
	return
}
