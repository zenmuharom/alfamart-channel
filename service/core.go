package service

import (
	"alfamart-channel/domain"
	"alfamart-channel/function"
	"alfamart-channel/helper"
	"alfamart-channel/repo"
	"alfamart-channel/tool"
	"alfamart-channel/util"
	"errors"
	"fmt"
	"reflect"

	"github.com/zenmuharom/zenlogger"
)

type DefaultCoreService struct {
	logger zenlogger.Zenlogger
}

type CoreService interface {
	ConstructBit61(productCode, billNumber string) (result string)
	// Request(process string, req dto.Trx) (response map[string]interface{}, err error)
	Request2(process string, config *domain.UserProduct, req map[string]interface{}) (middlewResponseIDs []int, response map[string]interface{}, err error)
	Parse(process string, tsResponse map[string]interface{}) (middlewResponseIDs []int, response map[string]interface{}, err error)
}

func NewCoreService(logger zenlogger.Zenlogger) CoreService {
	return &DefaultCoreService{logger: logger}
}

func (service *DefaultCoreService) ConstructBit61(productCode, billNumber string) (result string) {
	service.logger.Debug("ConstructBit61", zenlogger.ZenField{Key: "productCode", Value: productCode}, zenlogger.ZenField{Key: "billNumber", Value: billNumber})
	prefixSuffixRepo := repo.NewPrefixSuffixRepo(service.logger)
	prefixSuffixConfigs, err := prefixSuffixRepo.FindAllByProductCode(productCode)
	if err != nil {
		service.logger.Error(err.Error())
	}

	// if there is no prefix suffix config then skip it
	if len(prefixSuffixConfigs) == 0 {
		result = billNumber
		return
	}

	prefix := ""

	var defaultConfig domain.PrefixSuffixConfig
	foundSpecificConfig := false

	for _, preSufConfig := range prefixSuffixConfigs {

		if preSufConfig.Prefix.String == "" {
			defaultConfig = preSufConfig
			continue
		}

		prefix = billNumber[int(preSufConfig.FixNumberStart.Int64):int(preSufConfig.FixNumberEnd.Int64)]

		service.logger.Debug("ConstructBit61", zenlogger.ZenField{Key: "process", Value: "matching rule"}, zenlogger.ZenField{Key: "billNumber_prefix", Value: prefix}, zenlogger.ZenField{Key: "list_prefix", Value: preSufConfig.Prefix.String})

		service.logger.Debug("ConstructBit61", zenlogger.ZenField{Key: "logic", Value: fmt.Sprintf("ConstructBit61 %v == %v", prefix, preSufConfig.Prefix.String)})
		if prefix == preSufConfig.Prefix.String {

			// config suffix if suffix length not null
			if preSufConfig.SuffixLength.Int64 != 0 {
				suffix := billNumber[len(prefix):]
				suffix_padded := helper.PadLeft(suffix, int(preSufConfig.SuffixLength.Int64), preSufConfig.FixCharacter.String)

				if preSufConfig.PresufType.String == "0" {
					billNumber = helper.PadLeft(fmt.Sprintf("%v%v", prefix, suffix_padded), int(preSufConfig.Length.Int64), preSufConfig.FixCharacter.String)
				} else {
					billNumber = helper.PadRight(fmt.Sprintf("%v%v", prefix, suffix_padded), int(preSufConfig.Length.Int64), preSufConfig.FixCharacter.String)
				}

			} else {
				if preSufConfig.PresufType.String == "0" {
					billNumber = helper.PadLeft(billNumber, int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
				} else {
					billNumber = helper.PadRight(billNumber, int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
				}
			}

			foundSpecificConfig = true
			break
		}
	}

	if !foundSpecificConfig {
		service.logger.Debug("ConstructBit61", zenlogger.ZenField{Key: "foundSpecificConfig", Value: foundSpecificConfig})
		prefix = billNumber[int(defaultConfig.FixNumberStart.Int64):int(defaultConfig.FixNumberEnd.Int64)]
		service.logger.Debug("ConstructBit61", zenlogger.ZenField{Key: "process", Value: "do based on default config"}, zenlogger.ZenField{Key: "billNumber_prefix", Value: prefix}, zenlogger.ZenField{Key: "list_prefix", Value: defaultConfig.Prefix.String})

		// config suffix if suffix length not null
		if defaultConfig.SuffixLength.Int64 != 0 {
			suffix := billNumber[len(prefix):]
			suffix_padded := helper.PadLeft(suffix, int(defaultConfig.SuffixLength.Int64), defaultConfig.FixCharacter.String)
			if defaultConfig.PresufType.String == "0" {
				billNumber = helper.PadLeft(fmt.Sprintf("%v%v", prefix, suffix_padded), int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
			} else {
				billNumber = helper.PadRight(fmt.Sprintf("%v%v", prefix, suffix_padded), int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
			}
		} else {
			if defaultConfig.PresufType.String == "0" {
				billNumber = helper.PadLeft(billNumber, int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
			} else {
				billNumber = helper.PadRight(billNumber, int(defaultConfig.Length.Int64), defaultConfig.FixCharacter.String)
			}
		}

	}

	// billNumber = helper.PadLeft(prefix+suffix, 13, "0")

	return billNumber
}

func (service *DefaultCoreService) Request2(process string, config *domain.UserProduct, req map[string]interface{}) (middlewResponseIDs []int, response map[string]interface{}, err error) {
	service.logger.Debug("Request2", zenlogger.ZenField{Key: "req", Value: req})

	// send request to TS Adapter
	tsRepo := repo.NewTSRepo(service.logger)

	// find bit32
	bit32Key, bit32IValue := tool.FindFieldAs(service.logger, domain.MIDDLEWARE_REQUEST, fmt.Sprintf("ts_adapter|%v", process), "bit32", req)
	service.logger.Debug("bit32", zenlogger.ZenField{Key: "bit33", Value: bit32Key}, zenlogger.ZenField{Key: "value", Value: bit32IValue})
	if bit32IValue == nil || bit32IValue == "" {
		service.logger.Info("Request2", zenlogger.ZenField{Key: "process", Value: "Override bit32"}, zenlogger.ZenField{Key: "key", Value: bit32Key.Field}, zenlogger.ZenField{Key: "value", Value: config.Bit32.String})
		if bit32Key.Parent == "" {
			req[bit32Key.Field] = config.Bit32.String
		} else {
			// if it has parent key
			parentRCObject := req[bit32Key.Parent]
			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
			newParentRCObject := make(map[string]interface{})
			switch valueOfVariableParentRC.Kind() {
			case reflect.Map:
				iter := valueOfVariableParentRC.MapRange()
				for iter.Next() {
					if iter.Key().String() == bit32Key.Field {
						newParentRCObject[iter.Key().String()] = config.Bit32.String
					} else {
						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
					}
				}
			}
		}
	}

	// find bit33
	bit33Key, bit33IValue := tool.FindFieldAs(service.logger, domain.MIDDLEWARE_REQUEST, fmt.Sprintf("ts_adapter|%v", process), "bit33", req)
	service.logger.Debug("bit33", zenlogger.ZenField{Key: "bit33", Value: bit33Key}, zenlogger.ZenField{Key: "value", Value: bit33IValue})
	if bit33IValue == nil || bit33IValue == "" {
		service.logger.Info("Request2", zenlogger.ZenField{Key: "process", Value: "Override bit33"}, zenlogger.ZenField{Key: "key", Value: bit32Key.Field}, zenlogger.ZenField{Key: "value", Value: config.Bit33.String})
		if bit33Key.Parent == "" {
			req[bit33Key.Field] = config.Bit33.String
		} else {
			// if it has parent key
			parentRCObject := req[bit33Key.Parent]
			valueOfVariableParentRC := reflect.ValueOf(parentRCObject)
			newParentRCObject := make(map[string]interface{})
			switch valueOfVariableParentRC.Kind() {
			case reflect.Map:
				iter := valueOfVariableParentRC.MapRange()
				for iter.Next() {
					if iter.Key().String() == bit33Key.Field {
						newParentRCObject[iter.Key().String()] = config.Bit33.String
					} else {
						newParentRCObject[iter.Key().String()] = iter.Value().Interface()
					}
				}
			}
		}
	}

	// FOR DEBUG IN LOCAL
	// fmt.Println(service.logger.Debug("req", zenlogger.ZenField{Key: "req", Value: req}, zenlogger.ZenField{Key: "bit32", Value: bit32Key}, zenlogger.ZenField{Key: "bit33", Value: bit33Key}))

	tsRawRes, err := tsRepo.Request2(config.SwitchingUrl.String, int(config.TimeoutBiller.Int32), req)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	// parse response value
	middlewResponseIDs, response, err = service.Parse(process, tsRawRes)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	return
}

func (service *DefaultCoreService) Parse(process string, tsResponse map[string]interface{}) (middlewResponseIDs []int, response map[string]interface{}, err error) {

	middlewResponseIDs = make([]int, 0)

	service.logger.Debug("parsing", zenlogger.ZenField{Key: "middleware", Value: util.MIDDLEWARE_TS_ADAPTER}, zenlogger.ZenField{Key: "process", Value: process}, zenlogger.ZenField{Key: "tsResponse", Value: tsResponse})

	responseRepo := repo.NewMiddlewareResponseRepo(service.logger)
	config, err := responseRepo.FindAllByMiddlewareAndProcess(util.MIDDLEWARE_TS_ADAPTER, process)
	if err != nil {
		service.logger.Error(err.Error())
		return
	}

	if len(config) == 0 {
		err = errors.New("middleware response route not set yet by Finnet")
		service.logger.Error(err.Error(), zenlogger.ZenField{Key: "process", Value: process})
		return
	}

	for _, conf := range config {
		middlewResponseIDs = append(middlewResponseIDs, int(conf.Id))
	}

	response, err = function.MiddlewareResponseParse(config, tsResponse)
	if err != nil {
		err = errors.New("middleware response route not set yet by Finnet")
		service.logger.Error(err.Error(), zenlogger.ZenField{Key: "process", Value: process})
		return
	}

	return
}
