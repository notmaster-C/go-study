package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go-study/db"
	"net/http"
	"strings"
)

type apiAndDb struct {
	Db  string `json:"db"`
	Api string `json:"api"`
}

var (
	apiClient    = make(map[string]apiInfo)
	restyRequest = make(map[string]*resty.Request)
)

const (
	authTokenByHeaderBearer = "headerBearerToken"
	contentTypeJson         = "application/json"
)

// {"host":"http://202.115.158.71:8770", "authType":"headerBearerToken", "token":"PanguSuperAdmin@2020"}
type apiInfo struct {
	Host     string `json:"host"`
	AuthType string `json:"authType"`
	Token    string `json:"token"`
}
type statApiImpl struct {
	Uri         string `json:"uri"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Params      string `json:"params"`
}

func getClient(ck string) *resty.Request {
	request := resty.New().R()

	request.SetHeader("Authorization", "Bearer 9dfb67e6-d6b3-4d2f-b970-3dee2a58a2b8")

	request.URL = "http://202.115.158.71:8770/"

	return request

	return nil
}
func refreshStat(sid int64, value string) {
	err := db.GetDb("mysql").Model(&db.StatReport{}).Where("id = ?", sid).Update("value", value).Error
	if err != nil {
		fmt.Errorf("failed to update stat-report statId:%v value:%v error:%v", sid, value, err)
	}
	fmt.Errorf("stat-report data refresh successfly with statId:%v value:%v", sid, value)
}
func ExecuteStatApiDb(data *db.StatDataSource) error {
	sid := data.StatId
	var ad *apiAndDb
	errs := json.Unmarshal([]byte(data.ExtendKey), &ad)
	if errs != nil {
		return fmt.Errorf("json.Unmarshal failed:%v", ad)
	}
	statApi := getClient(ad.Api)
	if statApi == nil {
		return fmt.Errorf("invaild api-client key:%v", data.ExtendKey)
	}
	var (
		err error
		api statApiImpl
	)
	err = json.Unmarshal([]byte(data.Param), &api)
	if err != nil {
		fmt.Errorf("failed to json unmarshal api-stat whith statDataSource:%+v error:%v", data.Param, err)
		return err
	}
	url := fmt.Sprintf("%v/%v", statApi.URL, api.Uri)
	if api.ContentType == contentTypeJson {
		statApi.SetHeader("Content-Type", contentTypeJson)
		statApi.SetBody(api.Params)
	} else {
		statApi.SetQueryString(api.Params)
	}
	var resp *resty.Response
	switch strings.ToUpper(api.Method) {
	case resty.MethodGet:
		resp, err = statApi.Get(url)
	case resty.MethodPost:
		resp, err = statApi.Post(url)
	case resty.MethodPut:
		resp, err = statApi.Put(url)
	case resty.MethodHead:
		resp, err = statApi.Head(url)
	case resty.MethodDelete:
		resp, err = statApi.Delete(url)
	}

	var xjt1 db.XJT1
	err = json.Unmarshal(resp.Body(), &xjt1)
	if err != nil {
		fmt.Errorf("failed to api-stat url:%v method:%v params:%v error:%v", url, api.Method, api.Params, err)
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		statDb := db.GetDb("mysql")

		for _, record := range xjt1.Result.Records {
			statDb.Migrator().AutoMigrate(&record)
			statDb.Save(&record)
		}
		go refreshStat(sid, resp.String())
	}

	return nil
}
