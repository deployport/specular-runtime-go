package client

import (
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
)

// HTTPClientResult is the result of a HTTP client
type HTTPClientResult struct {
	ServerResult *HTTPResult
	Multipart    *multipart.Reader
}

// ParseHTTPClientResult parses a HTTP client result
// error is returned if there was an error parsing the result or an error in the transport
// HTTPResult is returned if the result was parsed successfully
// note that HTTPResult can also have a Struct representing the error of the spec
func ParseHTTPClientResult(pk *Package, res *http.Response, err error) (*HTTPClientResult, error) {
	if err != nil {
		// TODO: some HTTP status contain embedded responses in errors
		return nil, err
	}
	contentTypeHeader := res.Header.Get("Content-Type")
	// log.Printf("content type header: %v", contentTypeHeader)
	mediaType, params, err := mime.ParseMediaType(contentTypeHeader)
	if err != nil {
		return nil, fmt.Errorf("error parsing media type: %v", err)
	}
	if mediaType == HTTPResultMimeTypeHeartbeat.String() {
		return &HTTPClientResult{
			ServerResult: &HTTPResult{
				Heartbeat: true,
			},
		}, nil
	}
	if mediaType == HTTPResultMimeTypeError.String() {
		dec := json.NewDecoder(res.Body)
		he := &HTTPError{}
		if err := dec.Decode(&he); err != nil {
			return nil, fmt.Errorf("failed to decode JSON error, %w", err)
		}
		return &HTTPClientResult{
			ServerResult: &HTTPResult{
				Err: he,
			},
		}, err
	}
	if mediaType == HTTPResultMimeTypeStruct.String() {
		dec := json.NewDecoder(res.Body)
		st := NewContent()
		if err := dec.Decode(&st); err != nil {
			return nil, err
		}
		resTypeFQTN, err := TypeFQTNFromString(st.GetStruct())
		if err != nil {
			return nil, fmt.Errorf("failed to parse response type FQDN: %w", err)
		}

		resType, err := pk.TypeByFQDN(resTypeFQTN)
		if err != nil {
			return nil, err
		}
		sst := resType.TypeBuilder()()
		if err := sst.Hydrate(NewHydratationContext(pk, st)); err != nil {
			return nil, fmt.Errorf("failed to hydrate output of operation, %w", err)
		}
		return &HTTPClientResult{
			ServerResult: &HTTPResult{
				Struct: sst,
			},
		}, nil
	}
	if mediaType == "multipart/mixed" {
		boundary, ok := params["boundary"]
		if !ok {
			return nil, fmt.Errorf("no boundary found")
		}
		mr := multipart.NewReader(res.Body, boundary)
		return &HTTPClientResult{
			Multipart: mr,
		}, nil
	}
	return nil, fmt.Errorf("server returned unknown media type: %s", mediaType)
}
