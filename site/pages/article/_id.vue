<template>
  <section class="main">
    <div class="container">
      <div class="left-main-container">
        <div class="m-left">
          <article class="article-item">
            <div class="article-item-left">
              <a
                :href="'/user/' + article.user.id"
                :title="article.user.nickname"
                target="_blank"
              >
                <div
                  class="avatar is-rounded"
                  :style="{'backgroundImage':'url(' + article.user.avatar + ')'}"
                />
              </a>
            </div>

            <div class="article-item-right">
              <div class="article-title">
                <a :href="'/article/' + article.articleId">{{ article.title }}</a>
              </div>

              <div class="article-meta">
                <span class="article-meta-item">
                  <a :href="'/user/' + article.user.id">{{ article.user.nickname }}</a>
                </span>

                <span v-if="article.category" class="article-meta-item">
                  <a
                    :href="'/articles/cat/' + article.category.categoryId"
                  >{{ article.category.categoryName }}</a>
                </span>

                <span v-for="tag in article.tags" :key="tag.tagId" class="article-meta-item">
                  <a :href="'/articles/tag/' + tag.tagId">{{ tag.tagName }}</a>
                </span>

                <span class="article-meta-item">
                  <time itemprop="datePublished">{{ article.createTime | prettyDate }}</time>
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
                    <i class="iconfont icon-favorite" />{{ favorited ? '已收藏' : '收藏' }}
                  </a>
                </span>
              </div>
            </div>

            <!--
              <div v-if="article.share" class="article-content content">
                <a
                  class="article-share-summary"
                  :href="'/article/redirect/' + article.articleId"
                  target="_blank"
                  v-html="article.summary"
                />
                <a :href="'/article/redirect/' + article.articleId" target="_blank">点击阅读原文>></a>
                <ins
                  class="adsbygoogle"
                  style="display:block"
                  data-ad-format="fluid"
                  data-ad-layout-key="-ig-s+1x-t-q"
                  data-ad-client="ca-pub-5683711753850351"
                  data-ad-slot="4728140043"
                />
                <script>
                  (adsbygoogle = window.adsbygoogle || []).push({});
                </script>
              </div>
              <div v-else class="article-content content">
                <p v-highlight v-html="article.content" />
                <ins
                  class="adsbygoogle"
                  style="display:block"
                  data-ad-format="fluid"
                  data-ad-layout-key="-ig-s+1x-t-q"
                  data-ad-client="ca-pub-5683711753850351"
                  data-ad-slot="4728140043"
                />
                <script>
                  (adsbygoogle = window.adsbygoogle || []).push({});
                </script>
              </div>
              -->
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
                (adsbygoogle = window.adsbygoogle || []).push({});
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
                  <li v-if="article.user.type != 1 && article.sourceUrl">
                    <strong>原文地址：</strong>
                    <a href="javascript:void(0)">{{ article.sourceUrl }}</a>
                  </li>
                  <li>
                    <strong>免责声明：</strong>
                    我们尊重原创，也注重分享。版权原作者所有，如有侵犯您的权益请及时联系（
                    <a
                      href="mailto:mlog1@qq.com"
                    >mlog1@qq.com</a>），我们将在24小时之内删除。
                  </li>
                </ul>
              </blockquote>
            </div>
          </article>

          <!-- 评论 -->
          <comment entity-type="article" :entity-id="article.articleId" :show-ad="false" />

          <div class="columns article-related">
            <div class="column">
              <div v-if="newestArticles && newestArticles.length" class="widget">
                <div class="header">
                  最新文章
                </div>
                <div class="content">
                  <ul>
                    <li v-for="a in newestArticles" :key="a.articleId">
                      <a
                        :href="'/article/' + a.articleId"
                        :title="a.title"
                        target="_blank"
                      >{{ a.title }}</a>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <div v-if="relatedArticles && relatedArticles.length" class="column">
              <div class="widget">
                <div class="header">
                  相关文章
                </div>
                <div class="content">
                  <ul>
                    <li v-for="a in relatedArticles" :key="a.articleId">
                      <a
                        :href="'/article/' + a.articleId"
                        :title="a.title"
                        target="_blank"
                      >{{ a.title }}</a>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="m-right">
          <!-- 展示广告190x190 -->
          <ins
            class="adsbygoogle"
            style="display:inline-block;width:190px;height:190px"
            data-ad-client="ca-pub-5683711753850351"
            data-ad-slot="5685455263"
          />
          <script>
            (adsbygoogle = window.adsbygoogle || []).push({});
          </script>

          <div v-if="article.toc" class="toc widget">
            <div class="header">
              目录
            </div>
            <div class="content" v-html="article.toc" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
