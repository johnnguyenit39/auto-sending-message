package httprequest

import (
	"messenging_test/common"
	"net/http"
)

type HttpRequestRepository interface {
	DoRequest(request common.HttpRequestModel) (*http.Response, error)
}
