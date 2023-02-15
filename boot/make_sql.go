package boot

import (
	"github.com/zander-84/seagull/boot/internal/tpl"
	"strings"
)

func makeMySql(conf conf) (string, error) {
	mysqlTpl := tpl.Mysql
	mysqlTpl = strings.ReplaceAll(mysqlTpl, "${tableName}", conf.Entity.Name)
	mysqlTpl = strings.ReplaceAll(mysqlTpl, "${comment}", conf.Entity.Label)

	fields, err := _makeMysqlFields(conf)
	if err != nil {
		return "", err
	}
	mysqlTpl = strings.ReplaceAll(mysqlTpl, "${fields}", fields)

	keys := _makeMysqlKeys(conf)
	mysqlTpl = strings.ReplaceAll(mysqlTpl, "${keys}", keys)

	return mysqlTpl, nil
}

func _makeMysqlKeys(conf conf) string {
	out := strings.Join(conf.Entity.Keys, ",\r\n")
	return strings.TrimSuffix(out, ",\r\n")
}

func _makeMysqlFields(conf conf) (string, error) {
	fields := ""
	for _, v := range conf.Entity.Fields {
		fd, err := _makeMysqlField(conf, v)
		if err != nil {
			return "", err
		}
		fd += ",\r\n"
		fields += "		" + fd
	}
	return strings.TrimSuffix(fields, "\r\n"), nil
}

func _makeMysqlField(conf conf, v field) (string, error) {
	fd := ""

	// 字段名
	fd += "`" + v.GetName() + "`   "

	lines, err := v.GetMysql()
	if err != nil {
		return "", err
	}

	fd += lines

	return fd, err
}
