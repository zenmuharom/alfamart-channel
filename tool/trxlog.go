package tool

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"alfamart-channel/util"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/zenmuharom/zenlogger"
)

type DefaultTrxLog struct {
	logger zenlogger.Zenlogger
	db     *sqlx.DB
}

type TrxLog interface {
	Write(trx *domain.Trx) (err error)
}

func NewTrxLog(logger zenlogger.Zenlogger) TrxLog {
	dbMysql := util.GetDB()

	return &DefaultTrxLog{
		logger: logger,
		db:     dbMysql,
	}
}

func (this *DefaultTrxLog) Write(trx *domain.Trx) (err error) {

	dataToInsert := domain.Trx{
		Pid:            this.logger.GetPid(),
		AmendmentDate:  sql.NullTime{Time: trx.AmendmentDate.Time, Valid: true},
		SourceCode:     sql.NullString{String: trx.SourceCode.String, Valid: true},
		SourceMerchant: sql.NullString{String: trx.SourceMerchant.String, Valid: true},
		TargetProduct:  sql.NullString{String: trx.TargetProduct.String, Valid: true},
		TargetNumber:   sql.NullString{String: trx.TargetNumber.String, Valid: true},
		Bit61:          sql.NullString{String: trx.Bit61.String, Valid: true},
		Amount:         sql.NullFloat64{Float64: trx.Amount.Float64, Valid: true},
		Rc:             sql.NullString{String: trx.Rc.String, Valid: true},
		RcDesc:         sql.NullString{String: trx.RcDesc.String, Valid: true},
		Status:         sql.NullString{String: trx.Status.String, Valid: true},
		ElapsedTime:    sql.NullInt64{Int64: trx.ElapsedTime.Int64, Valid: true},
	}

	trxRepo := repo.NewTrxRepo(this.logger, this.db)
	_, err = trxRepo.Upsert(dataToInsert)
	if err != nil {
		this.logger.Error("Write", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call _, err := trxRepo.insertTrx(dataToInsert)"})
	}

	return
}
