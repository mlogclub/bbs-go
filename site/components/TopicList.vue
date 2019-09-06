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
            <div class="avatar avatar-size-45 is-rounded" :style="{backgroundImage:'url('+ topic.user.avatar +')'}" />
          </div>
          <div class="center">
            <a :href="'/topic/' + topic.topicId" :title="topic.title">
              <div class="topic-title">{{ topic.title }}</div>
            </a>

            <div class="topic-meta">
              <span><a :href="'/user/' + topic.user.id">{{ topic.user.nickname }}</a></span>
              <span>{{ topic.lastCommentTime | prettyDate }}</span>
              <span v-for="tag in topic.tags" :key="tag.tagId" class="tag">
                <a :href="'/topics/tag/' + tag.tagId + '/1'">{{ tag.tagName }}</a></span>
            </div>
          </div>
          <div class="right">
            <span class="view-count">{{ topic.viewCount }}</span>
          </div>
        </div>
      </li>
    </template>
  </ul>
</template>

<script>
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
    transition: background .2s;

    &:hover {
      background: #f3f6f9;
    }

    // &:not(:last-child) {
    //   border-bottom: 1px dashed #f2f2f2;
    // }

    .topic-item {
      display: flex;

      .left {
        min-width: 45px;
        min-height: 45px;
      }

      .center {
        width: 100%;
        margin-left: 5px;

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
          color: #bbb;
          margin-top: 6px;

          span {
            font-size: 12px;

            &:not(:last-child) {
              margin-right: 3px;
            }

            &.tag {
              height: auto !important;
            }

            &.btn a {
              color: #3273dc;
            }
          }

          a {
            color: #778087;
          }
        }
      }

      .right {
        min-width: 65px;
        max-width: 65px;
        text-align: right;
        padding-right: 8px;

        span.view-count {
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
