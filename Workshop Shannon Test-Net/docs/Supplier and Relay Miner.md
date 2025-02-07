# Supply Set-Up: Supplier and Relay Miner

**REQUIREMENTS**: A running full node. See [Full Node Setup](./Full%20Node%20Setup.md)

You will need to have an account to run a Servicer, so lets create one, execute:
`docker exec -it full-node poktrolld keys add beta-relayminer-rawthil`
and fund it (see [Account Generation](./Account%20Generation.md))

Edit the file `stake_configs/supplier_stake_config_example.yaml` and set:
- `owner_address: <REPLACE WITH YOUR SUPPLIER ADDRESS>` 
- `service_id: <THE ID OF YOUR SERVICE>`
- `publicly_exposed_url: <YOUR PUBLIC IP/URL>:8545`

Then stake the supplier account
```bash
docker exec -it full-node poktrolld tx supplier stake-supplier \
  --config=/poktroll/stake_configs/supplier_stake_config_example.yaml \
  --from=beta-relayminer-rawthil \
  --gas=auto \
  --gas-prices=1upokt \
  --gas-adjustment=1.5 \
  --chain-id=pocket-beta \
  --yes
```


### Starting the Relay Miner

Using the data displayed when you created the account, complete the `.env` file (replace with your own):
- `SUPPLIER_MNEMONIC="<REPLACE WITH THOSE 24 WORDS>"`
- `SUPPLIER_ADDR="<REPLACE WITH YOUR SUPPLIER ADDRESS>"`

Now edit the relay miner config in `relayminer/config/relayminer_config.yaml`:
- Replace the `service_id` with the ID of the service you created (if you chose to).
- Replace the `backend_url` with the url of the service in your backend. In my case, an OpenAI compatible service, it will be `http://localhost:9900/`.
- In `publicly_exposed_endpoints`, set the same public IP/URL you used during the stake, in `stake_configs/supplier_stake_config_example.yaml`.

Now make sure that your relay miner port `8545` is open for the internet, you will need it to serve requests.

Finally, start the relay miner: 
`docker compose up -d relayminer`
(you can track logs with `docker compose logs -f --tail 100 relayminer`)


##### Example files

`stake_configs/supplier_stake_config_example.yaml`:
```yaml
stake_amount: 1000000upokt
owner_address: pokt1r90ujjku55rldjpxsuwx0s2cg7yp5uphxnaa5l
services:
  - service_id: A100
    endpoints:
      - publicly_exposed_url: http://rawthil.zapto.org:8545
        rpc_type: json_rpc

```

`relayminer/config/relayminer_config.yaml`:
```yaml
default_signing_key_names:
  - supplier
smt_store_path: /home/pocket/.poktroll/smt
pocket_node:
  query_node_rpc_url: tcp://full-node:26657
  query_node_grpc_url: tcp://full-node:9090
  tx_node_rpc_url: tcp://full-node:26657
suppliers:
  - service_id: "A100"
    service_config:
      backend_url: "http://localhost:9900/"
      publicly_exposed_endpoints:
        - rawthil.zapto.org
    listen_url: http://0.0.0.0:8545 
metrics:
  enabled: true
  addr: :9090
pprof:
  enabled: false
  addr: :6060

```