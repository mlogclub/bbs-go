<template>
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
          <span class="score-score">
            <i class="iconfont icon-score" />
            <span>{{ scoreLog.score }}</span></span
          >
          <span class="score-description">{{ scoreLog.description }}</span>
          <span class="score-time"
            >@{{ scoreLog.createTime | formatDate }}</span
          >
        </li>
      </ul>
      <pagination :page="scoreLogsPage.page" url-prefix="/user/scores?p=" />
    </div>
  </div>
</template>

<script>
export default {
  layout: 'ucenter',
  middleware: 'authenticated',
  async asyncData({ $axios, query }) {
    const [scoreLogsPage] = await Promise.all([
      $axios.get('/api/user/scorelogs?page=' + (query.p || 1)),
    ])
    return {
      scoreLogsPage,
    }
  },
}
</script>

<style lang="scss" scoped>
.score-logs {
  // margin-top: 10px;
  font-size: 1rem;
  li {
    line-height: 2;
    margin-bottom: 10px;

    .score-type {
      color: green;
    }

    .score-score {
      margin: 0 3px;
      span {
        font-weight: bold;
      }
    }

    .score-time {
      color: var(--text-color3);
    }

    .score-description {
      color: var(--text-color3);
    }

    &.plus {
      .score-type {
        color: red;
      }
    }
  }
}
</style>
