<template>
  <section class="main">
    <div v-if="loading" class="loading modal is-active">
      <div class="modal-background" />
      <div class="modal-content" style="text-align: center; margin-top: -100px">
        <div class="loading-animation" />
        <span class="loading-text">提交中，请稍后...</span>
      </div>
    </div>
    <div class="container">
      <el-card class="box-card">
        <ul class="msg-box">
          <li>
            <h4>积分购买（1:500）</h4>
          </li>
          <li>
            <h4 style="margin-bottom: 15px">购买金额</h4>
            <el-radio-group v-model="amountVal" @change="amountChange">
              <el-radio border :label="'' + 10">10￥</el-radio>
              <el-radio border :label="'' + 20">20￥</el-radio>
              <el-radio border :label="'' + 50">50￥</el-radio>
              <el-radio border :label="'' + 100">100￥</el-radio>
            </el-radio-group>
          </li>

          <li>
            <el-radio
              v-model="payType"
              style="margin-left: 90%"
              border
              label="1"
              >支付宝</el-radio
            >
          </li>
        </ul>
        <div style="text-align: center; margin-top: 30px">
          <el-button type="primary" @click="surePay">确认支付</el-button>
          <!-- <nuxt-link class="button" to="/pay/vip">网站赞助</nuxt-link> -->
        </div>
      </el-card>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      amountVal: '',
      loading: false,
      disabled: false,
      payType: '1',
      code: '',
    }
  },
  methods: {
    amountChange(val) {
      this.amountVal = val
      if (val === '') {
        this.disabled = false
      } else {
        this.disabled = true
      }
    },
    surePay() {
      if (this.amountVal === '') {
        this.$message.warning('请输入金额')
      }
      this.loading = true
      this.$axios
        .get('/api/pay/url', {
          params: {
            price: this.amountVal,
          },
        })
        .then((data) => {
          window.open(data, '_self')
        })
        .catch((rsp) => {
          this.$message.error(rsp.message)
        })
    },
  },
}
</script>

<style scoped>
/* 信息列表样式 */
.msg-box > li {
  list-style: none;
  border-bottom: 1px solid #c5c5c5;
  padding: 20px 10px;
}
</style>
