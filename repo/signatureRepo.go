package repo

import (
	"alfamart-channel/dto"
	"crypto/md5"
	"crypto/sha256"
	"fmt"

	"github.com/zenmuharom/zenlogger"
)

type DefaultSignatureRepo struct {
	logger zenlogger.Zenlogger
}

type SignatureRepo interface {
	Check(request dto.Trx, password string) (valid bool)
}

func NewSignatureRepo(logger zenlogger.Zenlogger) SignatureRepo {
	return &DefaultSignatureRepo{
		logger: logger,
	}
}

func (repo *DefaultSignatureRepo) Check(request dto.Trx, password string) (valid bool) {
	repo.logger.Debug("Signature Raw", zenlogger.ZenField{Key: "request", Value: request}, zenlogger.ZenField{Key: "password", Value: password})

	hash := sha256.New()

	sha256 := request.TimeStamp + password
	hash.Write([]byte(sha256))

	sha256res := fmt.Sprintf("%x", hash.Sum(nil))

	signature := request.UserName + request.ProductCode + request.Terminal + request.TransactionType + request.BillNumber + request.TraxId + sha256res

	signInq := fmt.Sprintf("%x", md5.Sum([]byte(signature)))

	repo.logger.Debug("Signature Processed", zenlogger.ZenField{Key: "sha256", Value: sha256}, zenlogger.ZenField{Key: "sha256res", Value: sha256res}, zenlogger.ZenField{Key: "signature", Value: signature}, zenlogger.ZenField{Key: "signInqRes", Value: signInq})

	if string(signInq[:]) == request.Signature {
		valid = true
	}

	return
}
