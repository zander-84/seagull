package project

var UsecaseTpl = `package usecase

import (
	"${project}/apps/${server}/internal/usecase/root"
)

var (

	RootUseCase root.UseCase
)

func InitUseCase() {


	RootUseCase = root.NewService()
}

`

var UsecaserootTpl = `package root

type Service struct {
}

func NewService() UseCase {
	out := new(Service)

	return out
}
`
var UsecaserootInterfaceTpl = `package root

type UseCase interface {
}`
