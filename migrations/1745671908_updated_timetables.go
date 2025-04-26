package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("tjuajguccuqmony")
		if err != nil {
			return err
		}

		// update
		edit_exam_timetables := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "mfyrawfn",
			"name": "exam_timetables",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "u4dlxvkuasevu69",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": null
			}
		}`), edit_exam_timetables); err != nil {
			return err
		}
		collection.Schema.AddField(edit_exam_timetables)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("tjuajguccuqmony")
		if err != nil {
			return err
		}

		// update
		edit_exam_timetables := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "mfyrawfn",
			"name": "exam_timetable",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "u4dlxvkuasevu69",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_exam_timetables); err != nil {
			return err
		}
		collection.Schema.AddField(edit_exam_timetables)

		return dao.SaveCollection(collection)
	})
}
