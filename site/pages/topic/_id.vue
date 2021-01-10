<template>
  <div>
    <section class="main">
      <div class="container main-container left-main size-360">
        <div class="left-container">
          <div class="main-content no-padding">
            <article
              class="topic-detail"
              itemscope
              itemtype="http://schema.org/BlogPosting"
            >
              <div class="topic-header">
                <div class="topic-header-left">
                  <avatar :user="topic.user" size="45" />
                </div>
                <div class="topic-header-center">
                  <div class="topic-nickname" itemprop="headline">
                    <a
                      itemprop="author"
                      itemscope
                      itemtype="http://schema.org/Person"
                      :href="'/user/' + topic.user.id"
                      >{{ topic.user.nickname }}</a
                    >
                  </div>
                  <div class="topic-meta">
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
                <h1 v-if="topic.title" class="topic-title" itemprop="headline">
                  {{ topic.title }}
                </h1>
                <div
                  v-lazy-container="{ selector: 'img' }"
                  class="topic-content-detail"
                  v-html="topic.content"
                ></div>
                <ul
                  v-if="topic.imageList && topic.imageList.length"
                  v-viewer
                  class="topic-image-list"
                >
                  <li v-for="(image, index) in topic.imageList" :key="index">
                    <div class="image-item">
                      <img :src="image.preview" :data-src="image.url" />
                    </div>
                  </li>
                </ul>
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
                <avatar
                  v-for="likeUser in likeUsers"
                  :key="likeUser.id"
                  :user="likeUser"
                  :round="true"
                  size="30"
                />
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
                <a
                  class="action"
                  :class="{ disabled: liked }"
                  @click="like(topic)"
                >
                  <i
                    class="action-icon iconfont icon-like"
                    :class="{ 'checked-icon': liked }"
                  />
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
                      'checked-icon': favorited,
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
                :mode="topic.type === 1 ? 'text' : 'markdown'"
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
import Vue from 'vue'
import Viewer from 'v-viewer'
import Comment from '~/components/Comment'
import UserInfo from '~/components/UserInfo'
import TopicManageMenu from '~/components/topic/TopicManageMenu'
import Avatar from '~/components/Avatar'
import 'viewerjs/dist/viewer.css'

Vue.use(Viewer, {
  defaultOptions: {
    zIndex: 9999,
    navbar: false,
    title: false,
    tooltip: false,
    movable: false,
    scalable: false,
    url: 'data-src',
  },
})

export default {
  components: {
    Comment,
    UserInfo,
    TopicManageMenu,
    Avatar,
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

    const [liked, favorited, commentsPage, likeUsers] = await Promise.all([
      $axios.get('/api/like/liked', {
        params: {
          entityType: 'topic',
          entityId: params.id,
        },
      }),
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
      liked: liked.liked,
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
          this.$message.success('已取消收藏！')
        } else {
          await this.$axios.get('/api/topic/favorite/' + topicId)
          this.favorited = true
          this.$message.success('收藏成功')
        }
      } catch (e) {
        console.error(e)
        this.$message.error('收藏失败：' + (e.message || e))
      }
    },
    async like(topic) {
      try {
        if (this.liked) {
          return
        }
        await this.$axios.post('/api/topic/like/' + topic.topicId)
        this.liked = true
        topic.likeCount++
        this.likeUsers = this.likeUsers || []
        this.likeUsers.unshift(this.$store.state.user.current)
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.liked = true
          this.$message.error(e.message || e)
        }
      }
    },
  },
  head() {
    return {
      title: this.$topicSiteTitle(this.topic),
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

<style lang="scss" scoped>
@import './assets/styles/variable.scss';

.topic-detail {
  margin-bottom: 20px;

  .topic-header {
    //padding: 10px;
    display: flex;
    margin: 0 10px;
    border-bottom: 1px solid $line-color-base;

    @media screen and (max-width: 1024px) {
      .topic-header-right {
        display: none;
      }
    }

    .topic-header-left {
      margin: 10px 10px 0 0;
      //width: 50px;
      //height: 50px;
    }

    .topic-header-center {
      margin: 10px 10px 0 0;
      width: 100%;

      .topic-nickname a {
        color: #555;
        font-size: 14px;
        font-weight: bold;
        overflow: hidden;
      }

      .topic-meta {
        position: relative;
        font-size: 12px;
        line-height: 24px;
        color: #70727c;
        margin-top: 5px;

        span.meta-item {
          font-size: 12px;

          &:not(:last-child) {
            margin-right: 8px;
          }
        }
      }
    }

    .topic-header-right {
      min-width: 42px;
    }
  }

  .topic-content,
  .topic-tags,
  .topic-like-users,
  .topic-actions {
    margin: 20px 12px;
  }

  .topic-content {
    font-size: 15px;
    color: #000;
    white-space: normal;
    word-break: break-all;
    word-wrap: break-word;
    padding-top: 0 !important;
    margin: 0 12px;

    .topic-title {
      font-weight: 700;
      font-size: 20px;
      word-wrap: break-word;
      word-break: normal;
    }

    .topic-content-detail {
      font-size: 16px;
      line-height: 24px;
    }

    .topic-image-list {
      margin-left: 0;
      margin-top: 10px;

      li {
        cursor: pointer;
        border: 1px dashed #ddd;
        text-align: center;

        // 图片尺寸
        $image-size: 120px;

        display: inline-block;
        vertical-align: middle;
        width: $image-size;
        height: $image-size;
        line-height: $image-size;
        margin: 0 8px 8px 0;
        background-color: #e8e8e8;
        background-size: 32px 32px;
        background-position: 50%;
        background-repeat: no-repeat;
        overflow: hidden;
        position: relative;

        .image-item {
          display: block;
          width: $image-size;
          height: $image-size;
          overflow: hidden;
          transform-style: preserve-3d;

          & > img {
            width: 100%;
            height: 100%;
            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }
      }
    }

    nav {
      background-color: #fdfdfd;
      border: 1px solid #f6f6f6;
      padding: 10px 0;
      font-size: 14px;

      ul {
        list-style: disc outside;
        margin-left: 2em;
        margin-top: 0;
      }
    }
  }

  .topic-tags {
    .topic-tag {
      height: 25px;
      padding: 0 8px;
      display: inline-flex;
      justify-content: center;
      align-items: center;
      border-radius: 12.5px;
      margin-right: 10px;
      background: #f7f7f7;
      border: 1px solid #f7f7f7;
      color: #777;
      font-size: 12px;

      &:hover {
        color: #1878f3;
        background: #fff;
        border: 1px solid #1878f3;
      }
    }
  }

  .topic-like-users {
    width: 80%;
    margin: 0 auto;
    display: flex;
    flex-wrap: wrap;
    justify-content: center;

    .avatar-a {
      margin-right: 3px;
    }
  }

  .topic-actions {
    margin: 20px auto;
    padding: 0 25px;
    display: flex;
    justify-content: space-between;

    .action {
      background: #ffffff;
      cursor: pointer;
      flex: 1;
      display: flex;
      align-items: center;
      flex-direction: column;
      color: #8590a6;

      .checked-icon {
        color: red;
      }

      &.disabled {
        cursor: not-allowed;

        &:hover {
          color: #8590a6;

          > .action-icon {
            fill: #8590a6;
          }
        }
      }

      > .action-icon {
        font-size: 30px;
        fill: #8590a6;
      }

      &:hover {
        color: #1878f3;

        > .action-icon {
          fill: #1878f3;
        }
      }

      > .content {
        margin-top: 10px;
        font-size: 12px;
      }
    }
  }
}
</style>
