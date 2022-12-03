<template>
  <section class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        文章正在审核中
      </div>
    </div>
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <article
          class="article-item article-detail"
          itemscope
          itemtype="http://schema.org/BlogPosting"
        >
          <div class="main-content">
            <div class="article-header">
              <div class="article-title-wrapper">
                <h1 class="article-title" itemprop="headline">
                  {{ article.title }}
                </h1>
                <div class="article-manage-menu">
                  <article-manage-menu v-model="article" />
                </div>
              </div>
              <div class="article-meta">
                <span class="article-meta-item">
                  由
                  <nuxt-link
                    :to="'/user/' + article.user.id"
                    class="article-author"
                    itemprop="author"
                    itemscope
                    itemtype="http://schema.org/Person"
                    ><span itemprop="name">{{
                      article.user.nickname
                    }}</span></nuxt-link
                  >发布于
                  <time
                    :datetime="
                      article.createTime | formatDate('yyyy-MM-ddTHH:mm:ss')
                    "
                    itemprop="datePublished"
                    >{{ article.createTime | prettyDate }}</time
                  >
                </span>
              </div>
            </div>

            <div
              v-lazy-container="{ selector: 'img' }"
              class="article-content content line-numbers"
              itemprop="articleBody"
              v-html="article.content"
            ></div>

            <!--节点、标签-->
            <div class="article-tags">
              <nuxt-link
                v-for="tag in article.tags"
                :key="tag.tagId"
                :to="'/articles/' + tag.tagId"
                class="article-tag"
                >#{{ tag.tagName }}</nuxt-link
              >
            </div>
          </div>
        </article>

        <!-- 评论 -->
        <comment
          :entity-id="article.articleId"
          :comments-page="commentsPage"
          entity-type="article"
        />
      </div>
      <div class="right-container">
        <user-info :user="article.user" />
        <div
          v-if="relatedArticles && relatedArticles.length"
          class="widget no-margin"
        >
          <div class="widget-header">相关文章</div>
          <div class="widget-content article-related">
            <ul>
              <li v-for="a in relatedArticles" :key="a.articleId">
                <nuxt-link
                  :to="'/article/' + a.articleId"
                  :title="a.title"
                  class="article-related-title"
                  target="_blank"
                  >{{ a.title }}</nuxt-link
                >
              </li>
            </ul>
          </div>
        </div>

        <div v-if="nearlyArticles && nearlyArticles.length" class="widget">
          <div class="widget-header">近期文章</div>
          <div class="widget-content article-related">
            <ul>
              <li v-for="a in nearlyArticles" :key="a.articleId">
                <nuxt-link
                  :to="'/article/' + a.articleId"
                  :title="a.title"
                  class="article-related-title"
                  target="_blank"
                  >{{ a.title }}</nuxt-link
                >
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import Prism from 'prismjs'
import UserHelper from '~/common/UserHelper'

