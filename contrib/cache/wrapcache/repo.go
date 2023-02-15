package wrapcache

import (
	"errors"
	"fmt"
	"github.com/zander-84/seagull/contract/def"
	"golang.org/x/exp/constraints"
)

type repository interface {
	CachePrefix() string
	Name() string
	Version() int
}

func RepoKey[T constraints.Integer | string](r repository, in T) def.K {
	return def.K{
		Key:   fmt.Sprintf("%s:%s:%d,%v", r.CachePrefix(), r.Name(), r.Version(), in),
		Alias: []any{in},
	}
}

func RepoKeys[T constraints.Integer | string](r repository, in []T) []def.K {
	out := make([]def.K, 0, len(in))

	for _, v := range in {
		out = append(out, RepoKey(r, v))
	}
	return out
}

func RepoDBKeys(in []def.K) ([]any, error) {
	ids := make([]any, 0, len(in))
	for _, v := range in {
		if len(v.Alias) != 1 {
			return nil, errors.New("lost id")
		}
		ids = append(ids, v.Alias[0])
	}
	return ids, nil
}
