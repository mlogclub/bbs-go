<template>
  <section class="main">
    <div class="container-wrapper main-container left-main">
      <div class="left-container">
        <article class="article-item article-detail">
          <div class="article-header">
            <div class="article-item-left">
              <a
                :href="'/user/' + article.user.id"
                :title="article.user.nickname"
                target="_blank"
              >
                <div
                  :style="{
                    backgroundImage: 'url(' + article.user.avatar + ')'
                  }"
                  class="avatar is-rounded"
                />
              </a>
            </div>
            <div class="article-item-right">
              <div class="article-title">{{ article.title }}</div>

              <div class="article-meta">
                <span class="article-meta-item">
                  由
                  <a :href="'/user/' + article.user.id" class="article-author"
                    >&nbsp;{{ article.user.nickname }}&nbsp;</a
                  >发布于
                  <time itemprop="datePublished">{{
                    article.createTime | prettyDate
                  }}</time>
                </span>

                <span v-if="article.category" class="article-meta-item">
                  <span class="article-tag tag">
                    <a :href="'/articles/cat/' + article.category.categoryId">{{
                      article.category.categoryName
                    }}</a>
                  </span>
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
                    <a :href="'/articles/tag/' + tag.tagId" class>{{
                      tag.tagName
                    }}</a>
                  </span>
                </span>
              </div>

              <div class="article-tool">
                <span v-if="isOwner">
                  <a @click="deleteArticle(article.articleId)">
                    <i class="iconfont icon-delete" />删除
                  </a>
                </span>
                <span v-if="isOwner">
                  <a :href="'/article/edit/' + article.articleId">
                    <i class="iconfont icon-edit" />修改
                  </a>
                </span>
                <span>
                  <a @click="addFavorite(article.articleId)">
                    <i class="iconfont icon-favorite" />{{
                      favorited ? '已收藏' : '收藏'
                    }}
                  </a>
                </span>
              </div>
            </div>
          </div>

          <div class="article-content content">
            <ins
              class="adsbygoogle"
              style="display:block"
              data-ad-format="fluid"
              data-ad-layout-key="-ig-s+1x-t-q"
              data-ad-client="ca-pub-5683711753850351"
              data-ad-slot="4728140043"
            />
            <script>
              ;(adsbygoogle = window.adsbygoogle || []).push({})
            </script>

            <p v-highlight v-html="article.content" />
          </div>

          <div class="article-footer">
            <blockquote v-if="article.share">
              <ul>
                <li v-if="article.user.type == 1">
                  <strong>转载自公众号：</strong>
                  <a href="javascript:void(0)">{{ article.user.nickname }}</a>
                </li>
                <li>
                  <strong>免责声明：</strong>
                  我们尊重原创，也注重分享。版权原作者所有，如有侵犯您的权益请及时联系（
                  <a href="mailto:mlog1@qq.com">mlog1@qq.com</a
                  >），我们将在24小时之内删除。
                </li>
              </ul>
            </blockquote>
          </div>
        </article>

        <!-- 评论 -->
        <comment
          :entity-id="article.articleId"
          :comments-page="commentsPage"
          :show-ad="true"
          entity-type="article"
        />

        <div class="columns article-related">
          <div class="column">
            <div v-if="newestArticles && newestArticles.length" class="widget">
              <div class="widget-header">最新文章</div>
              <div class="widget-content">
                <ul>
                  <li v-for="a in newestArticles" :key="a.articleId">
                    <a
                      :href="'/article/' + a.articleId"
                      :title="a.title"
                      target="_blank"
                      >{{ a.title }}</a
                    >
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div class="column">
            <div
              v-if="relatedArticles && relatedArticles.length"
              class="widget"
            >
              <div class="widget-header">相关文章</div>
              <div class="widget-content">
                <ul>
                  <li v-for="a in relatedArticles" :key="a.articleId">
                    <a
                      :href="'/article/' + a.articleId"
                      :title="a.title"
                      target="_blank"
                      >{{ a.title }}</a
                    >
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="right-container">
        <weixin-gzh />

        <!-- 展示广告190x190 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:190px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="5685455263"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>

        <div ref="toc" v-if="article.toc" class="widget no-bg toc">
          <div class="widget-header">目录</div>
          <div v-html="article.toc" class="widget-content" />
        </div>

        <!-- 展示广告190x90 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:90px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="9345305153"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>

        <!-- 展示广告190x480 -->
        <ins
          class="adsbygoogle"
          style="display:inline-block;width:190px;height:480px"
          data-ad-client="ca-pub-5683711753850351"
          data-ad-slot="3438372357"
        />
        <script>
          ;(adsbygoogle = window.adsbygoogle || []).push({})
        </script>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import Comment from '~/components/Comment'
import WeixinGzh from '~/components/WeixinGzh'

export default {
  components: {
    Comment,
    WeixinGzh
  },
  async asyncData({ $axios, params, error }) {
    let article
    try {
      article = await $axios.get('/api/article/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '文章不存在，或已被删除'
      })
      return
    }
    const [
      commentsPage,
      favorited,
      newestArticles,
      relatedArticles
    ] = await Promise.all([
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'article',
          entityId: article.articleId
        }
      }),
      $axios.get('/api/favorite/favorited', {
        params: {
          entityType: 'article',
          entityId: params.id
        }
      }),
      $axios.get('/api/article/user/newest/' + article.user.id),
      $axios.get('/api/article/related/' + article.articleId)
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
      newestArticles,
      relatedArticles,
      commentsPage,
      keywords,
      description
    }
  },
  computed: {
    isOwner() {
      return (
        this.$store.state.user.current &&
        this.article &&
        this.$store.state.user.current.id === this.article.user.id
      )
    }
  },
  mounted() {
    utils.handleToc(this.$refs.toc)
  },
  methods: {
    async deleteArticle(articleId) {
      try {
        await this.$axios.post('/api/article/delete/' + articleId)
        this.$toast.success('删除成功！', {
          duration: 1000,
          onComplete() {
            utils.linkTo('/articles')
          }
        })
      } catch (e) {
        this.$toast.error('删除失败：' + (e.message || e))
      }
    },
    async addFavorite(articleId) {
      try {
        if (this.favorited) {
          await this.$axios.get('/api/favorite/delete', {
            params: {
              entityType: 'article',
              entityId: articleId
            }
          })
          this.favorited = false
          this.$toast.success('已取消收藏！')
        } else {
          await this.$axios.post('/api/article/favorite/' + articleId)
          this.favorited = true
          this.$toast.success('收藏成功！')
        }
      } catch (e) {
        console.error(e)
        this.$toast.error('收藏失败：' + (e.message || e))
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle(this.article.title),
      meta: [
        { hid: 'description', name: 'description', content: this.description },
        { hid: 'keywords', name: 'keywords', content: this.keywords }
      ]
    }
  }
}
</script>

<style lang="scss" scoped></style>
