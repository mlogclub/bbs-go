<template>
  <ul class="topic-list2">
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
          <div class="topic-user-info">
            <a :href="'/user/' + topic.user.id">{{ topic.user.nickname }}</a>
          </div>
          <div class="topic-time">
            发布于{{ topic.createTime | prettyDate }}
          </div>
        </div>
        <div class="topic-content">
          <template v-if="topic.type === 0">
            <div class="topic-title">
              发布了：<a :href="'/topic/' + topic.topicId">{{ topic.title }}</a>
            </div>
          </template>
        </div>
        <div class="topic-handlers">
          <div class="btn" @click="like(topic)">
            <i class="iconfont icon-like" />{{ topic.liked ? '已赞' : '赞' }}
            <span v-if="topic.likeCount > 0">{{ topic.likeCount }}</span>
          </div>
          <div class="btn" @click="toTopicDetail(topic.topicId)">
            <i class="iconfont icon-comments" />评论
            <span v-if="topic.commentCount > 0">{{ topic.commentCount }}</span>
          </div>
          <div class="btn" @click="toTopicDetail(topic.topicId)">
            <i class="iconfont icon-view" />查看
            <span v-if="topic.viewCount > 0">{{ topic.viewCount }}</span>
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

<style lang="scss" scoped>
.topic-list2 {
  .topic-item {
    padding: 12px 12px 12px 8px;
    display: flex;
    //padding: 20px 20px 10px 20px;
    position: relative;
    overflow: hidden;
    transition: background 0.5s;
    border-bottom: 1px solid #f2f2f2;
    //cursor: pointer;
    //&:hover {
    //  background: #f3f6f9;
    //}

    .topic-main-content {
      flex: 1;
      margin-left: 15px;

      .topic-top {
        margin-bottom: 8px;

        .topic-user-info {
          display: inline-flex;
          align-items: center;

          a {
            font-weight: 700;
            font-size: 16px;
            color: rgb(51, 51, 51);
            display: flex;
            max-width: 250px;
            overflow: hidden;
          }
        }

        @media screen and (max-width: 768px) {
          .topic-time {
            float: none !important;
            margin-top: 8px;
          }
        }

        .topic-time {
          color: #8590a6;
          font-size: 12px;
          float: right;
        }
      }

      .topic-content {
        .topic-title {
          display: inline-block;
          font-size: 15px;
          margin-bottom: 6px;
          word-wrap: break-word;
          word-break: normal;
          width: 100%;

          a {
            color: #3273dc;

            &:hover {
              color: #3273dc;
              text-decoration: underline;
            }
          }
        }
      }

      .topic-handlers {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-top: 6px;
        font-size: 12px;
        flex: 1;

        .btn {
          color: #8590a6;
          cursor: pointer;
          min-width: 100px;

          &:hover {
            color: #1878f3;
          }

          i {
            margin-right: 3px;
            font-size: 20px;
            position: relative;
            top: 2px;
          }
        }
      }
    }
  }
}
</style>
