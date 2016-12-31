// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package controllers

import (
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
	"encoding/json"
)

type rowHistoryPage struct {
	Alert        string
	SignData     string
	ShowSignData bool
	CountSignArr []int
	Lang         map[string]string
	History    []map[string]string
	WalletId     int64
	CitizenId    int64
	TableName    string
	StateId      int64
	Global       string
	Columns               map[string]string
}

func (c *Controller) RowHistory() (string, error) {

	var history []map[string]string
	rbId := utils.StrToInt64(c.r.FormValue("rbId"))
	if rbId < 1 {
		return "", utils.ErrInfo(`Incorrect rbId`)
	}
	var tableName string
	if utils.CheckInputData(c.r.FormValue("tableName"), "string") {
		tableName = c.r.FormValue("tableName")
	}

	global := c.r.FormValue("global")
	prefix := c.StateIdStr
	if global == "1" {
		prefix = "global"
	} else {
		global = "0"
	}
	columns, err := c.GetMap(`SELECT data.* FROM "`+prefix+`_tables", jsonb_each_text(columns_and_permissions->'update') as data WHERE name = ?`, "key", "value", tableName)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	columns["id"] = ""
	columns["block_id"] = ""
	for i:=0; i<100; i++ {
		data, err := c.OneRow(`SELECT data, block_id FROM "rollback" WHERE rb_id = ?`, rbId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		var messageMap map[string]string
		json.Unmarshal([]byte(data["data"]), &messageMap)
		rbId = utils.StrToInt64(messageMap["rb_id"])
		messageMap["block_id"] = data["block_id"]
		history = append(history, messageMap)
		if rbId == 0 {
			break
		}
	}

	TemplateStr, err := makeTemplate("row_history", "rowHistory", &rowHistoryPage{
		Alert:        c.Alert,
		Lang:         c.Lang,
		ShowSignData: c.ShowSignData,
		SignData:     "",
		WalletId:     c.SessWalletId,
		History:         history,
		CitizenId:    c.SessCitizenId,
		CountSignArr: c.CountSignArr,
		Columns:      columns,
		StateId:      c.SessStateId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
