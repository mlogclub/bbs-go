<template>
  <div class="topics-nav">
    <a :href="createTopicUrl" class="button is-primary post">
      <i class="iconfont icon-topic" />&nbsp;
      <strong>发表主题</strong>
    </a>
    <ul>
      <li>
        <a :class="{ active: currentNodeId === 0 }" href="/topics/node/newest">
          <i class="iconfont icon-topic" />
          <span>最新</span>
        </a>
      </li>
      <li>
        <a
          :class="{ active: currentNodeId === -1 }"
          href="/topics/node/recommend"
        >
          <i class="iconfont icon-topic" />
          <span>推荐</span>
        </a>
      </li>
      <li v-for="node in nodes" :key="node.nodeId" class="topics-nav-item">
        <a
          :class="{ active: currentNodeId == node.nodeId }"
          :href="'/topics/node/' + node.nodeId"
        >
          <i class="iconfont icon-topic" />
          <span>{{ node.name }}</span>
        </a>
      </li>
    </ul>
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
  min-width: 150px;
  margin-right: 10px;
  padding: 5px;

  ul {
    li {
      a {
        display: block;
        width: 100%;
        padding: 8px 0;
        color: #667d99;
        font-weight: 500;
        font-size: 15px;

        i {
          margin-right: 8px;
        }

        &:hover {
          color: #f36322;
          text-decoration: none;
        }

        &.active {
          color: rgb(75, 147, 209);
          font-weight: 900;
        }
      }
    }
  }
}
</style>
