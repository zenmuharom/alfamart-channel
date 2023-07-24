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
	// return true

	hash := sha256.New()

	// find userName
	userNameKey, userNameIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, "/", "userName", req)
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

	// find timeStamp
	timeStampKey, timeStampIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "timeStamp", req)
	if timeStampKey.Parent == "" && timeStampKey.Field == "" {
		valid = false
		err = errors.New("timestamp route not set yet by Finnet.")
		return
	}
	timeStamp := fmt.Sprintf("%v", timeStampIValue)

	sha256timeStampAndPassword := timeStamp + user.Password.String
	hash.Write([]byte(sha256timeStampAndPassword))

	sha256res := fmt.Sprintf("%x", hash.Sum(nil))

	// find channel
	channelKey, channelIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "channel", req)
	if channelKey.Parent == "" && channelKey.Field == "" {
		valid = false
		err = errors.New("channel route not set yet by Finnet.")
		return
	}
	channel := fmt.Sprintf("%v", channelIValue)

	// find terminal
	terminalKey, terminalIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "terminal", req)
	if terminalKey.Parent == "" && terminalKey.Field == "" {
		valid = false
		err = errors.New("terminal route not set yet by Finnet.")
		return
	}
	terminal := fmt.Sprintf("%v", terminalIValue)

	// find transactionType
	transactionTypeKey, transactionTypeIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "transactionType", req)
	if transactionTypeKey.Parent == "" && transactionTypeKey.Field == "" {
		valid = false
		err = errors.New("transactionType route not set yet by Finnet.")
		return
	}
	transactionType := fmt.Sprintf("%v", transactionTypeIValue)

	// find billNumber
	billNumberKey, billNumberIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "billNumber", req)
	if billNumberKey.Parent == "" && billNumberKey.Field == "" {
		valid = false
		err = errors.New("bill number route not set yet by Finnet.")
		return
	}
	billNumber := fmt.Sprintf("%v", billNumberIValue)

	// find traxId
	traxIdKey, traxIdIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "traxId", req)
	if traxIdKey.Parent == "" && traxIdKey.Field == "" {
		valid = false
		err = errors.New("traxId route not set yet by Finnet.")
		return
	}
	traxId := fmt.Sprintf("%v", traxIdIValue)

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

	rawSignature := userName + productCode + channel + terminal + transactionType + billNumber + traxId + sha256res

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
