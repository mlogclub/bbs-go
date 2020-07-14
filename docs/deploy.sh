# 生成静态文件
vuepress build

# 进入生成的文件夹
cd .vuepress/dist

# 如果是发布到自定义域名
echo 'docs.bbs-go.com' > CNAME

git init
git add -A
git config user.name 'mlogclub'
git config user.email 'mlog1@qq.com'
git commit -m 'deploy docs'

git push -f git@github.com:mlogclub/bbs-go-docs.git master
