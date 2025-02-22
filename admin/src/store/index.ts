import { createPinia } from 'pinia';
import useAppStore from './modules/app';
import useUserStore from './modules/user';
import useTabBarStore from './modules/tab-bar';
import useDictStore from './modules/dict';

const pinia = createPinia();

export { useAppStore, useUserStore, useTabBarStore, useDictStore };
export default pinia;
