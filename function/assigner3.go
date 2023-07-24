package function

import (
	"alfamart-channel/domain"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/zenmuharom/zenfunction"
	"github.com/zenmuharom/zenlogger"
)

type AssignVariableValue struct {
	Key     string
	VarType string
	Value   interface{}
}

type DefaultAssigner2 struct {
	logger   zenlogger.Zenlogger
	function zenfunction.Zenfunction
}

type Assigner2 interface {
	ServerResponseParse(srvConfs []domain.ServerResponseValue, middlewareResponseValue map[string]interface{}) (parsed map[string]interface{}, errs []error)
	ServerResponseConstruct(srConfs []domain.ServerResponse, values map[string]interface{}) (constructed map[string]interface{}, errs []error)
	MiddlewareRequestParse(mrvConfs []domain.MiddlewareRequestValue, middlewareRequestValue map[string]interface{}) (parsed map[string]interface{}, errs []error)
	MiddlewareRequestConstruct(mrConfs []domain.MiddlewareRequest, values map[string]interface{}) (constructed map[string]interface{}, errs []error)
	ParseServerResponseValue(conf domain.ServerResponseValue, valueReference interface{}) (nVal interface{}, err error)
	ParseMiddlewareRequestValue(conf domain.MiddlewareRequestValue, valueReference interface{}) (nVal interface{}, err error)
}

func NewAssigner2(logger zenlogger.Zenlogger) Assigner2 {
	return &DefaultAssigner2{
		logger:   logger,
		function: zenfunction.New(logger),
	}
}

func (assigner *DefaultAssigner2) ServerResponseConstruct(srConfs []domain.ServerResponse, values map[string]interface{}) (constructed map[string]interface{}, errs []error) {
	constructed = make(map[string]interface{}, 0)
	assigner.logger.Debug("ServerResponseConstruct", zenlogger.ZenField{Key: "values", Value: values})
	for key, config := range srConfs {
		assigner.logger.Debug("ServerResponseConstruct", zenlogger.ZenField{Key: "key", Value: key}, zenlogger.ZenField{Key: "config", Value: config}, zenlogger.ZenField{Key: "ParentId", Value: config.ParentId.Int64})

		var parentValue interface{}
		var valueToAssignValue interface{}

		// if key exist in value assign it otherwise set to empty string
		if _, ok := values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)]; ok {
			parentValue = values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)]
		} else {
			// assign empty value based on it's type
			switch config.FieldParentType.String {
			case "string":
				parentValue = ""
			case "integer":
				parentValue = 0
			case "boolean":
				parentValue = false
			case "object":
				parentValue = map[string]interface{}{}
			case "arrayString":
				parentValue = []string{}
			case "arrayInteger":
				parentValue = []int{}
			case "arrayBoolean":
				parentValue = []bool{}
			case "arrayObject":
				parentValue = []interface{}{}
			default:
				parentValue = ""
			}
		}

		// if key exist in value assign it otherwise set to empty string
		if _, ok := values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]; ok {
			valueToAssignValue = values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]
		} else {
			// assign empty value based on it's type
			switch config.Type.String {
			case "string":
				valueToAssignValue = ""
			case "integer":
				valueToAssignValue = 0
			case "boolean":
				valueToAssignValue = false
			case "object":
				valueToAssignValue = map[string]interface{}{}
			case "arrayString":
				valueToAssignValue = []string{}
			case "arrayInteger":
				valueToAssignValue = []int{}
			case "arrayBoolean":
				valueToAssignValue = []bool{}
			case "arrayObject":
				valueToAssignValue = []interface{}{}
			default:
				valueToAssignValue = ""
			}
		}

		valueConfig := domain.ValueConfig{
			FieldName:          sql.NullString{String: config.Field.String},
			ConditionFieldId:   sql.NullInt64{Int64: config.ConditionFieldId.Int64},
			ConditionFieldName: sql.NullString{String: config.ConditionFieldName.String},
			ConditionOperator:  sql.NullString{String: config.ConditionOperator.String},
			ConditionValue:     sql.NullString{String: config.ConditionValue.String},
		}

		if config.ParentId.Int64 != 0 {
			parent := AssignVariableValue{
				Key:     config.FieldParent.String,
				VarType: config.FieldParentType.String,
				Value:   parentValue,
			}

			valueToAssign := AssignVariableValue{
				Key:     config.Field.String,
				VarType: config.Type.String,
				Value:   valueToAssignValue,
			}

			// so the next iterate can get value from this
			if assigner.ValidityConditionField(valueConfig, values[fmt.Sprintf("%v|%v", config.ConditionFieldId.Int64, config.ConditionFieldName.String)]) {
				// fmt.Println(fmt.Sprintf("gua assign jg lo ajg %v", config.Field.String)) // FOR DEBUG
				values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)] = assigner.AssignValue(parent, valueToAssign)
			} else {
				// fmt.Println(fmt.Sprintf("ga masuk bos %v", config.Field.String)) // FOR DEBUG
			}
		} else {
			// constructed[config.Field.String] = values[config.Field.String]
			if assigner.ValidityConditionField(valueConfig, values[fmt.Sprintf("%v|%v", config.ConditionFieldId.Int64, config.ConditionFieldName.String)]) {
				constructed[config.Field.String] = values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]
			}
		}

		assigner.logger.Debug("ServerResponseConstruct", zenlogger.ZenField{Key: "constructed", Value: values})

	}

	return
}

