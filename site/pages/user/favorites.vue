<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <user-profile :user="currentUser" />
        <div class="widget">
          <div class="widget-header">
            <nav class="breadcrumb">
              <ul>
                <li><a href="/">首页</a></li>
                <li>
                  <a :href="'/user/' + currentUser.id">{{
                    currentUser.nickname
                  }}</a>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">收藏列表</a>
                </li>
              </ul>
            </nav>
          </div>

          <div class="widget-content">
            <ul v-if="favorites && favorites.length" class="article-list">
              <li v-for="favorite in favorites" :key="favorite.favoriteId">
                <article v-if="favorite.deleted" class="article-item">
                  <div class="article-summary">
                    收藏内容失效!!!
                  </div>
                </article>
                <article v-else class="article-item">
                  <div class="article-title">
                    <a :href="favorite.url">{{ favorite.title }}</a>
                  </div>
                  <div class="article-summary">
                    {{ favorite.content }}
                  </div>
                  <div class="article-meta">
                    <span class="article-meta-item"
                      ><a :href="'/user/' + favorite.user.id">{{
                        favorite.user.nickname
                      }}</a></span
                    >
                    <span class="article-meta-item"
                      ><time>{{ favorite.createTime | prettyDate }}</time></span
                    >
                  </div>
                </article>
              </li>
              <li v-if="hasMore" class="more">
                <a @click="list">查看更多&gt;&gt;</a>
              </li>
            </ul>
            <div v-else class="notification is-primary">
              暂无收藏
            </div>
          </div>
        </div>
      </div>
      <user-center-sidebar :user="currentUser" />
    </div>
  </section>
</template>

<script>
import UserProfile from '~/components/UserProfile'
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  middleware: 'authenticated',
  components: {
    UserProfile,
    UserCenterSidebar,
  },
  async asyncData({ $axios, params }) {},
  data() {
    return {
      favorites: [],
      cursor: 0,
      hasMore: true,
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
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
.article-list {
  margin: 0 !important;

  li {
    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;

    &:not(:last-child) {
      border-bottom: 1px solid #f2f2f2;
    }

    &.more {
      text-align: center;

      a {
        font-size: 15px;
        font-weight: bold;
      }
    }
  }

  .article-item {
    overflow: hidden;
    zoom: 1;
    line-height: 24px;
  }
}

article {
  .article-title {
    a {
      color: #999;
      font-weight: normal;
      overflow: hidden;
      text-overflow: ellipsis;
      font-size: 18px;
      line-height: 30px;
    }
  }

  .article-summary {
    color: #000;
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

  .article-meta {
    display: inline-block;
    font-size: 13px;
    padding-top: 6px;

    .article-meta-item {
      padding: 0 6px 0 0;
    }

    a {
      color: #3273dc;
    }

    span {
      color: #999;
    }
  }
}
</style>
