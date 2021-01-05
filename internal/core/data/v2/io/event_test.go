package io

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	v2 "github.com/edgexfoundry/go-mod-core-contracts/v2"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	dto "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/fxamacker/cbor/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ExampleUUID            = "82eb2e26-0f24-48aa-ae4c-de9dac3fb9bc"
	TestDeviceName         = "TestDevice"
	TestOriginTime         = 1600666185705354000
	TestDeviceResourceName = "TestDeviceResourceName"
	TestDeviceProfileName  = "TestDeviceProfileName"
	TestReadingValue       = "45"
)

var expectedEventId = uuid.New().String()

var testReading = dtos.BaseReading{
	DeviceName:   TestDeviceName,
	ResourceName: TestDeviceResourceName,
	ProfileName:  TestDeviceProfileName,
	Origin:       TestOriginTime,
	ValueType:    v2.ValueTypeUint8,
	SimpleReading: dtos.SimpleReading{
		Value: TestReadingValue,
	},
}

var testAddEvent = dto.AddEventRequest{
	BaseRequest: common.BaseRequest{
		RequestId: ExampleUUID,
	},
	Event: dtos.Event{
		Id:          expectedEventId,
		DeviceName:  TestDeviceName,
		ProfileName: TestDeviceProfileName,
		Origin:      TestOriginTime,
		Readings:    []dtos.BaseReading{testReading},
	},
}

func newRequestWithContentType(contentType string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("Test body"))
	req.Header.Set(clients.ContentType, contentType)
	return req
}

func TestNewEventRequestReader(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		expectedType interface{}
	}{
		{"Get Json Reader", clients.ContentTypeJSON, jsonEventReader{}},
		{"Get Cbor Reader", clients.ContentTypeCBOR, cborEventReader{}},
		{"Get Json Reader when content-type is unknown", "Unknown-Type", jsonEventReader{}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req := newRequestWithContentType(testCase.contentType)
			reader := NewEventRequestReader(req)
			assert.IsType(t, testCase.expectedType, reader, "unexpected reader type")
		})
	}
}

func TestJsonSerialization(t *testing.T) {
	tests := []struct {
		name          string
		targetDTO     interface{}
		errorExpected bool
	}{
		{"Valid", []dto.AddEventRequest{testAddEvent}, false},
		{"Invalid", []string{"string1", "string2"}, true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req := newRequestWithContentType(clients.ContentTypeJSON)
			jsonReader := NewEventRequestReader(req)
			byteArray, err := json.Marshal(testCase.targetDTO)
			require.NoError(t, err, "error occurs during json marshalling")
			r := ioutil.NopCloser(bytes.NewBuffer(byteArray))
			_, err = jsonReader.ReadAddEventRequest(r)
			if testCase.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCborSerialization(t *testing.T) {
	tests := []struct {
		name          string
		targetDTO     interface{}
		errorExpected bool
	}{
		{"Valid", []dto.AddEventRequest{testAddEvent}, false},
		{"Invalid", []string{"string1", "string2"}, true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req := newRequestWithContentType(clients.ContentTypeCBOR)
			cborReader := NewEventRequestReader(req)
			byteArray, err := cbor.Marshal(testCase.targetDTO)
			require.NoError(t, err, "error occurs during cbor marshalling")
			r := ioutil.NopCloser(bytes.NewBuffer(byteArray))
			_, err = cborReader.ReadAddEventRequest(r)
			if testCase.errorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
