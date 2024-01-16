<template>
  <div class="widget">
    <div class="widget-header">积分记录</div>
    <div class="widget-content">
      <load-more-async v-slot="{ results }" url="/api/user/score_logs">
        <ul class="score-logs">
          <li
            v-for="scoreLog in results"
            :key="scoreLog.id"
            :class="{ plus: scoreLog.type === 0 }"
          >
            <span class="score-type">{{
              scoreLog.type === 0 ? "获得积分" : "减少积分"
            }}</span>
            <span class="score-score">
              <i class="iconfont icon-score" />
              <span>{{ scoreLog.score }}</span></span
            >
            <span class="score-description">{{ scoreLog.description }}</span>
            <span class="score-time"
              >@{{ useFormatDate(scoreLog.createTime) }}</span
            >
          </li>
        </ul>
      </load-more-async>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  middleware: ["auth"],
  layout: "ucenter",
});
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
