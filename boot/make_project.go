package boot

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/zander-84/seagull/boot/internal/tpl/project"
	"os"
	"path/filepath"
	"strings"
)

type projectConf struct {
	Project  string // 项目
	Server   string // 服务
	HttpPort int    `json:"http_port"`
	GrpcPort int    `json:"grpc_port"`

	confPath string // 配置文件路径
	savePath string // 保存路径

}

func makeProject() error {
	cf, err := makeProjectInit()
	if err != nil {
		return err
	}

	// 1. 创建目录
	root := cf.savePath + "/" + cf.Server

	root_cmd := root + "/" + "cmd"
	root_cmd_api := root_cmd + "/" + "api"
	root_cmd_config := root_cmd + "/" + "config"

	root_internal := root + "/" + "internal"
	root_internal_config := root_internal + "/" + "config"
	root_internal_endpoint := root_internal + "/" + "endpoint"
	root_internal_endpoint_hello := root_internal_endpoint + "/" + "hello"

	root_internal_infrastructure := root_internal + "/" + "infrastructure"
	root_internal_infrastructure_repository := root_internal_infrastructure + "/" + "repository"

	root_internal_pkg := root_internal + "/" + "pkg"
	root_internal_transport := root_internal + "/" + "transport"
	root_internal_usecase := root_internal + "/" + "usecase"
	root_internal_usecase_root := root_internal_usecase + "/" + "root"

	var dirs = []string{
		root,
		root_cmd, root_cmd_api, root_cmd_config,
		root_internal,
		root_internal_config, root_internal_endpoint, root_internal_endpoint_hello, root_internal_infrastructure, root_internal_infrastructure_repository, root_internal_pkg,
		root_internal_transport, root_internal_usecase, root_internal_usecase_root,
	}

	for _, v := range dirs {
		if err := makeDir(v); err != nil {
			return err
		}
	}

	//root_cmd_api
	mainTpl := project.MainTpl
	mainTpl = replaceProjectGlobalVar(mainTpl, cf)
	if _, err := save(root_cmd_api, "main.go", mainTpl); err != nil {
		return err
	}

	bsTpl := project.BsTpl
	bsTpl = replaceProjectGlobalVar(bsTpl, cf)
	if _, err := save(root_cmd_api, "bs.go", bsTpl); err != nil {
		return err
	}

	// root_internal_config
	confTpl := project.ConfTpl
	confTpl = replaceProjectGlobalVar(confTpl, cf)
	if _, err := save(root_internal_config, "config.go", confTpl); err != nil {
		return err
	}

	//root_internal_pkg
	PkgCorsTpl := project.PkgCorsTpl
	PkgCorsTpl = replaceProjectGlobalVar(PkgCorsTpl, cf)
	if _, err := save(root_internal_pkg, "cors.go", PkgCorsTpl); err != nil {
		return err
	}

	PkgEtcdTpl := project.PkgEtcdTpl
	PkgEtcdTpl = replaceProjectGlobalVar(PkgEtcdTpl, cf)
	if _, err := save(root_internal_pkg, "etcd.go", PkgEtcdTpl); err != nil {
		return err
	}
	PkgMongoTpl := project.PkgMongoTpl
	PkgMongoTpl = replaceProjectGlobalVar(PkgMongoTpl, cf)
	if _, err := save(root_internal_pkg, "mongo.go", PkgMongoTpl); err != nil {
		return err
	}
	PkgMysqlTpl := project.PkgMysqlTpl
	PkgMysqlTpl = replaceProjectGlobalVar(PkgMysqlTpl, cf)
	if _, err := save(root_internal_pkg, "mysql.go", PkgMysqlTpl); err != nil {
		return err
	}
	PkgRedisTpl := project.PkgRedisTpl
	PkgRedisTpl = replaceProjectGlobalVar(PkgRedisTpl, cf)
	if _, err := save(root_internal_pkg, "redis.go", PkgRedisTpl); err != nil {
		return err
	}
	PkgInitTpl := project.PkgInitTpl
	PkgInitTpl = replaceProjectGlobalVar(PkgInitTpl, cf)
	if _, err := save(root_internal_pkg, "init.go", PkgInitTpl); err != nil {
		return err
	}

	//root_internal_usecase
	UsecaseTpl := project.UsecaseTpl
	UsecaseTpl = replaceProjectGlobalVar(UsecaseTpl, cf)
	if _, err := save(root_internal_usecase, "usecase.go", UsecaseTpl); err != nil {
		return err
	}
	//root_internal_usecase_root
	UsecaserootTpl := project.UsecaserootTpl
	UsecaserootTpl = replaceProjectGlobalVar(UsecaserootTpl, cf)
	if _, err := save(root_internal_usecase_root, "service.go", UsecaserootTpl); err != nil {
		return err
	}

	UsecaserootInterfaceTpl := project.UsecaserootInterfaceTpl
	UsecaserootInterfaceTpl = replaceProjectGlobalVar(UsecaserootInterfaceTpl, cf)
	if _, err := save(root_internal_usecase_root, "interface.go", UsecaserootInterfaceTpl); err != nil {
		return err
	}

	//root_internal_transport
	RmcTpl := project.RmcTpl
	RmcTpl = replaceProjectGlobalVar(RmcTpl, cf)
	if _, err := save(root_internal_transport, "rmc.go", RmcTpl); err != nil {
		return err
	}

	// root_internal_endpoint_hello
	HelloTpl := project.HelloTpl
	HelloTpl = replaceProjectGlobalVar(HelloTpl, cf)
	if _, err := save(root_internal_endpoint_hello, "hello.go", HelloTpl); err != nil {
		return err
	}

	HelloCodecTpl := project.HelloCodecTpl
	HelloCodecTpl = replaceProjectGlobalVar(HelloCodecTpl, cf)
	if _, err := save(root_internal_endpoint_hello, "hello_codec.go", HelloCodecTpl); err != nil {
		return err
	}
	return err
}

