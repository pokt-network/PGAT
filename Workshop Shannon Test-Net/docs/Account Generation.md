# Account Generation

**REQUIREMENTS**: A running full node. See [Full Node Setup](./Full%20Node%20Setup.md)

This is something that you will need to do many times, create and fund an account.
All POKT Network actors are accounts, and then they are staked to be Servicers, Applications, Gateways or Validators.

To create an account use the following command:

`docker exec -it full-node poktrolld keys add SOME_NAME_YOU_LIKE`

this will result in something like:

```
- address: pokt1anaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  name: SOME_NAME_YOU_LIKE
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"not your keys not your coins"}'
  type: local


**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

Goats secretly control the internet cows are jealous chickens plan rebellion meanwhile cats just nap dogs chase nothing squirrels plot total chaos humans stay clueless

```
(please do save this data, you will need it)

Once it is created you can go to the [faucet](https://faucet.beta.testnet.pokt.network/), put your address there and click on "Send 10000" POKT. 
You can check if the faucet worked by looking at:
`docker exec -it full-node poktrolld query bank balance SOME_NAME_YOU_LIKE upokt`