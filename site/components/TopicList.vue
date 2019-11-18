<template>
  <ul class="topic-list topic-wrap">
    <template v-for="(topic, index) in topics">
      <li v-if="showAd && index === 3" :key="'ad-' + index">
        <div class="ad">
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
                <span v-for="tag in topic.tags" :key="tag.tagId" class="tag">
                  <a :href="'/topics/tag/' + tag.tagId + '/1'">{{
                    tag.tagName
                  }}</a>
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

<style lang="scss" scoped>
.topic-list {
  margin: 0 0 10px 0 !important;

  li {
    padding: 8px 0 8px 8px;
    position: relative;
    overflow: hidden;
    border-radius: 4px;
    transition: background 0.5s;

    &:hover {
      background: #f3f6f9;
      border-bottom: none;
    }

    // &:not(:last-child) {
    //   border-bottom: 1px dashed #f2f2f2;
    // }
  }
}
</style>
