[English](README.en-US.md) | [中文](README.md)


> Thanks for your support. If you like this open-source project, please give it a star. Your support keeps us moving forward.

## Official Links

- Website: [https://bbs-go.com](https://bbs-go.com)
- Community: [https://bbs.bbs-go.com](https://bbs.bbs-go.com)
- GitHub: [https://github.com/mlogclub/bbs-go](https://github.com/mlogclub/bbs-go)
- Gitee: [https://gitee.com/mlogclub/bbs-go](https://gitee.com/mlogclub/bbs-go)

## Demo

- Frontend: https://demo.bbs-go.com
- Admin: https://demo.bbs-go.com/admin
- Account: `admin / 123456`

## Why Choose bbs-go

- **Ready out of the box**: Core community capabilities including signup/login, posting, commenting, likes, favorites, follows, and notifications.
- **Growth loop built-in**: Tasks, points, levels, and badges to improve user activity and retention.
- **Operations friendly**: Content governance, user governance, permission governance, and system settings for long-term operations.
- **Bilingual support**: Built-in `en-US` and `zh-CN` for multilingual community scenarios.

## Feature Map

```mermaid
graph LR
  A((bbs-go))

  subgraph L1[Community Capabilities]
    direction TB
    U[User Side]
    U1[Signup and Login]
    U2[Profile]
    U3[Notifications]
    U4[Follow and Fans]
    U5[Points Ranking]
    U --> U1
    U --> U2
    U --> U3
    U --> U4
    U --> U5

    C[Content Side]
    C1[Topics and Feeds]
    C2[Articles]
    C3[Comments and Replies]
    C4[Likes and Favorites]
    C5[Tags and Nodes]
    C6[Search]
    C --> C1
    C --> C2
    C --> C3
    C --> C4
    C --> C5
    C --> C6
  end

  subgraph R1[Growth and Operations]
    direction TB
    G[Growth]
    G1[Daily Check-in]
    G2[Task System]
    G3[Points and EXP]
    G4[Level Progression]
    G5[Badges]
    G --> G1
    G --> G2
    G --> G3
    G --> G4
    G --> G5

    O[Operations]
    O1[User Management]
    O2[Content Governance]
    O3[Reports and Forbidden Words]
    O4[Roles and Permissions]
    O5[System Settings]
    O6[Operation Logs]
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

## Core Features

### User Side

- Account registration and login (multiple login methods supported)
- User profile management and personal homepage
- Follow/fan relationship management
- In-site notifications and interaction reminders
- Point records and leaderboards

### Content Side

- Publish and edit topics, feeds, and articles
- Complete interaction loop with comments, replies, likes, and favorites
- Tags and nodes for better content organization and discovery
- Interactive features like voting and hidden content
- In-site search for faster content retrieval

### Growth Side

- Daily check-in for ongoing activity incentives
- Task system (new user, daily, achievement)
- Points and EXP reward mechanisms
- Level progression configuration
- Badge and honor system

### Operations Side

- Unified governance for users, topics, comments, and articles
- Report handling and forbidden-word management
- Role, menu, and API permission management
- System parameter and site configuration management
- Operation logs and audit trails

## Typical Use Cases

- Developer and technical communities
- Hobby and interest-based communities
- Product user communities
- Internal enterprise knowledge communities
- Content membership communities

## Contact

QQ Group:
![BBS-GO QQ Group](docs/images/qq.png)

WeChat:
![WeChat](docs/images/wechat.png)

## Commercial Services

Paid services help sustain long-term development while the project remains open source.

| Service | Price | Description |
| ------- | ----- | ----------- |
| Commercial License | CNY 1628 | Commercial usage license for bbs-go |
| Feature Customization | Negotiable | Custom feature development based on your needs |

## What Is bbs-go

`bbs-go` is an open-source community system that helps you quickly build an operable and growth-oriented content community.

In one sentence: **Publish easily, engage deeply, govern effectively, and grow continuously.**

## Contributors

<a href="https://github.com/mlogclub/bbs-go/graphs/contributors"><img src="https://opencollective.com/bbs-go/contributors.svg?width=890&button=false" /></a>
