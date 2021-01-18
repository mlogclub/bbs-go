<template>
  <nav class="dock-nav">
    <ul>
      <li :class="{ active: currentNodeId === 0 }">
        <a href="/topics/node/newest">最新</a>
      </li>
      <li :class="{ active: currentNodeId === -1 }">
        <a href="/topics/node/recommend">推荐</a>
      </li>
      <li class="dock-nav-divider"></li>
      <li
        v-for="node in nodes"
        :key="node.nodeId"
        :class="{ active: currentNodeId === node.nodeId }"
      >
        <a :href="'/topics/node/' + node.nodeId">{{ node.name }}</a>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  props: {
    currentNodeId: {
      type: Number,
      default: 0,
    },
    nodes: {
      type: Array,
      default() {
        return []
      },
    },
  },
}
</script>

<style lang="scss" scoped>
@import './assets/styles/variable';
.dock-nav {
  display: block;
  position: -webkit-sticky;
  position: sticky;
  top: 10px;

  width: 150px;
  border-radius: 2px;
  background-color: #fff;
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
      color: #2e3135;
      //padding: 0 12px;
      border-radius: 3px;
      transition: background-color 0.2s, color 0.2s;
      font-weight: 500;

      &:not(:first-child) {
        margin-top: 10px;
      }

      &.active {
        background-color: #ea6f5a;
        color: #fff;
        a {
          color: #fff;
        }
      }

      &:not(.active):hover {
        background-color: hsla(0, 0%, 94.9%, 0.6);
      }

      a {
        text-decoration: none;
        cursor: pointer;
        color: #909090;
        width: 100%;
        height: 100%;
        text-align: center;
        line-height: 30px;
      }
    }
    li.dock-nav-divider {
      height: 15px;
      border-bottom: 1px solid $border-color-base;
    }
  }
}
</style>
