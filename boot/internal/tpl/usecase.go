package tpl

var UseCaseInterfaceMysql = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"${project}/apps/${server}/internal/entity"
)

type Reader interface {
	Exist${EntityName}(key string, val any) (bool, error)
	Get${EntityName}(id int64) (*entity.${EntityName}, error)
	BatchGet${EntityName}(ids []int64) ([]entity.${EntityName}, error)
	Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error)
}

type Writer interface {
	Create${EntityName}(${entityName} *entity.${EntityName}) error
	Update${EntityName}(id int64, version int, entity *entity.${EntityName}) error
	Update${EntityName}Map(id int64, version int, data map[string]any) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get${EntityName}(ctx context.Context, id int64) (*entity.${EntityName}, error)
	BatchGet${EntityName}(ctx context.Context, ids []int64) ([]entity.${EntityName}, error)
	Create${EntityName}(ctx context.Context, ${entityName} *entity.${EntityName}) error
	Update${EntityName}(ctx context.Context, id int64, version int, ${entityName} *entity.${EntityName}) error
	Update${EntityName}Map(ctx context.Context, id int64, version int, ${entityName} map[string]any) error
	Search${EntityName}(ctx context.Context, searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error)
}`
var UseCaseInterfaceMongo = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"${project}/apps/${server}/internal/entity"
)

type Reader interface {
	Exist(key string, val any) (bool, error)
	Get${EntityName}(id string) (*entity.${EntityName}, error)
	BatchGet${EntityName}(ids []string) ([]entity.${EntityName}, error)
	Search${EntityName}(searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error)
}

type Writer interface {
	Create${EntityName}(${entityName} *entity.${EntityName}) error
	Update${EntityName}(id string, version int, entity *entity.${EntityName}) error
	Update${EntityName}Map(id string, version int, data map[string]any) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get${EntityName}(ctx context.Context, id string) (*entity.${EntityName}, error)
	BatchGet${EntityName}(ctx context.Context, ids []string) ([]entity.${EntityName}, error)
	Create${EntityName}(ctx context.Context, ${entityName} *entity.${EntityName}) error
	Update${EntityName}(ctx context.Context, id string, version int, ${entityName} *entity.${EntityName}) error
	Search${EntityName}(ctx context.Context, searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error)
}`

var UseCaseServerMysql = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"${project}/apps/${server}/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	out := new(Service)
	out.repo = repo

	return out
}

func (s *Service) Get${EntityName}(ctx context.Context, id int64) (*entity.${EntityName}, error) {
	out, err := s.repo.Get${EntityName}(id)
	if err != nil {
		if think.IsErrNotFound(err) {
			return nil, err
		}
		return nil, think.ErrSystemSpace(err.Error())
	}
	return out, err
}

func (s *Service) BatchGet${EntityName}(ctx context.Context, ids []int64) ([]entity.${EntityName}, error) {
	out, err := s.repo.BatchGet${EntityName}(ids)
	if err != nil {
		return nil, think.ErrSystemSpace(err.Error())
	}
	return out, err
}

func (s *Service) Create${EntityName}(ctx context.Context, ${entityName} *entity.${EntityName}) error {
	if err := s.repo.Create${EntityName}(${entityName}); err != nil {
		return think.ErrSystemSpace(err.Error())
	}
	return nil
}
func (s *Service) Update${EntityName}(ctx context.Context, id int64, version int, ${entityName} *entity.${EntityName}) error {
	if err := s.repo.Update${EntityName}(id, version, ${entityName}); err != nil {
		return think.ErrSystemSpace(err.Error())
	}
	return nil
}

func (s *Service) Update${EntityName}Map(ctx context.Context, id int64, version int, ${entityName} map[string]any) error {
	if err := s.repo.Update${EntityName}Map(id, version, ${entityName}); err != nil {
		return think.ErrSystemSpace(err.Error())
	}
	return nil
}

func (s *Service) Search${EntityName}(ctx context.Context, searchMeta contract.SearchMeta, searchParams contract.MysqlBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	data, cnt, err = s.repo.Search${EntityName}(searchMeta, searchParams)
	if err != nil {
		err = think.ErrSystemSpace(err.Error())
	}
	return data, cnt, err
}
`
var UseCaseServerMongo = `package ${useCasePkg}

import (
	"context"
	"github.com/zander-84/seagull/contract"
	"github.com/zander-84/seagull/think"
	"${project}/apps/${server}/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) UseCase {
	out := new(Service)
	out.repo = repo

	return out
}

func (s *Service) Get${EntityName}(ctx context.Context, id string) (*entity.${EntityName}, error) {
	out, err := s.repo.Get${EntityName}(id)
	if err != nil {
		if think.IsErrNotFound(err) {
			return nil, err
		}
		return nil, think.ErrSystemSpace(err.Error())
	}
	return out, err
}

func (s *Service) BatchGet${EntityName}(ctx context.Context, ids []string) ([]entity.${EntityName}, error) {
	out, err := s.repo.BatchGet${EntityName}(ids)
	if err != nil {
		return nil, think.ErrSystemSpace(err.Error())
	}
	return out, err
}

func (s *Service) Create${EntityName}(ctx context.Context, ${entityName} *entity.${EntityName}) error {
	if err := s.repo.Create${EntityName}(${entityName}); err != nil {
		return think.ErrSystemSpace(err.Error())
	}
	return nil
}
func (s *Service) Update${EntityName}(ctx context.Context, id string, version int, ${entityName} *entity.${EntityName}) error {
	if err := s.repo.Update${EntityName}(id, version, ${entityName}); err != nil {
		return think.ErrSystemSpace(err.Error())
	}
	return nil
}



func (s *Service) Search${EntityName}(ctx context.Context, searchMeta contract.SearchMeta, searchParams contract.MongoBuilder) (data []entity.${EntityName}, cnt int64, err error) {
	data, cnt, err = s.repo.Search${EntityName}(searchMeta, searchParams)
	if err != nil {
		err = think.ErrSystemSpace(err.Error())
	}
	return data, cnt, err
}`
