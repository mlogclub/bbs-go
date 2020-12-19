<template>
  <div>
    <section class="main">
      <div class="container main-container left-main">
        <div class="left-container">
          <div class="main-content no-padding">
            <article
              class="topic-detail topic-wrap"
              itemscope
              itemtype="http://schema.org/BlogPosting"
            >
              <div class="topic-header">
                <div class="topic-header-left">
                  <a
                    :href="'/user/' + topic.user.id"
                    :title="topic.user.nickname"
                  >
                    <img :src="topic.user.smallAvatar" class="avatar size-45" />
                  </a>
                </div>
                <div class="topic-header-center">
                  <h1 class="topic-title" itemprop="headline">
                    {{ topic.title }}
                  </h1>
                  <div class="topic-meta">
                    <span
                      class="meta-item"
                      itemprop="author"
                      itemscope
                      itemtype="http://schema.org/Person"
                    >
                      <a :href="'/user/' + topic.user.id" itemprop="name">{{
                        topic.user.nickname
                      }}</a>
                    </span>
                    <span class="meta-item">
                      发布于
                      <time
                        :datetime="
                          topic.lastCommentTime
                            | formatDate('yyyy-MM-ddTHH:mm:ss')
                        "
                        itemprop="datePublished"
                        >{{ topic.lastCommentTime | prettyDate }}</time
                      >
                    </span>
                  </div>
                </div>
                <div class="topic-header-right">
                  <topic-manage-menu :topic="topic" />
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

              <!--内容-->
              <div class="topic-content content" itemprop="articleBody">
                <div
                  v-lazy-container="{ selector: 'img' }"
                  v-html="topic.content"
                ></div>
              </div>

              <!--节点、标签-->
              <div class="topic-tags">
                <a
                  v-if="topic.node"
                  :href="'/topics/node/' + topic.node.nodeId"
                  class="topic-tag"
                  >{{ topic.node.name }}</a
                >
                <a
                  v-for="tag in topic.tags"
                  :key="tag.tagId"
                  :href="'/topics/tag/' + tag.tagId"
                  class="topic-tag"
                  >#{{ tag.tagName }}</a
                >
              </div>

              <!-- 点赞用户列表 -->
              <div
                v-if="likeUsers && likeUsers.length"
                class="topic-like-users"
              >
                <a
                  v-for="likeUser in likeUsers"
                  :key="likeUser.id"
                  :href="'/user/' + likeUser.id"
                  :alt="likeUser.nickname"
                  target="_blank"
                  class="like-user"
                >
                  <img
                    :src="likeUser.smallAvatar"
                    :alt="likeUser.nickname"
                    class="avatar"
                  />
                </a>
              </div>

              <!-- 功能按钮 -->
              <div class="topic-actions">
                <a class="action disabled">
                  <i class="action-icon iconfont icon-read" />
                  <span class="content">
                    <span>浏览</span>
                    <span v-if="topic.viewCount > 0">
                      ({{ topic.viewCount }})
                    </span>
                  </span>
                </a>
                <a class="action" @click="like(topic)">
                  <i class="action-icon iconfont icon-like" />
                  <span class="content">
                    <span>点赞</span>
                    <span v-if="topic.likeCount > 0">
                      ({{ topic.likeCount }})
                    </span>
                  </span>
                </a>
                <a class="action" @click="addFavorite(topic.topicId)">
                  <i
                    class="action-icon iconfont"
                    :class="{
                      'icon-has-favorite': favorited,
                      'icon-favorite': !favorited,
                    }"
                  />
                  <span class="content">
                    <span>收藏</span>
                  </span>
                </a>
              </div>

              <!-- 评论 -->
              <comment
                :entity-id="topic.topicId"
                :comments-page="commentsPage"
                :comment-count="topic.commentCount"
                :show-ad="false"
                entity-type="topic"
              />
            </article>
          </div>
        </div>
        <div class="right-container">
          <user-info :user="topic.user" />

          <div class="ad">
            <!-- 展示广告 -->
            <adsbygoogle ad-slot="1742173616" />
          </div>

          <div v-if="topic.toc" ref="toc" class="widget no-bg toc">
            <div class="widget-header">
              目录
            </div>
            <div class="widget-content" v-html="topic.toc" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
import utils from '~/common/utils'
import Comment from '~/components/Comment'
import UserInfo from '~/components/UserInfo'
import TopicManageMenu from '~/components/topic/TopicManageMenu'

export default {
  components: {
    Comment,
    UserInfo,
    TopicManageMenu,
  },
  async asyncData({ $axios, params, error }) {
    let topic
    try {
      topic = await $axios.get('/api/topic/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '话题不存在',
      })
      return
    }

    const [favorited, commentsPage, likeUsers] = await Promise.all([
      $axios.get('/api/favorite/favorited', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/comment/list', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
      $axios.get('/api/topic/recentlikes/' + params.id),
    ])

    return {
      topic,
      commentsPage,
      favorited: favorited.favorited,
      likeUsers,
    }
  },
  computed: {
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
    async addFavorite(topicId) {
      try {
        if (this.favorited) {
          await this.$axios.get('/api/favorite/delete', {
            params: {
              entityType: 'topic',
              entityId: topicId,
            },
          })
          this.favorited = false
          this.$toast.success('已取消收藏！')
        } else {
          await this.$axios.get('/api/topic/favorite/' + topicId)
          this.favorited = true
          this.$toast.success('收藏成功')
        }
      } catch (e) {
        console.error(e)
        this.$toast.error('收藏失败：' + (e.message || e))
      }
    },
    async like(topic) {
      try {
        await this.$axios.post('/api/topic/like/' + topic.topicId)
        topic.liked = true
        topic.likeCount++
        this.likeUsers = this.likeUsers || []
        this.likeUsers.unshift(this.$store.state.user.current)
      } catch (e) {
        if (e.errorCode === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              },
            },
          })
        } else {
          topic.liked = true
          this.$toast.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$siteTitle(this.topic.title),
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
