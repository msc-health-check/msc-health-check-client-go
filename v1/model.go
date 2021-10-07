package v1

type ProjectCheck struct {
	ID          string
	URL         string
	AppName     string
	ChecksOut   []map[string]interface{}
	LiveSignals []map[string]interface{}
	Errors      []map[string]interface{}
}

type ProjectCheckRequest struct {
	URL         string
	AppName     string
}
