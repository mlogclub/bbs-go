<template>
  <div>
    <section class="main">
      <div class="container">
        <div class="left-main-container">
          <div class="m-left">
            <div class="topic-detail">
              <div class="header">
                <div class="left">
                  <a :href="'/user/' + topic.user.id" :title="topic.user.nickname">
                    <div class="avatar radius8" :style="{backgroundImage:'url(' + topic.user.avatar + ')'}" />
                  </a>
                </div>
                <div class="right">
                  <div class="topic-title">
                    {{ topic.title }}
                  </div>
                  <div class="topic-meta">
                    <span><a :href="'/user/' + topic.user.id">{{ topic.user.nickname }}</a></span>
                    <span>{{ topic.lastCommentTime | prettyDate }}</span>
                    <span>点击:{{ topic.viewCount }}</span>
                    <span v-for="tag in topic.tags" :key="tag.tagId" class="tag"><a :href="'/topics/tag/' + tag.tagId">{{ tag.tagName }}</a></span>
                    <span v-if="isOwner" class="act">
                      <a @click="deleteTopic(topic.topicId)">
                        <i class="iconfont icon-delete" />&nbsp;删除
                      </a>
                    </span>
                    <span v-if="isOwner" class="act">
                      <a :href="'/topic/edit/' + topic.topicId">
                        <i class="iconfont icon-edit" />&nbsp;修改
                      </a>
                    </span>
                    <span class="act">
                      <a @click="addFavorite(topic.topicId)">
                        <i class="iconfont icon-favorite" />&nbsp;{{ favorited ? '已收藏' : '收藏' }}
                      </a>
                    </span>
                  </div>
                </div>
              </div>

              <div v-highlight class="content" v-html="topic.content" />

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

            <!-- 评论 -->
            <comment entity-type="topic" :entity-id="topic.topicId" :show-ad="true" />
          </div>
          <div class="m-right">
            <!-- 展示广告190x90 -->
            <ins
              class="adsbygoogle"
              style="display:inline-block;width:190px;height:90px"
              data-ad-client="ca-pub-5683711753850351"
              data-ad-slot="9345305153"
            />
            <script>
              (adsbygoogle = window.adsbygoogle || []).push({});
            </script>

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

            <div v-if="topic.toc" ref="toc" class="toc widget">
              <div class="header">
                目录
              </div>
              <div class="content" v-html="topic.toc" />
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
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
      return this.currentUser && this.topic && this.currentUser.id === this.topic.user.id
    }
  },
  async asyncData({ $axios, params, error }) {
    let topic
    try {
      topic = await $axios.get('/api/topic/' + params.id)
    } catch (e) {
      error({
        statusCode: 404,
        message: '话题不存在'
      })
      return
    }

    const currentUser = await $axios.get('/api/user/current')
    const favorited = await $axios.get('/api/favorite/favorited', {
      params: {
        entityType: 'topic',
        entityId: params.id
      }
    })

    return {
      currentUser: currentUser,
      topic: topic,
      favorited: favorited.favorited
    }
  },
  mounted() {
    utils.handleToc()
  },
  head() {
    return {
      title: this.$siteTitle(this.topic.title)
    }
  },
  methods: {
    async addFavorite(topicId) {
      try {
        if (this.favorited) {
          await this.$axios.get('/api/favorite/delete', {
            params: {
              entityType: 'topic',
              entityId: topicId
            }
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
      try {
        await this.$axios.post('/api/topic/delete/' + topicId)
        this.$toast.success('删除成功', {
          duration: 2000,
          onComplete: function () {
            utils.linkTo('/topics')
          }
        })
      } catch (e) {
        console.error(e)
        this.$toast.error('删除失败：' + (e.message || e))
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  .aside-avatar {
    width: 150px;
    height: 150px;
  }
  .topic-detail {
  margin-bottom: 20px;

  .header {
    display: flex;
    border-bottom: 1px dashed #f4f4f5;
    padding-bottom: 5px;

    .left {
      margin-right: 10px;
    }

    .right {
      .topic-title {
        color: #555;
        font-size: 16px;
      }

      .topic-meta {
        span {
          font-size: 12px;
          color: #778087;

          &.act {
            a {
              font-size: 12px;
              color: #3273dc;
              margin-left: 10px;

              i {
                font-size: 12px;
                color: #000;
              }
            }
          }

          a {
            font-size: 12px;
            color: #778087;
          }

          &:not(:last-child) {
            margin-right: 3px;
          }

          &.tag {
            height: auto !important;
          }
        }
      }
    }
  }

  .content {
    padding-top: 10px;
    font-size: 15px;
    color: #000;
  }
}
</style>
