package etl

import (
	"fmt"
	"reflect"
	"sync"

	conversion "etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
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

const (
	OrderByRidDESC = "rid DESC"
	OrderByRidASC  = "rid ASC"
)

type ETLQuery struct {
	OrderBy string
	Where   []interface{}
	Limit   int
}

type ETLEngine struct {
	etlModelPairs     []*types.ModelPair
	customConversions []func(scheme *conversion.Scheme)
	convertScheme     *conversion.Scheme
	queries           []ETLQuery
}

func (ng *ETLEngine) Register(fromDB *cmsql.Database, fromModel types.Model, toDB *cmsql.Database, toModel types.Model) *ETLEngine {
	ng.etlModelPairs = append(ng.etlModelPairs, types.NewModelPair(fromDB, fromModel, toDB, toModel))
	return ng
}

func (ng *ETLEngine) RegisterConversion(funcz func(*conversion.Scheme)) *ETLEngine {
	ng.customConversions = append(ng.customConversions, funcz)
	return ng
}

func (ng *ETLEngine) RegisterQuery(query ETLQuery) *ETLEngine {
	ng.queries = append(ng.queries, query)
	return ng
}

func (ng *ETLEngine) Bootstrap() {
	ng.convertScheme = conversion.Build(append(ng.customConversions)...)
}

/*
 * Get (20) latest RIDs from Destination Database
 * Get (20) latest RIDs from Source Database
 * Compare two results (above), then choose the RIDs that appear in Source Database but not in Destination Database
 * Reply on result RIDs (above), query records (with limit) from lowest RID (source database) and RIDs (above)
 * Transform data
 * Load data
 */
func (ng *ETLEngine) Run() {
	ng.Bootstrap()

	var wg sync.WaitGroup

	for modelPairIndex, etlModelPair := range ng.etlModelPairs {
		wg.Add(1)
		go func(index int, modelPair *types.ModelPair) {
			defer wg.Done()

			for {
				plistSrc := reflect.New(reflect.TypeOf(modelPair.Source.Model).Elem())
				plistDst := reflect.New(reflect.TypeOf(modelPair.Target.Model).Elem())

				var rid dot.ID = 0
				var excludedRIDs []dot.ID
				mapDstRIDs := make(map[dot.ID]bool)
				mapSrcRIDs := make(map[dot.ID]bool)

				// Get (20) latest RIDs from Destination Database
				latestDstRIDs, err := ng.scanRIDs(modelPair.Target.DB, plistDst.Interface().(types.Model), ETLQuery{
					OrderBy: OrderByRidDESC,
					Where:   []interface{}{},
					Limit:   20,
				})
				must(err)
				for _, _rid := range latestDstRIDs {
					mapDstRIDs[_rid] = true
				}
				if len(latestDstRIDs) != 0 {
					rid = latestDstRIDs[0]
				}

				// Get (20) latest RIDs from Source Database
				latestSrcRIDs, err := ng.scanRIDs(modelPair.Source.DB, plistSrc.Interface().(types.Model), ETLQuery{
					OrderBy: OrderByRidDESC,
					Where:   append(ng.queries[index].Where, sq.NewExpr("rid <= ?", rid.Int64())),
					Limit:   20,
				})
				must(err)
				for _, _rid := range latestSrcRIDs {
					mapSrcRIDs[_rid] = true
				}

				// Compare two results (above), then choose the RIDs that appear in Source Database but not in Destination Database
				for i := len(latestSrcRIDs) - 1; i >= 0; i-- {
					if !mapDstRIDs[latestSrcRIDs[i]] {
						rid = latestSrcRIDs[i]
						break
					}
				}

				for i := len(latestDstRIDs) - 1; i >= 0; i-- {
					if mapSrcRIDs[latestDstRIDs[i]] {
						excludedRIDs = append(excludedRIDs, latestDstRIDs[i])
					}
				}

				// Reply on result RIDs (above), query records (with limit) from lowest RID (source database) and RIDs (above)
				var condition sq.WriterTo
				if len(excludedRIDs) > 0 {
					condition = sq.NotIn("rid", excludedRIDs)
				} else {
					condition = sq.NewExpr("TRUE")
				}
				err = ng.scanModels(modelPair.Source.DB, plistSrc.Interface().(types.Model), ETLQuery{
					OrderBy: ng.queries[index].OrderBy,
					Where:   append(ng.queries[index].Where, sq.NewExpr("rid >= ?", rid.Int64()), condition),
					Limit:   ng.queries[index].Limit,
				})
				must(err)

				if plistSrc.Elem().Len() == 0 {
					break
				}

				// Transform data
				err = ng.transform(plistSrc.Elem().Interface(), plistDst.Interface())
				must(err)

				// Load data
				err = ng.loadModels(modelPair.Target.DB, plistDst.Interface())
				must(err)
			}

			return
		}(modelPairIndex, etlModelPair)
	}

	wg.Wait()
}

func (ng *ETLEngine) RunEveryDay() {
	ng.Bootstrap()

	var wg sync.WaitGroup

	for modelPairIndex, etlModelPair := range ng.etlModelPairs {
		go func(index int, modelPair *types.ModelPair) {
			wg.Add(1)
			defer wg.Done()

			var rid dot.ID = 0

			for {
				plistSrc := reflect.New(reflect.TypeOf(modelPair.Source.Model).Elem())
				plistDst := reflect.New(reflect.TypeOf(modelPair.Target.Model).Elem())

				err := ng.scanModels(modelPair.Source.DB, plistSrc.Interface().(types.Model), ETLQuery{
					OrderBy: OrderByRidASC,
					Where:   append(ng.queries[index].Where, sq.NewExpr("rid > ?", rid.Int64())),
					Limit:   1000,
				})
				must(err)

				valSrc := reflect.ValueOf(plistSrc.Interface()).Elem()
				if valSrc.Len() == 0 {
					break
				}

				mapRidSrc := make(map[dot.ID]reflect.Value)
				ridsSrc := make([]dot.ID, 0, valSrc.Len())
				var maxRid dot.ID
				for i := 0; i < valSrc.Len(); i++ {
					_rid := dot.ID(valSrc.Index(i).Elem().FieldByName("Rid").Int())
					ridsSrc = append(ridsSrc, _rid)
					mapRidSrc[_rid] = valSrc.Index(i).Elem()
					if maxRid < _rid {
						maxRid = _rid
					}
				}
				rid = maxRid

				err = ng.scanModels(modelPair.Target.DB, plistDst.Interface().(types.Model), ETLQuery{
					OrderBy: OrderByRidASC,
					Where:   []interface{}{sq.In("rid", ridsSrc)},
					Limit:   1000,
				})
				must(err)

				valDst := reflect.ValueOf(plistDst.Interface()).Elem()
				mapRidDst := make(map[dot.ID]reflect.Value)
				ridsDst := make([]dot.ID, 0, valDst.Len())
				for i := 0; i < valDst.Len(); i++ {
					_rid := valDst.Index(i).Elem().FieldByName("Rid")
					ridsDst = append(ridsDst, _rid.Interface().(dot.ID))
					mapRidDst[_rid.Interface().(dot.ID)] = valDst.Index(i).Elem()
				}

				var rids, deletedRids []dot.ID
				for k, v := range mapRidSrc {
					if _, ok := mapRidDst[k]; !ok {
						rids = append(rids, k)
						continue
					}

					if !isEqual(v, mapRidDst[k]) {
						rids = append(rids, k)
						deletedRids = append(deletedRids, k)
					}
				}

				if len(deletedRids) != 0 {
					err = deleteModelsByRids(modelPair.Target.DB, plistDst.Interface().(types.Model), deletedRids)
					must(err)
				}

				if len(rids) == 0 {
					break
				}

				err = ng.scanModels(modelPair.Source.DB, plistSrc.Interface().(types.Model), ETLQuery{
					OrderBy: OrderByRidASC,
					Where:   append(ng.queries[index].Where, sq.In("rid", rids)),
					Limit:   1000,
				})
				must(err)

				err = ng.transform(plistSrc.Elem().Interface(), plistDst.Interface())
				must(err)

				err = ng.loadModels(modelPair.Target.DB, plistDst.Interface())
				must(err)
			}
		}(modelPairIndex, etlModelPair)
	}

	wg.Wait()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (ng *ETLEngine) scanModels(db *cmsql.Database, plistSrc core.IFind, query ETLQuery) error {
	err := db.
		OrderBy(query.OrderBy).
		Where(query.Where).
		Limit(uint64(query.Limit)).
		Find(plistSrc)
	return err
}

func (ng *ETLEngine) scanRID(db *cmsql.Database, tableName string, query ETLQuery) (rid dot.ID, err error) {
	err = db.
		Select("rid").
		From(tableName).
		OrderBy(query.OrderBy).
		Where(query.Where).
		Limit(uint64(query.Limit)).
		Scan(&rid)
	return
}

func (ng *ETLEngine) scanRIDs(db *cmsql.Database, plistSrc core.IFind, query ETLQuery) (rids []dot.ID, err error) {
	var id dot.ID

	rows, err := db.NewQuery().
		Select("rid").
		From(plistSrc.SQLTableName()).
		Where(query.Where).
		Limit(uint64(query.Limit)).
		OrderBy(query.OrderBy).
		Query()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		rids = append(rids, id)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return
}

func deleteModelsByRids(db *cmsql.Database, plistSrc core.IFind, rids []dot.ID) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM %q WHERE rid in (%s)", plistSrc.SQLTableName(), convertIntArrayToString(rids)))
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

func convertIntArrayToString(args []dot.ID) (ids string) {
	for _, id := range args {
		ids += fmt.Sprintf("%d,", id.Int64())
	}
	ids = ids[:len(ids)-1]
	return
}

func isEqual(obj, otherObj reflect.Value) bool {
	objType := obj.Type()
	objValue := obj
	mapObjFieldIndex := make(map[string]int, objType.NumField())
	for i := 0; i < objType.NumField(); i++ {
		mapObjFieldIndex[objType.Field(i).Name] = i
	}

	otherObjType := otherObj.Type()
	otherObjValue := otherObj
	mapOtherObjectField := make(map[string]int, otherObjType.NumField())
	for i := 0; i < otherObjType.NumField(); i++ {
		mapOtherObjectField[otherObjType.Field(i).Name] = i
	}

	var differentFields []string
	for otherObjectField := range mapOtherObjectField {
		if _, ok := mapObjFieldIndex[otherObjectField]; ok {
			differentFields = append(differentFields, otherObjectField)
		}
	}

	for _, field := range differentFields {
		objField := objValue.Field(mapObjFieldIndex[field])
		otherObjectField := otherObjValue.Field(mapOtherObjectField[field])

		if !reflect.DeepEqual(objField.Interface(), otherObjectField.Interface()) {
			return false
		}
	}

	return true
}
