package boot

import (
	"github.com/zander-84/seagull/boot/internal/tpl"
	"strings"
)

func makeRepositoryBasic(conf conf) (string, error) {
	repoTpl := ""
	if conf.Repository.Typ == "" || conf.Repository.Typ == "mysql" {
		if conf.Repository.Cache.Enable {
			repoTpl = tpl.RepositoryMysqlCache
		} else {
			repoTpl = tpl.RepositoryMysql
		}
	} else if conf.Repository.Typ == "mongo" {
		if conf.Repository.Cache.Enable {
			repoTpl = tpl.RepositoryMongoCache
		} else {
			repoTpl = tpl.RepositoryMongo
		}

	}
	repoTpl = replaceGlobalVar(repoTpl, conf)
	repoTpl = strings.ReplaceAll(repoTpl, "${shortEntityName}", conf.shortEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${entityName}", conf.publicEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${version}", conf.Repository.Version)
	repoTpl = strings.ReplaceAll(repoTpl, "${cacheGetOrSetDuration}", conf.Repository.Cache.GetOrSetDuration)
	repoTpl = strings.ReplaceAll(repoTpl, "${cacheSetDuration}", conf.Repository.Cache.SetDuration)
	return repoTpl, nil
}

func makeRepository(conf conf) (string, error) {
	repoTpl := ""
	if conf.Repository.Typ == "" || conf.Repository.Typ == "mysql" {
		if conf.Repository.Cache.Enable {
			repoTpl = tpl.RepoMysqlCache
		} else {
			repoTpl = tpl.RepoMysql
		}

	} else if conf.Repository.Typ == "mongo" {
		if conf.Repository.Cache.Enable {
			repoTpl = tpl.RepoMongoCache
		} else {
			repoTpl = tpl.RepoMongo
		}
	}
	repoTpl = replaceGlobalVar(repoTpl, conf)
	repoTpl = strings.ReplaceAll(repoTpl, "${shortEntityName}", conf.shortEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${entityName}", conf.publicEntityName())
	repoTpl = strings.ReplaceAll(repoTpl, "${version}", conf.Repository.Version)
	repoTpl = strings.ReplaceAll(repoTpl, "${cacheGetOrSetDuration}", conf.Repository.Cache.GetOrSetDuration)
	repoTpl = strings.ReplaceAll(repoTpl, "${cacheSetDuration}", conf.Repository.Cache.SetDuration)
	return repoTpl, nil
}
