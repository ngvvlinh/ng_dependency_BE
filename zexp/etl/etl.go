package etl

import (
	"context"
	"reflect"

	conversion "etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/backend/zexp/etl/types"
	"etop.vn/capi/dot"
)

func NewETLEngine(etlModelPairs []*types.ModelPair) *ETLEngine {
	if etlModelPairs == nil {
		etlModelPairs = []*types.ModelPair{}
	}
	return &ETLEngine{
		etlModelPairs:     etlModelPairs,
		customConversions: []func(scheme *conversion.Scheme){},
	}
}

type ETLEngine struct {
	etlModelPairs     []*types.ModelPair
	customConversions []func(scheme *conversion.Scheme)
	convertScheme     *conversion.Scheme
}

func (ng *ETLEngine) Register(fromDB *cmsql.Database, fromModel types.Model, toDB *cmsql.Database, toModel types.Model) *ETLEngine {
	ng.etlModelPairs = append(ng.etlModelPairs, types.NewModelPair(fromDB, fromModel, toDB, toModel))
	return ng
}

func (ng *ETLEngine) RegisterConversion(funcz func(*conversion.Scheme)) *ETLEngine {
	ng.customConversions = append(ng.customConversions, funcz)
	return ng
}

func (ng *ETLEngine) Bootstrap() {
	ng.convertScheme = conversion.Build(append(ng.customConversions)...)
}

func (ng *ETLEngine) Run(context context.Context) {
	ng.Bootstrap()
	for _, modelPair := range ng.etlModelPairs {
		plistSrc := reflect.New(reflect.TypeOf(modelPair.Source.Model).Elem())
		plistDst := reflect.New(reflect.TypeOf(modelPair.Target.Model).Elem())

		err := ng.scanModels(modelPair.Source.DB, plistSrc.Interface().(types.Model), 0, 2)
		must(err)

		err = ng.transform(plistSrc.Elem().Interface(), plistDst.Interface())
		must(err)

		err = ng.loadModels(modelPair.Target.DB, plistDst.Interface())
		must(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (ng *ETLEngine) scanModels(db *cmsql.Database, plistSrc core.IFind, fromID dot.ID, limit int) error {
	err := db.
		OrderBy("id").
		Where("id >= ?", fromID).
		Limit(uint64(limit)).
		Find(plistSrc)
	return err
}

func (ng *ETLEngine) loadModels(db *cmsql.Database, data interface{}) error {
	_data := reflect.ValueOf(data)
	return db.ShouldUpsert(_data.Interface().(core.IUpsert))
}

func (ng *ETLEngine) transform(src, dst interface{}) error {
	err := ng.convertScheme.Convert(unwrap(src), unwrapPtr(dst))
	return err
}
