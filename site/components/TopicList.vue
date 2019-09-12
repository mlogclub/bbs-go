<template>
  <ul class="topic-list">
    <template v-for="(topic, index) in topics">
      <li v-if="showAd && (index === 3)" :key="'ad-' + index ">
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
            (adsbygoogle = window.adsbygoogle || []).push({});
          </script>
        </div>
      </li>
      <li :key="topic.topicId">
        <div class="topic-item">
          <div class="left">
            <div
              class="avatar avatar-size-45 is-rounded"
              :style="{backgroundImage:'url('+ topic.user.avatar +')'}"
            />
          </div>
          <div class="center">
            <a :href="'/topic/' + topic.topicId" :title="topic.title">
              <div class="topic-title">{{ topic.title }}</div>
            </a>

            <div class="topic-meta">
              <div class="meta-item">
                <a :href="'/user/' + topic.user.id">{{ topic.user.nickname }}</a>
              </div>
              <div class="meta-item">
                {{ topic.lastCommentTime | prettyDate }}
              </div>
              <div class="meta-item">
                <span v-for="tag in topic.tags" :key="tag.tagId" class="tag">
                  <a :href="'/topics/tag/' + tag.tagId + '/1'">{{ tag.tagName }}</a>
                </span>
              </div>
            </div>
          </div>
          <div class="right">
            <div class="like">
              <span class="like-btn" :class="{'liked': topic.liked}" @click="like(topic)">
                <i class="iconfont icon-like" />
              </span>
              <span v-if="topic.likeCount" class="like-count">{{ topic.likeCount }}</span>
            </div>
            <span class="count">{{ topic.commentCount }}&nbsp;/&nbsp;{{ topic.viewCount }}</span>
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
      default: function () {
        return null
      },
      required: true
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

    .topic-item {
      display: flex;

      .left {
        min-width: 48px;
        min-height: 48px;

        .avatar {
          min-width: 48px;
          min-height: 48px;
          width: 48px;
          height: 48px;
        }
      }

      .center {
        width: 100%;
        margin-left: 8px;

        .topic-title {
          color: #555;
          font-size: 16px;
          line-height: 21px;
          font-weight: normal;
          overflow: hidden;
          word-break: break-all;
          -webkit-line-clamp: 2;
          text-overflow: ellipsis;
          -webkit-box-orient: vertical;
          display: -webkit-box;
        }

        .topic-meta {
          position: relative;
          font-size: 12px;
          line-height: 24px;
          color: #bbb;
          margin-top: 3px;
          display: flex;

          .meta-item {
            font-size: 12px;

            &:not(:last-child) {
              margin-right: 8px;
            }

            .tag {
              height: auto !important;

              &:not(:last-child) {
                margin-right: 3px;
              }
            }
          }

          a {
            color: #778087;
          }
        }
      }

      .right {
        min-width: 100px;
        max-width: 100px;
        text-align: right;
        padding-right: 8px;

        .like {
          font-size: 12px;

          .like-btn {
            transition: all 0.5s;
            cursor: pointer;

            background-color: rgba(126, 107, 1, 0.08);
            color: #e7672e;
            // padding: 0;
            // display: inline-block;
            // height: 20px;
            // width: 20px;
            // line-height: 20px;
            // border-radius: 10px;
            // text-align: center;
            padding: 4px;
            border-radius: 50%;

            &:hover, &.liked  {
              color: #fff;
              background-color: #e7672e;
            }
          }

          .like-count {
            font-weight: bold;
            color: #e7672e;

            &::before {
              content: " x ";
            }
          }
        }

        .count {
          font-size: 12px;
          color: #fff;

          background: #aab0c6;
          padding: 2px 10px;
          border-radius: 6px;
          font-weight: 700;
        }
      }
    }
  }
}
</style>
