# CloudFlare DDNS

## GetStarted

### Manual

1. Download `ddns_client` in [Last Release](https://github.com/dingdangdog/cloudflare_ddns/releases), or [build ddns_client](#Package_ddns_client) yourself;
2. Download `config.json` in [Last Release](https://github.com/dingdangdog/cloudflare_ddns/releases);
3. Please fill in your configuration information in `config.json`, remember to delete the comments;
4. Put `ddns_client` and `config.json` in the same folder;
5. Set `ddns_client` as an executable file and run.

## ddns

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

## ip

A simple service that gets the client's public IP and returns it.

> deploy on the server, `ddns client` requests this interface of the server, obtains the client's real IP and returns it.

### Package whoisme_server

> **Tip**: In the `Windows` development environment, please execute the command in the `CMD` terminal command line, not in `PowerShell`, otherwise the packaged result may not be executed normally

```bash
# Packaging Linux binary files in Windows development environment
cd ip
SET GOOS=linux
SET GOARCH=amd64

go build -o whoisme_server whoisme.go
```
