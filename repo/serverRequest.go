package repo

import (
	"alfamart-channel/domain"
	"alfamart-channel/util"

	"github.com/zenmuharom/zenlogger"
)

type DefaultServerRequestRepo struct {
	logger zenlogger.Zenlogger
}

type ServerRequestRepo interface {
	FindAllByEndpoint(endpoint string) ([]domain.ServerRequest, error)
	FindRequestQueryByEndpoint(endpoint string) ([]domain.ServerRequest, error)
	FindByEndpointAndFieldAs(endpoint, fieldAs string) (*domain.ServerRequest, error)
}

func NewServerRequestRepo(logger zenlogger.Zenlogger) ServerRequestRepo {
	return &DefaultServerRequestRepo{
		logger: logger,
	}
}

func (repo *DefaultServerRequestRepo) FindAllByEndpoint(endpoint string) ([]domain.ServerRequest, error) {
	serverRequestConfs := make([]domain.ServerRequest, 0)

	sqlStmt := `SELECT t.id, t.endpoint, t.field, t.type, t.length_min, t.length_max, t.required, CASE WHEN t.parent_id = 0 THEN NULL ELSE t.parent_id END AS parent_id, t.field_as, t.created_at, t.updated_at, t.deleted_at FROM ( SELECT *, ROW_NUMBER() OVER ( PARTITION BY endpoint, field, type, length, required, parent_id, field_as ORDER BY CASE WHEN productCode <> 000000 THEN 0 ELSE 1 END, productCode DESC, id ) AS rn FROM server_request ) t WHERE t.rn = 1 AND t.id <> 0 AND t.endpoint = ? ORDER BY t.id ASC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "endpoint", Value: endpoint})

	err := util.GetDB().Select(&serverRequestConfs, sqlStmt, endpoint)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return serverRequestConfs, nil
}

func (repo *DefaultServerRequestRepo) FindRequestQueryByEndpoint(endpoint string) ([]domain.ServerRequest, error) {
	serverRequestConfs := make([]domain.ServerRequest, 0)

	sqlStmt := `SELECT t.id, t.endpoint, t.field, t.type, t.length_min, t.length_max, t.required, CASE WHEN t.parent_id = 0 THEN NULL ELSE t.parent_id END AS parent_id, t.field_as, t.created_at, t.updated_at, t.deleted_at FROM ( SELECT sr.*, ROW_NUMBER ( ) OVER ( PARTITION BY sr.endpoint, sr.field, sr.type, sr.length, sr.required, sr.parent_id, sr.field_as ORDER BY CASE WHEN sr.productCode <> 000000 THEN 0 ELSE 1 END, sr.productCode DESC, id ) AS rn FROM server_request sr JOIN server_request srBody ON srBody.id = sr.parent_id AND srBody.field_as = 'requestQuery' ) t WHERE t.rn = 1 AND t.id <> 0 AND t.endpoint = ? ORDER BY t.id ASC`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "endpoint", Value: endpoint})

	err := util.GetDB().Select(&serverRequestConfs, sqlStmt, endpoint)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return serverRequestConfs, nil
}

func (repo *DefaultServerRequestRepo) FindByEndpointAndFieldAs(endpoint, fieldAs string) (*domain.ServerRequest, error) {

	var data domain.ServerRequest

	sqlStmt := `SELECT t.id, t.endpoint, t.field, t.type, t.length_min, t.length_max, t.required, CASE WHEN t.parent_id = 0 THEN NULL ELSE t.parent_id END AS parent_id, t.field_as, t.productCode, t.created_at, t.updated_at, t.deleted_at FROM ( SELECT *, ROW_NUMBER ( ) OVER ( PARTITION BY endpoint, field, type, length, required, parent_id, field_as ORDER BY CASE WHEN productCode <> 000000 THEN 0 ELSE 1 END, productCode DESC, id ) AS rn FROM server_request ) t WHERE t.rn = 1 AND t.id <> 0 AND t.endpoint = ? AND t.field_as = ? ORDER BY t.id`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "endpoint", Value: endpoint}, zenlogger.ZenField{Key: "field_as", Value: fieldAs})

	err := util.GetDB().Get(&data, sqlStmt, endpoint, fieldAs)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	return &data, nil
}
