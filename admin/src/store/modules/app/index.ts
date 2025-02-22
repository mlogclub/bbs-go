import { defineStore } from 'pinia';
import type { RouteRecordNormalized } from 'vue-router';
import { getMenuList, convertMenus } from '@/api/user';
import { AppState } from './types';

const useAppStore = defineStore('app', {
  state: (): AppState => ({
    title: 'GO-ADMIN',
    theme: 'light',
    colorWeak: false,
    navbar: true,
    menu: true,
    topMenu: false,
    hideMenu: false,
    menuCollapse: false,
    footer: false,
    themeColor: '#165DFF',
    menuWidth: 220,
    globalSettings: false,
    device: 'desktop',
    tabBar: true,
    menuFromServer: true,
    serverMenu: [],
    routeLoaded: false,

    table: {
      size: 'medium',
      bordered: { wrapper: true, cell: true },
      // bordered: {
      //   cell: true,
      // },
    },
  }),

  getters: {
    appCurrentSetting(state: AppState): AppState {
      return { ...state };
    },
    appDevice(state: AppState) {
      return state.device;
    },
    appAsyncMenus(state: AppState): RouteRecordNormalized[] {
      const list = convertMenus(state.serverMenu);
      return list as unknown as RouteRecordNormalized[];
    },
  },

  actions: {
    // Update app settings
    updateSettings(partial: Partial<AppState>) {
      // @ts-ignore-next-line
      this.$patch(partial);
    },

    // Change theme color
    toggleTheme(dark: boolean) {
      if (dark) {
        this.theme = 'dark';
        document.body.setAttribute('arco-theme', 'dark');
      } else {
        this.theme = 'light';
        document.body.removeAttribute('arco-theme');
      }
    },
    toggleDevice(device: string) {
      this.device = device;
    },
    toggleMenu(value: boolean) {
      this.hideMenu = value;
    },
    async fetchServerMenuConfig() {
      try {
        this.serverMenu = await getMenuList();
      } catch (error) {
        useHandleError(error);
      }
    },
    clearServerMenu() {
      this.serverMenu = [];
    },
    setRouteLoaded(value: boolean) {
      this.routeLoaded = value;
    },
  },
});

export default useAppStore;
