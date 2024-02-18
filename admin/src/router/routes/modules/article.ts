import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/article',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '文章管理',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 4,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'Article',
      component: () => import('@/views/pages/article/index.vue'),
      meta: {
        title: '文章管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
