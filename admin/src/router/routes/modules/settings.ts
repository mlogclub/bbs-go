import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/settings',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '系统设置',
    requiresAuth: true,
    icon: 'icon-settings',
    order: 7,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'Settings',
      component: () => import('@/views/pages/settings/index.vue'),
      meta: {
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
