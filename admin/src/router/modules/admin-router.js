/** When your routing table is too long, you can split it into small modules **/

import Layout from "@/layout";

const adminRouter = [
  {
    path: "/users",
    component: Layout,
    redirect: "/users",
    children: [
      {
        name: "Users",
        path: "",
        component: () => import("@/views/pages/users/index"),
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
        name: "Topics",
        path: "topics",
        component: () => import("@/views/pages/topics/index"),
        meta: {
          title: "话题",
          icon: "iconfont icon-topic",
        },
      },
      {
        name: "Articles",
        path: "articles",
        component: () => import("@/views/pages/articles/index"),
        meta: {
          title: "文章",
          icon: "iconfont icon-article",
        },
      },
      {
        name: "Comments",
        path: "comments",
        component: () => import("@/views/pages/comments/index"),
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
        name: "Nodes",
        path: "nodes",
        component: () => import("@/views/pages/topics/nodes"),
        meta: {
          title: "节点",
          icon: "iconfont icon-tag",
        },
      },
      {
        name: "Tags",
        path: "tags",
        component: () => import("@/views/pages/tags/index"),
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
        name: "Links",
        path: "",
        component: () => import("@/views/pages/links/index"),
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
