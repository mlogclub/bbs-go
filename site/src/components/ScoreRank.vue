<template>
  <div v-if="scoreRank && scoreRank.length" class="widget">
    <div class="widget-header">
      <span class="widget-title">积分排行</span>
    </div>
    <div class="widget-content">
      <ul class="score-rank">
        <li v-for="user in scoreRank" :key="user.id">
          <my-avatar :user="user" :size="35" />
          <div class="score-user-info">
            <nuxt-link :to="'/user/' + user.id" class="score-nickname">{{
              user.nickname
            }}</nuxt-link>
            <p class="score-desc">
              {{ user.topicCount }} 帖子 • {{ user.commentCount }} 评论
            </p>
          </div>
          <div class="score-rank-info">
            <span class="score-user-score">
              <i class="iconfont icon-score" /><span>{{ user.score }}</span>
            </span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const { data: scoreRank } = useAsyncData(() =>
  useMyFetch("/api/user/score/rank")
);
</script>

<style scoped lang="scss">
.score-rank {
  li {
    display: flex;
    justify-content: flex-start;
    align-items: center;
    list-style: none;
    font-size: 13px;
    position: relative;
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .score-user-info {
      width: 100%;
      margin-left: 9px;
      line-height: 1.4;
      font-size: 12px;
      .score-nickname {
        font-size: 14px;
        color: var(--text-color);
        line-height: 20px;

        &:hover {
          color: rgba(0, 166, 244, 0.8);
        }
      }
      .score-desc {
        font-size: 11px;
        color: var(--text-color3);
        line-height: 20px;
        display: block;
      }
    }

    .score-rank-info {
      width: 120px;
      .score-user-score {
        float: right;
        border-radius: 12px;
        color: var(--text-color3);
        height: 21px;
        line-height: 21px;
        padding: 0 6px;
        text-shadow: 0 0 1px #fff;
        background-color: var(--bg-color2);
        font-size: 0.75rem;
        align-items: center;
        i {
          margin-right: 3px;
        }
      }
    }
  }
}
</style>
