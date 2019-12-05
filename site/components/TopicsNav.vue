<template>
  <div class="topics-nav">
    <ul class="topics-nav-list">
      <li :class="{ active: currentTagId <= 0 }" class="topics-nav-item">
        <a href="/topics">全部</a>
      </li>
      <li
        v-for="tag in bbsNavTags"
        :key="tag.tagId"
        :class="{ active: currentTagId == tag.tagId }"
        class="topics-nav-item"
      >
        <a :href="'/topics/tag/' + tag.tagId">{{ tag.tagName }}</a>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    currentTagId: {
      type: Number,
      default: 0
    }
  },
  computed: {
    bbsNavTags() {
      return this.$store.state.config.config.bbsNavTags
    }
  }
}
</script>

<style lang="scss" scoped>
.topics-nav {
  margin-bottom: 8px;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #e6ecf0 !important;

  .topics-nav-list {
    .topics-nav-item {
      display: inline-block;
      padding: 0 20px;

      a {
        position: relative;
        display: inline-block;
        padding: 14px 0;
        line-height: 22px;
        text-align: center;
        font-weight: 500;
        color: #8590a6;
      }

      &.active {
        a {
          color: #343a40;
          font-weight: 700;

          &:after {
            position: absolute;
            right: 0;
            bottom: -1px;
            left: 0;
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
