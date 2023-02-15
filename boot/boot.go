package boot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakeProject() error {
	return makeProject()
}
func Make() error {
	action, filepath, savePath, err := ParseCmd()
	_ = savePath
	if err != nil {
		return err
	}
	if action == "entity" {
		if _, err := MakeEntity(filepath); err != nil {
			return fmt.Errorf("boot.MakeEntity err: %s", err.Error())
		}
	} else if action == "mysql" {
		if _, err := MakeMysql(filepath); err != nil {
			return fmt.Errorf("boot.MakeMysql err: %s", err.Error())
		}

	} else if action == "repository" {
		if _, err := MakeRepository(filepath); err != nil {
			return fmt.Errorf("boot.MakeRepository err: %s", err.Error())
		}
	} else if action == "usecase" {
		if _, err := MakeUseCase(filepath); err != nil {
			return fmt.Errorf("boot.MakeRepository err: %s", err.Error())
		}
	} else if action == "endpoint" {
		if _, err := MakeEndpoint(filepath); err != nil {
			return fmt.Errorf("boot.MakeEndpoint err: %s", err.Error())
		}
	} else {
		if _, err := MakeEntity(filepath); err != nil {
			return fmt.Errorf("boot.MakeEntity err: %s", err.Error())
		}
		if _, err := MakeMysql(filepath); err != nil {
			return fmt.Errorf("boot.MakeMysql err: %s", err.Error())
		}
		if _, err := MakeRepository(filepath); err != nil {
			return fmt.Errorf("boot.MakeRepository err: %s", err.Error())
		}
		if _, err := MakeUseCase(filepath); err != nil {
			return fmt.Errorf("boot.MakeRepository err: %s", err.Error())
		}
		if _, err := MakeEndpoint(filepath); err != nil {
			return fmt.Errorf("boot.MakeEndpoint err: %s", err.Error())
		}
	}

	return nil
}
func MakeEntity(fileName string) (string, error) {
	fileData, err := readFile(fileName)
	if err != nil {
		return "", err
	}
	cf := new(conf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return "", err
	}

	content, err := makeEntity(*cf)
	if err != nil {
		return "", err
	}

	filePath, err := save(filepath.Dir(fileName), "entity."+cf.Entity.Name+"_bee.go.tpl", content)
	if err != nil {
		return "", err
	}
	fmt.Printf("make entity success save on %s \n", filePath)
	return "", nil
}

func MakeMysql(fileName string) (string, error) {
	fileData, err := readFile(fileName)
	if err != nil {
		return "", err
	}
	cf := new(conf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return "", err
	}

	content, err := makeMySql(*cf)
	if err != nil {
		return "", err
	}

	filePath, err := save(filepath.Dir(fileName), cf.Entity.Name+".sql", content)
	if err != nil {
		return "", err
	}
	fmt.Printf("make mysql success save on %s \n", filePath)
	return "", nil
}
func MakeRepository(fileName string) (string, error) {
	fileData, err := readFile(fileName)
	if err != nil {
		return "", err
	}
	cf := new(conf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return "", err
	}

	content, err := makeRepositoryBasic(*cf)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", errors.New("content empty")
	}
	fn := ""
	if cf.RepositoryTypeIsMysql() {
		fn = "repository." + cf.Entity.Name + "_mysql.go.tpl"
	} else if cf.RepositoryTypeIsMongo() {
		fn = "repository." + cf.Entity.Name + "_mongo.go.tpl"
	} else {
		return "", errors.New("err RepositoryType")
	}
	filePath, err := save(filepath.Dir(fileName), fn, content)
	if err != nil {
		return "", err
	}
	fmt.Printf("make repository basic success save on %s \n", filePath)

	content2, err := makeRepository(*cf)
	if err != nil {
		return "", err
	}

	if content2 == "" {
		return "", errors.New("content empty")
	}
	fn2 := ""
	fn2 = "repository." + cf.Entity.Name + "_repository.go.tpl"

	filePath2, err := save(filepath.Dir(fileName), fn2, content2)
	if err != nil {
		return "", err
	}
	fmt.Printf("make repository success save on %s \n", filePath2)

	return "", nil
}

func MakeUseCase(fileName string) (string, error) {
	fileData, err := readFile(fileName)
	if err != nil {
		return "", err
	}
	cf := new(conf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return "", err
	}

	content, err := makeUseCaseInterface(*cf)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", errors.New("content empty")
	}
	fn := ""
	fn = "usecase.interface" + ".go.tpl"

	filePath, err := save(filepath.Dir(fileName), fn, content)
	if err != nil {
		return "", err
	}
	fmt.Printf("make usecase.interface  success save on %s \n", filePath)

	content2, err := makeUseCaseServer(*cf)
	if err != nil {
		return "", err
	}

	if content2 == "" {
		return "", errors.New("content empty")
	}
	fn2 := ""
	fn2 = "usecase.server" + ".go.tpl"

	filePath2, err := save(filepath.Dir(fileName), fn2, content2)
	if err != nil {
		return "", err
	}
	fmt.Printf("make usecase.server  success save on %s \n", filePath2)

	return "", nil

}

func MakeEndpoint(fileName string) (string, error) {
	fileData, err := readFile(fileName)
	if err != nil {
		return "", err
	}
	cf := new(conf)
	if err := json.Unmarshal([]byte(fileData), cf); err != nil {
		return "", err
	}

	content, err := makeEndpoint(*cf)
	if err != nil {
		return "", err
	}

	filePath, err := save(filepath.Dir(fileName), "endpoint."+cf.Entity.Name+".go.tpl", content)
	if err != nil {
		return "", err
	}
	fmt.Printf("make endpiont success save on %s \n", filePath)

	content2, err := makeEndpointCodec(*cf)
	if err != nil {
		return "", err
	}

	filePath2, err := save(filepath.Dir(fileName), "endpoint."+cf.Entity.Name+"_codec.go.tpl", content2)
	if err != nil {
		return "", err
	}
	fmt.Printf("make endpiont codec success save on %s \n", filePath2)
	return "", nil
}

func save(dir, filename string, content string) (string, error) {
	filePath := strings.TrimSuffix(dir, "/") + "/" + filename
	return filePath, os.WriteFile(filePath, []byte(content), 0644)
}

func replaceGlobalVar(content string, conf conf) string {
	content = strings.ReplaceAll(content, "${project}", conf.Project)
	content = strings.ReplaceAll(content, "${server}", conf.Server)
	content = strings.ReplaceAll(content, "${entityName}", conf.privateEntityName())
	content = strings.ReplaceAll(content, "${EntityName}", conf.publicEntityName())
	content = strings.ReplaceAll(content, "${shortEntityName}", conf.shortEntityName())
	content = strings.ReplaceAll(content, "${tableName}", conf.packageName())
	content = strings.ReplaceAll(content, "${useCasePkg}", conf.UseCasePackage())

	if conf.isMongoEntity() {
		content = strings.ReplaceAll(content, "${mongoPkg}", "\"go.mongodb.org/mongo-driver/bson/primitive\"")
	} else {
		content = strings.ReplaceAll(content, "${mongoPkg}", "")
	}
	return content
}
