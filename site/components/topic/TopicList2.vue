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
          <div class="topic-userinfo">
            <avatar class="topic-inline-avatar" :user="topic.user" size="20" />
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
          <template v-if="topic.type === 1">
            <div v-if="topic.content" class="topic-summary">
              {{ topic.content }}
            </div>
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
            <i class="iconfont icon-read" />浏览
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
    padding: 12px 12px;
    display: flex;
    position: relative;
    overflow: hidden;
    transition: background 0.5s;
    border-bottom: 1px solid #f2f2f2;

    //&:hover {
    //  background: #f3f6f9;
    //}

    .topic-main-content {
      flex: 1;
      margin-left: 12px;

      .topic-top {
        margin-bottom: 8px;

        .topic-userinfo {
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

          .topic-inline-avatar {
            display: none;
            margin-right: 5px;
          }
        }

        .topic-time {
          color: #8590a6;
          font-size: 12px;
          float: right;
        }

        @media screen and (max-width: 1024px) {
          .topic-time {
            float: none;
            margin-top: 8px;
          }
        }
      }

      .topic-content {
        .topic-title {
          display: inline-block;
          font-size: 15px;
          margin-bottom: 6px;
          word-wrap: break-word;
          word-break: break-all;
          width: 100%;

          a {
            color: #3273dc;

            &:hover {
              color: #3273dc;
              text-decoration: underline;
            }
          }
        }

        .topic-summary {
          display: inline-block;
          font-size: 15px;
          margin-bottom: 6px;
          word-wrap: break-word;
          word-break: break-all;
          width: 100%;
          color: #17181a;
        }

        .topic-image-list {
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
      }

      .topic-handlers {
        display: flex;
        align-items: center;
        //justify-content: space-between;
        margin-top: 6px;
        font-size: 12px;
        flex: 1;

        .btn {
          color: #8590a6;
          cursor: pointer;

          &:not(:last-child) {
            margin-right: 20px;
          }

          &:hover {
            color: #1878f3;
          }

          i {
            margin-right: 3px;
            font-size: 12px;
            position: relative;
          }
        }
      }
    }

    @media screen and (max-width: 768px) {
      .topic-avatar {
        display: none;
      }

      .topic-main-content {
        margin-left: 0;

        .topic-inline-avatar {
          display: block !important;
        }
      }
    }
  }
}
</style>
