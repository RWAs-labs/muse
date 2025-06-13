## Usage 

### Fetching the Inbound Ballot Identifier

### Command
```shell
musetool get-ballot [inboundHash] [chainID] --config <filename.json>
```
### Example
```shell
musetool get-ballot 0x61008d7f79b2955a15e3cb95154a80e19c7385993fd0e083ff0cbe0b0f56cb9a 1
{"level":"info","time":"2025-01-20T11:30:47-05:00","message":"ballot identifier: 0xae189ab5cd884af784835297ac43eb55deb8a7800023534c580f44ee2b3eb5ed"}
```

- `inboundHash`: The inbound hash of the transaction for which the ballot identifier is to be fetched
- `chainID`: The chain ID of the chain to which the transaction belongs
- `config`: [Optional] The path to the configuration file. When not provided, the configuration in the file is user. A sample config is provided at `cmd/musetool/config/sample_config.json`

The Config contains the rpcs needed for the tool to function,
if not provided the tool automatically uses the default rpcs.It is able to fetch the rpc needed using the chain ID

The command returns a ballot identifier for the given inbound hash.