import Comment from '~/components/Comment'
export default {
  components: {
    Comment
  },
  computed: {
    isOwner: function () {
      return this.currentUser && this.article && this.currentUser.id === this.article.user.id
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
  },
  async asyncData({ $axios, params, error }) {
    let article
    try {
      article = await $axios.get('/api/article/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '文章不存在'
      })
      return
    }
    const currentUser = await $axios.get('/api/user/current')
    const favorited = await $axios.get('/api/favorite/favorited', {
      params: {
        entityType: 'article',
        entityId: params.id
      }
    })
    const [newestArticles, relatedArticles] = await Promise.all([
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
      currentUser: currentUser,
      article: article,
      favorited: favorited.favorited,
      newestArticles: newestArticles,
      relatedArticles: relatedArticles,
      keywords: keywords,
      description: description
    }
  },
  mounted() {
    utils.handleToc()
  },
  methods: {
    async deleteArticle(articleId) {
      try {
        await this.$axios.post('/api/article/delete/' + articleId)
        this.$toast.success('删除成功！', {
          duration: 1000,
          onComplete: function () {
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
  }
}
</script>

<style lang="scss" scoped>
.article-item-left {
  width: 50px;
  height: 50px;
  float: left;
  vertical-align: middle;
}

.article-item-right {
  margin-left: 50px;
  padding-left: 10px;
  vertical-align: middle;
}

article {
  .article-title {
    a {
      // color: #999;
      color: #0f0f0f;
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

  .article-tool {
    display: inline-block;
    font-size: 13px;

    &>span{
      margin-left: 10px;
    }

    a {
      font-size: 12px;
      i {
        font-size: 12px;
        color: #000
      }
    }
  }
}

.article-content {
  margin-top: 20px;

  a.article-share-summary {
    color: #4a4a4a;
  }
}

.article-footer {
  margin-top: 20px;
  word-break: break-all;

  blockquote {
    padding: 10px 15px;
    margin: 0 0 20px;
    border: 1px dotted #eeeeee;
    border-left: 3px solid #eeeeee;
    background-color: #fbfbfb;
  }

  blockquote p:last-child,
  blockquote ul:last-child,
  blockquote ol:last-child {
    margin-bottom: 0;
  }

  blockquote footer,
  blockquote small,
  blockquote .small {
    display: block;
    font-size: 80%;
    line-height: 1.42857;
    color: #777777;
  }

  blockquote footer:before,
  blockquote small:before,
  blockquote .small:before {
    content: '\2014 \00A0';
  }

  blockquote.pull-right {
    padding-right: 15px;
    padding-left: 0;
    border-right: 5px solid #eeeeee;
    border-left: 0;
    text-align: right;
  }

  blockquote.pull-right footer:before,
  blockquote.pull-right small:before,
  blockquote.pull-right .small:before {
    content: '';
  }

  blockquote.pull-right footer:after,
  blockquote.pull-right small:after,
  blockquote.pull-right .small:after {
    content: '\00A0 \2014';
  }

}

.article-related {
  margin-top: 20px;

  .widget > .header {
    color: #000;
    font-size: 20px;
  }
}
</style>
