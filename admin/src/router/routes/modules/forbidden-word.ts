import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/forbidden-word',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '违禁词',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 5,
    hideChildrenInMenu: true,
  },
  children: [
    {
      path: '',
      name: 'ForbiddenWord',
      component: () => import('@/views/pages/forbidden-word/index.vue'),
      meta: {
        title: '违禁词',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
