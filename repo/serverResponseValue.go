package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/jmoiron/sqlx"
	"github.com/zenmuharom/zenlogger"
)

type DefaultServerResponseValueRepo struct {
	logger zenlogger.Zenlogger
}

type ServerResponseValueRepo interface {
	FindAllByProductCode(productCode string) ([]domain.ServerResponseValue, error)
	FindAllByServerResponseID(productCode string, FieldIDs []int, middlewResponseIDs []int) ([]domain.ServerResponseValue, error)
}

func (repo *DefaultServerResponseValueRepo) FindAllByServerResponseID(productCode string, fieldIDs []int, middlewResponseIDs []int) ([]domain.ServerResponseValue, error) {
	middlewareResponseValueConfs := make([]domain.ServerResponseValue, 0)

	if productCode == "" {
		productCode = "000000"
	}

	sqlStmt := "SELECT srv.id, srv.product_code, srv.field_id, sr.field AS field_name, sr.type AS field_type, srv.middleware_response_id, mr.field AS middleware_response_field, srv.args, srv.condition_field_id, src.field AS condition_field_name, srv.condition_operator, srv.condition_value FROM ( SELECT *, ROW_NUMBER ( ) OVER ( PARTITION BY field_id, middleware_response_id, condition_field_id, condition_operator, condition_value ORDER BY CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END, product_code DESC, id ) AS rn FROM server_response_value WHERE deleted_at IS NULL AND ( product_code = '000000' OR product_code = :productCode ) ) srv LEFT JOIN middleware_response mr ON mr.id = srv.middleware_response_id LEFT JOIN server_response sr ON srv.field_id = sr.id LEFT JOIN server_response src ON srv.condition_field_id = src.id WHERE srv.rn = 1 AND srv.field_id IN ( :field_ids ) ORDER BY srv.id ASC"
	if len(middlewResponseIDs) > 0 {
		sqlStmt = "SELECT srv.id, srv.product_code, srv.field_id, sr.field AS field_name, sr.type AS field_type, srv.middleware_response_id, mr.field AS middleware_response_field, srv.args, srv.condition_field_id, src.field AS condition_field_name, srv.condition_operator, srv.condition_value FROM ( SELECT *, ROW_NUMBER () OVER ( PARTITION BY field_id, middleware_response_id, condition_field_id, condition_operator, condition_value ORDER BY CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END, product_code DESC, id ) AS rn FROM server_response_value WHERE deleted_at IS NULL AND (product_code = '000000' OR product_code = :productCode) ) srv LEFT JOIN middleware_response mr ON mr.id = srv.middleware_response_id LEFT JOIN server_response sr ON srv.field_id = sr.id LEFT JOIN server_response src ON srv.condition_field_id = src.id WHERE srv.rn = 1 AND srv.field_id IN (:field_ids) AND (srv.middleware_response_id IN (:middleware_response_ids) OR srv.middleware_response_id IS NULL) ORDER BY srv.id ASC"
	}

	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "fieldIDs", Value: fieldIDs}, zenlogger.ZenField{Key: "middlewResponseIDs", Value: middlewResponseIDs})

	arg := map[string]interface{}{
		"productCode": productCode,
		"field_ids":   fieldIDs,
	}

	if len(middlewResponseIDs) > 0 {
		arg["middleware_response_ids"] = middlewResponseIDs
	}

	query, args, err := sqlx.Named(sqlStmt, arg)
	if err != nil {
		repo.logger.Error(err.Error())
	}

	// fills qids on query dynamically
	query, args, err = sqlx.In(query, args...)
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

func NewServerResponseValueRepo(logger zenlogger.Zenlogger) ServerResponseValueRepo {
	return &DefaultServerResponseValueRepo{
		logger: logger,
	}
}

func (repo *DefaultServerResponseValueRepo) FindAllByProductCode(productCode string) ([]domain.ServerResponseValue, error) {
	middlewareResponseValueConfs := make([]domain.ServerResponseValue, 0)

	sqlStmt := `SELECT srv.id, srv.product_code, srv.field_id, sr.field AS field_name, sr.type AS field_type, srv.middleware_response_id, mr.field AS middleware_response_field, srv.args, srv.condition_field_id, src.field AS condition_field_name, srv.condition_operator, srv.condition_value FROM server_response_value srv JOIN middleware_response mr ON mr.id = srv.middleware_response_id JOIN server_response sr ON srv.field_id = sr.id LEFT JOIN server_response src ON srv.condition_field_id = src.id WHERE srv.product_code = ? ORDER BY srv.id ASC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "product_code", Value: productCode})

	err := util.GetDB().Select(&middlewareResponseValueConfs, sqlStmt, productCode)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return middlewareResponseValueConfs, nil
}
