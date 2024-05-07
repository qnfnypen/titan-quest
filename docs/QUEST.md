[TOC]

## 登陆 

接口URL有变动, 其他不变

> POST /api/v1/user/login_before

> POST /api/v1/user/login



## 查询任务

> GET /api/v1/quest/query_missions

参数：


响应:

```
{
    "code": 0,
    "data": {
        "basic_mission": [
            {
                "id": 1001,
                "title": "Connect Wallet",
                "channel": "Wallet",
                "logo": "https://static.fansland.io/fans/wallet.png",
                "credit": 5,
                "status": 1,
                "open_url": "",
                "start_time": "2024-04-23T17:38:17+08:00",
                "end_time": "2024-04-23T17:38:17+08:00",
                "type": 1,
                "parent_id": 0,
                "created_at": "2024-04-23T17:38:17+08:00",
                "updated_at": "2024-04-23T17:38:17+08:00",
                "sub_mission": null
            },
            {
                "id": 1002,
                "title": "Follow Titan Twitter",
                "channel": "Twitter",
                "logo": "https://static.fansland.io/fans/twitter.png",
                "credit": 5,
                "status": 1,
                "open_url": "https://twitter.com/intent/user?screen_name=fansland_io",
                "start_time": "2024-04-23T17:40:54+08:00",
                "end_time": "2024-04-23T17:40:54+08:00",
                "type": 1,
                "parent_id": 0,
                "created_at": "2024-04-23T17:40:54+08:00",
                "updated_at": "2024-04-23T17:40:54+08:00",
                "sub_mission": null
            },
            {
                "id": 1003,
                "title": "\nRetweet a tweet",
                "channel": "Twitter",
                "logo": "https://static.fansland.io/fans/twitter.png",
                "credit": 5,
                "status": 1,
                "open_url": "https://twitter.com/intent/retweet?tweet_id=1777949881136722272&related=twitterapi,twittermedia,twitter,support",
                "start_time": "2024-04-23T17:41:27+08:00",
                "end_time": "2024-04-23T17:41:27+08:00",
                "type": 1,
                "parent_id": 0,
                "created_at": "2024-04-23T17:41:27+08:00",
                "updated_at": "2024-04-23T17:41:27+08:00",
                "sub_mission": null
            },
            {
                "id": 1004,
                "title": "Like this tweet",
                "channel": "Twitter",
                "logo": "https://static.fansland.io/fans/twitter.png",
                "credit": 5,
                "status": 1,
                "open_url": "https://twitter.com/intent/like?tweet_id=1777949881136722272",
                "start_time": "2024-04-23T17:42:01+08:00",
                "end_time": "2024-04-23T17:42:01+08:00",
                "type": 1,
                "parent_id": 0,
                "created_at": "2024-04-23T17:42:01+08:00",
                "updated_at": "2024-04-23T17:42:01+08:00",
                "sub_mission": null
            },
            {
                "id": 1005,
                "title": "Join Fansland on Discord",
                "channel": "Discord",
                "logo": "https://static.fansland.io/fans/discord.png",
                "credit": 5,
                "status": 1,
                "open_url": "https://discord.gg/fansland",
                "start_time": "2024-04-23T17:42:23+08:00",
                "end_time": "2024-04-23T17:42:23+08:00",
                "type": 1,
                "parent_id": 0,
                "created_at": "2024-04-23T17:42:23+08:00",
                "updated_at": "2024-04-23T17:42:23+08:00",
                "sub_mission": null
            }
        ],
        "daily_mission": [
            {
                "id": 1006,
                "title": "Quote this tweet & tag 3 friends on Twitter",
                "channel": "Twitter",
                "logo": "https://static.fansland.io/fans/twitter.png",
                "credit": 10,
                "status": 1,
                "open_url": "https://twitter.com/fansland_io/status/1777949881136722272",
                "start_time": "2024-04-23T17:49:28+08:00",
                "end_time": "2024-04-23T17:49:28+08:00",
                "type": 2,
                "parent_id": 0,
                "created_at": "2024-04-23T17:49:28+08:00",
                "updated_at": "2024-04-23T17:49:28+08:00",
                "sub_mission": null
            },
            {
                "id": 1007,
                "title": "Post a tweet with the specific content",
                "channel": "Twitter",
                "logo": "https://static.fansland.io/fans/twitter.png",
                "credit": 10,
                "status": 1,
                "open_url": "https://twitter.com/fansland_io/status/1767394060950802439",
                "start_time": "2024-04-23T17:50:15+08:00",
                "end_time": "2024-04-23T17:50:15+08:00",
                "type": 2,
                "parent_id": 0,
                "created_at": "2024-04-23T17:50:15+08:00",
                "updated_at": "2024-04-23T17:50:15+08:00",
                "sub_mission": null
            }
        ],
        "weekly_mission": [
            {
                "id": 1008,
                "title": "Invite friends to Discord",
                "channel": "Discord",
                "logo": "https://static.fansland.io/fans/discord.png",
                "credit": 0,
                "status": 1,
                "open_url": "",
                "start_time": "2024-04-23T17:50:45+08:00",
                "end_time": "2024-04-23T17:50:45+08:00",
                "type": 3,
                "parent_id": 0,
                "created_at": "2024-04-23T17:50:45+08:00",
                "updated_at": "2024-04-23T17:50:45+08:00",
                "sub_mission": [
                    {
                        "id": 1010,
                        "title": "5",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 5,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1011,
                        "title": "10",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 15,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1012,
                        "title": "30",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 60,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1013,
                        "title": "50",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 120,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1014,
                        "title": "100",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 300,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1015,
                        "title": "500",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 1000,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1008,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    }
                ]
            },
            {
                "id": 1009,
                "title": "Engage & Chat in Discord to level up, get weekly $Fans Points Airdrop!",
                "channel": "Discord",
                "logo": "https://static.fansland.io/fans/discord.png",
                "credit": 0,
                "status": 1,
                "open_url": "",
                "start_time": "2024-04-23T18:04:47+08:00",
                "end_time": "2024-04-23T18:04:47+08:00",
                "type": 3,
                "parent_id": 0,
                "created_at": "2024-04-23T18:04:47+08:00",
                "updated_at": "2024-04-23T18:04:47+08:00",
                "sub_mission": [
                    {
                        "id": 1016,
                        "title": "Lv10",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 100,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1009,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1017,
                        "title": "Lv20",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 150,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1009,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1018,
                        "title": "Lv30",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 200,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1009,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1019,
                        "title": "Lv40",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 250,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1009,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    },
                    {
                        "id": 1020,
                        "title": "Lv50",
                        "channel": "Discord",
                        "logo": "",
                        "credit": 300,
                        "status": 1,
                        "open_url": " ",
                        "start_time": "2024-04-23T18:40:53+08:00",
                        "end_time": "2024-04-23T18:40:53+08:00",
                        "type": 3,
                        "parent_id": 1009,
                        "created_at": "2024-04-23T18:40:53+08:00",
                        "updated_at": "2024-04-23T18:40:53+08:00"
                    }
                ]
            }
        ]
    },
    "success": true
}
```


