package main

import (
	"log"

	"gorm.io/gorm"
)

func DumpFloors(indexName string) {
	var holes Holes
	var floors Items
	result := DB.Where("hidden = false").
		FindInBatches(&holes, 1000, func(tx *gorm.DB, batch int) error {
			if len(holes) == 0 {
				return nil
			}
			holeIDs := make([]int, len(holes))
			for i, hole := range holes {
				holeIDs[i] = hole.ID
			}

			err := tx.
				Table("floor").
				Select("id", "content", "updated_at").
				Where("hole_id in (?) and deleted = 0 and ((is_actual_sensitive IS NOT NULL AND is_actual_sensitive = false) OR (is_actual_sensitive IS NULL AND is_sensitive = false))", holeIDs).
				Scan(&floors).Error
			if err != nil {
				return err
			}
			if len(floors) == 0 {
				return nil
			}

			err = BulkInsert(floors, indexName)
			if err != nil {
				return err
			}

			log.Printf("insert holes [%d, %d]\n", holes[0].ID, holes[len(holes)-1].ID)
			return nil
		})

	if result.Error != nil {
		log.Fatalf("dump err: %s", result.Error)
	}
}

func DumpTag() {
	var items Items
	err := DB.Table("tag").Select("id", "name", "updated_at").Scan(&items).Error
	if err != nil {
		log.Fatalf("dump err: %s", err)
		return
	}

	if len(items) == 0 {
		return
	}

	err = BulkInsert(items, IndexNameTag)
	if err != nil {
		log.Fatalf("dump err: %s", err)
	}
}

func DumpProject() {
	var holes Holes
	var projects Items
	var hole_projects []HoleProject
	result := DB.Where("hidden = false and approved = ?", true).
		FindInBatches(&holes, 1000, func(tx *gorm.DB, batch int) error {
			if len(holes) == 0 {
				return nil
			}
			holeIDs := make([]int, len(holes))
			for i, hole := range holes {
				holeIDs[i] = hole.ID
			}

			err := tx.Table("hole_projects").Where("hole_id in (?)", holeIDs).Scan(&hole_projects).Error
			if err != nil {
				return nil
			}

			projectIDs := make([]int, len(hole_projects))
			for i, hole_project := range hole_projects {
				projectIDs[i] = hole_project.ProjectID
			}

			err = tx.
				Table("project").
				Select("id", "CONCAT(content, description)", "updated_at").
				Where("id in (?)", projectIDs).
				Scan(&projects).Error
			if err != nil {
				return err
			}
			if len(projects) == 0 {
				return nil
			}

			err = BulkInsert(projects, IndexNameProject)
			if err != nil {
				return err
			}

			log.Printf("insert project holes, hole_id in [%d, %d]\n", holes[0].ID, holes[len(holes)-1].ID)
			return nil
		})

	if result.Error != nil {
		log.Fatalf("dump err: %s", result.Error)
	}
}
