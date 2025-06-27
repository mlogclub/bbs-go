import localeLogin from '@/views/login/locale/zh-CN';
import localeDashboard from '@/views/pages/dashboard/locale/zh-CN';
import localeUser from '@/views/pages/user/locale/zh-CN';
import localeTopic from '@/views/pages/topic/locale/zh-CN';
import localeTopicNode from '@/views/pages/topic-node/locale/zh-CN';
import localeArticle from '@/views/pages/article/locale/zh-CN';
import localeForbiddenWord from '@/views/pages/forbidden-word/locale/zh-CN';
import localeLink from '@/views/pages/link/locale/zh-CN';
import localeRole from '@/views/pages/system/role/locale/zh-CN';
import localeApi from '@/views/pages/system/api/locale/zh-CN';
import localeMenu from '@/views/pages/system/menu/locale/zh-CN';
import localePermission from '@/views/pages/system/permission/locale/zh-CN';
import localeDict from '@/views/pages/system/dict/locale/zh-CN';
import localeSettings from '@/views/pages/settings/locale/zh-CN';

export default {
  settings: {
    title: '页面配置',
    themeColor: '主题色',
    content: '内容区域',
    search: '搜索',
    language: '语言',
    navbar: {
      title: '导航栏',
      theme: {
        toLight: '点击切换为亮色模式',
        toDark: '点击切换为暗黑模式',
      },
      screen: {
        toFull: '点击切换全屏模式',
        toExit: '点击退出全屏模式',
      },
      alerts: '消息通知',
    },
    menuWidth: '菜单宽度 (px)',
    menu: '菜单栏',
    topMenu: '顶部菜单栏',
    tabBar: '多页签',
    footer: '底部',
    otherSettings: '其他设置',
    colorWeak: '色弱模式',
    alertContent:
      '配置之后仅是临时生效，要想真正作用于项目，点击下方的 "复制配置" 按钮，将配置替换到 settings.json 中即可。',
    copySettings: {
      title: '复制配置',
      message: '复制成功，请粘贴到 src/settings.json 文件中',
    },
    close: '关闭',
    color: {
      tooltip:
        '根据主题颜色生成的 10 个梯度色（将配置复制到项目中，主题色才能对亮色 / 暗黑模式同时生效）',
    },
    menuFromServer: '菜单来源于后台',
  },
  menu: {
    list: '列表页',
    result: '结果页',
    exception: '异常页',
    form: '表单页',
    profile: '详情页',
    visualization: '数据可视化',
    user: '个人中心',
    arcoWebsite: 'Arco Design',
    faq: '常见问题',
  },
  navbar: {
    docs: '文档中心',
    action: {
      locale: '切换为中文',
    },
  },
  messageBox: {
    logout: '退出登录',
  },
  router: {
    dashboard: '仪表盘',
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