export default {
  async asyncData({ $axios, params, error }) {
    let article
    try {
      article = await $axios.get('/api/article/' + params.id)
    } catch (e) {
      error({
        message: e.message,
      })
      return
    }
    const [commentsPage, nearlyArticles, relatedArticles] = await Promise.all([
      $axios.get('/api/comment/comments', {
        params: {
          entityType: 'article',
          entityId: article.articleId,
        },
      }),
      $axios.get('/api/article/nearly/' + article.articleId),
      $axios.get('/api/article/related/' + article.articleId),
    ])

    // 文章关键字
    let keywords = ''
    const keywordsArr = []
    if (article.tags && article.tags.length > 0) {
      article.tags.forEach((tag) => {
        keywordsArr.push(tag.tagName)
      })
      if (keywordsArr.length > 0) {
        keywords = keywordsArr.join(',')
      }
    }

    // 文章描述
    let description = ''
    if (article.summary && article.summary.length > 0) {
      description = article.summary.substr(0, 128)
      if (article.summary.length > 128) {
        description += '...'
      }
    }

    return {
      article,
      nearlyArticles,
      relatedArticles,
      commentsPage,
      keywords,
      description,
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.article.title),
      meta: [
        { hid: 'description', name: 'description', content: this.description },
        { hid: 'keywords', name: 'keywords', content: this.keywords },
      ],
    }
  },
  computed: {
    hasPermission() {
      return this.isOwner || this.isOwnerOrAdmin
    },
    isOwnerOrAdmin() {
      return UserHelper.isOwner(this.user) || UserHelper.isAdmin(this.user)
    },
    isPending() {
      return this.article.status === 2
    },
    user() {
      return this.$store.state.user.current
    },
  },
  mounted() {
    Prism.highlightAll()
  },
  methods: {
    deleteArticle(articleId) {
      if (!process.client) {
        return
      }
      const me = this
      this.$confirm('是否确认删除该文章？').then(function () {
        me.$axios
          .post('/api/article/delete/' + articleId)
          .then(() => {
            me.$msg({
              message: '删除成功',
              onClose() {
                me.$linkTo('/articles')
              },
            })
          })
          .catch((e) => {
            me.$message.error('删除失败：' + (e.message || e))
          })
      })
    },
    async addFavorite(articleId) {
      try {
        if (this.article.favorited) {
          await this.$axios.post('/api/favorite/delete', {
            params: {
              entityType: 'article',
              entityId: articleId,
            },
          })
          this.article.favorited = false
          this.$message.success('已取消收藏')
        } else {
          await this.$axios.post('/api/favorite/add', {
            entityType: 'article',
            entityId: articleId,
          })
          this.article.favorited = true
          this.$message.success('收藏成功')
        }
      } catch (e) {
        console.error(e)
        this.$message.error('收藏失败：' + (e.message || e))
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.article-list {
  margin: 0 !important;

  li {
    padding: 12px 12px;
    display: flex;
    position: relative;
    overflow: hidden;
    transition: background 0.5s;
    border-radius: 3px;
    background: var(--bg-color);

    &:not(:last-child) {
      margin-bottom: 10px;
    }
  }
}

.article-item {
  overflow: hidden;
  zoom: 1;
  line-height: 24px;

  .article-title {
    a {
      font-size: 18px;
      line-height: 30px;
      font-weight: 500;
      color: var(--text-color);
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  // 详情页
  &.article-detail {
    .article-header {
      padding: 10px 0;
      border-bottom: 1px solid var(--border-color);
    }

    .article-title-wrapper {
      display: flex;
      .article-title {
        width: 100%;
        color: var(--text-color);
        font-weight: normal;
        overflow: hidden;
        text-overflow: ellipsis;
        font-size: 18px;
        line-height: 30px;
      }
      .article-manage-menu {
        min-width: max-content;
      }
    }

    .article-tags {
      margin-top: 10px;
      .article-tag {
        height: 25px;
        padding: 0 8px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
        border-radius: 12.5px;
        margin-right: 10px;
        background: var(--bg-color2);
        border: 1px solid var(--border-color);
        color: var(--text-color3);
        font-size: 12px;

        &:hover {
          color: var(--text-link-color);
          background: var(--bg-color);
          border: 1px solid var(--border-hover-color);
        }
      }
    }
  }

  .article-summary {
    font-size: 14px;
    color: var(--text-color2);
    overflow: hidden;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 3;
    text-align: justify;
    padding-top: 6px;
    word-break: break-all;
    text-overflow: ellipsis;
  }

  .article-meta {
    display: inline-block;
    font-size: 13px;
    padding-top: 6px;

    .article-meta-item {
      padding: 0 6px 0 0;
      color: var(--text-color3);

      a {
        color: var(--text-link-color);

        &.article-author {
          font-weight: bold;
          padding: 0 3px;
        }
      }

      .article-tag {
        height: 25px;
        padding: 0 8px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
        border-radius: 12.5px;
        background: var(--bg-color2);
        border: 1px solid var(--border-color);
        color: var(--text-color3);
        font-size: 12px;

        &:hover {
          color: var(--text-link-color);
          background: var(--bg-color);
          border: 1px solid var(--border-hover-color);
        }

        &:not(:last-child) {
          margin-right: 10px;
        }
      }
    }
  }

  .article-tool {
    display: inline-block;
    margin-right: 5px;
    line-height: 32px;

    & > span {
      margin-left: 5px;

      a {
        font-size: 12px;
        color: var(--text-color3);
        font-weight: 700;

        &:hover {
          text-decoration: underline;
        }

        i {
          font-size: 12px;
          color: var(--text-color);
        }
      }
    }
  }

  .article-content {
    font-size: 15px;
    margin-top: 10px;
    margin-bottom: 10px;

    a.article-share-summary {
      color: var(--text-color);
    }
  }

  .article-footer {
    word-break: break-all;
    background: var(--bg-color);
    padding: 10px;

    &,
    a {
      color: var(--text-color);
      font-size: 14px;
    }
  }
}

.article-related {
  margin-top: 0 !important;

  li {
    // margin: 8px 0;
    padding: 5px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }
  }

  .article-related-title {
    overflow: hidden;
    word-break: break-all;
    text-overflow: ellipsis;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    display: -webkit-box;

    color: var(--text-color2);
    font-size: 14px;

    &:hover {
      color: var(--text-link-color);
    }
  }
}
</style>
