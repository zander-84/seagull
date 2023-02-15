package boot

import (
	"fmt"
	"github.com/zander-84/seagull/boot/internal/tpl"
	"strings"
)

func makeEndpoint(conf conf) (string, error) {
	endpointTpl := tpl.Endpoint
	endpointTpl = replaceGlobalVar(endpointTpl, conf)

	assignCreateFields, err := _makeEndpointAssignCreateFields(conf)
	if err != nil {
		return "", err
	}
	endpointTpl = strings.ReplaceAll(endpointTpl, "${assignCreateFields}", assignCreateFields)

	assignUpdateFields, err := _makeEndpointAssignUpdateFields(conf)
	if err != nil {
		return "", err
	}
	endpointTpl = strings.ReplaceAll(endpointTpl, "${assignUpdateFields}", assignUpdateFields)

	return endpointTpl, nil
}

func makeEndpointCodec(conf conf) (string, error) {
	endpointTpl := tpl.EndpointCodec
	endpointTpl = replaceGlobalVar(endpointTpl, conf)

	fields, err := _makeEndpointFields(conf)
	if err != nil {
		return "", err
	}
	endpointTpl = strings.ReplaceAll(endpointTpl, "${fields}", fields)

	idWithType := ""
	outId := ""
	if conf.isMysqlEntity() {
		idWithType = "Id int64"
		outId = `out.Id = conv.ShouldStringToInt64(httpCtx.Param("id"))`
	} else if conf.isMongoEntity() {
		idWithType = "Pk string"
		outId = `out.Id = httpCtx.Param("id")`
	}
	endpointTpl = strings.ReplaceAll(endpointTpl, "${IdWithType}", idWithType)
	endpointTpl = strings.ReplaceAll(endpointTpl, "${outId}", outId)

	idsWithType := ""
	if conf.isMysqlEntity() {
		idsWithType = "Ids []int64"
	} else if conf.isMongoEntity() {
		idsWithType = "Ids []string"
	}
	endpointTpl = strings.ReplaceAll(endpointTpl, "${idsWithType}", idsWithType)

	return endpointTpl, nil
}

func _makeEndpointAssignUpdateFields(conf conf) (string, error) {
	fields := ""
	for _, v := range conf.Entity.Fields {
		if strings.ToLower(v.Name) == "id" {
			if conf.isMysqlEntity() {
				continue
			}
		}
		if strings.ToLower(v.Name) == "created_at" || strings.ToLower(v.Name) == "updated_at" || strings.ToLower(v.Name) == "version" {
			continue
		}
		fd := fmt.Sprintf("%s.Update%s(in.%s)", conf.privateEntityName(), upperCamelCase(v.GetName()), upperCamelCase(v.GetName()))
		fd += "\r\n"
		fields += "		" + fd
	}
	return fields, nil
}

func _makeEndpointAssignCreateFields(conf conf) (string, error) {
	fields := ""
	for _, v := range conf.Entity.Fields {
		if strings.ToLower(v.Name) == "id" {
			if conf.isMysqlEntity() {
				continue
			}
		}
		if strings.ToLower(v.Name) == "created_at" || strings.ToLower(v.Name) == "updated_at" || strings.ToLower(v.Name) == "version" {
			continue
		}
		fd := fmt.Sprintf("%s.%s=in.%s", conf.privateEntityName(), upperCamelCase(v.GetName()), upperCamelCase(v.GetName()))
		fd += "\r\n"
		fields += "		" + fd
	}
	return fields, nil
}
func _makeEndpointFields(conf conf) (string, error) {
	fields := ""

	for _, v := range conf.Entity.Fields {
		if strings.ToLower(v.Name) == "id" {
			if conf.isMysqlEntity() {
				continue
			}
		}
		if strings.ToLower(v.Name) == "created_at" || strings.ToLower(v.Name) == "updated_at" {
			continue
		}
		fd, err := _makeEndpointField(conf, v)
		if err != nil {
			return "", err
		}
		fd += "\r\n"
		fields += "		" + fd
	}
	return fields, nil
}

func _makeEndpointField(conf conf, v field) (string, error) {
	fd := ""

	// 字段名
	fd += upperCamelCase(v.GetName()) + "   "

	//+ 类型
	gotype, err := v.GetGoTyp()
	if err != nil {
		return "", err
	}

	fd += gotype

	return fd, err
}
