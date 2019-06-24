> 本人日常开发的主力语言为Java，因为对golang的兴趣，所以才边学边写了该项目，完全是为了学习。

## 技术栈

- iris (https://github.com/kataras/iris) mvc框架
- gorm (http://gorm.io/) orm框架
- resty (https://github.com/go-resty/resty) 好用的http-client
- cron (https://github.com/robfig/cron) 定时任务
- goquery（https://github.com/PuerkitoBio/goquery）html dom元素解析
- Element-UI (https://element.eleme.cn) 饿了么开源的基于vue.js的前端库 


## startup

1. 打开根目录下的mlog.json，配置自己的数据库，只需要创建好库就可以，表会自动为你创建。
2. 根目录下执行:`go mod tidy`将所有需要依赖的库下载到本地。（依赖管理使用了go mod，不会用的请看这里：https://mlog.club/topic/9）
3. 执行`go run main.go`启动项目

## 注意事项

- 后台管理项目使用vue+element-ui编写，采用前后端分离方案，后台工程在`web/admin`请自行编辑。
