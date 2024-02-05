package function

import (
	"alfamart-channel/domain"
	"alfamart-channel/helper"
	"alfamart-channel/models"
	"alfamart-channel/variable"
	"encoding/json"
	"fmt"
	"strings"

	"reflect"
	"strconv"

	"github.com/zenmuharom/zenfunction"
	"github.com/zenmuharom/zenlogger"
)

type DefaultAssigner struct {
	logger   zenlogger.Zenlogger
	function zenfunction.Zenfunction
	hasher   helper.Hasher
}

type Assigner interface {
	AssignMiddlewareRequestValue(field domain.FieldValue, fields map[string]models.Variable) (assigned interface{})
	AssignServerResponseValue(field domain.FieldValue, fields map[string]models.Variable) (assigned interface{})
	ReadCommand(dtype, arg string) (result any, err error)
}

func NewAssigner(logger zenlogger.Zenlogger) Assigner {
	return &DefaultAssigner{
		logger:   logger,
		function: zenfunction.New(logger),
	}
}

func (assigner *DefaultAssigner) ReadCommand(dtype, arg string) (result any, err error) {
	res, err := assigner.function.ReadCommand(arg)
	if err != nil {
		assigner.logger.Error("ReadCommand", zenlogger.ZenField{Key: "error", Value: err.Error()})
	}
	result = assigner.ConvertValue(dtype, res)
	return
}

func (assigner *DefaultAssigner) ConvertValue(dtype, originValue any) (value any) {
	switch dtype {
	case "string":
		value = fmt.Sprintf("%v", originValue)
	case "integer":
		iVal, err := convertToInt64(originValue)
		if err != nil {
			assigner.logger.Error("ConvertValue", zenlogger.ZenField{Key: "error", Value: err.Error()})
			value = int64(iVal)
		} else {
			value = iVal
		}
	case "float":
		iVal, err := convertToFloat64(originValue)
		if err != nil {
			assigner.logger.Error("ConvertValue", zenlogger.ZenField{Key: "error", Value: err.Error()})
			value = int64(iVal)
		} else {
			value = iVal
		}
	case "boolean":
		// Convert the interface{} to an int using a type assertion
		bVal, err := strconv.ParseBool(fmt.Sprintf("%v", originValue))
		if err != nil {
			value = false
		} else {
			value = bVal
		}
	case "object":
		value = ""
	}

	return
}

func convertToInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

func convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

func (assigner *DefaultAssigner) ValidityConditionField(config domain.ValueConfig, valueToCompare interface{}) (valid bool) {

	assigner.logger.Debug("ValidityConditionField", zenlogger.ZenField{Key: "config", Value: config}, zenlogger.ZenField{Key: "valueToCompare", Value: valueToCompare})
	functionGenerated := ""

	if config.ConditionFieldId.Int64 == 0 {
		functionGenerated = "no condition"
		valid = true
	} else {
		functionGenerated, valid = assigner.validityConditionCore(config.ConditionOperator.String, valueToCompare, config.ConditionValue.String)
	}

	assigner.logger.Debug("ValidityConditionField", zenlogger.ZenField{Key: "fieldName", Value: config.FieldName.String}, zenlogger.ZenField{Key: "functionGenerated", Value: functionGenerated}, zenlogger.ZenField{Key: "valid", Value: valid})

	return
}

func (assigner *DefaultAssigner) ValidityConditionValue(config domain.FieldValue, valueToCompare interface{}) (valid bool) {
	assigner.logger.Debug("ValidityCondition", zenlogger.ZenField{Key: "config", Value: config}, zenlogger.ZenField{Key: "valueToCompare", Value: valueToCompare})
	functionGenerated := ""

	if config.ConditionFieldId.Int64 == 0 {
		functionGenerated = "no condition"
		valid = true
	} else {
		functionGenerated, valid = assigner.validityConditionCore(config.ConditionOperator.String, valueToCompare, config.ConditionValue.String)
	}

	assigner.logger.Debug("ValidityCondition", zenlogger.ZenField{Key: "field", Value: config.FieldName.String}, zenlogger.ZenField{Key: "functionGenerated", Value: functionGenerated}, zenlogger.ZenField{Key: "valueToCompare", Value: valueToCompare}, zenlogger.ZenField{Key: "valid", Value: valid})

	return
}

