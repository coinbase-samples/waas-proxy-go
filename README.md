# waas-proxy-go

Package providing an example server to utilize with Coinbase [Wallet as a Service](https://www.coinbase.com/cloud/products/waas). It is built with Golang 1.20, 
and runs a web server to communicate with a mobile app.

For more information about [Coinbase WaaS](https://docs.cloud.coinbase.com/waas/docs/welcome)

## Warning
This is a sample reference implementation, and as such is not built to be fully production-ready. 
Do not directly use this code in a customer facing application without subjecting it to significant load testing and a security review.

## Running waas-proxy-go

### Setup
From the [Coinbase Cloud Console](https://cloud.coinbase.com/) create a new waas project, download the api key, and save the key file.

Copy the sample.env to .env and fill in the values.

### Start Application
To start up the local server, run the following:
```
make start-local
```

This does the following:
* Starts the application server - default port is 8443

## License
This library is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file.

