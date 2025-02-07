# Sending a Relay using the Shannon SDK

**REQUIREMENTS**: A running full node. See [Full Node Setup](./Full%20Node%20Setup.md)
**REQUIREMENTS**: A Servicer. See [Supplier and Relay Miner](./Supplier%20and%20Relay%20Miner.md)
**REQUIREMENTS**: A staked Application. See [Applications and Gateways](./Applications%20and%20Gateways.md)

### Using Shanon SDK Example

So, you like to get your hands dirty? This is the way then...

This repository comes with an example called "send_relay", located at `./example code/GO relayer`. This example is a bare-bones implementation of the [Shanon GO SDK](https://github.com/pokt-network/shannon-sdk/) and **its sole intention** is to show what are the main parts required to send a relay through the POKT Network.

First make sure you are positioned there:
`cd example\ code/GO\ relayer`
To make the example work, you need the relay miner running and an staked application. Additionally, you will need to set the config file, for which we provide a sample:
`cp config.json.sample config.json`
And modify the following fields:
- `app_address`: Set it to the app address.
- `app_private_key_hex`: Set it to the hex representation of the App's private key.

Then just build and run it:
`go build send_relay.go && CONFIG_FILE=./config.json ./send_relay`
The output should be something like this (ignore the nonsense of an answer from the model...):
```
App pokt1jc3ttp3w5lku9cxh2k23w0f9uskvekgsdsqchr found, stake: 100000000upokt
Found 1 endponits!
Selected endpoint supplier: pokt1r90ujjku55rldjpxsuwx0s2cg7yp5uphxnaa5l-http://rawthil.zapto.org:8545
         URL: http://rawthil.zapto.org:8545
Relay Succeded!! RPC status: 200
Response: {"id":"cmpl-5a7ff6653fd44cabad6a233a2e1df4ae","object":"text_completion","created":1738939926,"model":"pocket_network","choices":[{"index":0,"text":"\n\nAre you a lawyer?\n\nAre you authorized to give legal advice?\n\nI am an AI chatbot assistant and cannot give","logprobs":null,"finish_reason":"length","stop_reason":null,"prompt_logprobs":null}],"usage":{"prompt_tokens":5,"total_tokens":30,"completion_tokens":25,"prompt_tokens_details":null}}
```

If you are not using an LLM backend, you can change the payload to whatever you desire, by means of the `payload` field in the `config.json`.

### Step by Step

The code tries to be self-explanatory, but these are the main steps that you will see.

First we read the config and instantiate a `LazyFullNode`, which we will use to query the network. This instance will abstract away all calls to cosmos, it is the one using the Shannon SDK in the back. You can see the code in the `./shannon` folder.
```go
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
```

Then we check if the Application that was passed is staked for the selected service ID, otherwise the relay will fail, as the relay miner node will reject it:

```go
// Check if the app is correctly staked for service
onchainApp, err := FullNode.GetApp(ctx, config.AppAddr)
if err != nil {
    fmt.Println(err)
    log.Fatal(fmt.Sprintf("Error getting onchain data for app %s", config.AppAddr))
}
if onchainApp == nil {
    log.Fatal(fmt.Sprintf("No data found for app %s", config.AppAddr))
}
fmt.Printf("App %s found, stake: %s\n", onchainApp.Address, onchainApp.Stake)

// Check if the app is staked for the requested service
if !shannon.AppIsStakedForService(types.ServiceID(config.ServiceID), onchainApp) {
    log.Fatal(fmt.Sprintf("App %s is not staked for service %s", config.AppAddr, config.ServiceID))
}
```

After we look for the current session of this app, remember that in POKT the available nodes get rotated from session to session, so you will need to do this each time to ensure that the target relay-miner node will respond. If the target relay-miner node is not in session with the application that you are using, it will reject the relay.
```go
// Get App session
appSession, err := FullNode.GetSession(types.ServiceID(config.ServiceID), config.AppAddr)
if err != nil {
    fmt.Println(err)
    log.Fatal(fmt.Sprintf("Error getting session data for app %s in service %s", config.AppAddr, config.ServiceID))
}
fmt.Println(appSession)
```

With the session data, we look for the actual endpoint that are available. If you recall when you staked the servicer, you added a URL to the stake command, that URL is written in the blockchain and we are going to retrieve that now.
```go
// Get all the endpoint available in this session
endpoints, err := shannon.EndpointsFromSession(appSession)
if err != nil {
    fmt.Println(err)
    log.Fatal(fmt.Sprintf("Failed getting endpoints for current session", config.AppAddr, config.ServiceID))
}
if len(endpoints) < 1 {
    log.Fatal("No endpoint found for the requested service. Are there Servicers staked?")
} else {
    fmt.Printf("Found %d endponits!\n", len(endpoints))
}
```

The result of this will be a map of relay-miner node addresses and endpoint data, containing among other things, the URL that we need to call to reach the relay-miner node.
At this point you might want to apply any kind of discrimination mechanism to ensure Quality of Service (QoS). The network itself wont provide you with any kind of data on the relay-miner node's service, so your custom gateway will probably need to source other data sources for this (like a private database on each node's reliability).
Since this is a **bare bones** example, we wont do any of that, and we will just pick one randomly:
```go
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
```

Now we have an endpoint to target so we must prepare the payload. A critical step here is to create a `signer` that will sign each transaction sent to the relay-miner node. The signer is created using the application private key and it will be the one signing the payload. If the payload sent to the relay-miner node is not signed, the node will reject the RPC.
```go
// Create a signer
signerApp := shannon.RelayRequestSigner{
    AccountClient: *FullNode.GetAccountClient(),
    PrivateKeyHex: config.AppPrivHex,
}
```
The payload is what actually reaches the relay-miner node's backend. In our example the `vLLM` engine. The data will contain the request data, the `Path` is the path at which we want to reach in the backend and the method is the method of the RPC call done to the backend.
```go
// Build the payload
thisPayload := types.Payload{
    Data:            config.Payload,
    Method:          "POST",
    Path:            config.Path,
    TimeoutMillisec: 10000,
}
```
Now that we have everything we need, we can actually send the relay:
```go
// Send the relay
response, err := shannon.SendRelay(thisPayload, selectedEndpoint, types.ServiceID(config.ServiceID), *FullNode, signerApp)
if err != nil {
    log.Fatal(fmt.Printf("relay: error sending relay for service %s endpoint %s: %w",
        config.ServiceID, selectedEndpoint.Addr(), err,
    ))
}
```
This function takes the endpoint, payload and signer, creates the Shannon-compatible RPC request and signs it using your app's private key. Then it just performs an RPC call to the selected endpoint URL. Before returning, this function will also validate the response and make sure it is signed by the selected relay-miner node, if this fails it means that we cannot be sure who is in the other side and cannot leverage their stake for trust.

The response from the relay miner is a serialized `http.Response` struct, which need to be deserialized to access the data of the backend:
```go
relayResponse, err := shannon.DeserializeRelayResponse(response.Payload)
if err != nil {
    log.Fatal(fmt.Printf("relay: error unmarshalling endpoint response into a POKTHTTP response for service %s endpoint %s: %w",
        config.ServiceID, selectedEndpoint.Addr(), err,
    ))
}
relayResponse.EndpointAddr = selectedEndpoint.Addr()
```
This deserialized data will contain the status of the backend in `relayResponse.HTTPStatusCode` and the payload in `relayResponse.Bytes`.

Congratulations, now you know how to send your own relays!
Now it is your turn to add what makes your gateway unique...