func (assigner *DefaultAssigner) validityConditionCore(operator string, valueToCompare interface{}, value interface{}) (functionGenerated string, valid bool) {
	switch operator {
	case "=":
		valueToCompareParsed := ""
		if valueToCompare != nil {
			valueToCompareParsed = valueToCompare.(string)
		}
		functionGenerated = fmt.Sprintf("%v = %v", valueToCompareParsed, value)
		valid = fmt.Sprintf("%v", valueToCompareParsed) == fmt.Sprintf("%v", value)
	case "<":
		typeValueToCompare := reflect.ValueOf(valueToCompare)
		valueToCompareInt := 0
		typeValue := reflect.ValueOf(value)
		valueInt := 0
		var e error

		// set to 0 if valueToCompare not setted
		if valueToCompare == nil {
			valueToCompare = 0
		}

		// convert valueToCompare to int
		if typeValueToCompare.Kind() != reflect.Int {
			valueToCompareInt, e = strconv.Atoi(fmt.Sprintf("%v", valueToCompare))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueToCompareInt = len(fmt.Sprintf("%v", valueToCompare))
			}
		}

		// set to 0 if value not setted
		if value == nil {
			valueInt = 0
		}

		// convert value to int
		if typeValue.Kind() != reflect.Int {
			valueInt, e = strconv.Atoi(fmt.Sprintf("%v", value))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueInt = len(fmt.Sprintf("%v", value))
			}
		}

		functionGenerated = fmt.Sprintf("%v < %v", valueToCompareInt, value)
		valid = valueToCompareInt < valueInt
	case "<=":
		typeValueToCompare := reflect.ValueOf(valueToCompare)
		valueToCompareInt := 0
		typeValue := reflect.ValueOf(value)
		valueInt := 0
		var e error

		// set to 0 if valueToCompare not setted
		if valueToCompare == nil {
			valueToCompare = 0
		}

		// convert valueToCompare to int
		if typeValueToCompare.Kind() != reflect.Int {
			valueToCompareInt, e = strconv.Atoi(fmt.Sprintf("%v", valueToCompare))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueToCompareInt = len(fmt.Sprintf("%v", valueToCompare))
			}
		}

		// set to 0 if value not setted
		if value == nil {
			valueInt = 0
		}

		// convert value to int
		if typeValue.Kind() != reflect.Int {
			valueInt, e = strconv.Atoi(fmt.Sprintf("%v", value))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueInt = len(fmt.Sprintf("%v", value))
			}
		}

		functionGenerated = fmt.Sprintf("%v <= %v", valueToCompareInt, value)
		valid = valueToCompareInt <= valueInt
	case ">":
		typeValueToCompare := reflect.ValueOf(valueToCompare)
		valueToCompareInt := 0
		typeValue := reflect.ValueOf(value)
		valueInt := 0
		var e error

		// set to 0 if valueToCompare not setted
		if valueToCompare == nil {
			valueToCompare = 0
		}

		// convert valueToCompare to int
		if typeValueToCompare.Kind() != reflect.Int {
			valueToCompareInt, e = strconv.Atoi(fmt.Sprintf("%v", valueToCompare))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueToCompareInt = len(fmt.Sprintf("%v", valueToCompare))
			}
		}

		// set to 0 if value not setted
		if value == nil {
			valueInt = 0
		}

		// convert value to int
		if typeValue.Kind() != reflect.Int {
			valueInt, e = strconv.Atoi(fmt.Sprintf("%v", value))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueInt = len(fmt.Sprintf("%v", value))
			}
		}

		functionGenerated = fmt.Sprintf("%v > %v", valueToCompareInt, value)
		valid = valueToCompareInt > valueInt
	case ">=":
		typeValueToCompare := reflect.ValueOf(valueToCompare)
		valueToCompareInt := 0
		typeValue := reflect.ValueOf(value)
		valueInt := 0
		var e error

		// set to 0 if valueToCompare not setted
		if valueToCompare == nil {
			valueToCompare = 0
		}

		// convert valueToCompare to int
		if typeValueToCompare.Kind() != reflect.Int {
			valueToCompareInt, e = strconv.Atoi(fmt.Sprintf("%v", valueToCompare))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueToCompareInt = len(fmt.Sprintf("%v", valueToCompare))
			}
		}

		// set to 0 if value not setted
		if value == nil {
			valueInt = 0
		}

		// convert value to int
		if typeValue.Kind() != reflect.Int {
			valueInt, e = strconv.Atoi(fmt.Sprintf("%v", value))
			if e != nil {
				assigner.logger.Debug(e.Error())
				valueInt = len(fmt.Sprintf("%v", value))
			}
		}

		functionGenerated = fmt.Sprintf("%v >= %v", valueToCompareInt, value)
		valid = valueToCompareInt >= valueInt
	case "!=":
		valueToCompareParsed := ""
		if valueToCompare != nil {
			valueToCompareParsed = valueToCompare.(string)
		}
		functionGenerated = fmt.Sprintf("%v != %v", valueToCompareParsed, value)
		valid = fmt.Sprintf("%v", valueToCompareParsed) != fmt.Sprintf("%v", value)
	}
	return
}

