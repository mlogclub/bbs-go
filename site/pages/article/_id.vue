<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <article
          class="article-item article-detail"
          itemscope
          itemtype="http://schema.org/BlogPosting"
        >
          <div class="main-content">
            <div class="article-header">
              <h1 class="article-title" itemprop="headline">
                {{ article.title }}
              </h1>
              <div class="article-meta">
                <span class="article-meta-item">
                  由
                  <a
                    :href="'/user/' + article.user.id"
                    class="article-author"
                    itemprop="author"
                    itemscope
                    itemtype="http://schema.org/Person"
                    ><span itemprop="name">{{ article.user.nickname }}</span></a
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
                    <a :href="'/articles/' + tag.tagId" class>{{
                      tag.tagName
                    }}</a>
                  </span>
                </span>

                <div class="article-tool">
                  <span v-if="hasPermission">
                    <a @click="deleteArticle(article.articleId)">
                      <i class="iconfont icon-delete" />&nbsp;删除
                    </a>
                  </span>
                  <span v-if="hasPermission">
                    <a :href="'/article/edit/' + article.articleId">
                      <i class="iconfont icon-edit" />&nbsp;修改
                    </a>
                  </span>
                  <span>
                    <a @click="addFavorite(article.articleId)">
                      <i class="iconfont icon-favorite" />&nbsp;{{
                        favorited ? '已收藏' : '收藏'
                      }}
                    </a>
                  </span>
                  <span v-if="isPending">
                    <a
                      href="javascript:void(0)"
                      style="cursor: default; text-decoration: none;"
                    >
                      <i class="iconfont icon-shenhe" />&nbsp;审核中
                    </a>
                  </span>
                </div>
              </div>
            </div>

            <div class="ad">
              <!-- 信息流广告 -->
              <adsbygoogle
                ad-slot="4980294904"
                ad-format="fluid"
                ad-layout-key="-ht-19-1m-3j+mu"
              />
            </div>

            <div
              v-lazy-container="{ selector: 'img' }"
              class="article-content content"
              itemprop="articleBody"
              v-html="article.content"
            ></div>

            <div class="ad">
              <!-- 展示广告 -->
              <adsbygoogle ad-slot="1742173616" />
            </div>

            <!-- 评论 -->
            <comment
              :entity-id="article.articleId"
              :comments-page="commentsPage"
              :show-ad="true"
              entity-type="article"
            />
          </div>
        </article>
      </div>
      <div class="right-container">
        <!-- 展示广告 -->
        <adsbygoogle ad-slot="1742173616" />

        <div
          v-if="relatedArticles && relatedArticles.length"
          class="widget no-margin"
        >
          <div class="widget-header">相关文章</div>
          <div class="widget-content article-related">
            <ul>
              <li v-for="a in relatedArticles" :key="a.articleId">
                <a
                  :href="'/article/' + a.articleId"
                  :title="a.title"
                  class="article-related-title"
                  target="_blank"
                  >{{ a.title }}</a
                >
              </li>
            </ul>
          </div>
        </div>

        <!-- 展示广告 -->
        <adsbygoogle ad-slot="1742173616" />

        <div v-if="nearlyArticles && nearlyArticles.length" class="widget">
          <div class="widget-header">近期文章</div>
          <div class="widget-content article-related">
            <ul>
              <li v-for="a in nearlyArticles" :key="a.articleId">
                <a
                  :href="'/article/' + a.articleId"
                  :title="a.title"
                  class="article-related-title"
                  target="_blank"
                  >{{ a.title }}</a
                >
              </li>
            </ul>
          </div>
        </div>

        <!-- 展示广告 -->
        <adsbygoogle ad-slot="1742173616" />
      </div>
    </div>
  </section>
</template>

<script>
import UserHelper from '~/common/UserHelper'
import Comment from '~/components/Comment'

export default {
  components: {
    Comment,
  },
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
    const [
      commentsPage,
      favorited,
      nearlyArticles,
      relatedArticles,
    ] = await Promise.all([
      $axios.get('/api/comment/list', {
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
  computed: {
    hasPermission() {
      return (
        this.isOwner ||
        UserHelper.isOwner(this.user) ||
        UserHelper.isAdmin(this.user)
      )
    },
    isOwner() {
      if (!this.user || !this.article) {
        return false
      }
      return this.user.id === this.article.user.id
    },
    isPending() {
      return this.article.status === 2
    },
    user() {
      return this.$store.state.user.current
    },
  },
  mounted() {
    this.initHighlight()
  },
  methods: {
    initHighlight() {
      if (process.client) {
        window.hljs.initHighlighting()
      }
    },
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
          this.$message.success('已取消收藏！')
        } else {
          await this.$axios.post('/api/article/favorite/' + articleId)
          this.favorited = true
          this.$message.success('收藏成功！')
        }
      } catch (e) {
        console.error(e)
        this.$message.error('收藏失败：' + (e.message || e))
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle(this.article.title),
      meta: [
        { hid: 'description', name: 'description', content: this.description },
        { hid: 'keywords', name: 'keywords', content: this.keywords },
      ],
      link: [
        {
          rel: 'stylesheet',
          href:
            '//cdn.staticfile.org/highlight.js/10.3.2/styles/github.min.css',
        },
      ],
      script: [
        {
          src: '//cdn.staticfile.org/highlight.js/10.3.2/highlight.min.js',
        },
      ],
    }
  },
}
</script>

<style lang="scss" scoped></style>
