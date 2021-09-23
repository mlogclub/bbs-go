/** When your routing table is too long, you can split it into small modules **/

import Layout from "@/layout";

const adminRouter = [
  {
    path: "/users",
    component: Layout,
    redirect: "/users",
    children: [
      {
        path: "",
        component: () => import("@/views/pages/users/index"),
        name: "users",
        meta: {
          title: "用户管理",
          icon: "iconfont icon-username",
        },
      },
    ],
  },
  {
    path: "/content",
    component: Layout,
    redirect: "/content/topics",
    name: "content",
    meta: {
      title: "内容管理",
      icon: "iconfont icon-topic",
    },
    children: [
      {
        path: "topics",
        component: () => import("@/views/pages/topics/index"),
        name: "topics",
        meta: {
          title: "话题",
          icon: "iconfont icon-topic",
        },
      },
      {
        path: "articles",
        component: () => import("@/views/pages/articles/index"),
        name: "articles",
        meta: {
          title: "文章",
          icon: "iconfont icon-article",
        },
      },
      {
        path: "comments",
        component: () => import("@/views/pages/comments/index"),
        name: "comments",
        meta: {
          title: "评论",
          icon: "iconfont icon-comments",
        },
      },
    ],
  },
  {
    path: "/category",
    component: Layout,
    redirect: "/category/nodes",
    name: "cocategoryntent",
    meta: {
      title: "分类管理",
      icon: "iconfont icon-tags",
    },
    children: [
      {
        path: "nodes",
        component: () => import("@/views/pages/topics/nodes"),
        name: "nodes",
        meta: {
          title: "节点",
          icon: "iconfont icon-tag",
        },
      },
      {
        path: "tags",
        component: () => import("@/views/pages/tags/index"),
        name: "tags",
        meta: {
          title: "标签",
          icon: "iconfont icon-tags",
        },
      },
    ],
  },
  {
    path: "/links",
    component: Layout,
    redirect: "/links",
    children: [
      {
        path: "",
        component: () => import("@/views/pages/links/index"),
        name: "links",
        meta: {
          title: "友情链接",
          icon: "iconfont icon-link",
        },
      },
    ],
  },
  {
    path: "/settings",
    component: Layout,
    redirect: "/settings",
    children: [
      {
        path: "",
        component: () => import("@/views/pages/settings/index"),
        name: "settings",
        meta: {
          title: "系统设置",
          icon: "iconfont icon-setting",
        },
      },
    ],
  },
];
export default adminRouter;
