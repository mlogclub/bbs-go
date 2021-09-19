/** When your routing table is too long, you can split it into small modules **/

import Layout from '@/layout'

const adminRouter = {
  path: '/topic',
  component: Layout,
  redirect: '/list',
  name: 'topic',
  meta: {
    title: '话题'
    // icon: 'table'
  },
  children: [
    {
      path: 'list',
      component: () => import('@/views/topic/index'),
      name: 'list',
      meta: {
        title: '列表',
        affix: true
      }
    }
  ]
}
export default adminRouter
