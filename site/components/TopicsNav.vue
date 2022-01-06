<template>
  <nav class="dock-nav">
    <ul>
      <li :class="{ active: currentNodeId === 0 }">
        <nuxt-link to="/topics/node/newest">
          <img class="node-logo" src="~/assets/images/new.png" />
          <span class="node-name">最新</span>
        </nuxt-link>
      </li>
      <li :class="{ active: currentNodeId === -1 }">
        <nuxt-link to="/topics/node/recommend">
          <img class="node-logo" src="~/assets/images/recommend2.png" />
          <span class="node-name">推荐</span>
        </nuxt-link>
      </li>
      <li :class="{ active: currentNodeId === -2 }">
        <nuxt-link to="/topics/node/feed">
          <img class="node-logo" src="~/assets/images/feed.png" />
          <span class="node-name">关注</span>
        </nuxt-link>
      </li>
      <li class="dock-nav-divider"></li>
      <li
        v-for="node in nodes"
        :key="node.nodeId"
        :class="{ active: currentNodeId === node.nodeId }"
      >
        <nuxt-link :to="'/topics/node/' + node.nodeId">
          <img v-if="node.logo" class="node-logo" :src="node.logo" />
          <img v-else class="node-logo" src="~/assets/images/node.png" />
          <span class="node-name">{{ node.name }}</span>
        </nuxt-link>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  props: {
    nodes: {
      type: Array,
      default() {
        return []
      },
    },
  },
  computed: {
    currentNodeId() {
      return this.$store.state.env.currentNodeId
    },
  },
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
