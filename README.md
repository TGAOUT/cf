<p align="center">
<img width="500" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022209168.png"><br><br>
<a href="https://github.com/teamssix/cf/stargazers"><img alt="GitHub stars" src="https://img.shields.io/github/stars/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/releases"><img alt="GitHub issues" src="https://img.shields.io/github/release/teamssix/cf"/></a>
<a href="https://github.com/teamssix/cf/blob/main/LICENSE"><img alt="License" src="https://img.shields.io/badge/License-Apache%202.0-blue.svg"/></a>
<a href="https://github.com/teamssix/cf/releases"><img alt="Downloads" src="https://img.shields.io/github/downloads/teamssix/cf/total?color=brightgreen"/></a>
<a href="https://goreportcard.com/report/github.com/teamssix/cf"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/teamssix/cf"/></a>
<a href="https://twitter.com/intent/tweet/?text=CF%2C%20an%20amazing%20cloud%20exploitation%20framework%0Ahttps%3A%2F%2Fgithub.com%2Fteamssix%2Fcf%0A%23cloud%20%23security%20%23cloudsecurity%20%23cybersecurtiy"><img alt="tweet" src="https://img.shields.io/twitter/url?url=https://github.com/teamssix/cf" /></a>
<a href="https://twitter.com/teamssix"><img alt="Twitter" src="https://img.shields.io/twitter/url/https/twitter.com/teamssix.svg?style=social&label=Follow" /></a>
<a href="https://github.com/teamssix"><img alt="Github" src="https://img.shields.io/github/followers/TeamsSix?style=social" /></a><br></br>
中文 | <a href="README_EN.md">English</a>
</p>


---

CF 是一个云环境利用框架，主要用来方便红队人员在获得云服务的访问凭证即 Access Key 的后续工作。

CF 下载地址：[github.com/teamssix/cf/releases](https://github.com/teamssix/cf/releases)

> 目前 CF 仅支持阿里云，后续会不断更新对其他云的支持

目前 CF 可以实现以下功能：

* 已实现

  - [x] 列出对象存储（包括存储桶大小和文件数量信息）
  - [x] 列出 ECS 实例
  - [x] 一键获得实例上的临时访问凭证
  - [x] 一键为所有实例执行三要素，方便 HVV
  - [x] 一键为实例反弹 Shell
  - [x] 支持阿里云
  - [x] 列出 RDS 云数据库实例
  - [x] 一键接管控制台
  - [x] 一键查看当前访问凭证所拥有的权限
  - [x] ……
  
* 预计实现

  - [ ] 云上痕迹清除
  
  - [ ] 自动检测当前运行环境是不是实例，如果是则一键扫描本地实例的凭证信息
  - [ ] 一键将获取到的临时凭证添加到工具中
  - [ ] 支持腾讯云等其他云厂商
  - [ ] ……

## 使用手册

使用手册请参见：[wiki.teamssix.com/cf](https://wiki.teamssix.com/cf)

[![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207112152449.png)](https://wiki.teamssix.com/cf)

## 简单上手

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207052307021.png)

配置 CF

```bash
cf configure
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241064.png)

一键列出当前访问凭证的云服务资源

```bash
cf ls
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207040107386.png)

一键列出当前访问凭证的权限

```bash
cf ls permissions
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207082108875.png)

一键接管控制台

```bash
cf console
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207082104422.png)

查看 CF 为实例执行命令的操作的帮助信息

```bash
cf ecs exec -h
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022215293.png)

一键为所有实例执行三要素，方便 HVV

```
cf ecs exec -b
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241381.png)

一键获取实例中的临时访问凭证数据

```bash
cf ecs exec -m
```

![img](https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202207022241672.png)

如果感觉还不错的话，师傅记得给个 Star 呀 ~，另外 CF 的更多使用方法可以参见使用文档：[wiki.teamssix.com/cf](https://wiki.teamssix.com/cf)

## 注意事项

* 本工具仅用于合法合规用途，严禁用于违法违规用途。
* 本工具中所涉及的风险点均属于租户责任，与云厂商无关。

## 更多

如果你对云安全比较感兴趣，可以看我的另外一个项目 [Awesome Cloud Security](https://github.com/teamssix/awesome-cloud-security)，这里收录了很多国内外的云安全资源。

另外在我的[云安全文库](https://wiki.teamssix.com/)里有大量的云安全方向的笔记和文章，最后，下面这个是我的个人微信公众号，欢迎关注 ~

<div align=center><img width="700" src="https://cdn.jsdelivr.net/gh/teamssix/BlogImages/imgs/202204152148071.png" div align=center/></div>

<div align=center><img src="https://api.star-history.com/svg?repos=teamssix/cf&type=Timeline" div align=center/></div>









