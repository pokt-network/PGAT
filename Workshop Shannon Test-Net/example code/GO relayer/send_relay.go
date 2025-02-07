package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"shannon_send_relay/shannon"
	"shannon_send_relay/types"
)

type Config struct {
	AppAddr    string           `json:"app_address"`
	AppPrivHex string           `json:"app_private_key_hex"`
	ServiceID  string           `json:"service_id"`
	RpcUrl     string           `json:"rpc_url"`
	GrcpConf   types.GRPCConfig `json:"grpc_config"`
	Path       string           `json:"path"`
	Payload    string           `json:"payload"`
}

func ReadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {

	// Create a new context
	ctx := context.Background()

	// Load the configuration from the config.json file
	config, err := ReadConfig(os.Getenv("CONFIG_FILE"))
	if err != nil {
		log.Fatal(err)
	}

	nodeConfig := types.FullNodeConfig{
		RpcURL:     config.RpcUrl,
		GRPCConfig: config.GrcpConf,
	}

	// Create a LazyFull node from the config
	FullNode, err := shannon.NewLazyFullNode(nodeConfig)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Failed to create Lazy Node")
	}

	// Check if the app is correctly staked for service
	onchainApp, err := FullNode.GetApp(ctx, config.AppAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal(fmt.Sprintf("Error getting on-chain data for app %s", config.AppAddr))
	}
	if onchainApp == nil {
		log.Fatal(fmt.Sprintf("No data found for app %s", config.AppAddr))
	}
	fmt.Printf("App %s found, stake: %s\n", onchainApp.Address, onchainApp.Stake)

	// Check if the app is staked for the requested service
	if !shannon.AppIsStakedForService(types.ServiceID(config.ServiceID), onchainApp) {
		log.Fatal(fmt.Sprintf("App %s is not staked for service %s", config.AppAddr, config.ServiceID))
	}

	// Get App session
	appSession, err := FullNode.GetSession(types.ServiceID(config.ServiceID), config.AppAddr)
	if err != nil {
		fmt.Println(err)
		log.Fatal(fmt.Sprintf("Error getting session data for app %s in service %s", config.AppAddr, config.ServiceID))
	}
	fmt.Println(appSession)

	// Get all the endpoint available in this session
	endpoints, err := shannon.EndpointsFromSession(appSession)
	if err != nil {
		fmt.Println(err)
		log.Fatal(fmt.Sprintf("Failed getting endpoints for current session", config.AppAddr, config.ServiceID))
	}
	if len(endpoints) < 1 {
		log.Fatal("No endpoint found for the requested service. Are there Servicers staked?")
	} else {
		fmt.Printf("Found %d endpoints!\n", len(endpoints))
	}

	// Pick one at random
	// Extract keys into a slice
	servicers := make([]types.EndpointAddr, 0, len(endpoints))
	for k := range endpoints {
		servicers = append(servicers, k)
	}
	// Pick a random servicer
	selectedServicer := servicers[rand.Intn(len(servicers))]
	selectedEndpoint := endpoints[selectedServicer]
	fmt.Printf("Selected endpoint supplier: %s\n\t URL: %s\n", selectedEndpoint.Addr(), selectedEndpoint.PublicURL())

	// Create a signer
	signerApp := shannon.RelayRequestSigner{
		AccountClient: *FullNode.GetAccountClient(),
		PrivateKeyHex: config.AppPrivHex,
	}

	// Build the payload
	thisPayload := types.Payload{
		Data:            config.Payload,
		Method:          "POST",
		Path:            config.Path,
		TimeoutMillisec: 10000,
	}

	// Send the relay
	response, err := shannon.SendRelay(thisPayload, selectedEndpoint, types.ServiceID(config.ServiceID), *FullNode, signerApp)
	if err != nil {
		log.Fatal(fmt.Printf("relay: error sending relay for service %s endpoint %s: %w",
			config.ServiceID, selectedEndpoint.Addr(), err,
		))
	}

	// The Payload field of the response received from the endpoint, i.e. the relay miner,
	// is a serialized http.Response struct. It needs to be deserialized into an HTTP Response struct
	// to access the Service's response body, status code, etc.
	relayResponse, err := shannon.DeserializeRelayResponse(response.Payload)
	if err != nil {
		log.Fatal(fmt.Printf("relay: error unmarshalling endpoint response into a POKTHTTP response for service %s endpoint %s: %w",
			config.ServiceID, selectedEndpoint.Addr(), err,
		))
	}
	relayResponse.EndpointAddr = selectedEndpoint.Addr()

	fmt.Printf("Relay Succeeded!! RPC status: %d\n", relayResponse.HTTPStatusCode)
	fmt.Printf("Response: %s\n", string(relayResponse.Bytes))
}
