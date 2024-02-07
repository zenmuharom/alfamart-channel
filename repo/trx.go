package repo

import (
	"alfamart-channel/domain"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zenmuharom/zenlogger"
)

type DefaultTrxRepo struct {
	db     *sqlx.DB
	logger zenlogger.Zenlogger
}

type TrxRepo interface {
	FindAll() ([]domain.Trx, error)
	FindByTargetNumber(targetNumber string) (trx domain.Trx, errRes error)
	FindByInquiryPayment(targetNumber string) (trx domain.Trx, errRes error)
	FindByCommit(targetNumber string) (trx domain.Trx, errRes error)
	Insert(data domain.Trx) (inserted *domain.Trx, errRes error)
	Upsert(data domain.Trx) (updated *domain.Trx, errRes error)
}

func NewTrxRepo(logger zenlogger.Zenlogger, db *sqlx.DB) TrxRepo {
	return &DefaultTrxRepo{
		logger: logger,
		db:     db,
	}
}

func (repo *DefaultTrxRepo) FindAll() ([]domain.Trx, error) {
	var merchants []domain.Trx

	sqlStmt := `SELECT id, amendment_date, source_merchant, target_product, pid, amount, status, rc, rc_desc, created_at, updated_at, deleted_at FROM trx WHERE deleted_at IS NULL`
	repo.logger.Debug("FindAll", zenlogger.ZenField{Key: "sqlStmt", Value: sqlStmt})

	err := repo.db.Select(&merchants, sqlStmt)
	if err != nil {

		repo.logger.Error("FindAll", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call err := repo.db.Select(&merchants, sqlStmt)"})
		return merchants, err
	}

	return merchants, nil
}

func (repo *DefaultTrxRepo) FindByTargetNumber(targetNumber string) (trx domain.Trx, errRes error) {

	sqlStmt := `SELECT pid, amendment_date, source_code, source_merchant, target_product, target_number, bit61, amount, status, rc, rc_desc FROM trx WHERE deleted_at IS NULL AND source_code = 'INQUIRY' AND status = 'approve' AND created_at >= NOW() - INTERVAL 3 MINUTE AND target_number = ? LIMIT 1`
	repo.logger.Debug("FindByTargetNumber", zenlogger.ZenField{Key: "sqlStmt", Value: sqlStmt}, zenlogger.ZenField{Key: "targetNumber", Value: targetNumber})

	err := repo.db.Get(&trx, sqlStmt, targetNumber)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			repo.logger.Error("FindAll", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call err := repo.db.Select(&trx, sqlStmt)"})
		}
		return trx, err
	}

	return trx, nil
}

func (repo *DefaultTrxRepo) FindByInquiryPayment(targetNumber string) (trx domain.Trx, errRes error) {

	sqlStmt := `SELECT pay.pid, pay.amendment_date, pay.source_code, pay.source_merchant, pay.target_product, pay.target_number, pay.bit61, pay.amount, pay.status, pay.rc, pay.rc_desc FROM trx inq JOIN trx pay ON inq.status = 'approve' AND pay.status = 'approve' AND inq.target_number = pay.target_number AND pay.source_code = 'PAYMENT' AND inq.status = pay.status AND inq.deleted_at IS NULL AND pay.deleted_at IS NULL AND inq.created_at >= NOW( ) - INTERVAL 3 MINUTE AND pay.created_at >= NOW( ) - INTERVAL 3 MINUTE WHERE inq.source_code = 'INQUIRY' AND inq.target_number = ? LIMIT 1`
	repo.logger.Debug("FindByInquiryPayment", zenlogger.ZenField{Key: "sqlStmt", Value: sqlStmt}, zenlogger.ZenField{Key: "targetNumber", Value: targetNumber})

	err := repo.db.Get(&trx, sqlStmt, targetNumber)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			repo.logger.Error("FindAll", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call err := repo.db.Select(&trx, sqlStmt)"})
		}
		return trx, err
	}

	return trx, nil
}

func (repo *DefaultTrxRepo) FindByCommit(targetNumber string) (trx domain.Trx, errRes error) {

	sqlStmt := `SELECT pay.pid, pay.amendment_date, pay.source_code, pay.source_merchant, pay.target_product, pay.target_number, pay.bit61, pay.amount, pay.status, pay.rc, pay.rc_desc FROM trx pay WHERE deleted_at IS NULL AND source_code = 'PAYMENT' AND status = 'approve' AND created_at >= NOW() - INTERVAL 3 MINUTE AND target_number = ? AND NOT EXISTS ( SELECT 1 FROM trx com WHERE com.target_number = pay.target_number AND com.source_code = 'COMMIT' AND com.status = 'approve' ) LIMIT 1`
	repo.logger.Debug("FindByCommit", zenlogger.ZenField{Key: "sqlStmt", Value: sqlStmt}, zenlogger.ZenField{Key: "targetNumber", Value: targetNumber})

	err := repo.db.Get(&trx, sqlStmt, targetNumber)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			repo.logger.Error("FindAll", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call err := repo.db.Select(&trx, sqlStmt)"})
		}
		return trx, err
	}

	return trx, nil
}

func (repo *DefaultTrxRepo) Insert(data domain.Trx) (inserted *domain.Trx, errRes error) {
	var err error
	sqlStmt := `INSERT INTO trx (amendment_date, source_merchant, target_product, pid, amount, status, rc, rc_desc) VALUES (:amendment_date, :source_merchant, :target_product, :pid, :amount, :status, :rc, :rc_desc)`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "data", Value: data})

	_, err = repo.db.NamedExec(sqlStmt, data)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	inserted = &data

	return
}

func (repo *DefaultTrxRepo) Upsert(data domain.Trx) (updated *domain.Trx, errRes error) {
	var err error
	sqlStmt := `
	INSERT INTO trx (pid, amendment_date, source_code, source_merchant, target_product, target_number, bit61, amount, rc, rc_desc, status, elapsed_time)
	VALUES (:pid, :amendment_date, :source_code, :source_merchant, :target_product, :target_number, :bit61, :amount, :rc, :rc_desc, :status, :elapsed_time)
	ON DUPLICATE KEY UPDATE
		source_code = :source_code,
		source_merchant = :source_merchant, 
		target_product = :target_product,
		target_number = :target_number,
		bit61 = :bit61,
		amount = :amount,
		rc = :rc, 
		rc_desc = :rc_desc,
		status = :status,
		elapsed_time = :elapsed_time
	`
	repo.logger.Debug(sqlStmt, zenlogger.ZenField{Key: "data", Value: data})

	_, err = repo.db.NamedExec(sqlStmt, data)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}

	updated = &data

	return
}
