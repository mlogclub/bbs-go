import type { RouteRecordNormalized } from 'vue-router';

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
  serverMenu: RouteRecordNormalized[];

  table: TableConfig;
  [key: string]: unknown;
}
