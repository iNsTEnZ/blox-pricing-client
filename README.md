# Token Pricing Client

How to use
```
go build
```
```
go run blox-client-service
```

Environment variables
```
PRICING_ADDRESS (default: localhost:8080) - defines the address of the token pricing server.
MAX_CACHE_SIZE (default: 2000) - defines the maximu allowed token data to be stored.
FETCH_INTERVAL_SECONDS (default: 60) - defines, in seconds, what would be the interval for fetching token prices.
REQUEST_TIMEOUT_SECONDS (default: 10) - defines, in seconds,  the timeout for fetching the data for all tokens.
CURRENCY (default: USD) - defines the currency for the token prices.
TOKENS (default: BTC,ETH,BNB) - a list of token the client would fetch pricing for.
```