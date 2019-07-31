import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import UserIndex from './views/user/Index';
import CategoryIndex from './views/category/Index';
import TagIndex from './views/tag/Index';
import TopicIndex from './views/topic/Index';
import ArticleIndex from './views/article/Index';
import CommentIndex from './views/comment/Index';
import OauthClientIndex from './views/oauth-client/Index'
import OauthTokenIndex from './views/oauth-token/Index'
import SysConfigIndex from './views/sys-config/Index'

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: '首页',
      hidden: true,
      redirect: {
        path: '/article/index'
      }
    },
    {
      path: '1',
      component: Home,
      name: '内容管理',
      iconCls: 'iconfont icon-article',
      children: [
        {
          path: '/article/index',
          component: ArticleIndex,
          name: '文章',
          iconCls: 'iconfont icon-article'
        },
        {
          path: '/topic/index',
          component: TopicIndex,
          name: '话题',
          iconCls: 'iconfont icon-topic'
        },
        {
          path: '/category/index',
          component: CategoryIndex,
          name: '分类',
          iconCls: 'iconfont icon-category'
        },
        {
          path: '/tag/index',
          component: TagIndex,
          name: '标签',
          iconCls: 'iconfont icon-tags'
        },
        {
          path: '/comment/index',
          component: CommentIndex,
          name: '评论',
          iconCls: 'iconfont icon-comment'
        },
        {
          path: '/sys-config/index',
          component: SysConfigIndex,
          name: '系统配置',
          iconCls: 'iconfont icon-setting',
        },
      ]
    },

    {
      path: '2',
      component: Home,
      name: '用户管理',
      iconCls: 'iconfont icon-user',
      children: [
        {
          path: '/user/index',
          component: UserIndex,
          name: '用户',
          iconCls: 'iconfont icon-user'
        }
      ]
    },

    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (about.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
    // },
  ],
});
