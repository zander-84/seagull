package boot

import (
	"fmt"
	"github.com/zander-84/seagull/boot/internal/tpl"
	"strings"
)

func makeEntity(conf conf) (string, error) {
	entityTpl := tpl.Entity
	entityTpl = replaceGlobalVar(entityTpl, conf)

	fields, err := _makeEntityFields(conf)
	if err != nil {
		return "", err
	}
	entityTpl = strings.ReplaceAll(entityTpl, "${fields}", fields)

	updatesFunc, err := _makeEntityUpdatesFunc(conf)
	if err != nil {
		return "", err
	}
	entityTpl = strings.ReplaceAll(entityTpl, "${updateFunc}", updatesFunc)

	getFieldName, err := _makeEntityFieldName(conf)
	if err != nil {
		return "", err
	}
	entityTpl = strings.ReplaceAll(entityTpl, "${getFieldNameFunc}", getFieldName)

	validateFiled, err := _makeEntityValidateFiled(conf)
	if err != nil {
		return "", err
	}
	entityTpl = strings.ReplaceAll(entityTpl, "${entityValidateField}", validateFiled)

	validate, err := _makeEntityValidate(conf)
	if err != nil {
		return "", err
	}
	entityTpl = strings.ReplaceAll(entityTpl, "${entityValidate}", validate)

	getVersion := _makeEntityGetVersion(conf)
	entityTpl = strings.ReplaceAll(entityTpl, "${getVersionFunc}", getVersion)

	updatedFields := _makeEntityUpdatedFields(conf)
	entityTpl = strings.ReplaceAll(entityTpl, "${updateFields}", updatedFields)

	return entityTpl, nil
}

func _makeEntityUpdatedFields(conf conf) string {
	out := tpl.EntityUpdatedFields
	out = strings.ReplaceAll(out, "${shortEntityName}", conf.shortEntityName())
	out = strings.ReplaceAll(out, "${EntityName}", conf.publicEntityName())
	return out
}

func _makeEntityGetVersion(conf conf) string {
	out := tpl.EntityGetVersion
	out = strings.ReplaceAll(out, "${shortEntityName}", conf.shortEntityName())
	out = strings.ReplaceAll(out, "${EntityName}", conf.publicEntityName())
	return out
}

func _makeEntityValidate(conf conf) (string, error) {
	out := tpl.EntityValidateFields
	validateFields := ""
	for _, v := range conf.Entity.Fields {
		_, err := v.GetGoTyp()
		if err != nil {
			return "", err
		}
		validateFields += fmt.Sprintf(`		if err:=%s.Validate%s();err!=nil{ return err }
`, conf.shortEntityName(), upperCamelCase(v.GetName()))
	}
	out = strings.ReplaceAll(out, "${shortEntityName}", conf.shortEntityName())
	out = strings.ReplaceAll(out, "${EntityName}", conf.publicEntityName())
	out = strings.ReplaceAll(out, "${validateFields}", validateFields)

	return out, nil
}
func _makeEntityValidateFiled(conf conf) (string, error) {
	out := ""
	for _, v := range conf.Entity.Fields {
		tmp := tpl.EntityValidateField
		gt, err := v.GetGoTyp()
		if err != nil {
			return "", err
		}
		minString := ""
		maxString := ""
		minNumber := ""
		maxNumber := ""

		if gt == "string" {
			if v.Min > 0 {
				minString = fmt.Sprintf(`
	if len(%s.%s) < %d {
		return errors.New("%s Min len  %d")
	}
`, conf.shortEntityName(), upperCamelCase(v.GetName()), v.Min, upperCamelCase(v.GetName()), v.Min)
			}

			if v.Max > 0 {
				maxString = fmt.Sprintf(`
	if len(%s.%s) > %d {
		return errors.New("%s  Max len %d")
	}
`, conf.shortEntityName(), upperCamelCase(v.GetName()), v.Max, upperCamelCase(v.GetName()), v.Max)
			}
		}
		tmp = strings.ReplaceAll(tmp, "${shortEntityName}", conf.shortEntityName())
		tmp = strings.ReplaceAll(tmp, "${EntityName}", conf.publicEntityName())
		tmp = strings.ReplaceAll(tmp, "${Field}", upperCamelCase(v.GetName()))
		tmp = strings.ReplaceAll(tmp, "${minString}", minString)
		tmp = strings.ReplaceAll(tmp, "${maxString}", maxString)
		tmp = strings.ReplaceAll(tmp, "${minNumber}", minNumber)
		tmp = strings.ReplaceAll(tmp, "${maxNumber}", maxNumber)
		out += tmp
	}
	return out, nil
}

func _makeEntityFieldName(conf conf) (string, error) {
	out := ""
	for _, v := range conf.Entity.Fields {
		tmp := tpl.EntityFieldName
		gt, err := v.GetGoTyp()
		if err != nil {
			return "", err
		}

		tmp = strings.ReplaceAll(tmp, "${shortEntityName}", conf.shortEntityName())
		tmp = strings.ReplaceAll(tmp, "${EntityName}", conf.publicEntityName())
		tmp = strings.ReplaceAll(tmp, "${Field}", upperCamelCase(v.GetName()))
		tmp = strings.ReplaceAll(tmp, "${field}", v.GetName())
		tmp = strings.ReplaceAll(tmp, "${type}", gt)
		out += tmp
	}
	return out, nil
}
func _makeEntityUpdatesFunc(conf conf) (string, error) {
	updatesFunc := ""
	for _, v := range conf.Entity.Fields {
		updateFunc := tpl.EntityUpdateFunc
		gt, err := v.GetGoTyp()
		if err != nil {
			return "", err
		}
		if strings.ToLower(v.GetName()) == "id" {
			continue
		}
		updateFunc = strings.ReplaceAll(updateFunc, "${shortEntityName}", conf.shortEntityName())
		updateFunc = strings.ReplaceAll(updateFunc, "${EntityName}", conf.publicEntityName())
		updateFunc = strings.ReplaceAll(updateFunc, "${Field}", upperCamelCase(v.GetName()))
		updateFunc = strings.ReplaceAll(updateFunc, "${field}", lowerCamelCase(v.GetName()))
		updateFunc = strings.ReplaceAll(updateFunc, "${type}", gt)
		updatesFunc += updateFunc
	}

	return updatesFunc, nil
}
func _makeEntityFields(conf conf) (string, error) {
	fields := ""
	if conf.isMongoEntity() {
		fd := "Pk primitive.ObjectID" + "`bson:\"_id,omitempty\"  gorm:\"-\"`"
		fd += "\r\n"
		fields += "		" + fd
	}
	for _, v := range conf.Entity.Fields {
		fd, err := _makeEntityField(conf, v)
		if err != nil {
			return "", err
		}
		fd += "\r\n"
		fields += "		" + fd
	}
	return fields, nil
}

func _makeEntityField(conf conf, v field) (string, error) {
	fd := ""

	// 字段名
	fd += upperCamelCase(v.GetName()) + "   "

	//+ 类型
	gotype, err := v.GetGoTyp()
	if err != nil {
		return "", err
	}

	fd += gotype

	// + tag  bson
	fd += v.GetMongoBson()

	// + 注释
	fd += " // " + v.Comment
	// + 注释mysql
	if conf.isMysqlEntity() {
		mysqlLine, err := v.GetMysql()
		if err != nil {
			return "", err
		}
		fd += " sql: " + mysqlLine
	}

	return fd, err
}
