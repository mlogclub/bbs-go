// import type { RouteRecordNormalized, RouteRecordRaw } from 'vue-router';
import type { MenuItem } from '@/api/user';

export interface TableConfig {
  size: 'mini' | 'medium' | 'large' | 'small';
  bordered:
  | boolean
  | import('@arco-design/web-vue/es/table/interface').TableBorder;
}

export interface AppState {
  theme: string;
  colorWeak: boolean;
  navbar: boolean;
  menu: boolean;
  topMenu: boolean;
  hideMenu: boolean;
  menuCollapse: boolean;
  footer: boolean;
  themeColor: string;
  menuWidth: number;
  globalSettings: boolean;
  device: string;
  tabBar: boolean;
  menuFromServer: boolean;
  serverMenu: MenuItem[];
  routeLoaded: boolean; // 路由是否已经加载

  table: TableConfig;
  [key: string]: unknown;
}
