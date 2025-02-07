# Full Node Setup

### Set Environment
`cp .env.sample .env`
Add  the address of your machine to the `NODE_HOSTNAME` entry

If you want to participate in the network, not only sync, you will need these ports open:
  - 26657
  - 26656
  - 9090

### Spin up the Full Node

Start your full node with
`docker compose up -d full-node`
(you can track logs with `docker compose logs -f --tail 100 full-node`)

Check the height using: `curl -s -X POST localhost:26657/status | jq`, and look for the field `latest_block_height`.
Or just watch it go up using `watch -n 1 "curl -s -X POST localhost:26657/block | jq '.result.block.header.height'"`.

This will take a while, and use more than 20GB of disk. At some point the height will be the same as the one observed in https://shannon.alpha.testnet.pokt.network/poktroll
Block `58120` will take HOURS to sync, it is huge, it was a load testing...
