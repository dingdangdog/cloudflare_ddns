# CloudFalre DDNS

简单的脚本快速实现基于Cloudflare的DDNS。

[简体中文] [English](./README.MD)

安全前，建议阅读 [服务介绍](#服务介绍) 了解两个模块(`ddns & whoiam`)的作用。

## 安装

> 🧨**提示**：无论自动还是手动安装，请自行修改配置 `config.json` ！如果使用其他获取IP服务，请自行修改源码（如果无需授权等，可以尝试不修改源码）

**安装方式**：

1. [自动安装脚本](#自动安装)（`ddns`/`whoiam`）；
2. [手动下载](#手动安装)（`ddns`/`whoiam`）；
3. 源码安装（`ddns`/`whoiam`）；这个就不介绍了，懂代码的一看就会了。
4. [Docker](#docker)（`ddns`/`whoiam`）；
5. [Cloudflare Workers](#cloudflareworkers)（`whoiam`）；

### 自动安装

- ddns_client

```sh
curl -sSL https://github.com/dingdangdog/cloudflare_ddns/raw/main/script/update_ddns.sh | bash
```

修改配置并运行 `ddns_client`。

- whoiam_server

```sh
curl -sSL https://github.com/dingdangdog/cloudflare_ddns/raw/main/script/update_whoiam.sh | bash
```

修改配置并运行 `whoiam_server`。

### 手动安装

1. 在 [Last Release](https://github.com/dingdangdog/cloudflare_ddns/releases) 中下载 `ddns_client` 或 `whoiam_server`，或者自行 [打包ddns_client](#打包ddns_client) / [打包whoiam_server](#打包whoiam_server)；
2. 下载 ddns_client [配置示例](./ddns/config_demo.json) 或者 whoiam_server [配置示例](./whoiam/config_demo.json) 到你的设备，重命名为 `config.json`，然后编辑你的配置信息；
3. 运行 `ddns_client` 或 `whoiam_server`。

### Docker

可以阅读对应的 `docker-compose.yaml` 文件内容，自行选择使用 `docker-compose` 还是 `docker run` 部署，这里就不再详细介绍了，请各位自行学习 `Docker` 吧！

- [ddns: docker-compose.yaml](./ddns/docker-compose.yaml)
- [whoiam: docker-compose.yaml](./whoiam/docker-compose.yaml)

### CloudflareWorkers

1. 在 `CloudFlare` 控制台中创建 `Worker`；
2. 复制 [whoiam.js](./whoiam/whoiam.js) 的全部代码，填入 `worker`；
3. 启动 `worker`；
4. 配置环境变量 `CLIENT_KEYS`，可以配置多个，使用 `,` 间隔！
5. 配置完成后点击 `部署` ，重新部署 `worker`。

## 服务介绍

### ddns

- [ddns.go](./ddns/ddns.go)

功能：想办法拿到真实 `IP`，拿到 `IP` 后调用 `Cloudflare API` 修改 `DNS` 记录。

用法：部署到你的动态公网IP主机上，修改 [配置信息](./ddns/config_demo.json) 运行即可。

#### 打包ddns_client

> **提示**：下面是在 `Windows` 开发环境中的示例，请使用 `CMD` 终端命令行执行命令，不要在 `PowerShell` 中执行，否则打包的结果可能无法正常执行

```bash
# 打包Linux二进制文件，在Windows开发环境下的打包方式
cd ddns
SET GOOS=linux
SET GOARCH=amd64

go build -o ddns_client ddns.go
```

#### 配置参数说明

- [config_demo.json](./ddns/config_demo.json)

```json
{
  "CLOUDFLARE": {
    "CF_API_TOKEN": "xx", // 请自行获取你的 Cloudflare 配置
    "CF_ZONE_ID": "xxx", // 请自行获取你的 Cloudflare 配置
    "CF_RECORD_ID": "xxx", // 请自行获取你的 Cloudflare 配置
    "DNS_TYPE": "A",
    "DNS_DOMAIN_NAME": "xxx.com", // 要设置的 DNS 域名
    "DNS_DOMAIN_CONTENT": "xxx.com", // 暂时没用
    "DNS_TTL": 1, // 填1会使用cloudflare默认值
    "DNS_PROXIED": false
  },
  "WHOIAM_API_URL": "http://xxx.com/whoiam", // 获取 真实公网IP 的接口
  "WHOIAM_CLIENT_ID": 0, // whoiam接口授权id
  "WHOIAM_CLIENT_KEY": "Test_Cilent_Key", // whoiam接口授权key
  "MODE": "development",
  "INTERVAL": 180
}
```

### whoiam

- [whoiam.go](./whoiam/whoiam.go)

功能：一个简单的获取客户端公网IP并返回的服务。

用法：部署到你的另一台公网服务器上，使用 `ddns_client` 发起请求到 `whoiam` 服务，获取到 `ddns_client` 所在主机的真实公网IP。

#### 打包whoiam_server

> **提示**：下面是在 `Windows` 开发环境中的示例，请在 `CMD` 终端命令行执行命令，不要在 `PowerShell` 中执行，否则打包的结果可能无法正常执行

```bash
# 打包Linux二进制文件，在Windows开发环境下的打包方式
cd whoiam
SET GOOS=linux
SET GOARCH=amd64

go build -o whoiam_server whoiam.go
```

### 公开免费的服务

如果不想自己部署 `whoiam`，可以使用公开免费的接口，但是注意：

1. 保证接口能够被你的服务器访问（部分接口可能被中国大陆墙）；
2. 保证接口返回值只有IP！如果想要使用其他格式，请自行修改代码！
3. 如果返回IPV6的IP，则无法使用！

| **接口地址**                                    | **使用限制**                  | **备注**                           |
| ------------------------------------------- | ----------------------------- | ---------------------------------- |
| <https://api.ipify.org>                       | 无限制，免费                  | 简单易用|
| <https://icanhazip.com>                       | 无限制，免费                  | 由 Cloudflare 运营，纯文本返回     |
| <https://api.seeip.org>                       | 无限制，免费                  | 支持 IPv4/IPv6，开放源代码         |
