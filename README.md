# AV is AirGo's V2bX

# 1、直接安装

```
bash <(curl -Ls https://raw.githubusercontent.com/ppoonk/V2bX/main/scripts/install.sh)
```

# 2、docker 安装

## 2-1
提前准备配置文件，参考项目根目录下 `config.json`

## 2-2
- 启动docker命令参考如下：

```
docker run -tid \
  -v $PWD/av/config.json:/etc/V2bX/config.json \
  --name airgo \
  --restart always \
  --net=host \
  --privileged=true \
  ppoiuty/av:latest
```

- docker compose参考如下：

```
version: '3'
services:
  AV:
    container_name: AV
    image: ppoiuty/av:latest
    network_mode: "host"
    restart: "always"
    privileged: true
    volumes:
      - ./config.json:/etc/V2bX/config.json
```
