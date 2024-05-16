<template>
  <div class="topic-search-items">
    <div v-for="item in results" :key="item.id" class="topic-search-item">
      <h1 class="topic-search-item-title">
        <a target="_blank" :href="'/topic/' + item.id" v-html="item.title" />
      </h1>
      <a
        target="_blank"
        :href="'/topic/' + item.id"
        class="topic-search-item-summary content"
      >
        <p v-html="item.summary"></p>
      </a>
      <div class="topic-mates">
        <!-- <span>{{ item.user.nickname }}</span> -->
        <span>{{ useFormatDate(item.createTime) }}</span>
        <span v-if="item.node">{{ item.node.name }}</span>
        <!-- <template v-if="item.tags && item.tags.length">
          <span v-for="tag in item.tags" :key="tag.id" class="tag">{{
            tag.name
          }}</span>
        </template> -->
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    results: {
      type: Object,
      default: null,
    },
  },
};
</script>

<style lang="scss" scoped>
.topic-search-items {
  .topic-search-item {
    padding: 10px;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .topic-search-item-title {
      margin-bottom: 10px;
      font-size: 16px;
      font-weight: 600;

      a {
        color: var(--text-color);

        &:hover {
          text-decoration: underline;
        }
      }
    }

    .topic-search-item-summary {
      margin-bottom: 10px;
      display: inline-block;
      font-size: 15px;
      width: 100%;
      text-decoration: none;
      color: var(--text-color2);
      word-wrap: break-word;

      overflow: hidden;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 3;
      text-align: justify;
      word-break: break-all;
      text-overflow: ellipsis;
    }

    .topic-mates {
      font-size: 13px;
      color: var(--text-color3);
      span {
        margin-right: 10px;
      }
    }
  }
}
</style>
