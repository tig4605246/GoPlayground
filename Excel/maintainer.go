package doc

import (
	"encoding/json"
	"net/http"
	"time"

	"dforcepro.com/api"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type MaintainerAPI bool

func (ma MaintainerAPI) Enable() bool {
	return bool(ma)
}

func (ma MaintainerAPI) GetAPIs() *[]*api.APIHandler {
	return &[]*api.APIHandler{
		&api.APIHandler{Path: "/v1/maintenance_provider", Next: ma.getEndPoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/maintenance_provider/{ID}", Next: ma.getByIDEndPoint, Method: "GET", Auth: true},
		&api.APIHandler{Path: "/v1/maintenance_provider", Next: ma.createEndPoint, Method: "POST", Auth: true},
		&api.APIHandler{Path: "/v1/maintenance_provider/{ID}", Next: ma.editEndPoint, Method: "PUT", Auth: true},
		&api.APIHandler{Path: "/v1/maintenance_provider/{ID}", Next: ma.deleteEndPoint, Method: "DELETE", Auth: true},
	}
}

type MaintainerDoc struct {
	ID              bson.ObjectId `json:"ID,omitempty" bson:"_id"`
	Name            string        `json:"name,omitempty" bson:"name,omitempty"`
	CompanyName     string        `json:"company_name,omitempty" bson:"company_name,omitempty"`
	Location        string        `json:"location,omitempty" bson:"location,omitempty"`
	Phone           string        `json:"phone,omitempty" bson:"phone,omitempty"`
	MaintenanceBoss string        `json:"maintenance_boss,omitempty" bson:"maintenance_boss,omitempty"`
	Performance     Performance   `json:"performance,omitempty"`
}
type Performance struct {
	Total    int       `json:"Total" `
	Finish   int       `json:"Finished_count" `
	ToFix    int       `json:"ToFixCount"`
	OverTime int       `json:"OverTimeCount"`
	ArrivalR float64   `json:"Arrival_rate"`
	FinishR  float64   `json:"Finished_rate" `
	LatestT  time.Time `json:"LatestRecord_time" `
}

const (
	Maintainer = "pana_maintenance"
)

func _NewMaintainerForm() *MaintainerDoc {
	return &MaintainerDoc{}
}

func (wf MaintainerAPI) getEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)

	// Give nothing to query all instances
	var queryStr = bson.M{}
	var results []JobDoc
	var mainresults []MaintainerDoc

	mongo := getMongo()
	queries := req.URL.Query()
	mtname, ok := queries["name"]
	s, ok1 := queries["start"]
	e, ok2 := queries["end"]

	if !ok1 {
		s = append(s, "2000-01-01")
	}
	if !ok2 {
		e = append(e, "2100-01-01")
	}
	shortForm := "2006-01-02"
	t1, _ := time.Parse(shortForm, s[0])
	t2, _ := time.Parse(shortForm, e[0])

	if ok {
		queryStr["name"] = mtname[0]
	}

	query := mongo.DB(lzw).C(panaMaintenaner).Find(queryStr)
	maintotal, err := query.Count()
	if err != nil {
		_di.Log.Err(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
		// TODO:回傳空資料
	}
	if maintotal == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if maintotal > 0 {
		query.All(&mainresults)
		var outputresult []interface{}
		var allfinish, allover, allarrival, allfix, allArrivalR, allFinishR = 0, 0, 0, 0, 0.0, 0.0
		var tofixCount, arrivalCount, finishedCount, overTime, latestT, total = 0, 0, 0, 0, time.Now(), 0
		for _, element := range mainresults {

			//result = append(result, element)
			tofixCount, arrivalCount, finishedCount, overTime, latestT = 0, 0, 0, 0, time.Now()
			queryStr = bson.M{"maintenance": element.Name, "build_time": bson.M{"$gte": t1, "$lt": t2}, "status": bson.M{"$ne": "取消派工"}}
			query = mongo.DB(lzw).C(assign).Find(queryStr).Sort("-build_time")
			total, err = query.Count()

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			if total > 0 {

				query.All(&results)
				if len(results) > 0 {
					latestT = results[0].BuildT
				}
				for _, elements := range results {

					if elements.Status == "未回報" {
						overTime++
					}
					if elements.Status != "已完工" {
						tofixCount++
					} else {
						finishedCount++
						if elements.StartT.Sub(elements.BuildT) <= time.Duration(elements.ArrivalSt)*time.Hour {
							arrivalCount++
						}
					}

				}

			}
			if finishedCount == 0 {
				element.Performance.ArrivalR = 0.0
			} else {
				element.Performance.ArrivalR = (float64(arrivalCount) / float64(finishedCount)) * 100
			}
			if finishedCount+overTime == 0 {
				element.Performance.FinishR = 0
			} else {
				element.Performance.FinishR = (float64(finishedCount) / float64(finishedCount+overTime)) * 100
			}
			element.Performance.Total = total
			element.Performance.Finish = finishedCount
			element.Performance.ToFix = tofixCount
			element.Performance.OverTime = overTime
			element.Performance.LatestT = latestT
			allarrival += arrivalCount
			allfinish += finishedCount
			allover += overTime
			allfix += tofixCount
			outputresult = append(outputresult, element)

		}
		if allfinish == 0 {
			allArrivalR = 0.0
		} else {
			allArrivalR = (float64(allarrival) / float64(allfinish)) * 100
		}
		if allfinish+allover == 0 {
			allFinishR = 0.0
		} else {
			allFinishR = (float64(allfinish) / float64(allfinish+allover)) * 100
		}

		responseJSON := queryPerformance{&outputresult, maintotal, allfix, allFinishR, allArrivalR}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseJSON)
	}

	_afterEndPoint(w, req)

}
func (wf MaintainerAPI) getByIDEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)

	// Give nothing to query all instances

	var queryStr = bson.M{}
	var results []JobDoc
	var mainresults []MaintainerDoc

	mongo := getMongo()
	queries := req.URL.Query()
	s, ok1 := queries["start"]
	e, ok2 := queries["end"]

	if !ok1 {
		s = append(s, "2000-01-01")
	}
	if !ok2 {
		e = append(e, "2100-01-01")
	}
	shortForm := "2006-01-02"
	t1, _ := time.Parse(shortForm, s[0])
	t2, _ := time.Parse(shortForm, e[0])

	vars := mux.Vars(req)
	queryStr["_id"] = bson.ObjectIdHex(vars["ID"])
	query := mongo.DB(lzw).C(panaMaintenaner).Find(queryStr)
	maintotal, err := query.Count()

	if err != nil {
		_di.Log.Err(err.Error())
		return
		// TODO:回傳空資料
	}
	if maintotal == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	query.All(&mainresults)
	var outputresult interface{}

	var tofixCount, arrivalCount, finishedCount, overTime, latestT, total = 0, 0, 0, 0, time.Now(), 0
	for _, element := range mainresults {

		//result = append(result, element)
		tofixCount, arrivalCount, finishedCount, overTime, latestT = 0, 0, 0, 0, time.Now()
		queryStr = bson.M{"maintenance": element.Name, "build_time": bson.M{"$gte": t1, "$lt": t2}, "status": bson.M{"$ne": "取消派工"}}
		query = mongo.DB(lzw).C(assign).Find(queryStr).Sort("-build_time")
		total, err = query.Count()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if total > 0 {

			query.All(&results)
			if len(results) > 0 {
				latestT = results[0].BuildT
			}
			for _, elements := range results {

				if elements.Status == "未回報" {
					overTime++
				}
				if elements.Status != "已完工" {
					tofixCount++
				} else {
					finishedCount++
					if elements.StartT.Sub(elements.BuildT) <= time.Duration(elements.ArrivalSt)*time.Hour {
						arrivalCount++
					}
				}

			}

		}
		if finishedCount == 0 {
			element.Performance.ArrivalR = 0.0
		} else {
			element.Performance.ArrivalR = (float64(arrivalCount) / float64(finishedCount)) * 100
		}
		if finishedCount+overTime == 0 {
			element.Performance.FinishR = 0
		} else {
			element.Performance.FinishR = (float64(finishedCount) / float64(finishedCount+overTime)) * 100
		}
		element.Performance.Total = total
		element.Performance.Finish = finishedCount
		element.Performance.ToFix = tofixCount
		element.Performance.OverTime = overTime
		element.Performance.LatestT = latestT

		outputresult = element

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(outputresult)

	_afterEndPoint(w, req)

}

