import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/link',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '友情链接',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 6,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'Link',
      component: () => import('@/views/pages/link/index.vue'),
      meta: {
        title: '友情链接',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
