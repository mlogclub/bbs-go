import localeLogin from '@/views/login/locale/en-US';
import localeDashboard from '@/views/pages/dashboard/locale/en-US';
import localeUser from '@/views/pages/user/locale/en-US';
import localeTopic from '@/views/pages/topic/locale/en-US';
import localeTopicNode from '@/views/pages/topic-node/locale/en-US';
import localeArticle from '@/views/pages/article/locale/en-US';
import localeForbiddenWord from '@/views/pages/forbidden-word/locale/en-US';
import localeLink from '@/views/pages/link/locale/en-US';
import localeRole from '@/views/pages/system/role/locale/en-US';
import localeApi from '@/views/pages/system/api/locale/en-US';
import localeMenu from '@/views/pages/system/menu/locale/en-US';
import localePermission from '@/views/pages/system/permission/locale/en-US';
import localeDict from '@/views/pages/system/dict/locale/en-US';
import localeSettings from '@/views/pages/settings/locale/en-US';

export default {
  settings: {
    title: 'Settings',
    themeColor: 'Theme Color',
    content: 'Content Setting',
    search: 'Search',
    language: 'Language',
    navbar: {
      title: 'Navbar',
      theme: {
        toLight: 'Click to use light mode',
        toDark: 'Click to use dark mode',
      },
      screen: {
        toFull: 'Click to switch to full screen mode',
        toExit: 'Click to exit the full screen mode',
      },
      alerts: 'alerts',
    },
    menuWidth: 'Menu Width (px)',
    menu: 'Menu',
    topMenu: 'Top Menu',
    tabBar: 'Tab Bar',
    footer: 'Footer',
    otherSettings: 'Other Settings',
    colorWeak: 'Color Weak',
    alertContent:
      'After the configuration is only temporarily effective, if you want to really affect the project, click the "Copy Settings" button below and replace the configuration in settings.json.',
    copySettings: {
      title: 'Copy Settings',
      message: 'Copy succeeded, please paste to file src/settings.json.',
    },
    close: 'Close',
    color: {
      tooltip: '10 gradient colors generated according to the theme color',
    },
    menuFromServer: 'Menu From Server',
  },
  menu: {
    list: 'List',
    result: 'Result',
    exception: 'Exception',
    form: 'Form',
    profile: 'Profile',
    visualization: 'Data Visualization',
    user: 'User Center',
    arcoWebsite: 'Arco Design',
    faq: 'FAQ',
  },
  navbar: {
    docs: 'Docs',
    action: {
      locale: 'Switch to English',
    },
  },
  messageBox: {
    logout: 'Logout',
  },
  router: {
    dashboard: 'Dashboard',
  },
  pages: {
    ...localeLogin,
    ...localeDashboard,
    ...localeUser,
    ...localeTopic,
    ...localeTopicNode,
    ...localeArticle,
    ...localeForbiddenWord,
    ...localeLink,
    ...localeRole,
    ...localeApi,
    ...localeMenu,
    ...localePermission,
    ...localeDict,
    ...localeSettings,
  },
};
