<template>
  <div v-if="activeTags && activeTags.length" class="widget">
    <div class="header">
      <span>今日热议</span>
      <div class="slot">
        <a href="/tags">查看更多>></a>
      </div>
    </div>
    <div class="content">
      <div class="tags are-medium">
        <span v-for="tag in activeTags" :key="tag.tagId" class="tag is-normal is-white">
          <a :href="'/articles/tag/' + tag.tagId" :title="tag.tagName">{{ tag.tagName }}</a>
        </span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    activeTags: {
      type: Array,
      default: function () {
        return []
      }
    }
  }
}
</script>

<style lang="scss" scoped>

</style>
