package main

import (
	"sync"

	"o.o/backend/cmd/fabo-server/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	cfg config.Config
	ll  = l.New()
	db  *cmsql.Database

	wg sync.WaitGroup
	mu sync.Mutex
)

func main() {
	var err error
	if cfg, err = config.Load(); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	postgres := cfg.Databases.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connection database", l.Error(err))
	}

	mapCreatedByAndFfmIDs, err := scanAllCreatedByAndFfmIDs()
	if err != nil {
		ll.S.Infof("Error scanAllCreatedByAndFfmIDs")
	}

	count, updatedCount, errCount := len(mapCreatedByAndFfmIDs), 0, 0

	for createdBy, ffmIDs := range mapCreatedByAndFfmIDs {
		go func(_createdBy dot.ID, _ffmIDs []dot.ID) {
			wg.Add(1)

			update := map[string]interface{}{
				"created_by": _createdBy,
			}
			err := db.Table("fulfillment").
				In("id", _ffmIDs).
				ShouldUpdateMap(update)

			mu.Lock()
			if err != nil {
				errCount += 1
			} else {
				updatedCount += 1
			}
			mu.Unlock()

			wg.Done()
		}(createdBy, ffmIDs)
	}

	wg.Wait()

	ll.S.Infof("Done: updated %v/%v", updatedCount, count)
	ll.S.Infof("Error %v/%v", errCount, count)
}

func scanAllCreatedByAndFfmIDs() (mapCreatedByAndFfmIDs map[dot.ID][]dot.ID, _ error) {
	var (
		fromID           dot.ID = 0
		hasResult        bool
		ffmID, createdBy dot.ID
	)

	mapCreatedByAndMapFfmID := make(map[dot.ID]map[dot.ID]bool)
	mapCreatedByAndFfmIDs = make(map[dot.ID][]dot.ID)

	for {
		hasResult = false
		rows, err := db.
			SQL(`SELECT f.id, o.created_by FROM fulfillment f JOIN "order" o ON f.order_id = o.id`).
			Where("f.created_by IS NULL").
			Where("f.order_id IS NOT NULL").
			Where("o.created_by IS NOT NULL").
			Where("f.id > ?", fromID.Int64()).
			OrderBy("f.id").
			Limit(1000).
			Query()
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			hasResult = true
			if err := rows.Scan(&ffmID, &createdBy); err != nil {
				return nil, err
			}

			if _, ok := mapCreatedByAndMapFfmID[createdBy]; !ok {
				mapCreatedByAndMapFfmID[createdBy] = make(map[dot.ID]bool)
			}
			mapCreatedByAndMapFfmID[createdBy][ffmID] = true
			fromID = ffmID
		}

		if err := rows.Close(); err != nil {
			return nil, err
		}

		if !hasResult {
			break
		}
	}

	for createdBy, mapFfmID := range mapCreatedByAndMapFfmID {
		ffmIDs := make([]dot.ID, 0, len(mapFfmID))
		for ffmID := range mapFfmID {
			ffmIDs = append(ffmIDs, ffmID)
		}
		mapCreatedByAndFfmIDs[createdBy] = ffmIDs
	}

	return mapCreatedByAndFfmIDs, nil
}
