<template>
  <nav class="dock-nav">
    <ul>
      <!-- <li class="dock-nav-divider"></li> -->
      <li
        v-for="node in nodes"
        :key="node.id"
        :class="{ active: envStore.currentNodeId === node.id }"
      >
        <nuxt-link :to="nodeUrl(node)">
          <img class="node-logo" :src="nodeLogo(node)" />
          <span class="node-name">{{ node.name }}</span>
        </nuxt-link>
      </li>
    </ul>
  </nav>
</template>

<script setup>
import iconNew from "~/assets/images/new.png";
import iconRecommend from "~/assets/images/recommend.png";
import iconFeed from "~/assets/images/feed.png";
import iconNode from "~/assets/images/node.png";

const envStore = useEnvStore();

const { data: nodes } = await useAsyncData("nodes", () =>
  useMyFetch(`/api/topic/node_navs`)
);

function nodeLogo(node) {
  if (node.logo) {
    return node.logo;
  }
  if (node.id === 0) {
    return iconNew;
  } else if (node.id === -1) {
    return iconRecommend;
  } else if (node.id === -2) {
    return iconFeed;
  }
  return iconNode;
}

function nodeUrl(node) {
  if (node.id > 0) {
    return `/topics/node/${node.id}`;
  } else if (node.id === 0) {
    return "/topics/node/newest";
  } else if (node.id === -1) {
    return "/topics/node/recommend";
  } else if (node.id === -2) {
    return "/topics/node/feed";
  }
}
</script>

<style lang="scss" scoped>
.dock-nav {
  display: block;
  position: -webkit-sticky;
  position: sticky;
  top: 10px;

  width: 150px;
  border-radius: 2px;
  background-color: var(--bg-color);
  transition: all 0.2s linear;

  ul {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px 12px;

    li:not(.dock-nav-divider) {
      position: relative;
      cursor: pointer;
      height: 30px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      color: var(--text-color);
      //padding: 0 12px;
      border-radius: 3px;
      transition: background-color 0.2s, color 0.2s;
      font-weight: 500;

      &:not(:first-child) {
        margin-top: 10px;
      }

      &.active {
        background-color: #ea6f5a;
        color: var(--text-color5);

        a {
          color: var(--text-color5);
        }
      }

      &:not(.active):hover {
        background-color: hsla(0, 0%, 94.9%, 0.6);
      }

      a {
        text-decoration: none;
        cursor: pointer;
        color: var(--text-color3);
        width: 100%;
        height: 100%;
        text-align: center;
        line-height: 30px;
        padding-left: 10px;

        display: flex;
        align-items: center;
        //justify-content: center;
        .node-logo {
          width: 24px;
          height: 24px;
          border-radius: 4px;
          margin-right: 10px;
          background-color: var(--bg-color);
        }
      }
    }

    li.dock-nav-divider {
      height: 15px;
      border-bottom: 1px solid var(--border-color);
    }
  }
}
</style>
