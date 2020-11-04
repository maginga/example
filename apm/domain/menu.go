package domain

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CreateMenu(tenantID string) error {
	url := fmt.Sprintf("%v", viper.Get("metadata.grandview-url"))
	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
	}
	defer tx.Rollback()

	contents := `[
		{
		  "name": "dashboard",
		  "permissions": [],
		  "subMenus": [
			{
			  "name": "dashboard.predictive-status",
			  "permissions": [],
			  "translate": "MENU.Predictive Status",
			  "type": "MENU_LINK",
			  "url": "/dashboard/map"
			},
			{
			  "name": "dashboard.asset-summary",
			  "permissions": [],
			  "translate": "MENU.Asset Summary",
			  "type": "MENU_LINK",
			  "url": "/dashboard/overview"
			},
			{
			  "name": "dashboard.alarm-status",
			  "permissions": [],
			  "translate": "MENU.Alarm Status",
			  "type": "MENU_LINK",
			  "url": "/dashboard/alarm"
			},
			{
			  "name": "dashboard.sensor-status",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "translate": "MENU.Sensor Status",
			  "type": "MENU_LINK",
			  "url": "/dashboard/sensor"
			},
			{
			  "name": "dashboard.health-status",
			  "permissions": [],
			  "translate": "MENU.Health Status",
			  "type": "MENU_LINK",
			  "url": "/dashboard/health"
			}
		  ],
		  "translate": "MENU.Dashboard",
		  "type": "MENU_VIEW"
		},
		{
		  "name": "exploration",
		  "permissions": [],
		  "subMenus": [
			{
			  "name": "exploration.asset",
			  "permissions": [],
			  "subMenus": [
				{
				  "name": "exploration.asset.overview",
				  "permissions": [],
				  "translate": "MENU.Overview",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/overview"
				},
				{
				  "name": "exploration.asset.trend",
				  "permissions": [],
				  "translate": "MENU.Parameter",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/parameter"
				},
				{
				  "name": "exploration.asset.alarm",
				  "permissions": [],
				  "translate": "MENU.Alarm",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/alarm"
				},
				{
				  "name": "exploration.asset.model",
				  "permissions": [
					"PERM_SYSTEM_MANAGE_APP",
					"PERM_SYSTEM_MANAGE_ASSET"
				  ],
				  "translate": "MENU.Model",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/model"
				},
				{
				  "name": "exploration.asset.work-order",
				  "permissions": [],
				  "translate": "MENU.WorkOrder",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/workorder"
				},
				{
				  "name": "exploration.asset.information",
				  "permissions": [],
				  "translate": "MENU.Information",
				  "type": "MENU_TAB",
				  "url": "/asset/detail/information"
				}
			  ],
			  "translate": "MENU.Asset",
			  "type": "MENU_LINK",
			  "url": "/asset/list"
			},
			{
			  "name": "exploration.alarm",
			  "permissions": [],
			  "translate": "MENU.Alarm",
			  "type": "MENU_LINK",
			  "url": "/asset/alarm"
			},
			{
			  "name": "exploration.work-order",
			  "permissions": [],
			  "redirectUrl": "/asset/workorder/history",
			  "subMenus": [
				{
				  "name": "exploration.work-order.history",
				  "permissions": [],
				  "translate": "MENU.WorkOrderHistory",
				  "type": "MENU_TAB",
				  "url": "/asset/workorder/history"
				},
				{
				  "name": "exploration.work-order.visit",
				  "permissions": [],
				  "translate": "MENU.VisitHistory",
				  "type": "MENU_TAB",
				  "url": "/asset/workorder/visit"
				},
				{
				  "name": "exploration.work-order.repair-cost",
				  "permissions": [],
				  "translate": "MENU.RepairCost",
				  "type": "MENU_TAB",
				  "url": "/asset/workorder/repaircost"
				}
			  ],
			  "translate": "MENU.WorkOrder",
			  "type": "MENU_LINK"
			}
		  ],
		  "translate": "MENU.Exploration",
		  "type": "MENU_VIEW"
		},
		{
		  "name": "asset analysis",
		  "permissions": [
			"PERM_SYSTEM_MANAGE_APP"
		  ],
		  "subMenus": [
			{
			  "name": "asset analysis.asset matching",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP"
			  ],
			  "translate": "ASSET MATCHING.Asset Matching",
			  "type": "MENU_LINK",
			  "url": "/asset/matching"
			}
		  ],
		  "translate": "MENU.Asset Analysis",
		  "type": "MENU_VIEW"
		},
		{
		  "name": "management",
		  "permissions": [
			"PERM_SYSTEM_MANAGE_APP",
			"PERM_SYSTEM_MANAGE_ASSET"
		  ],
		  "subMenus": [
			{
			  "name": "management.asset-template",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "translate": "MENU.Asset Template",
			  "type": "MENU_LINK",
			  "url": "/management/template/list"
			},
			{
			  "name": "management.sensor-management",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "translate": "MENU.Sensor",
			  "type": "MENU_LINK",
			  "url": "/management/sensor/list"
			},
			{
			  "name": "management.parameter-template",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "subMenus": [
				{
				  "name": "management.parameter-template.parameter",
				  "permissions": [
					"PERM_SYSTEM_MANAGE_APP",
					"PERM_SYSTEM_MANAGE_ASSET"
				  ],
				  "translate": "MANAGE.Parameter",
				  "type": "MENU_TAB",
				  "url": "/management/parameter/detail/parameter"
				}
			  ],
			  "translate": "MENU.Parameter Template",
			  "type": "MENU_LINK",
			  "url": "/management/parameter/list"
			},
			{
			  "name": "management.location",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "translate": "MENU.Asset Hierarchy",
			  "type": "MENU_LINK",
			  "url": "/management/hierarchy"
			},
			{
			  "name": "management.analysis-model",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP",
				"PERM_SYSTEM_MANAGE_ASSET"
			  ],
			  "translate": "MENU.Analysis Model",
			  "type": "MENU_LINK",
			  "url": "/management/model/list"
			}
		  ],
		  "translate": "MENU.Asset Management",
		  "type": "MENU_VIEW"
		},
		{
		  "name": "admin",
		  "permissions": [
			"PERM_SYSTEM_MANAGE_APP"
		  ],
		  "subMenus": [
			{
			  "name": "admin.authority",
			  "permissions": [
				"PERM_SYSTEM_MANAGE_APP"
			  ],
			  "translate": "MENU.Authority",
			  "type": "MENU_LINK",
			  "url": "/admin/group/list"
			}
		  ],
		  "translate": "MENU.System Management",
		  "type": "MENU_VIEW"
		}
	  ]`

	uid := uuid.New().String()
	stmt := "INSERT INTO menu " +
		"(id, version, contents, tenant_id, created_by, created_time, modified_by, modified_time) VALUES " +
		"(?,?,?,?,'admin',NOW(),'admin',NOW()) "

	_, err = tx.Exec(stmt, uid, 0, contents, tenantID)

	if err != nil {
		log.Panic(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	return err
}
