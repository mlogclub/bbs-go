[中文](README.md) | [English](README.en-US.md)


> 感谢您的支持与鼓励！如果您喜欢这个开源项目，不妨给它点个⭐️⭐️⭐️，您的星星是我们前进的动力 🙏🙏🙏

## 官网

- 官网：[https://bbs-go.com](https://bbs-go.com)
- 交流社区：[https://bbs.bbs-go.com](https://bbs.bbs-go.com)
- Github：[https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)
- Gitee：[https://gitee.com/mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)

## 演示

- 前台: https://demo.bbs-go.com
- 后台: https://demo.bbs-go.com/admin
- 账号密码: admin/123456

## 为什么选择 bbs-go

- **开箱可用**：注册登录、发帖评论、点赞收藏、关注消息等核心社区能力可直接使用。
- **增长闭环**：内置任务、积分、等级、勋章，支持用户活跃和长期留存。
- **运营友好**：提供内容治理、用户治理、权限治理与系统配置能力，方便持续运营。
- **双语支持**：内置 `en-US` / `zh-CN`，适合面向不同语言用户的社区场景。

## 功能地图

```mermaid
graph LR
  A((bbs-go))

  subgraph L1[社区能力]
    direction TB
    U[用户侧]
    U1[注册登录]
    U2[个人主页]
    U3[消息通知]
    U4[关注粉丝]
    U5[积分排行]
    U --> U1
    U --> U2
    U --> U3
    U --> U4
    U --> U5

    C[内容侧]
    C1[帖子动态]
    C2[文章发布]
    C3[评论回复]
    C4[点赞收藏]
    C5[标签节点]
    C6[站内搜索]
    C --> C1
    C --> C2
    C --> C3
    C --> C4
    C --> C5
    C --> C6
  end

  subgraph R1[增长与运营]
    direction TB
    G[增长侧]
    G1[每日签到]
    G2[任务体系]
    G3[积分经验]
    G4[等级成长]
    G5[勋章激励]
    G --> G1
    G --> G2
    G --> G3
    G --> G4
    G --> G5

    O[运营侧]
    O1[用户管理]
    O2[内容治理]
    O3[举报违禁词]
    O4[角色权限]
    O5[系统配置]
    O6[运营日志]
    O --> O1
    O --> O2
    O --> O3
    O --> O4
    O --> O5
    O --> O6
  end

  A --> U
  A --> C
  A --> G
  A --> O
```

## 核心功能

### 用户侧

- 账号注册与登录（支持多种登录方式）
- 用户资料维护与个人主页展示
- 关注/粉丝关系管理
- 站内消息与互动提醒
- 积分记录与排行榜

### 内容侧

- 支持帖子、动态、文章发布与编辑
- 评论、回复、点赞、收藏等完整互动链路
- 标签与节点管理，便于内容组织和发现
- 支持投票、隐藏内容等互动玩法
- 站内搜索能力，提升内容检索效率

### 增长侧

- 每日签到，持续活跃激励
- 任务体系（新手、每日、成就）
- 积分与经验奖励机制
- 等级成长配置
- 勋章与荣誉体系

### 运营侧

- 用户、帖子、评论、文章等统一治理
- 举报处理与违禁词管理
- 角色、菜单、接口权限分配
- 系统参数与站点配置管理
- 运营日志与行为留痕

## 适用场景

- 技术交流社区
- 兴趣爱好社群
- 产品用户社区
- 企业内部知识社区
- 内容型会员社区

## 联系我

QQ群：
![BBS-GO用户交流群](docs/images/qq.png)

微信：
![微信](docs/images/wechat.png)

## 付费服务

付费是为了项目能够更好的生存下去，请谅解。项目将一如既往的开源下去~

| 服务     | 价格   | 服务内容                                         |
| -------- | ------ | ------------------------------------------------ |
| 商用授权 | ￥1628 | 提供 bbs-go 商业使用授权                         |
| 功能定制 | 面议   | 接受各种功能定制，只有你想不到的没有我们做不到的 |

## bbs-go 是什么

`bbs-go` 是一个开源社区系统，帮助你快速搭建可运营、可增长的内容社区。

一句话概括：**发得出来、聊得起来、管得住、长得快**。

## Contributors

<a href="https://github.com/mlogclub/bbs-go/graphs/contributors"><img src="https://opencollective.com/bbs-go/contributors.svg?width=890&button=false" /></a>
