import { defineStore } from 'pinia';
// import { Notification } from '@arco-design/web-vue';
// import type { NotificationReturn } from '@arco-design/web-vue/es/notification/interface';
import type { RouteRecordNormalized } from 'vue-router';
import { getMenuList } from '@/api/user';
import { AppState } from './types';

const useAppStore = defineStore('app', {
  state: (): AppState => ({
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

    table: {
      size: 'medium',
      bordered: { wrapper: true, cell: false },
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
      return state.serverMenu as unknown as RouteRecordNormalized[];
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
      // let notifyInstance: NotificationReturn | null = null;
      try {
        // notifyInstance = Notification.info({
        //   id: 'menuNotice', // Keep the instance id the same
        //   content: 'loading',
        //   closable: true,
        // });
        const data = await getMenuList();
        this.serverMenu = data as unknown as RouteRecordNormalized[];
        // notifyInstance = Notification.success({
        //   id: 'menuNotice',
        //   content: 'success',
        //   closable: true,
        // });
      } catch (error) {
        useHandleError(error);
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
        // notifyInstance = Notification.error({
        //   id: 'menuNotice',
        //   content: 'error',
        //   closable: true,
        // });
      }
    },
    clearServerMenu() {
      this.serverMenu = [];
    },
  },
});

export default useAppStore;
