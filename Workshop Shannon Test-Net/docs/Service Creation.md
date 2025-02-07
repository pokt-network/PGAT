# Service Creation

**REQUIREMENTS**: A running full node. See [Full Node Setup](./Full%20Node%20Setup.md)


In case you want to provide a new service, you can create your own.
On Main-Net this will have a non-negligible cost, since we want to keep things simple.

To do this first create an account, this account will be the service owner:
`docker exec -it full-node poktrolld keys add beta-test-rawthil`
and fund it (see [Account Generation](./Account%20Generation.md))

Now that your address is funded, execute the `add-service` transaction:

```
poktrolld tx service add-service "A100" "test some AI endpoint" 17 \
    --fees 1upokt --from beta-test-rawthil --chain-id pocket-beta
```
This is an example, replace the following elements for your custom service:
- `"A100"` with any 4-digit hex number that is not a service already
- `"test some AI endpoint"` with a short description of your service.
- `17` with the amount of compute units that will cost each relay in the new service.
- `beta-test-rawthil` with the name of the account you created before.

The output of this example was:
```
auth_info:
  fee:
    amount:
    - amount: "1"
      denom: upokt
    gas_limit: "200000"
    granter: ""
    payer: ""
  signer_infos: []
  tip: null
body:
  extension_options: []
  memo: ""
  messages:
  - '@type': /poktroll.service.MsgAddService
    owner_address: pokt1fnzay6mtz26r5dtzmxsdakawkn7a6rc70k6ajv
    service:
      compute_units_per_relay: "17"
      id: A100
      name: test some AI endpoint
      owner_address: pokt1fnzay6mtz26r5dtzmxsdakawkn7a6rc70k6ajv
  non_critical_extension_options: []
  timeout_height: "0"
signatures: []
confirm transaction before signing and broadcasting [y/N]: y
code: 0
codespace: ""
data: ""
events: []
gas_used: "0"
gas_wanted: "0"
height: "0"
info: ""
logs: []
raw_log: ""
timestamp: ""
tx: null
txhash: 44B9B994297F0DB8DDA6434CF21FB6627D23F9C4CD72178EF43C376E1FEBBBC9
```
Check the tx status using:

'docker exec -it full-node poktrolld query tx --type=hash 44B9B994297F0DB8DDA6434CF21FB6627D23F9C4CD72178EF43C376E1FEBBBC9'

You can see the transaction state and the service information in POKTscan by means of the `txhash`, [this is the TX of this example](https://shannon-beta.poktscan.com/tx/44B9B994297F0DB8DDA6434CF21FB6627D23F9C4CD72178EF43C376E1FEBBBC9).
Once you see this is successful (`status` should be 0), you can move on to deploy a servicer.