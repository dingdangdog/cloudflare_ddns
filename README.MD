# CloudFalre DDNS

## ddns

想办法拿到真实 `IP`，拿到 `IP` 后调用 `Cloudflare API` 修改 `DNS` 记录。

### 打包ddns

> **提示**：在 `Windows` 开发环境中，请使用 `CMD` 终端命令行执行命令，不要在 `PowerShell` 中执行，否则打包的结果可能无法正常执行

```bash
# 打包Linux二进制文件，在Windows开发环境下的打包方式
cd ddns
SET GOOS=linux
SET GOARCH=amd64

go build -o ddns_client ddns.go
```

## ip

一个简单的获取客户端公网IP并返回的服务

> 用于部署到服务器上，`ddns client` 请求服务器的这个接口，获取客户端真实IP并返回。

### 打包ip服务

> **提示**：在 `Windows` 开发环境中，请在 `CMD` 终端命令行执行命令，不要在 `PowerShell` 中执行，否则打包的结果可能无法正常执行

```bash
# 打包Linux二进制文件，在Windows开发环境下的打包方式
cd ip
SET GOOS=linux
SET GOARCH=amd64

go build -o whoisme_server whoisme.go
```
