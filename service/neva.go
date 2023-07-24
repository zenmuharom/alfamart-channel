package service

import (
	"alfamart-channel/dto"
	"alfamart-channel/repo"
	"fmt"

	"github.com/zenmuharom/zenlogger"
)

type DefaultNevaService struct {
	sessionId string
	trx       dto.NevaDebet
	logger    zenlogger.Zenlogger
}

type NevaService interface {
	Debet(debet dto.NevaDebet) (status bool, err error)
	Reversal() (status bool, err error)
}

func NewNevaService(logger zenlogger.Zenlogger) NevaService {
	return &DefaultNevaService{
		logger: logger,
	}
}

func (service *DefaultNevaService) Debet(debet dto.NevaDebet) (status bool, err error) {

	nevaRepo := repo.NewNevaRepo(service.logger)

	service.trx = debet

	if service.sessionId == "" {

		// Do request token
		reqToken := dto.NevaInputToken{
			Type: "urn:getToken",
			Username: dto.NevaField{
				Type:  "xsd:string",
				Value: debet.Username,
			},
			Password: dto.NevaField{
				Type:  "xsd:string",
				Value: debet.Password,
			},
			MitraCo: dto.NevaField{
				Type:  "xsd:string",
				Value: debet.MitraCo,
			},
		}
		tokenRes, errToken := nevaRepo.GetToken(reqToken)
		if errToken != nil {
			err = errToken
			service.logger.Error(errToken.Error())
			return
		}

		service.sessionId = tokenRes.Body.GetTokenResponse.OutputToken.SessionId

	}

	// Do debet
	reqDebet := dto.NevaInputTransaction{
		Type: "urn:inputTransaction",
		Description: dto.NevaField{
			Type:  "xsd:string",
			Value: fmt.Sprintf("PAYMENT %s - %s - %s", debet.ProdCode, debet.BillingNo, debet.TraxId),
		},
		Dest1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.DestAcc,
		},
		Dest1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.Amount,
		},
		Phoneno: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.PhoneNo,
		},
		Sessionid: dto.NevaField{
			Type:  "xsd:string",
			Value: service.sessionId,
		},
		Source1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.SourceAcc,
		},
		Source1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.Amount,
		},
		TransactionType: dto.NevaField{
			Type:  "xsd:string",
			Value: "payment",
		},
		Mitraco: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.MitraCo,
		},
		Prodcode: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.ProdCode,
		},
		Billingno: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.BillingNo,
		},
		TraxId: dto.NevaField{
			Type:  "xsd:string",
			Value: debet.TraxId,
		},
	}
	debetRes, err := nevaRepo.DebetAccountV2(reqDebet)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	status = (debetRes.Body.DebetAccountV2Response.OutputTransaction.Status.ResultCode == "0")

	return
}

func (service *DefaultNevaService) Reversal() (status bool, err error) {

	nevaRepo := repo.NewNevaRepo(service.logger)

	if service.sessionId == "" {

		// Do request token
		reqToken := dto.NevaInputToken{
			Type: "urn:getToken",
			Username: dto.NevaField{
				Type:  "xsd:string",
				Value: service.trx.Username,
			},
			Password: dto.NevaField{
				Type:  "xsd:string",
				Value: service.trx.Password,
			},
			MitraCo: dto.NevaField{
				Type:  "xsd:string",
				Value: service.trx.MitraCo,
			},
		}
		tokenRes, err2 := nevaRepo.GetToken(reqToken)
		if err2 != nil {
			service.logger.Error(err.Error())
			return
		}

		service.sessionId = tokenRes.Body.GetTokenResponse.OutputToken.SessionId
	}

	// Do reversal
	reqDebet := dto.NevaInputTransaction{
		Type: "urn:inputTransaction",
		Description: dto.NevaField{
			Type:  "xsd:string",
			Value: fmt.Sprintf("REVERSAL %s - %s - %s", service.trx.ProdCode, service.trx.BillingNo, service.trx.TraxId),
		},
		Dest1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.SourceAcc,
		},
		Dest1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.Amount,
		},
		Phoneno: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.DestAcc,
		},
		Sessionid: dto.NevaField{
			Type:  "xsd:string",
			Value: service.sessionId,
		},
		Source1Acc: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.DestAcc,
		},
		Source1Amount: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.Amount,
		},
		TransactionType: dto.NevaField{
			Type:  "xsd:string",
			Value: "reversal",
		},
		Mitraco: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.MitraCo,
		},
		Prodcode: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.ProdCode,
		},
		Billingno: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.BillingNo,
		},
		TraxId: dto.NevaField{
			Type:  "xsd:string",
			Value: service.trx.TraxId,
		},
	}
	debetRes, err := nevaRepo.DebetAccountV2(reqDebet)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	status = (debetRes.Body.DebetAccountV2Response.OutputTransaction.Status.ResultCode == "0")

	return
}
