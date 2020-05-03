<template>
  <div class="topics-nav">
    <nav class="dock-nav">
      <ul>
        <li :class="{ active: currentNodeId === 0 }">
          <a href="/topics/node/newest">最新</a>
        </li>
        <li :class="{ active: currentNodeId === -1 }">
          <a href="/topics/node/recommend">推荐</a>
        </li>
        <li
          v-for="node in nodes"
          :key="node.nodeId"
          :class="{ active: currentNodeId === node.nodeId }"
        >
          <a :href="'/topics/node/' + node.nodeId">{{ node.name }}</a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script>
export default {
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    },
    nodes: {
      type: Array,
      default() {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.topics-nav {
  margin-right: 10px;
  width: 112px;
  min-height: 100px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);

  .dock-nav {
    width: 112px;
    border-radius: 2px;
    background-color: #fff;
    transition: all 0.2s linear;

    ul {
      height: 100%;
      display: flex;
      flex-direction: column;
      padding: 16px 12px;
      li {
        position: relative;
        cursor: pointer;
        height: 30px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        color: #2e3135;
        padding: 0 12px;
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
        }
      }
    }
  }
}
</style>
