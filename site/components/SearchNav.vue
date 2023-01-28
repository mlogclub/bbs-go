<template>
  <div class="topics-nav">
    <ul class="topics-nav-list">
      <li :class="{ active: currentNodeId === 0 }" class="topics-nav-item">
        <a @click="setNodeId(0)">全部</a>
      </li>
      <li :class="{ active: currentNodeId === -1 }" class="topics-nav-item">
        <a @click="setNodeId(-1)">推荐</a>
      </li>
      <li
        v-for="node in nodes"
        :key="node.nodeId"
        :class="{ active: currentNodeId === node.nodeId }"
        class="topics-nav-item"
      >
        <a @click="setNodeId(node.nodeId)">{{ node.name }}</a>
      </li>
    </ul>
    <div class="search-time-range">
      <div class="select is-small">
        <select v-model="timeRange" @change="setTimeRange">
          <option :value="0">时间不限</option>
          <option :value="1">一天内</option>
          <option :value="2">一周内</option>
          <option :value="3">一月内</option>
          <option :value="4">一年内</option>
        </select>
      </div>
    </div>
  </div>
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
  data() {
    return {
      timeRange: 0,
    }
  },
  computed: {
    currentNodeId() {
      return this.$store.state.search.nodeId
    },
  },
  methods: {
    setNodeId(nodeId) {
      this.$store.dispatch('search/changeNodeId', nodeId)
    },
    setTimeRange() {
      this.$store.dispatch('search/changeTimeRange', this.timeRange)
    },
  },
}
</script>

<style lang="scss" scoped>
.topics-nav {
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid var(--border-color) !important;
  display: flex;
  justify-content: space-between;
  align-items: center;

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
        color: var(--text-color3);

        &:hover {
          color: var(--text-color);
        }
      }

      &.active {
        a {
          color: var(--text-color);
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
  .search-time-range {
  }
}
</style>
