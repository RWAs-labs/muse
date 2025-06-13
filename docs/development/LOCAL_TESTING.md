# MuseChain Localnet Development & Testing Environment

MuseChain localnet development and testing environment is divided into three different directories:

- [localnet](../../contrib/localnet/README.md): a set of Docker images and script for spinning up a localnet.
- [e2e](../../e2e/README.md): a set of Go packages for end-to-end testing between MuseChain and other blockchains.
- [musee2e](../../cmd/musee2e/README.md): a CLI tool using `e2e` for running end-to-end tests.

A description of each directory is provided in the respective README files.

## Running E2E tests

### Build musenode
```bash
$ make musenode
```

This Makefile rule builds the musenode image. **Rebuild if musecored/museclientd code is updated**.
```bash
$ docker build -t musenode .
```

### Run e2e tests

Now we have built all the docker images, we can run the e2e test with make command:
```bash
make start-e2e-test
```

This uses `docker compose` to start the localnet and run standard e2e tests inside the orchestrator container. There are several parameters that the `Makefile` can provide to `docker compose` via environment variables:

- `LOCALNET_MODE`
	- `setup-only`: only setup the localnet, do not run the e2e tests
	- `upgrade`: run the upgrade tests
	- unset: run the e2e tests
- `E2E_ARGS`: arguments to provide to the `musee2e local` command
    - `--verbose` will give you verbose logs
	- `--test-filter` allows you to filter which tests to run by regular expression
- `UPGRADE_HEIGHT`: block height to upgrade at when `LOCALNET_MODE=upgrade`
- `MUSECORED_IMPORT_GENESIS_DATA`: path to genesis data to import before starting musecored
- `MUSECORED_START_PERIOD`: duration to tolerate `musecored` health check failures during startup

Here is an example of using `E2E_ARGS`:

```
export E2E_ARGS='--verbose --test-filter eth_deposit|eth_withdraw'
make start-e2e-test
```

This will execute `musee2e local` as:

```
musee2e local --verbose --test-filter eth_deposit|eth_withdraw --config config.yml
```

More options directly to `docker compose` via the `NODE_COMPOSE_ARGS` variable. This allows setting additional profiles or configuring an overlay. Example:

```
export NODE_COMPOSE_ARGS="--profile monitoring -f docker-compose-persistent.yml"
make start-e2e-test
```

This starts the e2e tests while enabling the monitoring stack and persistence (data is not deleted between test runs).

#### Run admin functions e2e tests

We define e2e tests allowing to test admin functionalities (emergency network pause for example).
Since these tests interact with the network functionalities, these can't be run concurrently with the regular e2e tests.
Moreover, these tests test scoped functionalities of the protocol, and won't be tested in the same pipeline as the regular e2e tests.
Therefore, we provide a separate command to run e2e admin functions tests:
```bash
make start-e2e-admin-test
```

### Run upgrade tests

Upgrade tests run the E2E tests with an older version, upgrade the nodes to the new version, and run the E2E tests again.
This allows testing the upgrade process with a populated state.

Before running the upgrade tests, the old version must be specified the Makefile.

NOTE: We only specify the major version for `NEW_VERSION` since we use major version only for chain upgrade. Semver is needed for `OLD_VERSION` because we use this value to fetch the release tag from the GitHub repository.

The upgrade tests can be run with the following command:
```bash
make start-upgrade-test
```

The test the upgrade script faster a light version of the upgrade test can be run with the following command:
```bash
make start-upgrade-test-light
```
This command will run the upgrade test with a lower height and will not populate the state.

### Run stress tests

Stress tests run the E2E tests with a larger number of nodes and clients to test the performance of the network.
It also stresses the network by sending a large number of cross-chain transactions.

The stress tests can be run with the following command:
```bash
make start-stress-test
```

### Test logs

For all tests, the most straightforward logs to observe are the orchestrator logs.
If everything works fine, it should finish without panic.

The logs can be observed with the following command:
```bash
# in node/contrib/localnet/orchestrator
$ docker logs -f orchestrator
```

### Stop tests

To stop the tests,
```bash
make stop-localnet
```

### Run monitoring setup

