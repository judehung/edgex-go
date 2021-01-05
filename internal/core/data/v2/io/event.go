//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package io

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/errors"
	dto "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/fxamacker/cbor/v2"
)

const maxEventSize = int64(25 * 1e6) // 25 MB

// EventReader unmarshals a request body into an Event type
type EventReader interface {
	ReadAddEventRequest(reader io.Reader) ([]dto.AddEventRequest, errors.EdgeX)
}

// NewRequestReader returns a BodyReader capable of processing the request body
func NewEventRequestReader(request *http.Request) EventReader {
	contentType := request.Header.Get(clients.ContentType)

	switch strings.ToLower(contentType) {
	case clients.ContentTypeCBOR:
		return NewCborReader()
	default:
		return NewJsonReader()
	}
}

// cborEventReader handles unmarshaling of a CBOR encoded request body payload
type cborEventReader struct{}

// NewCborReader creates a new instance of cborEventReader.
func NewCborReader() cborEventReader {
	return cborEventReader{}
}

// Read reads and converts the request's CBOR encoded event data into an Event struct
func (cborEventReader) ReadAddEventRequest(reader io.Reader) ([]dto.AddEventRequest, errors.EdgeX) {
	var addEvents []dto.AddEventRequest
	bytes, err := ioutil.ReadAll(io.LimitReader(reader, maxEventSize))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindIOError, "cbor event reading failed", err)
	}

	err = cbor.Unmarshal(bytes, &addEvents)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "cbor event decoding failed", err)
	}
	//TODO Shall consider to add bytes as part of AddEvents, so that bytes could be published to message bus later?
	return addEvents, nil
}

// jsonReader handles unmarshaling of a JSON request body payload
type jsonEventReader struct{}

// NewJsonReader creates a new instance of jsonReader.
func NewJsonReader() jsonEventReader {
	return jsonEventReader{}
}

// Read reads and converts the request's JSON event data into an Event struct
func (jsonEventReader) ReadAddEventRequest(reader io.Reader) ([]dto.AddEventRequest, errors.EdgeX) {
	var addEvents []dto.AddEventRequest
	err := json.NewDecoder(reader).Decode(&addEvents)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "event json decoding failed", err)
	}
	return addEvents, nil
}
