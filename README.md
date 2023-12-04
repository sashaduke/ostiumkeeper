# Ethereum Price Feed Keeper/Relayer System for Ostium

### Usage Instructions
- First make sure you have the latest versions of Go and Redis installed on your system (links below)
- Next run `redis-server` and check it is running on port `6379` (this is the default port)
- Finally, `cd` to the directory containing this repository and execute `go run .`
- The system should start pulling GBP/USD prices from the WebSocket feed and broadcasting them to the blockchain
- The contract state and transactions can be viewed on blockchain explorers - a link has been provided below

### Viewing the Result on the Blockchain
- Etherscan link to storage contract https://sepolia.etherscan.io/address/0x48eb2302cfec7049820b66fc91955c5d250b3ff9
- To view the contract state change, click on a recent 'Store' transaction, then switch to the 'State' tab at the top
- The middle of the three addresses should have a drop down arrow, which will reveal the state change when clicked
- Change the outputs to text and you will see both the original GBP/USD rate prior to that transaction, as well as the updated price

### Other Resources
- Go installation docs: https://go.dev/doc/install
- Redis installation docs: https://redis.io/docs/install/install-redis
