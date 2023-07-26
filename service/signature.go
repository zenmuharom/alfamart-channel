package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/zenmuharom/zenlogger"
)

type DefaultSignatureService struct {
	logger zenlogger.Zenlogger
}

type SignatureService interface {
	Check(endpoint string, req map[string]interface{}) (valid bool, err error)
}

func NewSignatureService(logger zenlogger.Zenlogger) SignatureService {
	return &DefaultSignatureService{
		logger: logger,
	}
}

func (service *DefaultSignatureService) Check(endpoint string, req map[string]interface{}) (valid bool, err error) {
	// BYPASS FOR DEBUG
	// valid = true
	// err = nil
	// return

	hash := sha256.New()

	// find userName
	userNameKey, userNameIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "userName", req)
	if userNameKey.Parent == "" && userNameKey.Field == "" {
		valid = false
		err = errors.New("username route not set yet by Finnet.")
		return
	}
	userName := fmt.Sprintf("%v", userNameIValue)

	// find productCode
	productCodeKey, productCodeIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "productCode", req)
	if productCodeKey.Parent == "" && productCodeKey.Field == "" {
		valid = false
		err = errors.New("product code route not set yet by Finnet.")
		return
	}
	productCode := fmt.Sprintf("%v", productCodeIValue)

	// find username on DB config
	userRepo := repo.NewUserRepo(service.logger)
	user, err := userRepo.Find(userName)
	if err != nil {
		service.logger.Info(err.Error())
		valid = false
		return
	}

	service.logger.Debug("Signature Raw", zenlogger.ZenField{Key: "request", Value: req}, zenlogger.ZenField{Key: "password", Value: user.Password.String})

	sha256timeStampAndPassword := user.Password.String
	hash.Write([]byte(sha256timeStampAndPassword))

	sha256res := fmt.Sprintf("%x", hash.Sum(nil))

	// find signature
	signatureKey, signatureIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "signature", req)
	if signatureKey.Parent == "" && signatureKey.Field == "" {
		valid = false
		err = errors.New("signature route not set yet by Finnet.")
		return
	}
	reqSignature := fmt.Sprintf("%v", signatureIValue)

	// do not remove this block code
	if reqSignature == "" {
		return
	}

	rawSignature := userName + productCode + sha256res

	generatedSignature := ""

	service.logger.Debug("Signature Processed", zenlogger.ZenField{Key: "sha256timeStampAndPassword", Value: sha256timeStampAndPassword}, zenlogger.ZenField{Key: "sha256res", Value: sha256res}, zenlogger.ZenField{Key: "rawSignature", Value: rawSignature}, zenlogger.ZenField{Key: "reqSignature", Value: reqSignature})
	hashAlogirthm := ""

	if tool.IsMD5Hash(reqSignature) {
		hashAlogirthm = "MD5"
		generatedSignature = fmt.Sprintf("%x", md5.Sum([]byte(rawSignature)))

	} else if tool.IsSHA256Hash(reqSignature) {
		hashAlogirthm = "SHA256"
		hash := sha256.New()
		hash.Write([]byte(rawSignature))
		generatedSignature = fmt.Sprintf("%x", hash.Sum(nil))
	} else {
		service.logger.Debug("Signature Validating", zenlogger.ZenField{Key: "reqSignature", Value: reqSignature}, zenlogger.ZenField{Key: "error", Value: "signature not matches any Hash Algorithm"})
	}

	service.logger.Debug("Signature Validating", zenlogger.ZenField{Key: "reqSignature", Value: reqSignature}, zenlogger.ZenField{Key: "result", Value: generatedSignature}, zenlogger.ZenField{Key: "hashAlgorithm", Value: hashAlogirthm})

	valid = reqSignature == generatedSignature
	return
}
