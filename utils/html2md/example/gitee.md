# 书籍阅读微信小程序



## 项目介绍


<code>BookStackWechat</code> - 书籍阅读微信小程序，不仅仅是[BookStack](https://gitee.com/TruthHun/BookStack)的配套小程序。

起初仅仅是想作为[BookStack](https://github.com/TruthHun/BookStack)的配套小程序，但是后来想了一下，小说书籍和文档阅读类的小程序，不外乎也就是那些功能，所以<code>BookStackWechat</code>会开发成为通用的书籍阅读类小程序。

由于 [DocHub](https://gitee.com/TruthHun/DocHub) 文库项目的开源，小程序的开发，暂停了一段时间，目前只开发实现了部分页面，当前更多页面还在继续开发中，源码届时会以 Apache 2.0 开源协议进行开源。在放出源码之前，你可以先对当前项目进行<code>star</code>和<code>watch</code>，以关注项目动态。

小程序采用CSS3的<code>flex</code>实现页面布局，预计会迭代发布三个版本：

- 
第一个版本： 纯模板页面，实现页面间的链接跳转

- 
第二个版本： 约定所有API接口，请求动态数据

- 
第三个版本： 套接<code>BookStack</code>的API接口，实现真正意义上的书籍阅读（届时会提前发布<code>BookStack</code> v2.x版本）


## 页面预览


由于没有原型，小程序的功能布局和页面布局，主要借(抄)鉴(袭)了3个手机APP，所以，请允许我免费给他们打个小广告：

- [iReader](http://www.zhangyue.com/) - 引领品质阅读
- [微信读书](http://weread.qq.com/) - 让阅读不再孤独
- [熊猫阅读](http://www.pandadushu.com/) - (没有slogan就是最好的slogan...)
- - -

好吧，他们都有一个响亮的Slogan，那么我们也要有一个：


<blockquote>
    <p>书栈阅读 - 让阅读，成为一种本能</p>
</blockquote>

- - -

页面一览：


### 首页

![](https://gitee.com/truthhun/BookStackWeChat/raw/master/screenshot/index1.png)![](https://gitee.com/truthhun/BookStackWeChat/raw/master/screenshot/index2.png)

### 分类

![](https://gitee.com/truthhun/BookStackWeChat/raw/master/screenshot/cate.png)

### 书架

![](https://gitee.com/truthhun/BookStackWeChat/raw/master/screenshot/bookshelf.png)

### 用户中心

![](https://gitee.com/truthhun/BookStackWeChat/raw/master/screenshot/me.png)