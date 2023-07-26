package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultServerResponseRepo struct {
	logger zenlogger.Zenlogger
}

type ServerResponseRepo interface {
	FindAllByProductCodeAndEndpoint(productCode, endpoint string) ([]domain.ServerResponse, error)
	FindAllByProductCodeAndEndpointParentOnly(productCode, endpoint string) ([]domain.ServerResponse, error)
	FindByEndpointAndFieldAs(endpoint, fieldAs string) (*domain.ServerResponse, error)
}

func NewServerResponseRepo(logger zenlogger.Zenlogger) ServerResponseRepo {
	return &DefaultServerResponseRepo{
		logger: logger,
	}
}

func (repo *DefaultServerResponseRepo) FindAllByProductCodeAndEndpoint(productCode, endpoint string) ([]domain.ServerResponse, error) {
	middlewareResponseValueConfs := make([]domain.ServerResponse, 0)

	sqlStmt := `SELECT sr.id, sr.order AS 'order', sr.field, sr.type, sr.parent_id, sr.product_code, srp.field AS field_parent, srp.type AS field_parent_type, sr.condition_field_id AS condition_field_id, srp2.field AS condition_field_name, sr.condition_operator AS condition_operator, sr.condition_value AS condition_value FROM ( SELECT *, ROW_NUMBER ( ) OVER ( PARTITION BY endpoint, field, parent_id, product_code, condition_field_id, condition_operator, condition_value ORDER BY CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END, product_code DESC, id ) AS rn FROM server_response WHERE deleted_at IS NULL AND ( product_code = '000000' OR product_code = ? ) ) sr LEFT JOIN server_response srp ON srp.id = sr.parent_id LEFT JOIN server_response srp2 ON srp2.id = sr.condition_field_id WHERE sr.rn = 1 AND sr.deleted_at IS NULL AND sr.endpoint = ? ORDER BY sr.order ASC, sr.parent_id DESC, sr.id DESC, srp.id DESC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "product_code", Value: productCode}, zenlogger.ZenField{Key: "endpoint", Value: endpoint})

	err := util.GetDB().Select(&middlewareResponseValueConfs, sqlStmt, productCode, endpoint)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return middlewareResponseValueConfs, nil
}

func (repo *DefaultServerResponseRepo) FindAllByProductCodeAndEndpointParentOnly(productCode, endpoint string) ([]domain.ServerResponse, error) {
	middlewareResponseValueConfs := make([]domain.ServerResponse, 0)

	sqlStmt := `SELECT sr.id, sr.ORDER AS 'order', sr.field, sr.type, sr.parent_id, sr.product_code, srp.field AS field_parent, srp.type AS field_parent_type, sr.condition_field_id AS condition_field_id, srp2.field AS condition_field_name, sr.condition_operator AS condition_operator, sr.condition_value AS condition_value FROM ( SELECT *, ROW_NUMBER ( ) OVER ( PARTITION BY endpoint, field, parent_id, product_code, condition_field_id, condition_operator, condition_value ORDER BY CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END, product_code DESC, id ) AS rn FROM server_response WHERE deleted_at IS NULL AND ( product_code = '000000' OR product_code = ? ) ) sr LEFT JOIN server_response srp ON srp.id = sr.parent_id LEFT JOIN server_response srp2 ON srp2.id = sr.condition_field_id WHERE sr.rn = 1 AND sr.deleted_at IS NULL AND sr.endpoint = ? AND sr.parent_id IS NULL ORDER BY sr.ORDER ASC, sr.parent_id DESC, sr.id DESC, srp.id DESC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "product_code", Value: productCode}, zenlogger.ZenField{Key: "endpoint", Value: endpoint})

	err := util.GetDB().Select(&middlewareResponseValueConfs, sqlStmt, productCode, endpoint)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return middlewareResponseValueConfs, nil
}

func (repo *DefaultServerResponseRepo) FindByEndpointAndFieldAs(endpoint, fieldAs string) (*domain.ServerResponse, error) {

	var data domain.ServerResponse

	sqlStmt := `SELECT t.id, t.endpoint, t.field, t.type, CASE WHEN t.parent_id = 0 THEN NULL ELSE t.parent_id END AS parent_id, t.field_as, t.product_code, t.created_at, t.updated_at, t.deleted_at FROM ( SELECT *, ROW_NUMBER ( ) OVER ( PARTITION BY endpoint, field, type, parent_id, field_as ORDER BY CASE WHEN product_code <> 000000 THEN 0 ELSE 1 END, product_code DESC, id ) AS rn FROM server_response ) t WHERE t.rn = 1 AND t.id <> 0 AND t.endpoint = ? AND t.field_as = ? ORDER BY t.id`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "field_as", Value: fieldAs})

	err := util.GetDB().Get(&data, sqlStmt, endpoint, fieldAs)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return &data, nil
}
