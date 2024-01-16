<script setup>
const props = defineProps({
  // 请求URL
  url: {
    type: String,
    required: true,
  },
  // 请求参数
  params: {
    type: Object,
    default() {
      return {};
    },
  },
});
// 是否正在加载中
const loading = ref(true);
const pageData = ref({
  cursor: "",
  results: [],
  hasMore: true,
});

const disabled = computed(() => {
  return loading.value || !pageData.value.hasMore;
});

const empty = computed(() => {
  return (
    pageData.value.hasMore === false && pageData.value.results.length === 0
  );
});

loadMore();

async function loadMore() {
  loading.value = true;
  try {
    const filters = Object.assign(props.params || {}, {
      cursor: pageData.value.cursor || "",
    });
    const data = await useHttpGet(props.url, {
      params: filters,
    });

    pageData.value.cursor = data.cursor;
    pageData.value.hasMore = data.hasMore;
    if (data.results && data.results.length) {
      data.results.forEach((item) => {
        pageData.value.results.push(item);
      });
    }
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
}
async function refresh() {
  pageData.value.cursor = "";
  pageData.value.results = [];
  pageData.value.hasMore = true;
  await loadMore();
}

function unshiftResults(item) {
  if (item && pageData.value && pageData.value.results) {
    pageData.value.results.unshift(item);
  }
}

defineExpose({
  refresh,
  unshiftResults,
});
</script>

<template>
  <div class="load-more">
    <slot v-if="empty" name="empty">
      <my-empty />
    </slot>
    <template v-else>
      <slot name="default" :results="pageData.results" />
      <div v-if="loading" class="loading">
        <el-skeleton :rows="3" animated />
      </div>
      <div class="has-more">
        <button
          class="button is-primary is-small"
          :disabled="disabled"
          @click="loadMore"
        >
          <span v-if="loading" class="icon">
            <i class="iconfont icon-loading" />
          </span>
          <span>{{ pageData.hasMore ? "查看更多" : "到底啦" }}</span>
        </button>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.load-more {
  .loading {
    background-color: var(--bg-color);
    padding: 10px;
  }
  .has-more {
    text-align: center;
    padding: 20px;
    button {
      width: 150px;
    }
  }

  .no-more {
    text-align: center;
    padding: 10px 0;
    color: var(--text-color3);
    font-size: 14px;
  }

  .icon-loading {
    animation: rotating 3s infinite linear;
  }
}
</style>
