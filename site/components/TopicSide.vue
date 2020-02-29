<template>
  <div class="right-container">
    <!--<post-btns :current-node-id="currentNodeId" />-->
    <site-notice />
    <div v-if="scoreRank && scoreRank.length" class="widget">
      <div class="widget-header">积分排行</div>
      <div class="widget-content">
        <ul class="score-rank">
          <li v-for="user in scoreRank" :key="user.id">
            <a :href="'/user/' + user.id" class="score-user-avatar">
              <img v-lazy="user.avatar" class="avatar" />
            </a>
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
    <div v-if="links && links.length" class="widget">
      <div class="widget-header">
        <span>友情链接</span>
        <span class="slot"><a href="/links">查看更多&gt;&gt;</a></span>
      </div>
      <div class="widget-content">
        <ul class="links">
          <li v-for="link in links" :key="link.linkId" class="link">
            <div class="link-logo">
              <img v-if="link.logo" :src="link.logo" />
              <img v-if="!link.logo" src="~/assets/images/net.png" />
            </div>
            <div class="link-content">
              <a
                :href="'/link/' + link.linkId"
                :title="link.title"
                class="link-title"
                target="_blank"
                >{{ link.title }}</a
              >
              <p class="link-summary">
                {{ link.summary }}
              </p>
            </div>
          </li>
        </ul>
      </div>
    </div>
    <!--
    <div class="ad">
      展示广告
      <adsbygoogle ad-slot="1742173616" />
    </div>
    -->
  </div>
</template>

<script>
// import PostBtns from '~/components/PostBtns'
import SiteNotice from '~/components/SiteNotice'

export default {
  components: {
    // PostBtns,
    SiteNotice
  },
  props: {
    currentNodeId: {
      type: Number,
      default: 0
    },
    links: {
      type: Array,
      default() {
        return null
      }
    },
    scoreRank: {
      type: Array,
      default() {
        return null
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.score-rank {
  li {
    display: flex;
    list-style: none;
    margin: 8px;
    font-size: 13px;
    position: relative;

    &:not(:last-child) {
      border-bottom: 1px solid #f7f7f7;
    }

    .score-user-avatar {
      min-width: 30px;
      .avatar {
        width: 30px;
        height: 30px;
      }
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
