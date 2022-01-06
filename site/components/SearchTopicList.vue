<template>
  <div class="topic-search-items">
    <div
      v-for="item in searchPage.results"
      :key="item.id"
      class="topic-search-item"
    >
      <a target="_blank" :href="'/topic/' + item.topicId">
        <h1 class="topic-search-item-title" v-html="item.title"></h1>
      </a>
      <p class="topic-search-item-summary content" v-html="item.summary"></p>
      <div class="topic-mates">
        <span v-html="item.nickname"></span>
        <span>{{ item.createTime | formatDate }}</span>
        <span v-if="item.node">{{ item.node.name }}</span>
        <!-- <template v-if="item.tags && item.tags.length">
          <span v-for="tag in item.tags" :key="tag.tagId" class="tag">{{
            tag.tagName
          }}</span>
        </template> -->
        <template v-if="item.tagNames && item.tagNames.length">
          <span
            v-for="(tagName, index) in item.tagNames"
            :key="index"
            class="tag"
            v-html="tagName"
          ></span>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    searchPage: {
      type: Object,
      default: null,
    },
  },
}
</script>

<style lang="scss" scoped>
.topic-search-items {
  background-color: var(--bg-color);
  .topic-search-item {
    padding: 10px;
    margin-bottom: 8px;
    background-color: var(--bg-color);
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    .topic-search-item-title {
      font-weight: bold;
      margin-bottom: 6px;
    }
    .topic-search-item-summary {
      font-size: 80%;
    }
    .topic-mates {
      font-size: 80%;
      color: var(--text-color3);
      span {
        margin-right: 10px;
      }
    }
  }
}
</style>