func (assigner *DefaultAssigner2) ServerResponseParse(srvConfs []domain.ServerResponseValue, middlewareResponseValue map[string]interface{}) (parsed map[string]interface{}, errs []error) {
	assigner.logger.Debug("ServerResponseParse", zenlogger.ZenField{Key: "srvConfs", Value: srvConfs}, zenlogger.ZenField{Key: "middlewareResponseValue", Value: middlewareResponseValue})
	fparse := make(map[string]interface{}, 0)
	for _, conf := range srvConfs {
		var err error

		valueConfig := domain.ValueConfig{
			FieldName:          sql.NullString{String: conf.FieldName.String},
			ConditionFieldId:   sql.NullInt64{Int64: conf.ConditionFieldId.Int64},
			ConditionFieldName: sql.NullString{String: conf.ConditionFieldName.String},
			ConditionOperator:  sql.NullString{String: conf.ConditionOperator.String},
			ConditionValue:     sql.NullString{String: conf.ConditionValue.String},
		}
		var valueToCompare interface{}
		if _, ok := fparse[fmt.Sprintf("%v|%v", conf.ConditionFieldId.Int64, conf.ConditionFieldName.String)]; ok {
			valueToCompare = fparse[fmt.Sprintf("%v|%v", conf.ConditionFieldId.Int64, conf.ConditionFieldName.String)]
		}

		// check condition based ON DB
		if assigner.ValidityConditionValue(valueConfig, valueToCompare) {

			// check if middleware_response set or not
			if conf.MiddlewareResponseId.Valid {
				// check if key exist in map fparse then set it otherwise do nothing
				if _, ok := middlewareResponseValue[conf.MiddlewareResponseField.String]; ok {
					assigner.logger.Debug(fmt.Sprintf("middlewareResponseValue[%v] found", conf.MiddlewareResponseField.String))
					fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)], err = assigner.ParseServerResponseValue(conf, middlewareResponseValue[conf.MiddlewareResponseField.String])
					if err != nil {
						errs = append(errs, err)
					}
				} else {
					assigner.logger.Debug(fmt.Sprintf("middlewareResponseValue[%v] not found", conf.MiddlewareResponseField.String))
				}
			} else {
				fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)], err = assigner.ParseServerResponseValue(conf, nil)
			}

		}
	}
	parsed = fparse
	assigner.logger.Debug("ServerResponseParsed", zenlogger.ZenField{Key: "parsed", Value: parsed})
	return
}

