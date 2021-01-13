<template>
  <section class="main">
    <div class="container main-container left-main size-320">
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
                  <avatar :user="message.from" size="40" :round="true" />
                </div>
                <div class="message-item-right">
                  <div class="message-item-meta">
                    <span v-if="message.from.id > 0" class="msg-nickname">
                      <a :href="'/user/' + message.from.id" target="_blank">{{
                        message.from.nickname
                      }}</a>
                    </span>
                    <span v-else class="msg-nickname">
                      <a href="javascript:void(0)" target="_blank">{{
                        message.from.nickname
                      }}</a>
                    </span>
                    <span class="msg-time">{{
                      message.createTime | prettyDate
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
import Avatar from '~/components/Avatar'

export default {
  middleware: 'authenticated',
  components: { UserProfile, UserCenterSidebar, Pagination, Avatar },
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
        span.msg-nickname {
          a {
            font-size: 16px;
            font-weight: 700;
            color: #1e70bf;
          }
        }

        span.msg-time {
          font-size: 13px;
          //color: #999;
        }

        span.msg-title {
          font-size: 16px;
          font-weight: 700;
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
          color: #000;
        }

        .message-quote {
          position: relative;
          padding: 8px 15px;
          border-radius: 1px;
          background: #f2f2f2;
          border: 1px solid #eaeaea;
          color: #292929;

          &:before {
            box-sizing: inherit;
            position: absolute;
            content: '';
            width: 0;
            height: 0;
            top: -7px;
            border-width: 0 7px 7px 7px;
            border-style: solid;
            border-color: transparent transparent #f2f2f2 transparent;
            z-index: 1;
          }
        }

        .message-show-more {
          a {
            &:hover {
              color: #3273dc;
              text-decoration: underline;
            }
          }
        }
      }
    }
  }
}
</style>
