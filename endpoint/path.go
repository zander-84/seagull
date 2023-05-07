package endpoint

import (
	"github.com/zander-84/seagull/transport"
	"strings"
)

type Path struct {
	path       string
	fullPath   string
	serverName string
	method     transport.Method
}

func NewPath(path string, method transport.Method) Path {
	p := Path{}
	p.path = path
	p.method = method
	p.fullPath = p.FullPath()
	p.serverName = p.ServerName()
	return p
}

func (p Path) FullPath() string {
	if p.fullPath != "" {
		return p.fullPath
	}
	return "/" + p.ServerName() + "/" + p.Method().String()
}
func (p Path) Path() string {
	return p.path
}
func (p Path) ServerName() string {
	if p.serverName != "" {
		return p.serverName
	}
	return strings.ReplaceAll(strings.Trim(p.path, "/"), "/", ".")
}

func (p Path) Method() transport.Method {
	return p.method
}
