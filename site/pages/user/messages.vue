<template>
  <section class="main">
    <div class="container main-container left-main">
      <div class="left-container">
        <user-profile :user="currentUser" />
        <div class="widget">
          <div class="widget-header">
            <nav class="breadcrumb">
              <ul>
                <li>
                  <a href="/">首页</a>
                </li>
                <li>
                  <a :href="'/user/' + currentUser.id">{{
                    currentUser.nickname
                  }}</a>
                </li>
                <li class="is-active">
                  <a href="#" aria-current="page">消息</a>
                </li>
              </ul>
            </nav>
          </div>

          <div class="widget-content">
            <ul
              v-if="messagesPage && messagesPage.results"
              class="message-list"
            >
              <li
                v-for="message in messagesPage.results"
                :key="message.messageId"
                class="message-item"
              >
                <div class="message-item-left">
                  <img :src="message.from.smallAvatar" class="avatar" />
                </div>
                <div class="message-item-right">
                  <div class="message-item-meta">
                    <span v-if="message.from.id > 0" class="nickname">
                      <a :href="'/user/' + message.from.id" target="_blank">{{
                        message.from.nickname
                      }}</a>
                    </span>
                    <span v-else class="nickname">
                      <a href="javascript:void(0)" target="_blank">{{
                        message.from.nickname
                      }}</a>
                    </span>
                    <span class="time">{{
                      message.createTime | prettyDate
                    }}</span>
                  </div>
                  <div class="content">
                    <div class="message-content">
                      {{ message.content }}
                      <a
                        v-if="message.detailUrl"
                        :href="message.detailUrl"
                        class="show-more"
                        target="_blank"
                        >点击查看详情&gt;&gt;</a
                      >
                    </div>
                    <blockquote class="message-quote">
                      {{ message.quoteContent }}
                    </blockquote>
                  </div>
                </div>
              </li>
            </ul>
            <div v-else class="notification is-primary">
              暂无消息
            </div>
            <pagination
              :page="messagesPage.page"
              url-prefix="/user/messages?p="
            />
          </div>
        </div>
      </div>
      <user-center-sidebar :user="currentUser" />
    </div>
  </section>
</template>

<script>
import UserProfile from '~/components/UserProfile'
import UserCenterSidebar from '~/components/UserCenterSidebar'
import Pagination from '~/components/Pagination'
export default {
  middleware: 'authenticated',
  components: { UserProfile, UserCenterSidebar, Pagination },
  async asyncData({ $axios, query }) {
    const [messagesPage] = await Promise.all([
      $axios.get('/api/user/messages?page=' + (query.p || 1)),
    ])
    return {
      messagesPage,
    }
  },
  data() {
    return {
      messages: [],
      cursor: 0,
      hasMore: true,
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
  },
}
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
      border-bottom: 1px solid #f2f2f2;
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
        span.nickname {
          font-size: 16px;
          font-weight: bold;
        }

        span.time {
          font-size: 13px;
          color: #999;
        }
      }

      .content {
        margin-top: 5px;
        margin-bottom: 0px;

        .message-content {
          font-size: 15px;
          font-weight: 400;
          color: #000;

          .show-more {
            text-align: right;
            margin-left: 5px;
          }
        }

        .message-quote {
          margin: 8px 0;
          padding: 8px;
          font-size: 13px;
          font-weight: 400;
          color: #4a4a4a;
        }
      }
    }
  }
}
</style>
