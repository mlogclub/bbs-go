<template>
  <div class="topics-nav">
    <ul class="topics-nav-list">
      <li :class="{ active: nodeId === 0 }" class="topics-nav-item">
        <a @click="setNodeId(0)">{{ $t("component.search.all") }}</a>
      </li>
      <li :class="{ active: nodeId === -1 }" class="topics-nav-item">
        <a @click="setNodeId(-1)">{{ $t("component.search.recommend") }}</a>
      </li>
      <li
        v-for="node in nodes"
        :key="node.id"
        :class="{ active: nodeId === node.id }"
        class="topics-nav-item"
      >
        <a @click="setNodeId(node.id)">{{ node.name }}</a>
      </li>
    </ul>
    <div class="search-time-range">
      <div class="select is-small">
        <select v-model="timeRange" @change="setTimeRange">
          <option :value="0">{{ $t("component.search.timeRange.all") }}</option>
          <option :value="1">{{ $t("component.search.timeRange.day") }}</option>
          <option :value="2">{{ $t("component.search.timeRange.week") }}</option>
          <option :value="3">{{ $t("component.search.timeRange.month") }}</option>
          <option :value="4">{{ $t("component.search.timeRange.year") }}</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup>
const { data: nodes } = await useMyFetch("/api/topic/nodes");

const route = useRoute();
const router = useRouter();
const nodeId = ref(parseInt(route.query.nodeId) || 0);
const timeRange = ref(parseInt(route.query.timeRange) || 0);

const setNodeId = (changeNodeId) => {
  nodeId.value = changeNodeId;
  setQuery("nodeId", changeNodeId);
};

const setTimeRange = () => {
  setQuery("timeRange", timeRange.value);
};

const setQuery = (key, value) => {
  const currentQuery = { ...route.query };

  const newQuery = {
    ...currentQuery,
  };

  newQuery[key] = value;

  router.push({
    path: "/search",
    query: newQuery,
  });
};
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
            content: "";
          }
        }
      }
    }
  }
  .search-time-range {
  }
}
</style>
