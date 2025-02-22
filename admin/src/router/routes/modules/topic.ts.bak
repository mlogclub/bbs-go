import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/topic',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '帖子管理',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 3,
  },
  children: [
    {
      path: 'topic-node',
      name: 'TopicNode',
      component: () => import('@/views/pages/topic-node/index.vue'),
      meta: {
        title: '节点管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'index',
      name: 'Topic',
      component: () => import('@/views/pages/topic/index.vue'),
      meta: {
        title: '帖子管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
