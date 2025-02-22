import type { Router, RouteRecordRaw } from 'vue-router';
import { useAppStore } from '@/store';
import { MenuItem } from '@/api/user';
import { DEFAULT_LAYOUT } from '@/router/routes/base';

// 匹配views里面所有的.vue文件
const modules = import.meta.glob('@/views/pages/**/*.vue');

export default async function loadRouters(router: Router) {
  const appStore = useAppStore();
  await appStore.fetchServerMenuConfig();
  addRoutes(appStore.serverMenu);
  appStore.setRouteLoaded(true);

  function addRoutes(menuList: MenuItem[]) {
    const routes = getRoutes([], menuList, 1);
    routes.forEach((route) => {
      router.addRoute(route);
    });
  }

  function getRoutes(
    routes: RouteRecordRaw[],
    menuList: MenuItem[],
    deep: number
  ): RouteRecordRaw[] {
    if (!menuList || !menuList.length) {
      return [];
    }
    if (!routes || !routes.length) {
      routes = [] as RouteRecordRaw[];
    }
    menuList.forEach((item) => {
      if (item.type !== 'menu') {
        return;
      }

      const currentRoute = {
        name: item.name,
        path: item.path,
        meta: {
          title: item.title,
          requiresAuth: true,
        },
        component: loadComponent(item.component || ''),
      };

      let route: RouteRecordRaw;
      if (deep === 1) {
        route = {
          path: item.path,
          component: DEFAULT_LAYOUT,
          meta: {
            title: item.title,
            requiresAuth: true,
          },
          children: [currentRoute],
        } as RouteRecordRaw;
      } else {
        route = currentRoute;
      }

      if (item.children && item.children.length) {
        route.children = getRoutes([], item.children, deep + 1);
      }

      routes.push(route);
    });
    return routes;
  }

  function loadComponent(view: string) {
    if (!view) {
      return DEFAULT_LAYOUT;
    }
    let res;
    // eslint-disable-next-line guard-for-in, no-restricted-syntax
    for (const path in modules) {
      const dir = path.split('views/pages/')[1].split('.vue')[0];
      if (dir === view) {
        res = () => modules[path]();
      }
    }
    return res;
  }
}
