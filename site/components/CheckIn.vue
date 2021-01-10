<template>
  <div class="widget">
    <div class="widget-header">签到</div>
    <div class="widget-content checkin">
      <div v-if="checkIn && checkIn.checkIn">
        今日已签到
        <div>
          已连续签到&nbsp;<b>{{ checkIn.consecutiveDays }}</b
          >&nbsp;天
        </div>
      </div>
      <div v-else>
        <a @click="doCheckIn">立即签到</a>
        <div style="color: #f14668;">签到可以获得积分哦!</div>
      </div>

      <div v-if="checkInRank && checkInRank.length" class="rank">
        <div class="rank-title">今日排行</div>
        <ul>
          <li v-for="rank in checkInRank" :key="rank.id">
            <avatar :user="rank.user" size="30" class="rank-user-avatar" />
            <div class="rank-user-info">
              <a :href="'/user/' + rank.user.id">{{ rank.user.nickname }}</a>
              <p>@{{ rank.updateTime | formatDate }}</p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import Avatar from '~/components/Avatar'
export default {
  components: { Avatar },
  // props: {
  //   checkInRank: {
  //     type: Array,
  //     default() {
  //       return null
  //     },
  //   },
  // },
  data() {
    return {
      checkIn: null,
      checkInRank: null,
    }
  },
  computed: {
    user() {
      return this.$store.state.user.current
    },
    isLogin() {
      return this.$store.state.user.current != null
    },
  },
  mounted() {
    this.getCheckIn()
    this.loadCheckInRank()
  },
  methods: {
    async getCheckIn() {
      try {
        this.checkIn = await this.$axios.get('/api/checkin/checkin')
      } catch (e) {
        console.log(e)
      }
    },
    async doCheckIn() {
      if (!this.isLogin) {
        this.$toSignin()
      }
      try {
        await this.$axios.post('/api/checkin/checkin')
        this.$message.success('签到成功')
        await this.getCheckIn()
        await this.loadCheckInRank()
      } catch (e) {
        console.error(e)
      }
    },
    async loadCheckInRank() {
      try {
        this.checkInRank = await this.$axios.get('/api/checkin/rank')
      } catch (e) {
        console.error(e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
@import './assets/styles/variable';
.checkin {
  .rank {
    border-top: 1px solid $border-color;
    margin-top: 10px;
    padding-top: 10px;
    .rank-title {
      font-size: 14px;
      font-weight: 600;
    }
    li {
      display: flex;
      list-style: none;
      margin: 8px 0;
      font-size: 13px;
      position: relative;

      &:not(:last-child) {
        border-bottom: 1px solid #f7f7f7;
      }

      .rank-user-avatar {
        min-width: 30px;
      }

      .rank-user-info {
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
    }
  }
}
</style>
