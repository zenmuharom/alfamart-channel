package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"crypto/sha1"
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

	// hash := sha256.New()
	hash := sha1.New()

	// find userName
	agentIDKey, agentIDIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "userName", req)
	if agentIDKey.Parent == "" && agentIDKey.Field == "" {
		valid = false
		err = errors.New("AgentID route not set yet by Finnet.")
		return
	}
	agentID := fmt.Sprintf("%v", agentIDIValue)

	// find agentPIN
	agentPINKey, agentPINIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "password", req)
	if agentPINKey.Parent == "" && agentPINKey.Field == "" {
		valid = false
		err = errors.New("AgentPIN route not set yet by Finnet.")
		return
	}
	agentPIN := fmt.Sprintf("%v", agentPINIValue)

	// find transactionId
	transactionIdKey, transactionIdIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "transactionId", req)
	if transactionIdKey.Parent == "" && transactionIdKey.Field == "" {
		valid = false
		err = errors.New("AgentTrxID route not set yet by Finnet.")
		return
	}
	transactionId := fmt.Sprintf("%v", transactionIdIValue)

	// find storeId
	storeIdKey, storeIdIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "storeId", req)
	if storeIdKey.Parent == "" && storeIdKey.Field == "" {
		valid = false
		err = errors.New("StoreID route not set yet by Finnet.")
		return
	}
	storeId := fmt.Sprintf("%v", storeIdIValue)

	// find productId
	productIdKey, productIdIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "productCode", req)
	if productIdKey.Parent == "" && productIdKey.Field == "" {
		valid = false
		err = errors.New("ProductID route not set yet by Finnet.")
		return
	}
	productId := fmt.Sprintf("%v", productIdIValue)

	// find customerId
	customerIdKey, customerIdIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "billNumber", req)
	if customerIdKey.Parent == "" && customerIdKey.Field == "" {
		valid = false
		err = errors.New("CustomerID route not set yet by Finnet.")
		return
	}
	customerId := fmt.Sprintf("%v", customerIdIValue)

	// find timeStamp
	timeStampKey, timeStampIValue := tool.FindFieldAs(service.logger, domain.SERVER_REQUEST, endpoint, "timeStamp", req)
	if timeStampKey.Parent == "" && timeStampKey.Field == "" {
		valid = false
		err = errors.New("CustomerID route not set yet by Finnet.")
		return
	}
	timeStamp := fmt.Sprintf("%v", timeStampIValue)

	// find username on DB config
	userRepo := repo.NewUserRepo(service.logger)
	user, err := userRepo.Find(agentID)
	if err != nil {
		service.logger.Info(err.Error())
		valid = false
		return
	}

	service.logger.Debug("Signature Raw", zenlogger.ZenField{Key: "request", Value: req}, zenlogger.ZenField{Key: "password", Value: user.Password.String}, zenlogger.ZenField{Key: "secret", Value: user.Secret})

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

	rawSignature := agentID + agentPIN + transactionId + storeId + productId + customerId + timeStamp + user.Secret.String

	hash.Write([]byte(rawSignature))
	generatedSignature := fmt.Sprintf("%x", hash.Sum(nil))

	service.logger.Debug("Signature Processed", zenlogger.ZenField{Key: "rawSignature", Value: rawSignature}, zenlogger.ZenField{Key: "reqSignature", Value: reqSignature})

	service.logger.Debug("Signature Validating", zenlogger.ZenField{Key: "reqSignature", Value: reqSignature}, zenlogger.ZenField{Key: "result", Value: generatedSignature}, zenlogger.ZenField{Key: "hashAlgorithm", Value: "sha1"})

	valid = reqSignature == generatedSignature
	return
}
