package function

import "alfamart-channel/domain"

func MiddlewareResponseParse(config []domain.MiddlewareResponse, response map[string]interface{}) (parsed map[string]interface{}, err error) {
	fparse := make(map[string]interface{}, 0)
	for _, conf := range config {
		fparse[conf.Field.String] = response[conf.Field.String]
	}
	parsed = fparse
	return
}
