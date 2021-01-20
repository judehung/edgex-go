//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/di"
	v2Clients "github.com/edgexfoundry/go-mod-core-contracts/v2/v2/clients/interfaces"
)

// MetadataDeviceClientName contains the name of the Metadata device client instance in the DIC.
var MetadataDeviceClientName = "V2MetadataDeviceClient"

// MetadataDeviceProfileClientName contains the name of the Metadata device profile client instance in the DIC.
var MetadataDeviceProfileClientName = "V2MetadataDeviceProfileClient"

// MetadataDeviceClientFrom helper function queries the DIC and returns the Metadata device client instance.
func MetadataDeviceClientFrom(get di.Get) v2Clients.DeviceClient {
	return get(MetadataDeviceClientName).(v2Clients.DeviceClient)
}

// MetadataDeviceProfileClientFrom helper function queries the DIC and returns the Metadata device profile client instance.
func MetadataDeviceProfileClientFrom(get di.Get) v2Clients.DeviceProfileClient {
	return get(MetadataDeviceProfileClientName).(v2Clients.DeviceProfileClient)
}