## 查询用户完成的任务情况

> GET /api/v1/quest/query_user_credits

**鉴权**

参数：


响应:


```
{
    "code": 0,
    "data": {
        "address": "0xe003B2Fb03F3126347afDBba460ED39e57F9588d",
        "credits": 5,
         "discord_user_id": "",
         "twitter_user_id": "1357906704566943745"
        "missions": {
            "basic_mission": [
                {
                    "id": 1000,
                    "username": "0xe003B2Fb03F3126347afDBba460ED39e57F9588d",
                    "mission_id": 1001,
                    "sub_mission_id": 0,
                    "type": 1,
                    "credit": 5,
                    "content": "0xe003B2Fb03F3126347afDBba460ED39e57F9588d",
                    "created_at": "2024-04-24T00:32:51+08:00",
                    "updated_at": "2024-04-24T00:32:51+08:00"
                }
            ],
            "daily_mission": null,
            "weekly_mission": null
        }
    },
    "success": true
}
```

## 推特OAUTH

> GET /api/v1/user/twitter/auth

**鉴权**

参数：
| 名称       | 类型     | 是否必须 | 描述                         |
| -------- | ------ | ---- | -------------------------- |
| callback | STRING | YES  |  回调URL                        |

示例:

```
/api/v1/user/twitter/auth?callback=http://192.168.0.120:8080/api/v1/twitter/callback
```

