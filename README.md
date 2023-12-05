# Ethereum Price Feed Keeper/Relayer System for Ostium

### Usage Instructions
- First make sure you have the latest versions of Go and Redis installed on your system (links below)
- Next run `redis-server` and check it is running on port `6379` (this is the default port)
- Finally, `cd` to the directory containing this repository and execute `go run .`
- Once running, the system should start pulling GBP/USD prices from the WebSocket feed and periodically broadcasting them to the blockchain automatically
- The contract state and transactions can be viewed on blockchain explorers - a link and explanation has been provided below

### Endpoints
- There are two REST API endpoints which will return the latest data cached in the database (/data) and the smart contract (/contracts)
- Once you have followed the above setup instructions and the daemon is running, you will be able to visit these URLs in your browser or use something like `curl` to query the endpoints:
- `curl http://localhost:8080/data`
- `curl http://localhost:8080/contracts`

### Viewing the Result on the Blockchain
- Etherscan link to storage contract https://sepolia.etherscan.io/address/0x48eb2302cfec7049820b66fc91955c5d250b3ff9
- To view the contract state change, click on a recent 'Store' transaction to the contract, then switch to the 'State' tab at the top next to 'Overview'
- The middle of the three addresses should have a drop down arrow on the left, which will reveal the state change when clicked
- Change the outputs from hex to text and you will see both the original GBP/USD rate prior to that transaction, as well as the new price it was updated to

### Environment Variables
- The system should work straight out of the box using its default configuration settings, but there are some environment variables that can be changed to update various parameters. The commands to modify these are as follows:
- `export CONTRACT_ADDR="48eB2302cfEc7049820b66FC91955C5d250b3fF9"` - storage smart contract address
- `export RPC_API_KEY="yourNodeProviderAPIToken"` - replace with your RPC provider API token (i.e. Infura)
- `export RPC_ENDPOINT="https://sepolia.infura.io/v3/"` - blockchain node RPC endpoint
- `export PRIVKEY_HEX="yourPrivateKeyHex"` - replace with your hex-encoded ECDSA private key
- `export WS_URL="wss://api.tiingo.com/fx"` - price feed WebSocket URL
- `export WS_API_KEY="yourPriceFeedAPIKey"` - replace with your API token for price feed
- `export FX_PAIR="gbpusd"` - can be changed to another FX pair i.e. eurusd
- `export WS_TIME_LAYOUT="2006-01-02T15:04:05.000000-07:00"` - timestamp layout for decoding responses
- `export CONTRACT_WRITE_FREQ="15"` (seconds) - frequency at which updates are written to contract 
- `export LOG_FILE="daemon_log.txt"` - file name for log output file

### Other Resources
- Go installation docs: https://go.dev/doc/install
- Redis installation docs: https://redis.io/docs/install/install-redis
