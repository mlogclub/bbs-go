<template>
  <ul class="topic-list topic-wrap">
    <template v-for="(topic, index) in topics">
      <li
        v-if="showAd && (index === 2 || index === 10 || index === 18)"
        :key="'topic-' + index"
      >
        <div class="ad">
          <!-- 信息流广告 -->
          <adsbygoogle
            ad-slot="4980294904"
            ad-format="fluid"
            ad-layout-key="-ht-19-1m-3j+mu"
          />
        </div>
      </li>
      <li :key="topic.topicId">
        <div class="topic-header">
          <div class="topic-header-left">
            <div
              :style="{ backgroundImage: 'url(' + topic.user.avatar + ')' }"
              class="avatar avatar-size-45 is-rounded"
            />
          </div>
          <div class="topic-header-center">
            <a :href="'/topic/' + topic.topicId" :title="topic.title">
              <div class="topic-title">{{ topic.title }}</div>
            </a>

            <div class="topic-meta">
              <span class="meta-item">
                <a :href="'/user/' + topic.user.id">{{
                  topic.user.nickname
                }}</a>
              </span>
              <span class="meta-item">
                {{ topic.lastCommentTime | prettyDate }}
              </span>
              <span class="meta-item">
                <a :href="'/topics/node/' + topic.node.nodeId" class="node">{{
                  topic.node.name
                }}</a>
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
                @click="like(topic)"
                class="like-btn"
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
        <!--
        <div class="topic-summary">
          <a :href="'/topic/' + topic.topicId" :title="topic.title">{{
            topic.summary
          }}</a>
        </div>
        -->
      </li>
    </template>
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
      required: false
    },
    showAd: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    async like(topic) {
      try {
        await this.$axios.get('/api/topic/like/' + topic.topicId)
        topic.liked = true
        topic.likeCount++
      } catch (e) {
        if (e.errorCode === 1) {
          this.$toast.info('请登录后点赞！！！', {
            action: {
              text: '去登录',
              onClick: (e, toastObject) => {
                utils.toSignin()
              }
            }
          })
        } else {
          this.$toast.error(e.message || e)
        }
      }
    }
  }
}
</script>

<style lang="scss" scoped></style>
