<template>
  <ul class="topic-list">
    <li v-for="topic in topics" :key="topic.topicId" class="topic-item">
      <div
        class="topic-avatar"
        :href="'/user/' + topic.user.id"
        :title="topic.user.nickname"
      >
        <avatar :user="topic.user" />
      </div>
      <div class="topic-main-content">
        <div class="topic-top">
          <div class="topic-userinfo">
            <avatar class="topic-inline-avatar" :user="topic.user" size="20" />
            <a :href="'/user/' + topic.user.id">{{ topic.user.nickname }}</a>
          </div>
          <div class="topic-time">
            发布于{{ topic.createTime | prettyDate }}
          </div>
        </div>
        <div class="topic-content" :class="{ 'topic-tweet': topic.type === 1 }">
          <template v-if="topic.type === 0">
            <h1 class="topic-title">
              <a :href="'/topic/' + topic.topicId">{{ topic.title }}</a>
            </h1>
            <a :href="'/topic/' + topic.topicId" class="topic-summary">{{
              topic.summary
            }}</a>
          </template>
          <template v-if="topic.type === 1">
            <a
              v-if="topic.content"
              :href="'/topic/' + topic.topicId"
              class="topic-summary"
              >{{ topic.content }}</a
            >
            <ul
              v-if="topic.imageList && topic.imageList.length"
              class="topic-image-list"
            >
              <li v-for="(image, index) in topic.imageList" :key="index">
                <a :href="'/topic/' + topic.topicId" class="image-item">
                  <img v-lazy="image.preview" />
                </a>
              </li>
            </ul>
          </template>
        </div>
        <div class="topic-bottom">
          <div class="topic-handlers">
            <div
              class="btn"
              :class="{ liked: topic.liked }"
              @click="like(topic)"
            >
              <i class="iconfont icon-like" />{{ topic.liked ? '已赞' : '赞' }}
              <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.topicId)">
              <i class="iconfont icon-comments" />评论
              <span v-if="topic.commentCount > 0">{{
                topic.commentCount
              }}</span>
            </div>
            <div class="btn" @click="toTopicDetail(topic.topicId)">
              <i class="iconfont icon-read" />浏览
              <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
            </div>
          </div>
          <div class="topic-tags">
            <span>
              <a
                v-if="topic.node"
                :href="'/topics/node/' + topic.node.nodeId"
                :alt="topic.node.name"
                >{{ topic.node.name }}</a
              >
            </span>
          </div>
        </div>
      </div>
    </li>
  </ul>
</template>

<script>
import Avatar from '~/components/Avatar'

export default {
  components: {
    Avatar,
  },
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
        this.$message.success('点赞成功')
      } catch (e) {
        if (e.errorCode === 1) {
          this.$msgSignIn()
        } else {
          this.$message.error(e.message || e)
        }
      }
    },
    toTopicDetail(topicId) {
      this.$linkTo(`/topic/${topicId}`)
    },
  },
}
</script>

<style lang="scss" scoped></style>
