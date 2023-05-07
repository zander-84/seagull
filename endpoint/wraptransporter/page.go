package wraptransporter

import (
	"github.com/zander-84/seagull/transport"
)

func SetToken(ts transport.Transporter, token string) {
	ts.Body().Set("_token", token)
}

func GetToken(ts transport.Transporter) string {
	data, ok := ts.Body().Get("_token")
	if !ok {
		return ""
	}
	res, _ := data.(string)
	return res
}
func SetTraceID(ts transport.Transporter, traceID string) {
	ts.Body().Set("_traceId", traceID)
}

func GetTraceID(ts transport.Transporter) string {
	data, ok := ts.Body().Get("_traceId")
	if !ok {
		return ""
	}
	res, _ := data.(string)
	return res
}

func SetUser(ts transport.Transporter, user any) {
	ts.Body().Set("_user", user)
}

// GetUser 每个项目需要再次wrap下
func GetUser(ts transport.Transporter) any {
	data, _ := ts.Body().Get("_user")
	return data
}

func SetPage(ts transport.Transporter, page int) {

	ts.Body().Set("_page", page)
}

func GetPage(ts transport.Transporter) int {
	data, ok := ts.Body().Get("_page")
	if !ok {
		return 0
	}
	res, _ := data.(int)
	return res
}

func SetPageSize(ts transport.Transporter, pageSize int) {
	ts.Body().Set("_pageSize", pageSize)
}

func GetPageSize(ts transport.Transporter) int {
	data, ok := ts.Body().Get("_pageSize")
	if !ok {
		return 0
	}
	res, _ := data.(int)

	return res
}

func SetCount(ts transport.Transporter, cnt int64) {
	ts.Body().Set("_cnt", cnt)
}

func GetCount(ts transport.Transporter) int64 {
	data, ok := ts.Body().Get("_cnt")
	if !ok {
		return 0
	}
	res, _ := data.(int64)
	return res
}

func SetData(ts transport.Transporter, data any) {
	ts.Body().Set("_data", data)
}

func GetData(ts transport.Transporter) any {
	data, _ := ts.Body().Get("_data")
	return data
}

func Set(ts transport.Transporter, key string, data any) {
	ts.Body().Set(key, data)
}

func Get(ts transport.Transporter, key string) (any, bool) {
	data, ok := ts.Body().Get(key)
	return data, ok
}
func ShouldGet(ts transport.Transporter, key string) any {
	data, _ := ts.Body().Get(key)
	return data
}
