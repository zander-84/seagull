package zap

import (
	"github.com/zander-84/seagull/contract/def"
	"go.uber.org/zap"
)

func (z *zLog) Debug(es ...def.E) {
	z.engine.Debug("", z.elements(es)...)
}

func (z *zLog) Info(es ...def.E) {
	z.engine.Info("", z.elements(es)...)
}

func (z *zLog) Error(es ...def.E) {
	z.engine.Error("", z.elements(es)...)
}

func (z *zLog) Panic(es ...def.E) {
	z.engine.Panic("", z.elements(es)...)
}

func (z *zLog) Fatal(es ...def.E) {
	z.engine.Fatal("", z.elements(es)...)
}

func (z *zLog) elements(es []def.E) []zap.Field {
	if len(es) < 1 {
		return nil
	}

	var fields = make([]zap.Field, 0)
	for _, v := range es {
		fields = append(fields, zap.Any(v.Key, v.Value))
	}
	return fields
}
