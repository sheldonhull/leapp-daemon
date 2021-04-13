package regional_endpoint

import (
  "fmt"
  "leapp_daemon/core/aws/region"
  http_error2 "leapp_daemon/infrastructure/http/http_error"
)

var regionalEndpoints = map[string]string {
	"af-south-1": "https://sts_client.af-south-1.amazonaws.com",
	"ap-east-1": "https://sts_client.ap-east-1.amazonaws.com",
	"ap-northeast-1": "https://sts_client.ap-northeast-1.amazonaws.com",
	"ap-northeast-2": "https://sts_client.ap-northeast-2.amazonaws.com",
	"ap-northeast-3": "https://sts_client.ap-northeast-3.amazonaws.com",
	"ap-south-1": "https://sts_client.ap-south-1.amazonaws.com",
	"ap-southeast-1": "https://sts_client.ap-southeast-1.amazonaws.com",
	"ap-southeast-2": "https://sts_client.ap-southeast-2.amazonaws.com",
	"ca-central-1": "https://sts_client.ca-central-1.amazonaws.com",
	"cn-north-1": "https://sts_client.cn-north-1.amazonaws.com.cn",
	"cn-northwest-1": "https://sts_client.cn-northwest-1.amazonaws.com.cn",
	"eu-central-1": "https://sts_client.eu-central-1.amazonaws.com",
	"eu-north-1": "https://sts_client.eu-north-1.amazonaws.com",
	"eu-south-1": "https://sts_client.eu-south-1.amazonaws.com",
	"eu-west-1": "https://sts_client.eu-west-1.amazonaws.com",
	"eu-west-2": "https://sts_client.eu-west-2.amazonaws.com",
	"eu-west-3": "https://sts_client.eu-west-3.amazonaws.com",
	"me-south-1": "https://sts_client.me-south-1.amazonaws.com",
	"sa-east-1": "https://sts_client.sa-east-1.amazonaws.com",
	"us-east-1": "https://sts_client.us-east-1.amazonaws.com",
	"us-east-2": "https://sts_client.us-east-2.amazonaws.com",
	"us-gov-east-1": "https://sts_client.us-gov-east-1.amazonaws.com",
	"us-gov-west-1": "https://sts_client.us-gov-west-1.amazonaws.com",
	"us-west-1": "https://sts_client.us-west-1.amazonaws.com",
	"us-west-2": "https://sts_client.us-west-2.amazonaws.com",
}

func GetRegionalEndpoint(regionName *string) (*string, error) {
	isRegionValid := region.IsRegionValid(*regionName)
	if !isRegionValid {
		return nil, http_error2.NewUnprocessableEntityError(fmt.Errorf("Region " + *regionName + " not valid"))
	}

	i, ok := regionalEndpoints[*regionName]
	if ok {
		return &i, nil
	} else {
		return &i, http_error2.NewNotFoundError(fmt.Errorf("No such endpoint for regionName " + *regionName))
	}
}
