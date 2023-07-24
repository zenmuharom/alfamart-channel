package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/zenmuharom/zenlogger"
)

type DefaultMiddlewareRequestValueRepo struct {
	logger zenlogger.Zenlogger
}

type MiddlewareRequestValueRepo interface {
	FindAllByAttributes(domain.MiddlewareRequestValue) ([]domain.MiddlewareRequestValue, error)
	FindAllByFieldIDs(fieldIDs []int) ([]domain.MiddlewareRequestValue, error)
}

func NewMiddlewareRequestValueRepo(logger zenlogger.Zenlogger) MiddlewareRequestValueRepo {
	return &DefaultMiddlewareRequestValueRepo{
		logger: logger,
	}
}

func (repo *DefaultMiddlewareRequestValueRepo) FindAllByAttributes(find domain.MiddlewareRequestValue) ([]domain.MiddlewareRequestValue, error) {
	middlewareResponseValueConfs := make([]domain.MiddlewareRequestValue, 0)

	sqlStmt := `SELECT srv.id, srv.product_code, srv.field_id, sr.field AS field_name, sr.type AS field_type, srv.server_request_id, mr.field AS server_request_field, mrp.field AS server_request_parent_field, srv.args, srv.condition_field_id, src.field AS condition_field_name, srv.condition_operator, srv.condition_value FROM middleware_request_value srv LEFT JOIN server_request mr ON mr.id = srv.server_request_id JOIN middleware_request sr ON srv.field_id = sr.id LEFT JOIN middleware_request src ON srv.condition_field_id = src.id LEFT JOIN server_request mrp ON mr.parent_id = mrp.id WHERE srv.ca_code = ? ORDER BY srv.id ASC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "ca_code", Value: find.CaCode.String})

	err := util.GetDB().Select(&middlewareResponseValueConfs, sqlStmt, find.CaCode.String)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return middlewareResponseValueConfs, nil
}

func (repo *DefaultMiddlewareRequestValueRepo) FindAllByFieldIDs(fieldIDs []int) ([]domain.MiddlewareRequestValue, error) {
	middlewareResponseValueConfs := make([]domain.MiddlewareRequestValue, 0)

	sqlStmt := "SELECT srv.id, srv.product_code, srv.field_id, sr.field AS field_name, sr.type AS field_type, srv.server_request_id, mr.field AS server_request_field, mrp.field AS server_request_parent_field, srv.args, srv.condition_field_id, src.field AS condition_field_name, srv.condition_operator, srv.condition_value FROM middleware_request_value srv LEFT JOIN server_request mr ON mr.id = srv.server_request_id JOIN middleware_request sr ON srv.field_id = sr.id LEFT JOIN middleware_request src ON srv.condition_field_id = src.id LEFT JOIN server_request mrp ON mr.parent_id = mrp.id AND mrp.id <> 0 WHERE srv.field_id IN ( ? ) ORDER BY srv.id ASC"
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "fieldIDs", Value: fieldIDs})

	if len(fieldIDs) == 0 {
		return middlewareResponseValueConfs, errors.New("Middleware request not set yet")
	}

	// fills qids on query dynamically
	query, args, err := sqlx.In(sqlStmt, fieldIDs)
	if err != nil {
		repo.logger.Error(err.Error())
	}

	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = util.GetDB().Rebind(query)

	err = util.GetDB().Select(&middlewareResponseValueConfs, query, args...)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return middlewareResponseValueConfs, nil

}
