/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const adminRouter = [
  {
    path: '/users',
    component: Layout,
    redirect: '/users',
    name: 'users',
    children: [{
      path: '',
      component: () => import('@/views/pages/users/index'),
      name: 'users',
      meta: {
        title: '用户管理'
      }
    }]
  },
  {
    path: '/content',
    component: Layout,
    redirect: '/content/topics',
    name: 'content',
    meta: {
      title: '内容管理'
      // icon: 'table'
    },
    children: [
      {
        path: 'topics',
        component: () => import('@/views/pages/topics/index'),
        name: 'topics',
        meta: {
          title: '话题'
        }
      },
      {
        path: 'articles',
        component: () => import('@/views/pages/articles/index'),
        name: 'articles',
        meta: {
          title: '文章'
        }
      },
      {
        path: 'comments',
        component: () => import('@/views/pages/comments/index'),
        name: 'comments',
        meta: {
          title: '评论'
        }
      }
    ]
  },
  {
    path: '/category',
    component: Layout,
    redirect: '/category/nodes',
    name: 'cocategoryntent',
    meta: {
      title: '分类管理'
      // icon: 'table'
    },
    children: [
      {
        path: 'nodes',
        component: () => import('@/views/pages/topics/nodes'),
        name: 'topics',
        meta: {
          title: '节点'
        }
      },
      {
        path: 'tags',
        component: () => import('@/views/pages/tags/index'),
        name: 'tags',
        meta: {
          title: '标签'
        }
      }
    ]
  },
  {
    path: '/links',
    component: Layout,
    redirect: '/links',
    name: 'links',
    children: [
      {
        path: '',
        component: () => import('@/views/pages/links/index'),
        name: 'links',
        meta: {
          title: '友情链接'
        }
      }
    ]
  },
  {
    path: '/settings',
    component: Layout,
    redirect: '/settings',
    name: 'settings',
    children: [
      {
        path: '',
        component: () => import('@/views/pages/settings/index'),
        name: 'settings',
        meta: {
          title: '系统设置'
        }
      }
    ]
  }
]
export default adminRouter
