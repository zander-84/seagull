package wraptransporter

import "github.com/zander-84/seagull/endpoint"

func SetToken(ts endpoint.Transporter, token string) {
	ts.Body().Set("_token", token)
}

func GetToken(ts endpoint.Transporter) string {
	data, ok := ts.Body().Get("_token")
	if !ok {
		return ""
	}
	res, _ := data.(string)
	return res
}
func SetTraceID(ts endpoint.Transporter, traceID string) {
	ts.Body().Set("_traceId", traceID)
}

func GetTraceID(ts endpoint.Transporter) string {
	data, ok := ts.Body().Get("_traceId")
	if !ok {
		return ""
	}
	res, _ := data.(string)
	return res
}

func SetUser(ts endpoint.Transporter, user any) {
	ts.Body().Set("_user", user)
}

// GetUser 每个项目需要再次wrap下
func GetUser(ts endpoint.Transporter) any {
	data, _ := ts.Body().Get("_user")
	return data
}

func SetPage(ts endpoint.Transporter, page int) {

	ts.Body().Set("_page", page)
}

func GetPage(ts endpoint.Transporter) int {
	data, ok := ts.Body().Get("_page")
	if !ok {
		return 0
	}
	res, _ := data.(int)
	return res
}

func SetPageSize(ts endpoint.Transporter, pageSize int) {
	ts.Body().Set("_pageSize", pageSize)
}

func GetPageSize(ts endpoint.Transporter) int {
	data, ok := ts.Body().Get("_pageSize")
	if !ok {
		return 0
	}
	res, _ := data.(int)

	return res
}

func SetCount(ts endpoint.Transporter, cnt int64) {
	ts.Body().Set("_cnt", cnt)
}

func GetCount(ts endpoint.Transporter) int64 {
	data, ok := ts.Body().Get("_cnt")
	if !ok {
		return 0
	}
	res, _ := data.(int64)
	return res
}

func SetData(ts endpoint.Transporter, data any) {
	ts.Body().Set("_data", data)
}

func GetData(ts endpoint.Transporter) any {
	data, _ := ts.Body().Get("_data")
	return data
}

func Set(ts endpoint.Transporter, key string, data any) {
	ts.Body().Set(key, data)
}

func Get(ts endpoint.Transporter, key string) (any, bool) {
	data, ok := ts.Body().Get(key)
	return data, ok
}
func ShouldGet(ts endpoint.Transporter, key string) any {
	data, _ := ts.Body().Get(key)
	return data
}
