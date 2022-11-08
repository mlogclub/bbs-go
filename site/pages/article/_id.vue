<template>
  <section class="main">
    <div v-if="isShowPendingStatus" class="container main-container">
      <div class="notification is-warning" style="width: 100%; margin: 20px 0">
        {{ $t('文章正在审核中，审核通过后展示') }}
      </div>
    </div>
    <div v-else class="container main-container left-main size-320">
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

                <span
                  v-if="article.tags && article.tags.length > 0"
                  class="article-meta-item"
                >
                  <span
                    v-for="tag in article.tags"
                    :key="tag.tagId"
                    class="article-tag tag"
                  >
                    <nuxt-link :to="'/articles/' + tag.tagId" class>{{
                      tag.tagName
                    }}</nuxt-link>
                  </span>
                </span>
              </div>
            </div>

            <div
              v-lazy-container="{ selector: 'img' }"
              class="article-content content line-numbers"
              itemprop="articleBody"
              v-html="article.content"
            ></div>
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
        statusCode: e.errorCode,
        message: e.message,
      })
      return
    }
    const [commentsPage, favorited, nearlyArticles, relatedArticles] =
      await Promise.all([
        $axios.get('/api/comment/comments', {
          params: {
            entityType: 'article',
            entityId: article.articleId,
          },
        }),
        $axios.get('/api/favorite/favorited', {
          params: {
            entityType: 'article',
            entityId: params.id,
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
      favorited: favorited.favorited,
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
    isShowPendingStatus() {
      if (this.article.status !== 2) {
        return false
      }
      // 用户没登陆
      if (!this.user) {
        return true
      }
      return this.user.id !== this.article.user.id && !this.isOwnerOrAdmin
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
        if (this.favorited) {
          await this.$axios.get('/api/favorite/delete', {
            params: {
              entityType: 'article',
              entityId: articleId,
            },
          })
          this.favorited = false
          this.$message.success('已取消收藏')
        } else {
          await this.$axios.post('/api/article/favorite/' + articleId)
          this.favorited = true
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

<style lang="scss" scoped></style>
