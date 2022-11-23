/** When your routing table is too long, you can split it into small modules **/

import Layout from "@/layout";

const adminRouter = [
  {
    path: "/user",
    component: Layout,
    redirect: "/user",
    children: [
      {
        name: "User",
        path: "",
        component: () => import("@/views/pages/user/index"),
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
    redirect: "/content/topic",
    name: "content",
    meta: {
      title: "内容管理",
      icon: "iconfont icon-topic",
    },
    children: [
      {
        name: "Topic",
        path: "topic",
        component: () => import("@/views/pages/topic/index"),
        meta: {
          title: "帖子管理",
          icon: "iconfont icon-topic",
        },
      },
      {
        name: "TopicReview",
        path: "topic/review",
        component: () => import("@/views/pages/topic/review"),
        meta: {
          title: "帖子审核",
          icon: "iconfont icon-audit",
        },
      },
      {
        name: "Article",
        path: "article",
        component: () => import("@/views/pages/article/index"),
        meta: {
          title: "文章管理",
          icon: "iconfont icon-article",
        },
      },
      {
        name: "ArticleReview",
        path: "article/review",
        component: () => import("@/views/pages/article/review"),
        meta: {
          title: "文章审核",
          icon: "iconfont icon-audit",
        },
      },
      {
        name: "ForbiddenWord",
        path: "forbidden-word",
        component: () => import("@/views/pages/forbidden-word/index"),
        meta: {
          title: "违禁词",
          icon: "iconfont icon-forbidden",
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
        component: () => import("@/views/pages/topic/nodes"),
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
