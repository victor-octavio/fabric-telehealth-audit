# fabric-telehealth-audit

A Hyperledger Fabric-based prototype for auditability of diagnostics in emergency telemedicine scenarios.

## Prerequisites

- Docker
- jq
- Go 1.21+

## Starting the local network
```bash
./network/friendly-microfab-setup.sh
source network/fabric-env.sh
peer channel list
```
