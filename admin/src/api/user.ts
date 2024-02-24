import axios from 'axios';
import type {
  RouteRecordNormalized,
  RouteRecordRaw,
  RouteMeta,
} from 'vue-router';
import { UserState } from '@/store/modules/user/types';

export interface LoginData {
  username: string;
  password: string;
  captchaId: string;
  captchaUrl: string;
  captchaCode: string;
}

export interface MenuItem {
  name: string;
  path: string;
  title: string;
  icon?: string;
  children?: Array<MenuItem>;
}

export interface LoginRes {
  token: string;
  redirect?: string;
  user?: UserState;
}

export function login(data: LoginData) {
  const formData = new FormData();
  formData.append('username', data.username);
  formData.append('password', data.password);
  formData.append('captchaId', data.captchaId);
  formData.append('captchaUrl', data.captchaUrl);
  formData.append('captchaCode', data.captchaCode);
  return axios.postForm<LoginRes>('/api/login/signin', formData);
}

export function logout() {
  axios.get('/api/login/signout');
}

export function getUserInfo() {
  return axios.get<UserState>('/api/user/current');
}

export async function getMenuList() {
  const ret = await axios.get<any, MenuItem[]>('/api/admin/menu/user_menus');
  return convertMenus(ret);
}

function convertMenus(items: MenuItem[]) {
  const menus: RouteRecordRaw[] = [];
  if (items && items.length) {
    items.forEach((item) => {
      const menu = convertMenuItem(item);
      if (item.children && item.children.length) {
        menu.children = convertMenus(item.children);
      }
      menus.push(menu);
    });
  }
  return menus;
}

function convertMenuItem(item: MenuItem) {
  const menu: RouteRecordRaw = {
    path: item.path,
    name: item.name,
    redirect: '',
    meta: {
      title: item.title,
      // icon: item.icon || 'icon-apps',
      icon: item.icon,
    } as RouteMeta,
  };
  return menu;
}
