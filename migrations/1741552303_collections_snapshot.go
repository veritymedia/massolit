package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2024-10-15 17:10:08.833Z",
				"updated": "2024-10-15 17:19:16.899Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [
								"image/jpg",
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					}
				],
				"indexes": [
					"CREATE INDEX ` + "`" + `__pb_users_auth__created_idx` + "`" + ` ON ` + "`" + `users` + "`" + ` (` + "`" + `created` + "`" + `)"
				],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "5k0uz7zn0m27i18",
				"created": "2024-10-16 17:43:29.048Z",
				"updated": "2024-12-14 17:02:58.230Z",
				"name": "books",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "gudbpp3u",
						"name": "title",
						"type": "text",
						"required": true,
						"presentable": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "cacwcquc",
						"name": "isbn",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": true
						}
					},
					{
						"system": false,
						"id": "88eznzr5",
						"name": "cover_url",
						"type": "url",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": [],
							"onlyDomains": []
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id != \"\"",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id != \"\"",
				"deleteRule": "@request.auth.id != \"\"",
				"options": {}
			},
			{
				"id": "zaenzsgxcbsif1x",
				"created": "2024-10-16 19:58:10.021Z",
				"updated": "2024-12-14 17:02:50.782Z",
				"name": "book_instances",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "r1q51d8l",
						"name": "book",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "5k0uz7zn0m27i18",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "szla0ipq",
						"name": "book_code",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id != \"\"",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id != \"\"",
				"deleteRule": "@request.auth.id != \"\"",
				"options": {}
			},
			{
				"id": "410vkrq314e2vl2",
				"created": "2024-10-16 20:00:24.947Z",
				"updated": "2024-12-14 17:03:39.590Z",
				"name": "rentals",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "kjp6dnoo",
						"name": "rented_to",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "rptqzw5t",
						"name": "book_instance",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "zaenzsgxcbsif1x",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id != \"\"",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id != \"\"",
				"deleteRule": "@request.auth.id != \"\"",
				"options": {}
			},
			{
				"id": "ob9ov162jg5mzye",
				"created": "2024-10-17 18:27:22.960Z",
				"updated": "2024-12-14 17:03:20.440Z",
				"name": "keys",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "1nyglzuu",
						"name": "app",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "glswqc03",
						"name": "key",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "zj818hrg1da8kgo",
				"created": "2024-12-11 19:31:38.619Z",
				"updated": "2024-12-14 17:08:01.123Z",
				"name": "behavior_notes",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "nmmtz01y",
						"name": "student_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "it6ludbd",
						"name": "first_name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "78rt4xk4",
						"name": "last_name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vaeokh6u",
						"name": "email",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "q0ds7lst",
						"name": "grade",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "sxrlmxtt",
						"name": "incident_time",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "hzjyovjc",
						"name": "behavior_type",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "jxrqx28m",
						"name": "notes",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "3xoknpyh",
						"name": "next_step",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "9s72zjbr",
						"name": "next_step_date",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "5ysqxzwt",
						"name": "author_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "fwtqckox",
						"name": "reported_by",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "qte02rx2",
						"name": "homeroom_advisor",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "mojdnvyg",
						"name": "visible_to_parents",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "7xt4qgkv",
						"name": "visible_to_students",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "flmhmzj9",
						"name": "created_at",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "cio3mmih",
						"name": "updated_at",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "adqiffxs",
						"name": "action_complete",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id != \"\"",
				"viewRule": "@request.auth.id != \"\"",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id != \"\"",
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "pu84v7pici1wp11",
				"created": "2024-12-11 19:34:12.006Z",
				"updated": "2024-12-12 22:44:51.333Z",
				"name": "config",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ctihfdr4",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vpwc94qg",
						"name": "value",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "02ats64dzd7u0ke",
				"created": "2024-12-13 19:49:39.468Z",
				"updated": "2024-12-14 17:03:27.683Z",
				"name": "mail_list",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "tgol0tzq",
						"name": "email",
						"type": "email",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "ggwqu9of",
						"name": "subs",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"behavior"
							]
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
