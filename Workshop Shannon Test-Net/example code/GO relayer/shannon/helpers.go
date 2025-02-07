package shannon

import (
	"shannon_send_relay/types"

	apptypes "github.com/pokt-network/poktroll/x/application/types"
	sessiontypes "github.com/pokt-network/poktroll/x/session/types"
	sdk "github.com/pokt-network/shannon-sdk"
)

func AppIsStakedForService(serviceID types.ServiceID, app *apptypes.Application) bool {
	for _, svcCfg := range app.ServiceConfigs {
		if types.ServiceID(svcCfg.ServiceId) == serviceID {
			return true
		}
	}

	return false
}

// endpointsFromSession returns the list of all endpoints from a Shannon session.
// It returns a map for efficient lookup, as the main/only consumer of this function uses
// the return value for selecting an endpoint for sending a relay.
func EndpointsFromSession(session sessiontypes.Session) (map[types.EndpointAddr]Endpoint, error) {
	sf := sdk.SessionFilter{
		Session: &session,
	}

	allEndpoints, err := sf.AllEndpoints()
	if err != nil {
		return nil, err
	}

	endpoints := make(map[types.EndpointAddr]Endpoint)
	for _, supplierEndpoints := range allEndpoints {
		for _, supplierEndpoint := range supplierEndpoints {
			endpoint := Endpoint{
				Supplier: string(supplierEndpoint.Supplier()),
				Url:      supplierEndpoint.Endpoint().Url,
				// Set the session field on the endpoint for efficient lookup when sending relays.
				Session: session,
			}
			endpoints[endpoint.Addr()] = endpoint
		}
	}

	return endpoints, nil
}