响应:

```
{
    "code": 0,
    "data": {
        "url": "https://api.twitter.com/oauth/authenticate?oauth_token=whL7EAAAAAABtRLFAAABjwwPYCE"
    },
    "success": true
}
```



## 验证任务 

> GET /api/v1/quest/check?mission_id=1002

**这个接口限制一分钟请求一次, 推特接口很贵的**

**鉴权**

参数：
| 名称       | 类型     | 是否必须 | 描述                         |
| -------- | ------ | ---- | -------------------------- |
| mission_id | STRING | YES  | 验证任务id                        |

示例:

```
 /api/v1/quest/check?mission_id=1002
```

响应:

```
{
    "code": 0,
    "data": {
        "missions": [
            {
                "id": 1008,
                "username": "0xe003B2Fb03F3126347afDBba460ED39e57F9588d",
                "mission_id": 1002,
                "sub_mission_id": 0,
                "type": 1,
                "credit": 5,
                "content": "1357906704566943745",
                "created_at": "2024-04-24T12:32:40+08:00",
                "updated_at": "2024-04-24T12:32:40+08:00"
            }
        ]
    },
    "success": true
}
```


## 提交推特链接 

> POST /api/v1/quest/twitter_link

**鉴权**

参数：
| 名称       | 类型     | 是否必须 | 描述                         |
| -------- | ------ | ---- | -------------------------- |
| mission_id | STRING | YES  | 验证任务id                        |
| link | STRING | YES  | 推特链接                        |

示例:

```
{
    "mission_id": 1106,
    "link": "https://x.com/shaqlin1/status/1783055657870213344"
}
```

响应:

```
{
    "code": 0,
    "data": null,
    "success": true
}
```

## Discord OAUTH

> GET /api/v1/user/discord/auth

**鉴权**

参数：


响应:

```
{
    "code": 0,
    "data": {
        "url": "https://discord.com/oauth2/authorize?client_id=1228239005540290602&redirect_uri=http%3A%2F%2F192.168.0.120%3A8080%2Fapi%2Fv1%2Fdiscord%2Fcallback&response_type=code&scope=identify+email&state=wfIAfhfvkoez"
    },
    "success": true
}
```


## 绑定KOL邀请码

> POST /api/v1/quest/kol_referral_code

**鉴权**

参数：
| 名称       | 类型     | 是否必须 | 描述                         |
| -------- | ------ | ---- | -------------------------- |
| code | STRING | YES  | KOL邀请码                        |

示例:

```
{
    "code": "brwJLw"
}
```

响应:

```
{
    "code": 0,
    "data": null,
    "success": true
}
```

## 浏览官网以及其回调
> GET /api/v1/quest/official_website/brows

**鉴权**

响应: 会携带code自动跳转到官网


> POST /api/v1/brows_official_website/callback

参数：
| 名称       | 类型     | 是否必须 | 描述                         |
| -------- | ------ | ---- | -------------------------- |
| code | STRING | YES  | 浏览官网时携带的code                        |

示例:

```
{
    "code": "brwJLw"
}
```

响应:

```
{
    "code": 0,
    "data": null,
    "success": true
}
```

> GET /api/v1/quest/official_website/verify

**鉴权**

响应:

```
{
    "code": 0,
    "data": {
        "verified": true
    },
    "success": true
}
```
