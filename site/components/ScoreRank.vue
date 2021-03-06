<template>
  <div v-if="scoreRank && scoreRank.length" class="widget">
    <div class="widget-header">
      <span class="widget-title">积分排行</span>
    </div>
    <div class="widget-content">
      <ul class="score-rank">
        <li v-for="user in scoreRank" :key="user.id">
          <avatar :user="user" size="30" />
          <div class="score-user-info">
            <a :href="'/user/' + user.id">{{ user.nickname }}</a>
            <p>{{ user.topicCount }} 帖子 • {{ user.commentCount }} 评论</p>
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

<script>
import Avatar from '~/components/Avatar'
export default {
  name: 'ScoreRank',
  components: { Avatar },
  props: {
    scoreRank: {
      type: Array,
      default() {
        return null
      },
    },
  },
}
</script>

<style scoped lang="scss">
.score-rank {
  li {
    display: flex;
    list-style: none;
    margin: 8px 0;
    font-size: 13px;
    position: relative;

    &:not(:last-child) {
      border-bottom: 1px solid #f7f7f7;
    }

    .score-user-info {
      width: 100%;
      margin-left: 5px;
      line-height: 1.4;
      font-size: 12px;
      a {
        font-weight: 700;
        &:hover {
          text-decoration: underline;
        }
      }
    }

    .score-rank-info {
      width: 120px;
      .score-user-score {
        float: right;
        border-radius: 12px;
        color: #778087;
        height: 21px;
        line-height: 21px;
        padding: 0 6px;
        text-shadow: 0 0 1px #fff;
        background-color: #f5f5f5;
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
