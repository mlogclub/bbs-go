<template>
  <div class="topics-nav">
    <ul class="topics-nav-list">
      <li :class="{ active: currentNodeId === 0 }" class="topics-nav-item">
        <a href="/topics/node/newest">最新</a>
      </li>
      <li :class="{ active: currentNodeId === -1 }" class="topics-nav-item">
        <a href="/topics/node/recommend">推荐</a>
      </li>
      <li
        v-for="node in nodes"
        :key="node.nodeId"
        :class="{ active: currentNodeId == node.nodeId }"
        class="topics-nav-item"
      >
        <a :href="'/topics/node/' + node.nodeId">{{ node.name }}</a>
      </li>
    </ul>
  </div>
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
.topics-nav {
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #e6ecf0 !important;

  .topics-nav-list {
    .topics-nav-item {
      display: inline-block;
      padding: 0 15px;

      a {
        position: relative;
        display: inline-block;
        padding: 14px 0;
        line-height: 22px;
        text-align: center;
        font-weight: 500;
        color: #8590a6;

        &:hover {
          color: #343a40;
        }
      }

      &.active {
        a {
          color: #343a40;
          font-weight: 700;

          &:after {
            position: absolute;
            right: -3px;
            bottom: -1px;
            left: -3px;
            height: 3px;
            background: #e0245e;
            content: '';
          }
        }
      }
    }
  }
}
</style>
