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
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      checkIn: null,
    }
  },
  mounted() {
    this.getCheckIn()
  },
  methods: {
    async getCheckIn() {
      try {
        this.checkIn = await this.$axios.get('/api/user/checkin')
      } catch (e) {
        console.log(e)
      }
    },
    async doCheckIn() {
      try {
        await this.$axios.post('/api/user/checkin')
        this.$toast.success('签到成功')
        await this.getCheckIn()
      } catch (e) {
        console.error(e)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.checkin {
}
</style>
