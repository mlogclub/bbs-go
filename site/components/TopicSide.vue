<template>
  <div class="right-container">
    <post-btns :current-node-id="currentNodeId" />
    <site-notice />
    <div v-if="scoreRank && scoreRank.length" class="widget">
      <div class="widget-header">积分排行</div>
      <div class="widget-content">
        <ul class="score-rank">
          <li v-for="user in scoreRank" :key="user.id">
            <a :href="'/user/' + user.id">
              <img :src="user.avatar" class="avatar" />
            </a>
            <div class="score-rank-info">
              <a :href="'/user/' + user.id">{{ user.nickname }}</a>
              <div>
                <i class="iconfont icon-score" /><span>{{ user.score }}</span>
              </div>
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
import PostBtns from '~/components/PostBtns'
import SiteNotice from '~/components/SiteNotice'

export default {
  components: { PostBtns, SiteNotice },
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

    .avatar {
      width: 30px;
      height: 30px;
    }

    .score-rank-info {
      margin-left: 5px;
      line-height: 1.4;
      font-size: 12px;
      a:hover {
        text-decoration: underline;
      }
      i {
        font-size: 12px;
      }
      span {
        margin-left: 3px;
        font-weight: 700;
        color: orangered;
      }
    }
  }
}
</style>
