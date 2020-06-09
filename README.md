# Go-Web-Nes

> 建议在局域网搭建

Go + jsnes 实现的网页版NES联机、观战服务

## 原理

玩家1将canvas转换为dataUrl广播给玩家2和观战玩家，玩家2将控制命令发送给玩家1。

假设fps为60，每帧dataUrl大小为10k，则每秒需要发送`广播玩家数量 * 60 * 10kb`的数据，因此建议在局域网搭建。

## 安装

```
git clone git@github.com:xiuxiubiu/go-web-nes.git

cd ./go-web-nes

go get ./...

go build

./go-web-nes
```

## 联机

* 浏览器打开`http://localhost:8181`

* 第一个进入的玩家为玩家1，第二个为玩家2，其余为观战玩家

* 玩家1退出，玩家2升级为玩家1，根据观战玩家加入顺序，提升为玩家2

<img alt="联机演示" src="https://github.com/xiuxiubiu/pictures/raw/master/go-web-nes/online.gif" width="100%" />