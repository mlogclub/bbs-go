import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

export default {
  path: '/permission',
  component: DEFAULT_LAYOUT,
  meta: {
    title: '权限管理',
    requiresAuth: true,
    icon: 'icon-apps',
    order: 8,
  },
  children: [
    {
      path: 'role',
      name: 'Role',
      component: () => import('@/views/pages/system/role/index.vue'),
      meta: {
        title: '角色管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'menu',
      name: 'Menu',
      component: () => import('@/views/pages/system/menu/index.vue'),
      meta: {
        title: '菜单管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'index',
      name: 'Permission',
      component: () => import('@/views/pages/system/permission/index.vue'),
      meta: {
        title: '权限管理',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
} as AppRouteRecordRaw;
