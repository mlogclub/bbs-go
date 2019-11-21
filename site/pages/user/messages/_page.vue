<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="columns">
        <div class="column is-9">
          <div class="main-body">
            <div class="widget">
              <div class="widget-header">
                <nav
                  class="breadcrumb"
                  aria-label="breadcrumbs"
                  style="margin-bottom: 0px;"
                >
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
                      <div
                        :style="{
                          backgroundImage: 'url(' + message.from.avatar + ')'
                        }"
                        class="avatar is-rounded has-border"
                      />
                    </div>
                    <div class="message-item-right">
                      <div class="message-item-meta">
                        <span v-if="message.from.id > 0" class="nickname">
                          <a
                            :href="'/user/' + message.from.id"
                            target="_blank"
                            >{{ message.from.nickname }}</a
                          >
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
                        {{ message.content }}
                        <span v-if="message.detailUrl" class="show-more">
                          <a :href="message.detailUrl" target="_blank"
                            >点击查看详情&gt;&gt;</a
                          >
                        </span>
                        <blockquote>{{ message.quoteContent }}</blockquote>
                      </div>
                    </div>
                  </li>
                </ul>
                <div
                  v-else
                  class="notification is-primary"
                  style="margin-top: 10px;"
                >
                  暂无消息
                </div>
                <pagination
                  :page="messagesPage.page"
                  url-prefix="/user/messages/"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="column is-3">
          <div class="main-aside">
            <user-center-sidebar
              :user="currentUser"
              :current-user="currentUser"
            />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import UserCenterSidebar from '~/components/UserCenterSidebar'
export default {
  middleware: 'authenticated',
  components: {
    UserCenterSidebar
  },
  async asyncData({ $axios, params }) {
    const [currentUser, messagesPage] = await Promise.all([
      $axios.get('/api/user/current'),
      $axios.get('/api/user/messages?page=' + (params.page || 1))
    ])
    return {
      currentUser,
      messagesPage
    }
  },
  data() {
    return {
      messages: [],
      cursor: 0,
      hasMore: true
    }
  },
  mounted() {
    // this.list()
  },
  methods: {
    async list() {
      const ret = await this.$axios.get('/api/user/messages', {
        params: {
          cursor: this.cursor
        }
      })
      if (ret.results && ret.results.length) {
        this.messages = this.messages.concat(ret.results)
      } else {
        this.hasMore = false
      }
      this.cursor = ret.cursor
    }
  }
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
        font-size: 14px;
        color: #4a4a4a;

        blockquote {
          margin: 0px;
        }
      }

      .show-more {
        text-align: right;

        a {
          font-size: 13px;
        }
      }
    }
  }
}
</style>
