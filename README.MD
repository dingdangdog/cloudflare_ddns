# CloudFlare DDNS

Simple script to quickly implement DDNS based on CloudFlare.

[简体中文](./README_ZH.MD) [English]

## GetStarted

### Auto

> 🧨**Tip**: Whether you install it automatically or manually, please modify the configuration `config.json` yourself! If you use other IP acquisition services, please modify the source code yourself (if you don't need authorization, you can try not to modify the source code)

- ddns_client

```sh
curl -sSL https://github.com/dingdangdog/cloudflare_ddns/raw/main/script/update_ddns.sh | bash
```

Modify the `config_demo.json` and run `ddns_client`.

- whoiam_server

```sh
curl -sSL https://github.com/dingdangdog/cloudflare_ddns/raw/main/script/update_whoiam.sh | bash
```

Modify the `config_demo.json` and run `whoiam_server`.

### Manual

1. Download `ddns_client` or `whoiam_server` from [Last Release](https://github.com/dingdangdog/cloudflare_ddns/releases), or [Package ddns_client](#Package_ddns_client) / [Package whoiam_server](#Package_whoiam_server) by yourself;
2. Download ddns_client [config_demo.json](./ddns/config demo.json) or whoiam_server [config_demo.json](./ip/config_demo.json) to your device, rename it to `config.json`, and then edit the configuration information。
3. Run `ddns_client` or `whoiam_server`。

## ddns

- [ddns.go](./ddns/ddns.go)

Find a way to get the real `IP`, and after getting the `IP`, call the `Cloudflare API` to modify the `DNS` record.

### Package_ddns_client

> **Tip**: In the `Windows` development environment, please execute the command in the `CMD` terminal command line, not in `PowerShell`, otherwise the packaged result may not be executed normally

```bash
# Packaging Linux binary files in Windows development environment
cd ddns
SET GOOS=linux
SET GOARCH=amd64

go build -o ddns_client ddns.go
```

### Config Options

- [ddns/config_demo.json](./ddns/config_demo.json)

```json
{
  "CLOUDFLARE": {
    "CF_API_TOKEN": "xx", // Please get your Cloudflare configuration
    "CF_ZONE_ID": "xxx", // Please get your Cloudflare configuration
    "CF_RECORD_ID": "xxx", // Please get your Cloudflare configuration
    "DNS_TYPE": "A",
    "DNS_DOMAIN_NAME": "xxx.com", // DNS domain name to be set
    "DNS_DOMAIN_CONTENT": "xxx.com", // Temporarily unused
    "DNS_TTL": 1, // Fill in 1 to use the default value of cloudflare
    "DNS_PROXIED": false
  },
  "IP_API_URL": "http://xxx.com/whoiam", // Interface to obtain the real public IP
  "DDD_CLIENT_ID": 0, // Whoisme interface authorization id
  "DDD_CLIENT_KEY": "Test_Cilent_Key", // Whoisme interface authorization key
  "MODE": "development",
  "INTERVAL": 180
}
```

## ip

- [whoiam.go](./ip/whoiam.go)

A simple service that gets the client's public IP and returns it.

> deploy on the server, `ddns client` requests this interface of the server, obtains the client's real IP and returns it.

### Dev Run

```sh
cd ip
go run whoiam.go
```

### Package_whoiam_server

> **Tip**: In the `Windows` development environment, please execute the command in the `CMD` terminal command line, not in `PowerShell`, otherwise the packaged result may not be executed normally

```bash
# Packaging Linux binary files in Windows development environment
cd ip
SET GOOS=linux
SET GOARCH=amd64

go build -o whoiam_server whoiam.go
```
