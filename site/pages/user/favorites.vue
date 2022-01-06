<template>
  <div class="widget">
    <div class="widget-header">
      <span>
        <i class="iconfont icon-favorite" />
        <span>收藏列表</span>
      </span>
    </div>

    <div class="widget-content">
      <ul v-if="favorites && favorites.length" class="favorite-list">
        <li
          v-for="favorite in favorites"
          :key="favorite.favoriteId"
          class="favorite-item"
        >
          <div v-if="favorite.deleted" class="favorite-item">
            <div class="favorite-summary">收藏内容失效</div>
          </div>
          <div v-else>
            <div class="favorite-title">
              <a :href="favorite.url" target="_blank">{{ favorite.title }}</a>
            </div>
            <div class="favorite-summary">
              {{ favorite.content }}
            </div>
            <div class="favorite-meta">
              <span class="favorite-meta-item"
                ><nuxt-link :to="'/user/' + favorite.user.id">{{
                  favorite.user.nickname
                }}</nuxt-link></span
              >
              <span class="favorite-meta-item"
                ><time>{{ favorite.createTime | prettyDate }}</time></span
              >
            </div>
          </div>
        </li>
        <li v-if="hasMore" class="favorite-item more">
          <a @click="list">查看更多&gt;&gt;</a>
        </li>
      </ul>
      <div v-else class="notification is-primary">暂无收藏</div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'ucenter',
  data() {
    return {
      favorites: [],
      cursor: 0,
      hasMore: true,
    }
  },
  head() {
    return {
      title: this.$siteTitle('收藏'),
    }
  },
  mounted() {
    this.list()
  },
  methods: {
    async list() {
      const ret = await this.$axios.get('/api/user/favorites', {
        params: {
          cursor: this.cursor,
        },
      })
      if (ret.results && ret.results.length) {
        this.favorites = this.favorites.concat(ret.results)
      } else {
        this.hasMore = false
      }
      this.cursor = ret.cursor
    },
  },
}
</script>

<style lang="scss" scoped>
.favorite-list {
  margin: 0 !important;

  .favorite-item {
    overflow: hidden;
    zoom: 1;
    line-height: 24px;

    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    &.more {
      text-align: center;

      a {
        font-size: 15px;
        font-weight: bold;
      }
    }

    .favorite-title {
      a {
        color: var(--text-color3);
        font-weight: normal;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 18px;
        line-height: 30px;
      }
    }

    .favorite-summary {
      color: var(--text-color);
      overflow: hidden;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 3;
      text-align: justify;
      padding-top: 6px;
      word-break: break-all;
      text-overflow: ellipsis;
      font-size: 14px;
    }

    .favorite-meta {
      display: inline-block;
      font-size: 13px;
      padding-top: 6px;

      .favorite-meta-item {
        padding: 0 6px 0 0;
      }

      a {
        color: var(--text-link-color);
      }

      span {
        color: var(--text-color3);
      }
    }
  }
}
</style>