func (assigner *DefaultAssigner2) MiddlewareRequestParse(mrvConfs []domain.MiddlewareRequestValue, middlewareRequestValue map[string]interface{}) (parsed map[string]interface{}, errs []error) {
	assigner.logger.Debug("MiddlewareRequestParse", zenlogger.ZenField{Key: "mrvConfs", Value: mrvConfs}, zenlogger.ZenField{Key: "middlewareRequestValue", Value: middlewareRequestValue})
	fparse := make(map[string]interface{}, 0)
	for _, conf := range mrvConfs {
		var err error

		valueConfig := domain.ValueConfig{
			FieldName:          sql.NullString{String: conf.FieldName.String},
			ConditionFieldId:   sql.NullInt64{Int64: conf.ConditionFieldId.Int64},
			ConditionFieldName: sql.NullString{String: conf.ConditionFieldName.String},
			ConditionOperator:  sql.NullString{String: conf.ConditionOperator.String},
			ConditionValue:     sql.NullString{String: conf.ConditionValue.String},
		}
		var valueToCompare interface{}
		if _, ok := fparse[fmt.Sprintf("%v|%v", conf.ConditionFieldId.Int64, conf.ConditionFieldName.String)]; ok {
			valueToCompare = fparse[fmt.Sprintf("%v|%v", conf.ConditionFieldId.Int64, conf.ConditionFieldName.String)]
		}

		// fmt.Println(fmt.Sprintf("valueToCompare fparse[%v|%v] = %v", conf.ConditionFieldId.Int64, conf.ConditionFieldName.String, valueToCompare)) // FOR DEBUG

		// check condition based ON DB
		if assigner.ValidityConditionValue(valueConfig, valueToCompare) {
			// check if the field refer to server request
			if conf.ServerRequestField.String != "" {
				// check if key exist in map fparse then set it otherwise do nothing
				// check if value reference to server request exist

				// check if server request field has parent
				if conf.ServerRequestParentField.String != "" {
					if parentField, ok := middlewareRequestValue[conf.ServerRequestParentField.String].(map[string]interface{}); ok {
						if _, ok := parentField[conf.ServerRequestField.String]; ok {
							fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)], err = assigner.ParseMiddlewareRequestValue(conf, parentField[conf.ServerRequestField.String])
							if err != nil {
								assigner.logger.Error("MiddlewareRequestParse", zenlogger.ZenField{Key: "error", Value: err.Error()})
								errs = append(errs, err)
							}
						} else {
							assigner.logger.Debug("MiddlewareRequestParse", zenlogger.ZenField{Key: "Assign", Value: "Map request not found"})
							fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)] = ""
						}
					}
				} else {
					// if server request doesn't have parent
					if _, ok := middlewareRequestValue[conf.ServerRequestField.String]; ok {
						fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)], err = assigner.ParseMiddlewareRequestValue(conf, middlewareRequestValue[conf.ServerRequestField.String])
						if err != nil {
							assigner.logger.Error("MiddlewareRequestParse", zenlogger.ZenField{Key: "error", Value: err.Error()})
							errs = append(errs, err)
						}
					} else {
						assigner.logger.Debug("MiddlewareRequestParse", zenlogger.ZenField{Key: "Assign", Value: "Map request not found"})
						fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)] = ""
					}

				}
			} else {
				assigner.logger.Debug("MiddlewareRequestParse", zenlogger.ZenField{Key: "Assign", Value: "No server request field set yet"})
				fparse[fmt.Sprintf("%v|%v", conf.FieldId.Int64, conf.FieldName.String)], err = assigner.ParseMiddlewareRequestValue(conf, "")
				if err != nil {
					assigner.logger.Error("MiddlewareRequestParse", zenlogger.ZenField{Key: "error", Value: err.Error()})
					errs = append(errs, err)
				}
			}

		}
	}
	parsed = fparse
	assigner.logger.Debug("MiddlewareRequestParse", zenlogger.ZenField{Key: "parsed", Value: parsed})
	return
}

