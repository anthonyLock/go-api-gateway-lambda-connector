package goapigatewaylambdaconnector

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var (
	errAPIGatewayV2UnexpectedRequest = errors.New("expected APIGatewayV2HTTPRequest event")
)

func newAPIGatewayV2HTTPRequest(ctx context.Context, payload []byte, opts *Options) (lambdaRequest, error) {
	var event events.APIGatewayV2HTTPRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return lambdaRequest{}, err
	}
	if event.Version != "2.0" {
		return lambdaRequest{}, errAPIGatewayV2UnexpectedRequest
	}

	//Add cookie to header.
	event.Headers["Cookie"] = strings.Join(event.Cookies, ";")

	req := lambdaRequest{
		HTTPMethod:                      event.RequestContext.HTTP.Method,
		Path:                            event.RequestContext.HTTP.Path,
		Headers:                         event.Headers,
		MultiValueQueryStringParameters: map[string][]string{},
		QueryStringParameters:           map[string]string{},
		Body:                            event.Body,
		IsBase64Encoded:                 event.IsBase64Encoded,
		SourceIP:                        event.RequestContext.HTTP.SourceIP,
		Context:                         newAPIGatewayV2HTTPRequestContext(ctx, event),
	}
	for k, v := range event.QueryStringParameters {
		if strings.Contains(v, ",") {
			req.MultiValueQueryStringParameters[k] = strings.Split(v, ",")
		} else {
			req.QueryStringParameters[k] = v

		}
	}

	if opts.UseProxyPath {
		req.Path = path.Join("/", event.PathParameters["proxy"])
	}

	return req, nil
}
