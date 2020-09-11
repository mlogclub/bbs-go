<template>
  <ul class="topic-list topic-wrap">
    <li
      v-for="(topic, index) in topics"
      :key="topic.topicId"
      class="topic-item"
    >
      <!-- 信息流广告 -->
      <adsbygoogle
        v-if="showAd && (index === 3 || index === 10 || index === 18)"
        ad-slot="4980294904"
        ad-format="fluid"
        ad-layout-key="-ht-19-1m-3j+mu"
      />
      <article itemscope itemtype="http://schema.org/BlogPosting">
        <div class="topic-header">
          <div v-if="showAvatar" class="topic-header-left">
            <a :href="'/user/' + topic.user.id" :title="topic.user.nickname">
              <img :src="topic.user.smallAvatar" class="avatar" />
            </a>
          </div>
          <div class="topic-header-center">
            <h1 class="topic-title" itemprop="headline">
              <a :href="'/topic/' + topic.topicId" :title="topic.title">{{
                topic.title
              }}</a>
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
                    topic.lastCommentTime | formatDate('yyyy-MM-ddTHH:mm:ss')
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
                <span v-for="tag in topic.tags" :key="tag.tagId" class="tag">
                  <a :href="'/topics/tag/' + tag.tagId">{{ tag.tagName }}</a>
                </span>
              </span>
            </div>
          </div>
          <div class="topic-header-right">
            <div class="like">
              <span
                :class="{ liked: topic.liked }"
                class="like-btn"
                @click="like(topic)"
              >
                <i class="iconfont icon-like" />
              </span>
              <span v-if="topic.likeCount" class="like-count">{{
                topic.likeCount
              }}</span>
            </div>
            <span class="count"
              >{{ topic.commentCount }}&nbsp;/&nbsp;{{ topic.viewCount }}</span
            >
          </div>
        </div>
      </article>
    </li>
  </ul>
</template>

<script>
import utils from '~/common/utils'
export default {
  props: {
    topics: {
      type: Array,
      default() {
        return []
      },
      required: false,
    },
    showAvatar: {
      type: Boolean,
      default: true,
    },
    showAd: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    async like(topic) {
      try {
        await this.$axios.post('/api/topic/like/' + topic.topicId)
        topic.liked = true
        topic.likeCount++
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
          this.$toast.error(e.message || e)
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped></style>
