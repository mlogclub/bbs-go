import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/views/Home.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      hidden: true,
      redirect: {
        path: '/topic/nodes'
      }
    },
    {
      path: '1',
      component: Home,
      meta: {
        title: '控制台',
        icon: 'iconfont icon-dashboard'
      },
      children: [
        {
          path: '/topic/nodes',
          component: () => import('@/views/topic/TopicNode.vue'),
          meta: {
            title: '节点',
            icon: 'iconfont icon-topic'
          }
        },
        {
          path: '/user/index',
          component: () => import('@/views/user/Index.vue'),
          meta: {
            title: '用户',
            icon: 'iconfont icon-username'
          }
        },
        {
          path: '/topic/index',
          component: () => import('@/views/topic/Index.vue'),
          meta: {
            title: '话题',
            icon: 'iconfont icon-topic'
          }
        },
        {
          path: '/article/index',
          component: () => import('@/views/article/Index.vue'),
          meta: {
            title: '文章',
            icon: 'iconfont icon-article'
          }
        },
        {
          path: '/tag/index',
          component: () => import('@/views/tag/Index.vue'),
          meta: {
            title: '标签',
            icon: 'iconfont icon-tags'
          }
        },
        {
          path: '/comment/index',
          component: () => import('@/views/comment/Index.vue'),
          meta: {
            title: '评论',
            icon: 'iconfont icon-comment'
          }
        },
        {
          path: '/link/index',
          component: () => import('@/views/link/Index.vue'),
          meta: {
            title: '链接',
            icon: 'iconfont icon-article'
          }
        },
        {
          path: '/sys-config/index',
          component: () => import('@/views/sys-config/Index.vue'),
          meta: {
            title: '配置',
            icon: 'iconfont icon-setting'
          }
        }
      ]
    }

    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (about.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import(/* webpackChunkName: "about" */ './views/About.vue'),
    // },
  ]
})