Before starting the monitoring setup, make sure the Musecore API is up at http://localhost:1317.
You can also add any additional ETH addresses to monitor in `muse-node/contrib/localnet/grafana/addresses.txt` file

```bash
make start-monitoring
```

### Grafana credentials and dashboards

The Grafana default credentials are admin:admin. The dashboards are located at http://localhost:3000.

### Stop monitoring setup

```bash
make stop-monitoring
```

## Interacting with the Localnet

In addition to running automated tests, you can also interact with the localnet directly for more specific testing.

The localnet can be started without running tests with the following command:

```bash
make start-localnet
```

The localnet takes a few minutes to start. Printing the logs of the orchestrator will show when the localnet is ready. Once setup, it will display:
```
✅ the localnet has been setup
```

### Interaction with MuseChain

MuseChain
The user can connect to the `musecore0` and directly use the node CLI with the musecored binary with a funded account:

The account is named `operator` in the keyring and has the address: `muse1amcsn7ja3608dj74xt93pcu5guffsyu2xfdcyp`

```bash
docker exec -it musecore0 sh
```

Performing a query:

```bash
musecored q bank balances muse1amcsn7ja3608dj74xt93pcu5guffsyu2xfdcyp
```

Sending a transaction:

```bash
musecored tx bank send operator muse172uf5cwptuhllf6n4qsncd9v6xh59waxnu83kq 5000amuse --from operator --fees 2000000000000000amuse
```

### Interaction with EVM

The user can interact with the local Ethereum node with the exposed RPC on `http://0.0.0.0:8545`. The following testing account is funded:

```
Address: 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC
Private key: d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263
```

Examples with the [cast](https://book.getfoundry.sh/cast/) CLI:

```bash
cast balance 0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC --rpc-url http://0.0.0.0:8545
98897999997945970464

cast send 0x9fd96203f7b22bCF72d9DCb40ff98302376cE09c --value 42 --rpc-url http://0.0.0.0:8545 --private-key "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
```

### Interaction using `musee2e`

`musee2e` CLI can also be used to interact with the localnet and test specific functionalities with the `run` command. The [local config](../../cmd/musee2e/config/local.yml) can be used to interact with the network. 

For local interactions, the Bitcoin account is not automatically setup. To set it up, the following command can be used:

```bash
musee2e setup-bitcoin cmd/musee2e/config/local.yml
```

The balances on the localnet can be checked with the following command:

```bash
musee2e balances cmd/musee2e/config/local.yml
```

Example of `run` command:

```dockerfile
musee2e run muse_deposit:2000000000000000000 eth_deposit:2000000000000000000 erc20_deposit:200000 --config cmd/musee2e/config/local.yml
```

## Useful data

- TSS Address (on ETH): 0xF421292cb0d3c97b90EEEADfcD660B893592c6A2

## Add more e2e tests

The e2e tests are located in the e2e/e2etests package. New tests can be added. The process:

1. Add a new test file in the e2e/e2etests package, the `test_` prefix should be used for the file name.
2. Implement a method that satisfies the interface:
```go
type E2ETestFunc func(*E2ERunner)
```
3. Add the test to list in the `e2e/e2etests/e2etests.go` file.

The test can interact with the different networks using the runned object:
```go
type E2ERunner struct {
	MEVMClient   *ethclient.Client
	EVMClient *ethclient.Client
	BtcRPCClient *rpcclient.Client

	CctxClient     crosschaintypes.QueryClient
	FungibleClient fungibletypes.QueryClient
	AuthClient     authtypes.QueryClient
	BankClient     banktypes.QueryClient
	ObserverClient observertypes.QueryClient
	MuseTxServer   txserver.MuseTxServer
	
	EVMAuth *bind.TransactOpts
	MEVMAuth   *bind.TransactOpts
	
	// ...
}
```

## Localnet Governance Proposals

Localnet can be used for testing the creation and execution of governance propoosals.

Exec into the `musecored0` docker container and run the script to automatically generate proposals in a variety of states and then extends the voting window to one hour, allowing you time to view a proposal in a pending state.
```
docker exec -it musecore0 bash
/root/test-gov-proposals.sh
```
