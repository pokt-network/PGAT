# Demand Set-Up: Applications and Gateways

**REQUIREMENTS**: A running full node. See [Full Node Setup](./Full%20Node%20Setup.md)

### Application

As before, we will need to create an account:
`docker exec -it full-node poktrolld keys add beta-application-rawthil`
and fund it (see [Account Generation](./Account%20Generation.md))

Edit the transaction config in `stake_configs/application_stake_config_example.yaml`, replacing the `F00C` with your custom service if you created one. In this example it was `A100`.

Now execute the stake transaction:
```bash
docker exec -it full-node poktrolld tx application stake-application \
  --config=/poktroll/stake_configs/application_stake_config_example.yaml \
  --from=beta-application-rawthil \
  --gas=auto \
  --gas-prices=1upokt \
  --gas-adjustment=1.5 \
  --chain-id=pocket-beta \
  --yes
```

### Gateway

If you want to test the awesome feature of delegation (or use PATH), you will need to stake a Gateway also.
So, you will need to create an account to use as gateway:
`docker exec -it full-node poktrolld keys add beta-gateway-rawthil`
and fund it (see [Account Generation](./Account%20Generation.md)).

Then proceed to stake the gateway (the config only states the stake amount):
```bash
docker exec -it full-node poktrolld tx gateway stake-gateway \
  --config=/poktroll/stake_configs/gateway_stake_config_example.yaml \
  --from=beta-gateway-rawthil \
  --gas=auto \
  --gas-prices=1upokt \
  --gas-adjustment=1.5 \
  --chain-id=pocket-beta \
  --yes
```

And delegate the Application to it:

delegate app to gateway
```bash
docker exec -it full-node tx application delegate-to-gateway <GATEWAY ACCOUNT ADDRESS> \
  --from=beta-application-rawthil \
  --gas=auto \
  --gas-prices=1upokt \
  --gas-adjustment=1.5 \
  --chain-id=pocket-beta \
  --yes
```
(replace `<GATEWAY ACCOUNT ADDRESS>` with the address of the gateway account)


##### Example files

`stake_configs/application_stake_config_example.yaml`:
```yaml
stake_amount: 100000000upokt
service_ids:
  - A100
```