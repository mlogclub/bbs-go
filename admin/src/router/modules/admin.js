/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const adminRouter = {
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
        title: '话题',
        affix: true
      }
    },
    {
      path: 'articles',
      component: () => import('@/views/pages/articles/index'),
      name: 'articles',
      meta: {
        title: '文章',
        affix: true
      }
    }
  ]
}
export default adminRouter