func (assigner *DefaultAssigner2) MiddlewareRequestConstruct(mrConfs []domain.MiddlewareRequest, values map[string]interface{}) (constructed map[string]interface{}, errs []error) {
	constructed = make(map[string]interface{}, 0)
	assigner.logger.Debug("MiddlewareRequestConstruct", zenlogger.ZenField{Key: "values", Value: values})
	for key, config := range mrConfs {
		assigner.logger.Debug("MiddlewareRequestConstruct", zenlogger.ZenField{Key: "key", Value: key}, zenlogger.ZenField{Key: "config", Value: config}, zenlogger.ZenField{Key: "ParentId", Value: config.ParentId.Int64})

		var parentValue interface{}
		var valueToAssignValue interface{}

		// if key exist in value assign it otherwise set to empty string
		if _, ok := values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)]; ok {
			parentValue = values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)]
		} else {
			parentValue = ""
		}

		// if key exist in value assign it otherwise set to empty string
		if _, ok := values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]; ok {
			valueToAssignValue = values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]
		} else {
			valueToAssignValue = ""
		}

		valueConfig := domain.ValueConfig{
			FieldName:          sql.NullString{String: config.Field.String},
			ConditionFieldId:   sql.NullInt64{Int64: config.ConditionFieldId.Int64},
			ConditionFieldName: sql.NullString{String: config.ConditionFieldName.String},
			ConditionOperator:  sql.NullString{String: config.ConditionOperator.String},
			ConditionValue:     sql.NullString{String: config.ConditionValue.String},
		}

		if config.ParentId.Int64 != 0 {
			parent := AssignVariableValue{
				Key:   config.FieldParent.String,
				Value: parentValue,
			}

			valueToAssign := AssignVariableValue{
				Key:   config.Field.String,
				Value: valueToAssignValue,
			}

			// so the next iterate can get value from this
			if assigner.ValidityConditionField(valueConfig, values[fmt.Sprintf("%v|%v", config.ConditionFieldId.Int64, config.ConditionFieldName.String)]) {
				values[fmt.Sprintf("%v|%v", config.ParentId.Int64, config.FieldParent.String)] = assigner.AssignValue(parent, valueToAssign)
			}
		} else {
			// constructed[config.Field.String] = values[config.Field.String]
			if assigner.ValidityConditionField(valueConfig, values[fmt.Sprintf("%v|%v", config.ConditionFieldId.Int64, config.ConditionFieldName.String)]) {
				constructed[config.Field.String] = values[fmt.Sprintf("%v|%v", config.Id, config.Field.String)]
			}
		}

		assigner.logger.Debug("MiddlewareRequestConstruct", zenlogger.ZenField{Key: "constructed", Value: values})

	}
	return
}