func (assigner *DefaultAssigner) AssignMiddlewareRequestValue(fieldValue domain.FieldValue, referenceFields map[string]models.Variable) (assigned interface{}) {

	assigner.logger.Debug("AssignValue", zenlogger.ZenField{Key: "fieldValue", Value: fieldValue}, zenlogger.ZenField{Key: "referenceFields", Value: referenceFields})

	if fieldValue.ConditionFieldId.Int64 != 0 && !assigner.ValidityConditionValue(fieldValue, referenceFields[assigner.hasher.GenerateHash(assigner.logger.GetPid(), fmt.Sprintf("%s", fieldValue.ConditionFieldId.Int64))]) {
		return
	}

	memRefFieldKey := assigner.hasher.GenerateHash(assigner.logger.GetPid(), fieldValue.FieldRefId.Int64)
	// if type is object then return
	if referenceFields[memRefFieldKey].Type == variable.TYPE_OBJECT {
		assigned = referenceFields[memRefFieldKey].Value
		return
	} else if fieldValue.Args.String == "$server_request_id" {
		assigned = referenceFields[memRefFieldKey].Value
		return
	}
	commandArg := strings.ReplaceAll(fieldValue.Args.String, "$server_request_id", fmt.Sprintf("%v", referenceFields[memRefFieldKey].Value))

	assigned, err := assigner.ReadCommand(fieldValue.FieldType, commandArg)
	if err != nil {
		assigner.logger.Error("AssignValue", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call ReadCommand"})
	}

	return
}

func (assigner *DefaultAssigner) AssignServerResponseValue(fieldValue domain.FieldValue, referenceFields map[string]models.Variable) (assigned interface{}) {

	assigner.logger.Debug("AssignValue", zenlogger.ZenField{Key: "fieldValue", Value: fieldValue}, zenlogger.ZenField{Key: "referenceFields", Value: referenceFields})

	if fieldValue.ConditionFieldId.Int64 != 0 && !assigner.ValidityConditionValue(fieldValue, referenceFields[assigner.hasher.GenerateHash(assigner.logger.GetPid(), fmt.Sprintf("%s", fieldValue.ConditionFieldId.Int64))]) {
		return
	}

	memRefFieldKey := assigner.hasher.GenerateHash(assigner.logger.GetPid(), fieldValue.FieldRefId.Int64)
	// if type is object then return
	if referenceFields[memRefFieldKey].Type == variable.TYPE_OBJECT {
		assigned = referenceFields[memRefFieldKey].Value
		return
	} else if fieldValue.Args.String == "$middleware_response_id" {
		assigned = referenceFields[memRefFieldKey].Value
		return
	}
	commandArg := strings.ReplaceAll(fieldValue.Args.String, "$middleware_response_id", fmt.Sprintf("%v", referenceFields[memRefFieldKey].Value))

	assigned, err := assigner.ReadCommand(fieldValue.FieldType, commandArg)
	if err != nil {
		assigner.logger.Error("AssignValue", zenlogger.ZenField{Key: "error", Value: err.Error()}, zenlogger.ZenField{Key: "addition", Value: "error while call ReadCommand"})
	}

	return
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func isMapStringInterface(value any) bool {
	_, ok := value.(map[string]any)
	return ok
}

func mapToJSON(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
