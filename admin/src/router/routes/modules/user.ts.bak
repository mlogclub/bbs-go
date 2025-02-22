import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/user',
  name: 'user',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '用户管理',
    requiresAuth: true,
    icon: 'icon-user',
    order: 2,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: 'index',
      name: 'User',
      component: () => import('@/views/pages/user/index.vue'),
      meta: {
        title: '用户管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