func (assigner *DefaultAssigner2) ValidityConditionField(config domain.ValueConfig, valueToCompare interface{}) (valid bool) {

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

func (assigner *DefaultAssigner2) ValidityConditionValue(config domain.ValueConfig, valueToCompare interface{}) (valid bool) {
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

func (assigner *DefaultAssigner2) validityConditionCore(operator string, valueToCompare interface{}, value interface{}) (functionGenerated string, valid bool) {
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

func (assigner *DefaultAssigner2) AssignValue(parent AssignVariableValue, valueToAssign AssignVariableValue) (assigned interface{}) {

	assigner.logger.Debug("AssignValue", zenlogger.ZenField{Key: "parent", Value: parent}, zenlogger.ZenField{Key: "valueToAssign", Value: valueToAssign})

	valueOfVariable := reflect.ValueOf(parent.Value)

	// if parent variable type is object
	switch parent.VarType {
	case "object":
		mapVariable := make(map[string]interface{}, 0)
		switch valueOfVariable.Kind() {
		case reflect.Map:
			iter := valueOfVariable.MapRange()
			for iter.Next() {
				mapVariable[iter.Key().String()] = iter.Value().Interface()
			}
			mapVariable[valueToAssign.Key] = valueToAssign.Value
			parent.Value = mapVariable
		default:
			mapVariable[valueToAssign.Key] = valueToAssign.Value
		}
		parent.Value = mapVariable
	case "arrayObject":
		// fmt.Println(fmt.Sprintf("field %v is %v to assign to parent %v which is %v", valueToAssign.Key, valueToAssign.VarType, parent.Key, parent.VarType)) // FOR DEBUG
		arrVariable := make([]interface{}, 0)
		switch valueOfVariable.Kind() {
		case reflect.Map:
			// fmt.Println("is map") // FOR DEBUG
			mapVariable := make(map[string]interface{}, 0)
			iter := valueOfVariable.MapRange()
			for iter.Next() {
				mapVariable[iter.Key().String()] = iter.Value().Interface()
			}
			mapVariable[valueToAssign.Key] = valueToAssign.Value
			parent.Value = mapVariable

		case reflect.Array, reflect.Slice:
			// fmt.Println("is array") // FOR DEBUG
			vRef := reflect.ValueOf(parent.Value)
			for i := 0; i < vRef.Len(); i++ {
				arrVariable = append(arrVariable, vRef.Index(i).Interface())
			}
			arrVariable = append(arrVariable, valueToAssign.Value)
			parent.Value = arrVariable
		default:
			// fmt.Println("default") // FOR DEBUG
			arrVariable = append(arrVariable, valueToAssign.Value)
		}
		parent.Value = arrVariable
	case "arrayString":

	case "arrayInteger":

	case "arrayBoolean":
	case "string":
		parent.Value = fmt.Sprintf("%v", valueToAssign.Value)
	case "integer":
		intVal, err := strconv.Atoi(fmt.Sprintf("%v", valueToAssign.Value))
		if err != nil {
			parent.Value = 0
		} else {

			parent.Value = intVal
		}
	case "boolean":
		parent.Value = valueToAssign.Value.(bool)
	default:
		parent.VarType = fmt.Sprintf("%v", valueToAssign.Value)
	}

	assigned = parent.Value

	return
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func (assigner *DefaultAssigner2) ParseServerResponseValue(conf domain.ServerResponseValue, valueReference interface{}) (nVal interface{}, err error) {
	assigner.logger.Debug("ParseServerResponseValue", zenlogger.ZenField{Key: "conf", Value: conf}, zenlogger.ZenField{Key: "valueReference", Value: valueReference})
	commandArg := ""
	if conf.Args.String != "" {

		if conf.MiddlewareResponseId.Int64 != 0 {

			valueStr := fmt.Sprintf("%v", valueReference)
			if !IsJSON(valueStr) {
				valueStr = strings.ReplaceAll(fmt.Sprintf("%v", valueReference), ",", `\,`)
			}

			commandArg = strings.ReplaceAll(conf.Args.String, "$middleware_response_id", valueStr)
		} else {
			commandArg = conf.Args.String
		}

		assigner.logger.Debug("ParseMiddlewareRequestValue", zenlogger.ZenField{Key: "field", Value: conf.FieldName.String}, zenlogger.ZenField{Key: "commandArg", Value: commandArg})

		nVal, err = assigner.function.ReadCommand(commandArg)
		if err != nil {
			assigner.logger.Error(err.Error())
			return
		}
	} else {
		// nVal = valueReference
		// fmt.Println(fmt.Sprintf("%v assign valueReference", conf.FieldName.String)) // FOR DEBUG
		switch conf.FieldType.String {
		case "string":
			nVal = ""
		case "integer":
			nVal = nil
		case "boolean":
			nVal = nil
		case "object":
			var emptyObject map[string]interface{}
			nVal = emptyObject
		case "arrayString":
			var emptyArray []string
			nVal = emptyArray
		case "arrayInteger":
			var emptyArray []int
			nVal = emptyArray
		case "arrayBoolean":
			var emptyArray []bool
			nVal = emptyArray
		case "arrayObject":
			var emptyArray []interface{}
			nVal = emptyArray
		default:
			nVal = nil
		}
	}

	return
}

func (assigner *DefaultAssigner2) ParseMiddlewareRequestValue(conf domain.MiddlewareRequestValue, valueReference interface{}) (nVal interface{}, err error) {
	assigner.logger.Debug("ParseMiddlewareRequestValue", zenlogger.ZenField{Key: "conf", Value: conf}, zenlogger.ZenField{Key: "valueReference", Value: valueReference})
	commandArg := ""
	if conf.Args.String != "" {

		if conf.ServerRequestId.Int64 != 0 {
			valueStr := fmt.Sprintf("%v", valueReference)
			if !IsJSON(valueStr) {
				valueStr = strings.ReplaceAll(fmt.Sprintf("%v", valueReference), ",", `\,`)
			}
			commandArg = strings.ReplaceAll(conf.Args.String, "$server_request_id", fmt.Sprintf("%v", valueStr))
		} else {
			commandArg = conf.Args.String
		}

		assigner.logger.Debug("ParseMiddlewareRequestValue", zenlogger.ZenField{Key: "field", Value: conf.FieldName.String}, zenlogger.ZenField{Key: "commandArg", Value: commandArg})

		nVal, err = assigner.function.ReadCommand(commandArg)
		if err != nil {
			assigner.logger.Error(err.Error())
			return
		}
	} else {
		// nVal = valueReference
		// fmt.Println(fmt.Sprintf("%v assign valueReference", conf.FieldName.String)) // FOR DEBUG
		switch conf.FieldType.String {
		case "string":
			nVal = ""
		case "integer":
			nVal = nil
		case "boolean":
			nVal = nil
		case "object":
			var emptyObject map[string]interface{}
			nVal = emptyObject
		case "arrayString":
			var emptyArray []string
			nVal = emptyArray
		case "arrayInteger":
			var emptyArray []int
			nVal = emptyArray
		case "arrayBoolean":
			var emptyArray []bool
			nVal = emptyArray
		case "arrayObject":
			var emptyArray []interface{}
			nVal = emptyArray
		default:
			nVal = nil
		}
	}

	return

}
