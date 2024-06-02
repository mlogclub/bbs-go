<template>
  <div class="widget">
    <div class="widget-header">
      <span>
        <i class="iconfont icon-message" />
        <span>消息</span>
      </span>
    </div>

    <div class="widget-content">
      <load-more-async v-slot="{ results }" url="/api/user/messages">
        <ul v-if="results && results.length" class="message-list">
          <li
            v-for="message in results"
            :key="message.messageId"
            class="message-item"
          >
            <div class="message-item-left">
              <my-avatar :user="message.from" :size="40" />
            </div>
            <div class="message-item-right">
              <div class="message-item-meta">
                <span v-if="message.from.id > 0" class="msg-nickname">
                  <nuxt-link :to="'/user/' + message.from.id" target="_blank">{{
                    message.from.nickname
                  }}</nuxt-link>
                </span>
                <span v-else class="msg-nickname">
                  <a href="javascript:void(0)" target="_blank">{{
                    message.from.nickname
                  }}</a>
                </span>
                <span class="msg-time">{{
                  usePrettyDate(message.createTime)
                }}</span>
                <span v-if="message.title" class="msg-title">
                  {{ message.title }}
                </span>
              </div>
              <div class="content">
                <div class="msg-attr message-quote">
                  {{ message.quoteContent }}
                </div>
                <div class="msg-attr message-content">
                  {{ message.content }}
                </div>
                <div
                  v-if="message.detailUrl"
                  class="msg-attr message-show-more"
                >
                  <a :href="message.detailUrl" target="_blank"
                    >点击查看详情&gt;&gt;</a
                  >
                </div>
              </div>
            </div>
          </li>
        </ul>
      </load-more-async>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: "ucenter",
  middleware: ["auth"],
});

useHead({
  title: useSiteTitle("消息"),
});
</script>

<style lang="scss" scoped>
.message-list {
  li.message-item {
    padding: 8px 0;
    zoom: 1;
    position: relative;
    overflow: hidden;
    display: flex;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .message-item-left {
      img {
        width: 50px;
        height: 50px;
        min-width: 50px;
        min-height: 50px;
        border-radius: 50%;
      }
    }

    .message-item-right {
      margin-left: 10px;
      width: 100%;

      .message-item-meta {
        display: flex;
        column-gap: 10px;
        align-items: center;
        span.msg-nickname {
          a {
            font-size: 16px;
            font-weight: 700;
            color: var(--text-link-color);
          }
        }

        span.msg-time {
          font-size: 13px;
          //color: #999;
        }

        span.msg-title {
          font-size: 14px;
          font-weight: 500;
        }
      }

      .content {
        font-size: 14px;
        margin-top: 5px;
        margin-bottom: 0;

        .msg-attr {
          margin: 10px 0 0;
        }

        .message-content {
          font-size: 15px;
          color: var(--text-color);
        }

        .message-quote {
          position: relative;
          padding: 8px 15px;
          border-radius: 1px;
          background: var(--bg-color4);
          border: 1px solid var(--border-color2);
          color: var(--text-color);

          &:before {
            box-sizing: inherit;
            position: absolute;
            content: "";
            width: 0;
            height: 0;
            top: -7px;
            border-width: 0 7px 7px 7px;
            border-style: solid;
            border-color: transparent transparent var(--border-color)
              transparent;
            z-index: 1;
          }
        }

        .message-show-more {
          a {
            &:hover {
              color: var(--text-link-color);
              text-decoration: underline;
            }
          }
        }
      }
    }
  }
}
</style>
