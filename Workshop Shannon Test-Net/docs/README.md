# The Docs, sweet docs...

IMPORTANT: All the documentation is implemented on the [docker compose example](https://github.com/pokt-network/poktroll-docker-compose-example) repo, so clone it before reading any docs:
`git clone https://github.com/pokt-network/poktroll-docker-compose-example.git`

The documentation is divided in some main tasks, you can read them in any order but we encourage you to do it in the following way:
1. [Full Node Setup](./Full%20Node%20Setup.md): Shows you how to set up a full node. You will need this to follow any other doc.
2. [Account Generation](./Account%20Generation.md): Presents how to create and fund an account in the test-net. You will need to do this many times.
3. [Service Creation](./Service%20Creation.md): How to create your own service. This guide will use a custom service, it is completelly optional if you want to use an existing service such as an ETH node (service ID : `F00C`).
4. [Supplier and Relay Miner](./Supplier%20and%20Relay%20Miner.md): How to setup a servicer account and deploy a relay miner pointing to your backend.
5. [Applications and Gateways](./Applications%20and%20Gateways.md): Here we show how to stake an application and a gateway and then how to delegate the application to the gateway so you can let it sign for you (awsome right?).
6. [Send a Relay with PATH](./Sending%20a%20Relay%20using%20PATH.md): With everything in place, we show you how to send a relay using PATH in delegated or centralized mode, super easy.
7. [Send a Relay using the SDK](./Sending%20a%20Relay%20using%20the%20Shannon%20SDK.md): If you like to get down to the core, here we present a really bare-bones example on how to use the Shannon SDK to send a relay. It is the same process that happens in PATH under the hood.

Happy coding!

Also, you can find the presentation slides [here](./Presentation%20Slides.pdf)