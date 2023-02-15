package boot

import (
	"fmt"
	"regexp"
	"strings"
)

var types = map[string]_field{
	"primitive.ObjectID": {_go: "primitive.ObjectID", _mysql: ""},
	"tinyint":            {_go: "int8", _mysql: "tinyint"},
	"uint8":              {_go: "uint8", _mysql: "tinyint unsigned"},
	"uint16":             {_go: "uint16", _mysql: "smallint unsigned"},
	"uint32":             {_go: "uint32", _mysql: "int unsigned"},
	"uint64":             {_go: "uint64", _mysql: "bigint unsigned"},
	"int8":               {_go: "int8", _mysql: "tinyint"},
	"int16":              {_go: "int16", _mysql: "smallint"},
	"int32":              {_go: "int32", _mysql: "int"},
	"int64":              {_go: "int64", _mysql: "bigint"},
	"bigint":             {_go: "int64", _mysql: "bigint"},
	"float32":            {_go: "float32", _mysql: "float"},
	"float64":            {_go: "float64", _mysql: "double"},
	"string":             {_go: "string", _mysql: "varchar"},
	"json":               {_go: "string", _mysql: "json"},
	"text":               {_go: "string", _mysql: "text"},
	"tinytext":           {_go: "string", _mysql: "tinytext"},
	"mediumtext":         {_go: "string", _mysql: "mediumtext"},
	"longtext":           {_go: "string", _mysql: "longtext"},
	"varchar":            {_go: "string", _mysql: "varchar"},
	"char":               {_go: "string", _mysql: "char"},
	"int":                {_go: "int", _mysql: "int"},
	"uint":               {_go: "uint", _mysql: "int unsigned"},
}

type _field struct {
	_go    string
	_mysql string
}
type field struct {
	Name       string
	Typ        string
	Max        int
	Min        int
	Constraint string
	Comment    string
}

func (f field) GetGoTyp() (string, error) {
	tpl := strings.ToLower(f.Typ)
	if out, ok := types[tpl]; ok {
		return out._go, nil
	}
	return tpl, fmt.Errorf("err type: %s", tpl)
}

func (f field) GetMysqlTyp() (string, error) {
	tpl := strings.ToLower(f.Typ)
	if out, ok := types[tpl]; ok {
		return out._mysql, nil
	}
	return tpl, fmt.Errorf("err type: %s", tpl)
}
func (f field) GetMongoBson() string {
	return " `bson:\"" + f.GetName() + "\"`"
}

func (f field) GetName() string {
	return strings.TrimSpace(f.Name)
}

func (f field) GetMysql() (string, error) {
	out, err := f.GetMysqlTyp()
	if err != nil {
		return "", err
	}
	if out == "varchar" || out == "char" {
		if f.Max < 1 {
			return "", fmt.Errorf("varchar or char type need max : %v", f)
		}
		out += fmt.Sprintf("(%d)", f.Max)
	}
	if f.Constraint != "" {
		out += " " + f.Constraint
	}
	if f.Comment != "" {
		out += " COMMENT '" + f.Comment + "'"
	}

	re := regexp.MustCompile(`'\s+'`)
	out = re.ReplaceAllString(out, "''")

	return out, err
}

type conf struct {
	Project string // 项目
	Server  string // 服务
	UseCase struct {
		Package string
	} `json:"use_case"`
	Repository struct {
		Typ     string
		Version string

		Cache struct {
			Enable           bool
			GetOrSetDuration string `json:"get_or_set_duration"`
			SetDuration      string `json:"set_duration"`
		}
	}
	Entity struct {
		Name         string
		Typ          string
		Label        string
		PrimaryField field `json:"primary_field"`
		Fields       []field
		Keys         []string
	}
}

func (c conf) UseCasePackage() string {
	if c.UseCase.Package != "" {
		return c.UseCase.Package
	}
	return c.packageName()
}
func (c conf) RepositoryTypeIsMysql() bool {
	if c.Repository.Typ == "" || c.Repository.Typ == "mysql" {
		return true
	}
	return false
}
func (c conf) RepositoryTypeIsMongo() bool {
	if c.Repository.Typ == "mongo" {
		return true
	}
	return false
}

func (c conf) GetEntityKeys() string {
	out := ""
	for _, v := range c.Entity.Keys {
		v = strings.TrimSuffix(v, ",")
		v = strings.TrimSuffix(v, "，")
		out += strings.TrimSpace(v) + ",\n"
	}
	out = strings.TrimSuffix(out, ",\n") + "\n"

	return out
}

func (c conf) isMysqlEntity() bool {
	return c.Entity.Typ == "mysql"
}

func (c conf) isMongoEntity() bool {
	return c.Entity.Typ == "mongo"
}
func (c conf) packageName() string {
	return c.Entity.Name
}

// 实体名缩写
func (c conf) shortEntityName() string {
	return strings.ToLower(string([]rune(c.Entity.Name)[0]))
}

// 实体名 大驼峰
func (c conf) publicEntityName() string {
	return upperCamelCase(c.Entity.Name)
}

// 实体私有名称
func (c conf) privateEntityName() string {
	names := strings.Split(c.Entity.Name, "_")
	out := ""
	for k, v := range names {
		if k == 0 {
			out += v
		} else {
			out += strFirstToUpper(v)
		}
	}
	return out
}
