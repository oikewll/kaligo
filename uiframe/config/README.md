## 结构说明
  1. `client`目录对应浏览器端的配置
  2. `server`目录对应服务端的配置,当ssr项目才需要配置此项

  > client端不能引入server端的配置，防止关键数据泄露

---
## ssr项目须知
> ps: 如果是ssr项目，由于server/api/** (node服务的api) 里面无法拿到process.server的状态，直接引入config，是拿不到server的,需要在用到的时候自己import对应环境
