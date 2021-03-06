package gremgo

import (
	"reflect"
	"testing"
)

/*
Dummy responses for mocking
*/

var dummySuccessfulResponse = []byte(`{"result":{"data":[{"id": 2,"label": "person","type": "vertex","properties": [
  {"id": 2, "value": "vadas", "label": "name"},
  {"id": 3, "value": 27, "label": "age"}]}
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":200,"attributes":{},"message":""}}`)

var dummyPartialResponse1 = []byte(`{"result":{"data":[{"id": 2,"label": "person","type": "vertex","properties": [
  {"id": 2, "value": "vadas", "label": "name"},
  {"id": 3, "value": 27, "label": "age"}]},
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":206,"attributes":{},"message":""}}`)

var dummyPartialResponse2 = []byte(`{"result":{"data":[{"id": 4,"label": "person","type": "vertex","properties": [
  {"id": 5, "value": "quant", "label": "name"},
  {"id": 6, "value": 54, "label": "age"}]},
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":200,"attributes":{},"message":""}}`)

var dummySuccessfulResponseMarshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      200,
	data:      "testData",
}

var dummyPartialResponse1Marshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      206, // Code 206 indicates that the response is not the terminating response in a sequence of responses
	data:      "testPartialData1",
}

var dummyPartialResponse2Marshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      200,
	data:      "testPartialData2",
}

// TestResponseHandling tests the overall response handling mechanism of gremgo
func TestResponseHandling(t *testing.T) {
	c := newClient()

	c.handleResponse(dummySuccessfulResponse)

	var expected []interface{}
	expected = append(expected, dummySuccessfulResponseMarshalled.data)

	if reflect.TypeOf(expected).String() != reflect.TypeOf(c.retrieveResponse(dummySuccessfulResponseMarshalled.requestid)).String() {
		t.Error("Expected data type does not match actual.")
	}
}

// TestResponseMarshalling tests the ability to marshal a response into a designated response struct for further manipulation
func TestResponseMarshalling(t *testing.T) {
	resp, err := marshalResponse(dummySuccessfulResponse)
	if err != nil {
		t.Error(err)
	}
	if dummySuccessfulResponseMarshalled.requestid != resp.requestid || dummySuccessfulResponseMarshalled.code != resp.code {
		t.Error("Expected requestid and code does not match actual.")
	} else if reflect.TypeOf(resp.data).String() != "[]interface {}" {
		t.Error("Expected data type does not match actual.")
	}
}

// TestResponseSortingSingleResponse tests the ability for sortResponse to save a response received from Gremlin Server
func TestResponseSortingSingleResponse(t *testing.T) {

	c := newClient()

	c.saveResponse(dummySuccessfulResponseMarshalled)

	var expected []interface{}
	expected = append(expected, dummySuccessfulResponseMarshalled.data)

	if reflect.DeepEqual(c.results[dummySuccessfulResponseMarshalled.requestid], expected) != true {
		t.Fail()
	}
}

// TestResponseSortingMultipleResponse tests the ability for the sortResponse function to categorize and group responses that are sent in a stream
func TestResponseSortingMultipleResponse(t *testing.T) {

	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled.data)
	expected = append(expected, dummyPartialResponse2Marshalled.data)

	if reflect.DeepEqual(c.results[dummyPartialResponse1Marshalled.requestid], expected) != true {
		t.Fail()
	}
}

// TestResponseRetrieval tests the ability for a requester to retrieve the response for a specified requestid generated when sending the request
func TestResponseRetrieval(t *testing.T) {
	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	resp := c.retrieveResponse(dummyPartialResponse1Marshalled.requestid)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled.data)
	expected = append(expected, dummyPartialResponse2Marshalled.data)

	if reflect.DeepEqual(resp, expected) != true {
		t.Fail()
	}
}

// TestResponseDeletion tests the ability for a requester to clean up after retrieving a response after delivery to a client
func TestResponseDeletion(t *testing.T) {
	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	c.deleteResponse(dummyPartialResponse1Marshalled.requestid)

	if len(c.results[dummyPartialResponse1Marshalled.requestid]) != 0 {
		t.Fail()
	}
}

var codes = []struct {
	code int
}{
	{200},
	{204},
	{206},
	{401},
	{407},
	{498},
	{499},
	{500},
	{597},
	{598},
	{599},
	{3434}, // Testing unknown error code
}

// Tests detection of errors and if an error is generated for a specific error code
func TestResponseErrorDetection(t *testing.T) {
	for _, co := range codes {
		err := responseDetectError(co.code)
		switch {
		case co.code == 200:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		case co.code == 204:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		case co.code == 206:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		default:
			if err == nil {
				t.Log("Unsuccessful response did not return error.")
			}
		}
	}
}
