package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultMiddlewareRequestRepo struct {
	logger zenlogger.Zenlogger
}

type MiddlewareRequestRepo interface {
	FindAllByAttributes(middlewareRequest domain.MiddlewareRequest) (configs []domain.MiddlewareRequest, err error)
	FindByProcessAndFieldAs(middleware, process, fieldAs string) (*domain.MiddlewareRequest, error)
}

func NewMiddlewareRequestRepo(logger zenlogger.Zenlogger) MiddlewareRequestRepo {
	return &DefaultMiddlewareRequestRepo{
		logger: logger,
	}
}

func (repo *DefaultMiddlewareRequestRepo) FindAllByAttributes(middlewareRequest domain.MiddlewareRequest) (configs []domain.MiddlewareRequest, err error) {

	sqlStmt := `
		SELECT
			id,
			field,
			type,
			middleware,
			condition_field_id,
			condition_operator,
			condition_value,
			product_code,
			created_at,
			updated_at,
			activated_at,
			deleted_at
		FROM
			(
				SELECT
					*,
					ROW_NUMBER() OVER (
						PARTITION BY 
							middleware, 
							process,
							field,
							type
						ORDER BY
							CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END,
							product_code DESC,
							id
					) AS rn
				FROM
					middleware_request
				WHERE
					deleted_at IS NULL
			) t
		WHERE
			t.rn = 1
			AND middleware = ? 
			AND process = ? 
			AND (product_code = '000000' OR product_code = ?)	
	`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "middleware", Value: middlewareRequest.Middleware.String}, zenlogger.ZenField{Key: "process", Value: middlewareRequest.Process.String}, zenlogger.ZenField{Key: "productCode", Value: middlewareRequest.ProductCode.String})

	err = util.GetDB().Select(&configs, sqlStmt, middlewareRequest.Middleware.String, middlewareRequest.Process.String, middlewareRequest.ProductCode.String)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return
}

func (repo *DefaultMiddlewareRequestRepo) FindByProcessAndFieldAs(middleware, process, fieldAs string) (*domain.MiddlewareRequest, error) {

	var data domain.MiddlewareRequest

	sqlStmt := `SELECT id, field, type FROM middleware_request WHERE deleted_at IS NULL AND middleware = ? AND process = ? AND field_as = ? LIMIT 1`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "middleware", Value: middleware}, zenlogger.ZenField{Key: "process", Value: process}, zenlogger.ZenField{Key: "field_as", Value: fieldAs})

	err := util.GetDB().Get(&data, sqlStmt, middleware, process, fieldAs)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return &data, nil
}