func (wf MaintainerAPI) createEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)
	maintainerForm := _NewMaintainerForm()
	_ = json.NewDecoder(req.Body).Decode(&maintainerForm)
	maintainerForm.ID = bson.NewObjectId()

	mongo := getMongo()

	err := mongo.DB(lzw).C(Maintainer).Insert(maintainerForm)
	if err != nil {
		// 寫 log 將 Document 寫進資料夾
		_di.Log.WriteFile("document/maintainer/create.log", maintainerForm.ToJsonByte())
		_di.Log.Err(err.Error())
		// 回傳錯誤息
		w.WriteHeader(http.StatusBadRequest)

	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(maintainerForm.ID.Hex()))
	}

	_afterEndPoint(w, req)
}

func (wf MaintainerAPI) editEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)

	//var queryStr = bson.M{}

	vars := mux.Vars(req)
	maintainerForm := _NewMaintainerForm()
	_ = json.NewDecoder(req.Body).Decode(&maintainerForm)

	//Check existence of the input parameters

	if len(vars["ID"]) <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRes{"ID error"})
		return

	}

	maintainerForm.ID = bson.ObjectIdHex(vars["ID"])

	mongo := getMongo()

	collection := mongo.DB(lzw).C(Maintainer)

	err := collection.Update(bson.M{"_id": bson.ObjectIdHex(vars["ID"])}, bson.M{"$set": maintainerForm})
	if err != nil {
		_di.Log.Err(err.Error())
	}

	if err != nil {
		// 寫 log 將 Document 寫進資料夾
		_di.Log.WriteFile("document/maintainer/edit.log", maintainerForm.ToJsonByte())
		_di.Log.Err(err.Error())
		// 回傳錯誤息
		w.WriteHeader(http.StatusBadRequest)

	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("success"))
	}

	_afterEndPoint(w, req)
}

func (wf MaintainerAPI) deleteEndPoint(w http.ResponseWriter, req *http.Request) {

	_beforeEndPoint(w, req)
	vars := mux.Vars(req)

	mongo := getMongo()

	//Check existence of the input parameters

	if len(vars["ID"]) <= 0 {

		json.NewEncoder(w).Encode(errorRes{"ID error"})
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	now := time.Now()
	unixNowTime := now.Unix()

	err := mongo.DB(lzw).C(Maintainer).Update(bson.M{"_id": bson.ObjectIdHex(vars["ID"])}, bson.M{"$set": bson.M{"_deleted": unixNowTime}})
	if err != nil {
		_di.Log.Err(err.Error())
	}

	if err != nil {
		// 寫 log 將 Document 寫進資料夾
		_di.Log.Err(err.Error())
		// 回傳錯誤息
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	_afterEndPoint(w, req)
}

func (w MaintainerDoc) ToJsonStr() string {
	return string(w.ToJsonByte())
}

func (w MaintainerDoc) ToJsonByte() []byte {
	jsonByte, _ := json.Marshal(w)
	return jsonByte
}
