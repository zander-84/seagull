package tpl

var EntityUpdateFunc = `
func (${shortEntityName} *${EntityName}) Update${Field}(${field} ${type}) { 
	${shortEntityName}.updateFields.Update(${shortEntityName}.FieldName${Field}(), ${field})
	${shortEntityName}.${Field} = ${field}
}`
var EntityFieldName = `
func (${shortEntityName} *${EntityName}) FieldName${Field}() string{ return "${field}"}`

var EntityValidateField = `
func (${shortEntityName} *${EntityName})  Validate${Field}() error{${minString}${maxString}${minNumber}${maxNumber}  return nil}`

var EntityValidateFields = `
func (${shortEntityName} *${EntityName}) Validate() error{ 
${validateFields}
	return nil
}`

var EntityGetVersion = `
func (${shortEntityName} *${EntityName}) GetVersion() int{ 
	return ${shortEntityName}.Version
}`

var EntityUpdatedFields = `
func (${shortEntityName} *${EntityName}) UpdatedFields() map[string]any{ 
	return ${shortEntityName}.updateFields.Get()
}`

var Entity = `package entity

import (
	"errors"
	"github.com/zander-84/seagull/contract"
	${mongoPkg}
)

type ${EntityName} struct {
${fields}
	updateFields contract.UpdateFields
}

func New${EntityName} (preCheck *${EntityName}) (*${EntityName},error) {
    err := preCheck.Validate()
    return preCheck, err
}

${entityValidateField}


${entityValidate}


${updateFunc}

${updateFields}

${getVersionFunc}


${getFieldNameFunc}

`
