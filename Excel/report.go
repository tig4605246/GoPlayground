package doc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"dforcepro.com/api"
	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

type ReportAPI bool

func (ra ReportAPI) Enable() bool {
	return bool(ra)
}

func (ra ReportAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/v1/report/{Start}/{End}", Next: ra.getEndPoint, Method: "GET", Auth: true},
	}
}

type ReportDoc struct {
	ID           bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	AssignNumber string        `json:"assign_number,omitempty" bson:"assign_number,omitempty"`
	Branch       string        `json:"branch,omitempty" bson:"branch,omitempty"`
	Maintenance  string        `json:"maintenance,omitempty" bson:"maintenance,omitempty"`
	AssignGuy    string        `json:"assign_guy,omitempty" bson:"assign_guy,omitempty"`
	BuildTime    time.Time     `json:"build_time,omitempty" bson:"build_time,omitempty"`
	Equipment    string        `json:"equipment,omitempty" bson:"equipment,omitempty"`
	FinishTime   time.Time     `json:"finished_time,omitempty" bson:"finished_time,omitempty"`
	FixTime      string        `json:"fix_time,omitempty" bson:"fix_time,omitempty"`
	Status       string        `json:"status,omitempty" bson:"status,omitempty"`
}

type BillDoc struct {
	ID          bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Type        string        `json:"type,omitempty" bson:"type,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
}

type BranchDoc struct {
	ID    bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Block string        `json:"block,omitempty" bson:"block,omitempty"`
}

type InputDoc struct {
	Start string `json:"start,omitempty" bson:"start,omitempty"`
	End   string `json:"end,omitempty" bson:"end,omitempty"`
}

const (
	panaBranch     = "SalesOffice"
	panaBill       = "Bill"
	panaAssignment = "Job"
	timeLayout     = "2006-01-02"
)

func _NewInputDoc() *InputDoc {
	return &InputDoc{}
}

func (wf ReportAPI) getEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)

	queries := req.URL.Query()
	inputForm := _NewInputDoc()
	st, ok1 := queries["Start"]
	et, ok2 := queries["End"]

	var queryStr = bson.M{}
	var reportResult []ReportDoc
	var reportResultUnfinish []ReportDoc
	var billResult BillDoc
	var branchResult BranchDoc
	mongo := getMongo()

	//Variables for excel

	var file *xlsx.File
	var sheetFinish *xlsx.Sheet
	var sheetUnfinish *xlsx.Sheet

	//Check existence of the input parameters

	if !ok1 {

		_di.Log.Err("Start time error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else {
		inputForm.Start = st[0]
	}

	if !ok2 {

		_di.Log.Err("End time error")
		w.WriteHeader(http.StatusBadRequest)
		return

	} else {
		inputForm.End = et[0]

	}

	//Format datetime

	startTime, err := time.Parse(timeLayout, inputForm.Start)
	if err != nil {

		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	endTime, err := time.Parse(timeLayout, inputForm.End)
	if err != nil {

		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	//Check if startTime >= endTime

	if !startTime.Before(endTime) {

		_di.Log.Err("Start is smaller than End")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Set up query command

	queryStr["build_time"] = bson.M{"$gte": startTime, "$lt": endTime}
	queryStr["status"] = "已完工"

	//Issue query command
	//If we find cache excel file first can speed up the response time

	query := mongo.DB(lzw).C(panaAssignment).Find(queryStr)
	total, err := query.Count()

	var limit, page = 200, 1

	//Debug print out

	_di.Log.Debug(strconv.Itoa(total))

	if err != nil {

		//Return nothing when error occurs

		_di.Log.Err(err.Error())

	} else if total != 0 {

		//Create Excel file

		file = xlsx.NewFile()

		//
		// Create sheet for DailyRepair
		//

		sheetFinish, err = file.AddSheet("DailyRepair")
		if err != nil {
			fmt.Printf(err.Error())
			_di.Log.Err(err.Error())
		}

		AddRowToSheet(sheetFinish, "代碼", "店舖名稱", "保修商", "叫修人", "叫修日期", "叫修時間", "叫修設備", "完畢日期", "完畢時間", "維修方式", "維修回報內容")

		//Query limitation: 200
		//This should be revised for better performance and agility
		query.Limit(limit).Skip(page - 1).All(&reportResult)

		for _, element := range reportResult {

			//Query for type and description

			_ = mongo.DB(lzw).C(panaBill).Find(bson.M{"assign_number": element.AssignNumber}).One(&billResult)

			_ = mongo.DB(lzw).C(panaBranch).Find(bson.M{"branch": element.Branch}).One(&branchResult)

			AddRowToSheet(sheetFinish, branchResult.Block, element.Branch, element.Maintenance,
				element.AssignGuy, element.BuildTime.String()[:10], element.BuildTime.String()[11:16], element.Equipment, element.FinishTime.String()[:10],
				element.FinishTime.String()[11:16], billResult.Type, billResult.Description)

		}

		//
		// Create sheet for Fail
		//

		//Revise query command

		queryStr["status"] = "未回報"

		//Issue query command

		query := mongo.DB(lzw).C(panaAssignment).Find(queryStr)

		total, _ = query.Count()

		//Total number of data

		sheetUnfinish, err = file.AddSheet("Fail")
		if err != nil {
			_di.Log.Err(err.Error())
		}

		AddRowToSheet(sheetUnfinish, "代碼", "店舖名稱", "保修商", "叫修人", "叫修日期", "叫修時間", "預計完修日期", "預計完修時間", "筆數")

		if total > 0 {

			//Query limitation: 200
			//This should be revised for better performance and agility
			query.Limit(limit).Skip(page - 1).All(&reportResultUnfinish)

			for _, element := range reportResultUnfinish {

				//Counting 筆數

				caseNum, _ := mongo.DB(lzw).C(panaAssignment).Find(bson.M{"branch": element.Branch}).Count()

				//Get Block

				_ = mongo.DB(lzw).C(panaBranch).Find(bson.M{"branch": element.Branch}).One(&branchResult)

				//Get fixed_time

				ft := element.FixTime + "h"

				//Calculate estimated completion time

				ftDuration, _ := time.ParseDuration(ft)

				bt := element.BuildTime

				btPlusFt := bt.Add(ftDuration)

				//fmt.Println("Above is ok bracnhResult ", branchResult)

				AddRowToSheet(sheetUnfinish, branchResult.Block, element.Branch, element.Maintenance,
					element.AssignGuy, element.BuildTime.String()[:10], element.BuildTime.String()[11:16], element.BuildTime.String()[:10],
					btPlusFt.String()[11:16], strconv.Itoa(caseNum))

			}

		}

		//Save file

		err = file.Save("/data/xls/lzw_Business_Report.xlsx")
		if err != nil {
			_di.Log.Err(err.Error())
		}

		//File done
		//Making it to be an async file stream can improve response time

		fileName := fmt.Sprintf("attachment; filename=\"lzw_Business_Report_%s-%s.xlsx\"", startTime.String()[:10], endTime.String()[:10])

		w.Header().Set("Content-Disposition", fileName)
		w.Header().Set("Content-Type", req.Header.Get("Content-Type"))
		http.ServeFile(w, req, "/data/xls/lzw_Business_Report.xlsx")

	}

	//Finish

	_di.Log.Debug("Generate report done")

	_afterEndPoint(w, req)

}

func (w ReportDoc) ToJsonStr() string {
	return string(w.ToJsonByte())
}

func (w ReportDoc) ToJsonByte() []byte {
	jsonByte, _ := json.Marshal(w)
	return jsonByte
}

/*

	Function Name: AddRowToSheet
	Purpose: Add a row of data to the sheet
	Input:

		sheet:  Pointer to xlsx.Sheet
		param:  Content of cells

	Output:

		 0 for success
		-1 for fail

*/

func AddRowToSheet(sheet *xlsx.Sheet, param ...string) int {

	var rows *xlsx.Row
	var cells *xlsx.Cell

	//Check availability of sheet

	if sheet == nil {
		fmt.Print("sheet is NULL!\n")
		return -1
	}

	//Create new row

	rows = sheet.AddRow()

	for _, paramName := range param {

		//Insert data to cells

		cells = rows.AddCell()
		cells.Value = paramName

	}
	return 0

}
