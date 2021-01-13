<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <user-profile :user="currentUser" />
        <div class="widget">
          <div class="widget-header">积分记录</div>
          <div class="widget-content">
            <ul class="score-logs">
              <li
                v-for="scoreLog in scoreLogsPage.results"
                :key="scoreLog.id"
                :class="{ plus: scoreLog.type === 0 }"
              >
                <span class="score-type">{{
                  scoreLog.type === 0 ? '获得积分' : '减少积分'
                }}</span>
                <span class="score-score"
                  ><i class="iconfont icon-score" /><span>{{
                    scoreLog.score
                  }}</span></span
                >
                <span class="score-description">{{
                  scoreLog.description
                }}</span>
                <span class="score-time"
                  >@{{ scoreLog.createTime | formatDate }}</span
                >
              </li>
            </ul>
            <pagination
              :page="scoreLogsPage.page"
              url-prefix="/user/scores?p="
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
    const [scoreLogsPage] = await Promise.all([
      $axios.get('/api/user/scorelogs?page=' + (query.p || 1)),
    ])
    return {
      scoreLogsPage,
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
.score-logs {
  // margin-top: 10px;
  font-size: 1rem;
  li {
    line-height: 2;

    .score-type {
      color: green;
    }

    .score-score {
      span {
        font-weight: bold;
      }
    }

    .score-time {
      color: #777;
    }

    .score-description {
      color: #777;
    }

    &.plus {
      .score-type {
        color: red;
      }
    }
  }
}
</style>
