import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/dashboard',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '仪表盘',
    requiresAuth: true,
    icon: 'icon-dashboard',
    order: 1,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'Dashboard',
      component: () => import('@/views/pages/dashboard/index.vue'),
      meta: {
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
