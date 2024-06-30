# cryptodashboard

This is a simple cryptocurrency dashboard that displays the current price of ETH, BTC and ADA in USD.
The prices are streamed in using the CoinGecko API and SSE (Server-Sent Events).

## Building the project

To build the project, run the following commands:

```bash
make
```

## Running the project

```bash
make run
```

## Configuration

The project uses a configuration file located at `config.yaml`. The configuration file contains the following fields:

```yaml
  coingecko_api_key: "your-api-key"
  update_interval_in_sec: 5
  target_currency: "usd"
  crypto_currency_ids:
    - ethereum
    - bitcoin
    - cardano
  host: "127.0.0.1"
  port: 8080
```

The configuration can also be loaded through environment variables. The names are just uppercase versions of the fields in the configuration file.
