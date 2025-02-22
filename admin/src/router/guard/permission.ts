import type { Router } from 'vue-router';
import NProgress from 'nprogress'; // progress bar

import usePermission from '@/hooks/permission';
import { useUserStore, useAppStore } from '@/store';
import loadRouters from './loadRouters';
import { appRoutes } from '../routes';
import { NOT_FOUND } from '../constants';

export default async function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const appStore = useAppStore();
    const userStore = useUserStore();
    if (userStore.id > 0) {
      if (appStore.menuFromServer) {
        if (!appStore.routeLoaded) {
          await loadRouters(router);
          next({
            path: to.path,
            replace: true,
          });
        } else {
          next();
        }
      } else {
        const Permission = usePermission();
        const permissionsAllow = Permission.accessRouter(to);
        // eslint-disable-next-line no-lonely-if
        if (permissionsAllow) {
          next();
        } else {
          const destination =
            Permission.findFirstPermissionRoute(appRoutes, userStore.role) ||
            NOT_FOUND;
          next(destination);
        }
      }
    } else {
      next();
    }
    NProgress.done();
  });
}
