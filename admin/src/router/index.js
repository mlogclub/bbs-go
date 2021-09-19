import Vue from 'vue';
import VueRouter from 'vue-router';

import Layout from '@/components/layout/Layout.vue';

Vue.use(VueRouter);

const routes = [
  // {
  //   path: '/',
  //   hidden: true,
  //   redirect: {
  //     path: '/home',
  //   },
  // },
  {
    path: '/',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/Home.vue'),
        meta: {
          title: '首页',
          closable: false,
          icon: 'el-icon-s-home',
        },
      },
    ],
  },
  {
    path: '/topic',
    component: Layout,
    meta: {
      title: '话题管理',
    },
    children: [
      {
        path: '/about',
        component: () => import('@/views/About.vue'),
        meta: {
          title: '关于',
          icon: 'iconfont icon-article',
        },
      },
    ],
  },
];

const router = new VueRouter({
  routes,
});

export default router;
