import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import UserIndex from './views/user/Index';
import CategoryIndex from './views/category/Index';
import TagIndex from './views/tag/Index';
import ArticleIndex from './views/article/Index';
import CommentIndex from './views/comment/Index';
import OauthClientIndex from './views/oauth-client/Index'
import OauthTokenIndex from './views/oauth-token/Index'

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: '首页',
      hidden: true,
      redirect: {
        path: '/category/index'
      }
    },

    {
      path: '1',
      component: Home,
      name: '文章管理',
      children: [
        {
          path: '/article/index',
          component: ArticleIndex,
          name: '文章',
          iconCls: 'el-icon-menu'
        },
        {
          path: '/category/index',
          component: CategoryIndex,
          name: '分类',
          iconCls: 'el-icon-menu'
        },
        {
          path: '/tag/index',
          component: TagIndex,
          name: '标签',
          iconCls: 'el-icon-menu'
        },
        {
          path: '/comment/index',
          component: CommentIndex,
          name: '评论',
          iconCls: 'el-icon-menu'
        },
      ]
    },

    {
      path: '2',
      component: Home,
      name: '用户管理',
      children: [
        {
          path: '/user/index',
          component: UserIndex,
          name: '用户',
          iconCls: 'el-icon-menu'
        },
        {
          path: '/oauth-client/index',
          component: OauthClientIndex,
          name: 'OauthClient',
          iconCls: 'el-icon-menu'
        },
        {
          path: '/oauth-token/index',
          component: OauthTokenIndex,
          name: 'OauthToken',
          iconCls: 'el-icon-menu'
        },
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
