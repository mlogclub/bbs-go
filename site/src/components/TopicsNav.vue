<template>
  <nav class="dock-nav">
    <ul>
      <li
        v-for="node in nodes"
        :key="node.id"
        :class="{ active: envStore.currentNodeId === node.id }"
      >
        <nuxt-link :to="nodeUrl(node)">
          <i
            class="node-logo"
            :style="'background-image: url(' + nodeLogo(node) + ')'"
          ></i>
          <div class="node-name">{{ node.name }}</div>
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
  useHttpGet(`/api/topic/node_navs`)
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

<style>
[data-theme="light"],
.theme-light {
  --topics-nav-color1: #f7f9ff;
  --topics-nav-color2: rgba(43, 89, 255, 0.06);
  --topics-nav-color3: #fff7f7;
}

[data-theme="dark"],
.theme-dark {
  --topics-nav-color1: #0d1a29;
  --topics-nav-color2: #3a2c2c;
  --topics-nav-color3: #3a2c2c;
}
</style>
<style lang="scss" scoped>
.dock-nav {
  display: block;
  position: -webkit-sticky;
  position: sticky;
  top: calc(52px + 1rem);

  width: 180px;
  border-radius: 12px;
  background-color: var(--bg-color);
  transition: all 0.2s linear;

  ul {
    padding: 16px 0;

    li {
      position: relative;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 12px;
      cursor: pointer;

      &:hover {
        &:before {
          visibility: visible;
        }
      }

      &:before {
        visibility: hidden;
        position: absolute;
        content: "";
        top: -2px;
        left: 0;
        right: 0;
        bottom: -2px;
        background-color: var(--topics-nav-color1);
        box-shadow: 0px 4px 4px var(--topics-nav-color2);
        transition: all 0.1s ease-out 0.05s;
      }

      &.active {
        background-color: var(--topics-nav-color3);
      }

      a {
        padding: 12px 24px;
        position: relative;
        z-index: 2;
        display: flex;
        // color: #16181f;
        color: var(--text-color);
        align-items: center;

        .node-logo {
          flex-shrink: 0;
          width: 24px;
          height: 24px;
          margin: 0 8px 0 0;
          background-position: center;
          background-repeat: no-repeat;
          background-size: 100% 100%;
          background-color: #fff;
          border-radius: 4px;
        }

        .node-name {
          height: 12px;
          line-height: 12px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
}
</style>