func replaceProjectGlobalVar(content string, conf *projectConf) string {
	content = strings.ReplaceAll(content, "${project}", conf.Project)
	content = strings.ReplaceAll(content, "${server}", conf.Server)
	content = strings.ReplaceAll(content, "${httpPort}", fmt.Sprintf("%d", conf.HttpPort))
	content = strings.ReplaceAll(content, "${grpcPort}", fmt.Sprintf("%d", conf.GrpcPort))
	return content
}

func makeDir(d string) error {
	fmt.Println("make dir :" + d)
	if err := os.Mkdir(d, 0755); err != nil {
		return err
	}
	return nil
}
func makeProjectInit() (*projectConf, error) {
	var in string
	flag.StringVar(&in, "i", "", "入参")
	in = strings.TrimSpace(in)

	var out string
	flag.StringVar(&out, "o", "", "保存位置")
	out = strings.TrimSpace(out)

	flag.Parse()
	if in == "" {
		return nil, errors.New("i 入参不能为空")
	}

	if out == "" {
		out = filepath.Dir(in)
	}
	if !isDir(out) {
		return nil, fmt.Errorf("save path err: 【%s】", out)
	}

	fileData, err := readFile(in)
	if err != nil {
		return nil, err
	}
	cf := new(projectConf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return nil, err
	}
	if cf.Project == "" || cf.Server == "" || cf.HttpPort < 1 || cf.GrpcPort < 1 {
		return nil, fmt.Errorf("config err,exist empty val: 【%v】", cf)
	}
	cf.confPath = in
	cf.savePath = strings.TrimSuffix(out, "/")
	if !strings.HasSuffix(cf.savePath, cf.Project+"/"+"apps") {
		return nil, fmt.Errorf("The file directory must be under %s \n", cf.Project+"/"+"apps")
	}
	if ok, err := fileExist(cf.savePath + "/" + cf.Server); err != nil {
		return nil, err
	} else if ok {
		return nil, fmt.Errorf("dist file exist: %s/%s", cf.savePath, cf.Server)
	}
	return cf, err

}
