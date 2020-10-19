<template>
  <div>
    <section class="main">
      <div class="container main-container left-main">
        <div class="left-container">
          <div class="main-content">
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
                      <time
                        :datetime="
                          topic.lastCommentTime
                            | formatDate('yyyy-MM-ddTHH:mm:ss')
                        "
                        itemprop="datePublished"
                        >{{ topic.lastCommentTime | prettyDate }}</time
                      >
                    </span>
                    <span class="meta-item">
                      <a
                        v-if="topic.node"
                        :href="'/topics/node/' + topic.node.nodeId"
                        class="node"
                        >{{ topic.node.name }}</a
                      >
                    </span>
                    <span class="meta-item">
                      <span
                        v-for="tag in topic.tags"
                        :key="tag.tagId"
                        class="tag"
                      >
                        <a :href="'/topics/tag/' + tag.tagId">{{
                          tag.tagName
                        }}</a>
                      </span>
                    </span>
                    <span v-if="hasPermission" class="meta-item act">
                      <a @click="deleteTopic(topic.topicId)">
                        <i class="iconfont icon-delete" />&nbsp;删除
                      </a>
                    </span>
                    <span v-if="hasPermission" class="meta-item act">
                      <a :href="'/topic/edit/' + topic.topicId">
                        <i class="iconfont icon-edit" />&nbsp;修改
                      </a>
                    </span>
                    <span class="meta-item act">
                      <a @click="addFavorite(topic.topicId)">
                        <i class="iconfont icon-favorite" />&nbsp;{{
                          favorited ? '已收藏' : '收藏'
                        }}
                      </a>
                    </span>
                  </div>
                </div>
                <div class="topic-header-right">
                  <span class="count"
                    >{{ topic.commentCount }}&nbsp;/&nbsp;{{
                      topic.viewCount
                    }}</span
                  >
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

              <div class="content topic-content" itemprop="articleBody">
                <div
                  v-lazy-container="{ selector: 'img' }"
                  v-html="topic.content"
                ></div>
              </div>
            </article>
            <div class="main-content-footer">
              <div class="topic-toolbar">
                <template v-if="likeUsers && likeUsers.length">
                  <a
                    v-for="likeUser in likeUsers"
                    :key="likeUser.id"
                    :href="'/user/' + likeUser.id"
                    :alt="likeUser.nickname"
                    target="_blank"
                    class="avatar-link"
                  >
                    <img
                      :src="likeUser.smallAvatar"
                      :alt="likeUser.nickname"
                      class="avatar"
                    />
                  </a>
                </template>
                <span class="like-desc">有{{ topic.likeCount }}人点赞</span>

                <div class="action-buttons">
                  <div
                    :class="{ active: topic.liked }"
                    class="action like"
                    @click="like(topic)"
                  >
                    <i class="iconfont icon-like" />
                  </div>
                  <div
                    :class="{ active: favorited }"
                    class="action favorite"
                    @click="addFavorite(topic.topicId)"
                  >
                    <i class="iconfont icon-favorite" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 评论 -->
          <comment
            :entity-id="topic.topicId"
            :comments-page="commentsPage"
            :comment-count="topic.commentCount"
            :show-ad="false"
            entity-type="topic"
          />
        </div>
        <div class="right-container">
          <a
            class="button is-success"
            href="/topic/create"
            style="width: 100%;"
          >
            <span class="icon"><i class="iconfont icon-topic" /></span>
            <span>发表话题</span>
          </a>

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
import UserHelper from '~/common/UserHelper'

export default {
  components: {
    Comment,
    UserInfo,
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
    hasPermission() {
      return (
        this.isOwner ||
        UserHelper.isOwner(this.user) ||
        UserHelper.isAdmin(this.user)
      )
    },
    isOwner() {
      if (!this.user || !this.topic) {
        return false
      }
      return this.user.id === this.topic.user.id
    },
    user() {
      return this.$store.state.user.current
    },
  },
  mounted() {},
  methods: {
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
    async deleteTopic(topicId) {
      if (process.client && !window.confirm('是否确认删除该话题？')) {
        return
      }
      try {
        await this.$axios.post('/api/topic/delete/' + topicId)
        this.$toast.success('删除成功', {
          duration: 2000,
          onComplete() {
            utils.linkTo('/topics')
          },
        })
      } catch (e) {
        console.error(e)
        this.$toast.error('删除失败：' + (e.message || e))
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
    }
  },
}
</script>

<style lang="scss" scoped></style>
