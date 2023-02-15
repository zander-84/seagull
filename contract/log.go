package contract

import "github.com/zander-84/seagull/contract/def"

var DefaultLogElement = "msg"

// LogM default msg key
func LogM(msg string) def.E {
	return def.E{Key: DefaultLogElement, Value: msg}
}

func LogE(key string, val any) def.E { return def.E{Key: key, Value: val} }

type Logger interface {
	Debug(es ...def.E)
	Info(es ...def.E)
	Error(es ...def.E)
	Panic(es ...def.E)
	Fatal(es ...def.E)
}
