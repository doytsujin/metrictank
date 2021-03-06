package expr

import (
	"fmt"

	"github.com/grafana/metrictank/api/models"
	"github.com/grafana/metrictank/consolidation"
)

type FuncConsolidateBy struct {
	in GraphiteFunc
	by string
}

func NewConsolidateBy() GraphiteFunc {
	return &FuncConsolidateBy{}
}

func NewConsolidateByConstructor(by string) func() GraphiteFunc {
	return func() GraphiteFunc {
		return &FuncConsolidateBy{by: by}
	}
}

func (s *FuncConsolidateBy) Signature() ([]Arg, []Arg) {
	if s.by != "" {
		return []Arg{
			ArgSeriesList{val: &s.in},
		}, []Arg{ArgSeriesList{}}
	}
	return []Arg{
		ArgSeriesList{val: &s.in},
		ArgString{val: &s.by, validator: []Validator{IsConsolFunc}},
	}, []Arg{ArgSeriesList{}}
}

func (s *FuncConsolidateBy) Context(context Context) Context {
	context.consol = consolidation.FromConsolidateBy(s.by)
	return context
}

func (s *FuncConsolidateBy) Exec(cache map[Req][]models.Series) ([]models.Series, error) {
	series, err := s.in.Exec(cache)
	if err != nil {
		return nil, err
	}
	consolidator := consolidation.FromConsolidateBy(s.by)
	for i, serie := range series {
		series[i].Target = fmt.Sprintf("consolidateBy(%s,\"%s\")", serie.Target, s.by)
		series[i].QueryPatt = fmt.Sprintf("consolidateBy(%s,\"%s\")", serie.QueryPatt, s.by)
		series[i].Consolidator = consolidator
		series[i].QueryCons = consolidator
	}
	return series, nil
}
