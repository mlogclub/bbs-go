import type { Router, LocationQueryRaw } from 'vue-router';
import NProgress from 'nprogress'; // progress bar

import { useUserStore } from '@/store';
import {
  DEFAULT_ROUTE_NAME,
  NOT_FOUND_ROUTE_NAME,
  WHITE_LIST,
} from '@/router/constants';

export default async function setupUserLoginInfoGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    NProgress.start();
    const userStore = useUserStore();
    const user = await userStore.info();
    if (user && user.id > 0) {
      if (user.id > 0) {
        next();
      } else {
        try {
          next();
        } catch (error) {
          await userStore.logout();
          next({
            name: 'login',
            query: {
              redirect:
                to.name !== NOT_FOUND_ROUTE_NAME ? to.name : DEFAULT_ROUTE_NAME,
              ...to.query,
            } as LocationQueryRaw,
          });
        }
      }
    } else {
      if (WHITE_LIST.find((el) => el.name === to.name)) {
        next();
        return;
      }
      next({
        name: 'login',
        query: {
          redirect:
            to.name !== NOT_FOUND_ROUTE_NAME ? to.name : DEFAULT_ROUTE_NAME,
          ...to.query,
        } as LocationQueryRaw,
      });
    }
  });
}
