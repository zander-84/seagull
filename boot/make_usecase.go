package boot

import (
	"github.com/zander-84/seagull/boot/internal/tpl"
	"strings"
)

func makeUseCaseInterface(conf conf) (string, error) {
	repoTpl := ""
	if conf.Repository.Typ == "" || conf.Repository.Typ == "mysql" {
		repoTpl = tpl.UseCaseInterfaceMysql

	} else if conf.Repository.Typ == "mongo" {
		repoTpl = tpl.UseCaseInterfaceMongo
	}
	repoTpl = replaceGlobalVar(repoTpl, conf)
	repoTpl = strings.ReplaceAll(repoTpl, "${shortEntityName}", conf.shortEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${entityName}", conf.publicEntityName())
	return repoTpl, nil
}

func makeUseCaseServer(conf conf) (string, error) {
	repoTpl := ""
	if conf.Repository.Typ == "" || conf.Repository.Typ == "mysql" {
		repoTpl = tpl.UseCaseServerMysql

	} else if conf.Repository.Typ == "mongo" {
		repoTpl = tpl.UseCaseServerMongo
	}
	repoTpl = replaceGlobalVar(repoTpl, conf)
	repoTpl = strings.ReplaceAll(repoTpl, "${shortEntityName}", conf.shortEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${entityName}", conf.publicEntityName())
	return repoTpl, nil
}
