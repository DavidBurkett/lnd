## Lightning Network Daemon for Litecoin

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ltcsuite/lnd/blob/master/LICENSE)
[![Godoc](https://godoc.org/github.com/ltcsuite/lnd?status.svg)](https://godoc.org/github.com/ltcsuite/lnd)

The Lightning Network Daemon for Litecoin (`lndltc`) - is a complete implementation of a
[Lightning Network](https://lightning.network) node.

It is adapted from the
[`lnd`](https://github.com/ltcsuite/lnd) project by Lightning Labs. This fork
focuses on maintaining complete support and bugfixes related to Litecoin. Users should create issues in this repo if
they are to find issues specific to Litecoin.

`lndltc` has several pluggable back-end
chain services including [`ltcd`](https://github.com/ltcsuite/ltcd) (a
full-node), [`litecoind`](https://github.com/litecoin-project/litecoin), and
[`neutrino`](https://github.com/ltcsuite/neutrino) (a new experimental light client). The project's codebase uses the
[ltcsuite](https://github.com/ltcsuite/) set of Litecoin libraries, and also
exports a large set of isolated re-usable Lightning Network related libraries
within it. In the current state `lndltc` is capable of:

- Creating channels.
- Closing channels.
- Completely managing all channel states (including the exceptional ones!).
- Maintaining a fully authenticated+validated channel graph.
- Performing path finding within the network, passively forwarding incoming payments.
- Sending outgoing [onion-encrypted payments](https://github.com/ltcsuite/lightning-onion)
  through the network.
- Updating advertised fee schedules.
- Automatic channel management ([`autopilot`](https://github.com/ltcsuite/lnd/tree/master/autopilot)).

## Lightning Network Specification Compliance

`lndltc` _fully_ conforms to the [Lightning Network specification
(BOLTs)](https://github.com/lightningnetwork/lightning-rfc). BOLT stands for:
Basis of Lightning Technology. The specifications are currently being drafted
by several groups of implementers based around the world including the
developers of `lndltc`. The set of specification documents as well as our
implementation of the specification are still a work-in-progress. With that
said, the current status of `lndltc`'s BOLT compliance is:

- [x] BOLT 1: Base Protocol
- [x] BOLT 2: Peer Protocol for Channel Management
- [x] BOLT 3: Bitcoin Transaction and Script Formats
- [x] BOLT 4: Onion Routing Protocol
- [x] BOLT 5: Recommendations for On-chain Transaction Handling
- [x] BOLT 7: P2P Node and Channel Discovery
- [x] BOLT 8: Encrypted and Authenticated Transport
- [x] BOLT 9: Assigned Feature Flags
- [x] BOLT 10: DNS Bootstrap and Assisted Node Location
- [x] BOLT 11: Invoice Protocol for Lightning Payments

## Developer Resources

The daemon has been designed to be as developer friendly as possible in order
to facilitate application development on top of `lndltc`. Two primary RPC
interfaces are exported: an HTTP REST API, and a [gRPC](https://grpc.io/)
service. The exported API's are not yet stable, so be warned: they may change
drastically in the near future.

An automatically generated set of documentation for the RPC APIs can be found
at [api.lightning.community](https://api.lightning.community). A set of developer
resources including guides, articles, example applications and community resources can be found at:
[docs.lightning.engineering](https://docs.lightning.engineering).

Finally, we also have an active
[Slack](https://lightning.engineering/slack.html) where protocol developers, application developers, testers and users gather to
discuss various aspects of `lndltc` and also Lightning in general. Discussions specific to the development of `lndltc` is done
in the #ltc channel.

## Installation

In order to build from source, please see [the installation
instructions](docs/INSTALL.md).

## Docker

To run lnd from Docker, please see the main [Docker instructions](docs/DOCKER.md)

## Safety

When operating a mainnet `lnd` node, please refer to our [operational safety
guidelines](docs/safety.md). It is important to note that `lnd` is still
**beta** software and that ignoring these operational guidelines can lead to
loss of funds.

## Security

The developers of `lnd` and `lndltc` take security _very_ seriously. The disclosure of
security vulnerabilities helps us secure the health of `lndltc`, privacy of our
users, and also the health of the Lightning Network as a whole. If you find
any issues regarding security or privacy, please disclose the information
responsibly by sending an email to loshan1212 at gmail dot com,
preferably encrypted using our designated PGP key
(`C0921846FED0BF4CF28BE1D73B2A6315CD51A673`) which can be found
[here](https://raw.githubusercontent.com/ltcsuite/lnd/master/scripts/keys/losh11.asc).

## Further reading

- [Step-by-step send payment guide with docker](https://github.com/ltcsuite/lnd/tree/master/docker)
- [Contribution guide](https://github.com/ltcsuite/lnd/blob/master/docs/code_contribution_guidelines.md)
