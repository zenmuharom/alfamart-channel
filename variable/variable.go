package variable

import (
	"alfamart-channel/models"
	"errors"
)

var InvalidFormat = models.ErrorMsg{
	GoCode: 1706,
	Err:    errors.New("Invalid Format"),
}

var UnknownServiceType = models.ErrorMsg{
	GoCode: 1707,
	Err:    errors.New("Unknown service type"),
}

const REQUEST_BODY = "requestBody"
const TYPE_OBJECT = "object"

const ROUTE_HTTP_RECEIVER = "HTTP_RECEIVER"
const ROUTE_HTTP_CLIENT = "HTTP_CLIENT"
const ROUTE_HTTP_SENDER = "HTTP_SENDER"
