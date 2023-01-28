<template>
  <section class="main">
    <div class="container">
      <el-card class="box-card">
        <ul class="msg-box">
          <li>
            <h4>网站赞助</h4>
          </li>
          <li>
            <h4 style="margin-bottom: 15px">赞助金额</h4>
            <el-radio-group v-model="amountVal" @change="amountChange">
              <el-radio border :label="'' + 99">VIP1</el-radio>
              <el-radio border :label="'' + 199">VIP2</el-radio>
              <el-radio border :label="'' + 299">VIP3</el-radio>
              <el-radio border :label="'' + 599">VIP4</el-radio>
              <el-radio border :label="'' + 799">VIP5</el-radio>
              <el-radio border :label="'' + 1888">VIP6</el-radio>
              <input
                v-model="rechargeParams.totalAmt"
                :disabled="true"
                class="input"
                style="width: 150px"
              />
            </el-radio-group>
          </li>

          <li>
            <el-radio-group
              v-model="rechargeParams.paymentType"
              style="margin-left: 90%"
            >
              <el-radio border :label="'' + 1">支付宝</el-radio>
            </el-radio-group>
          </li>
        </ul>
        <div style="text-align: center; margin-top: 30px">
          <el-button type="primary" @click="surePay">确认支付</el-button>
          <nuxt-link class="button" to="/pay">积分购买</nuxt-link>
        </div>
      </el-card>
    </div>
  </section>
</template>

<script>
export default {
  data() {
    return {
      amountVal: '99',
      rechargeParams: {
        totalAmt: '99',
        paymentType: '1',
      },
      code: '',
    }
  },

  methods: {
    amountChange(val) {
      this.rechargeParams.totalAmt = val
      if (val === '') {
        this.disabled = false
      } else {
        this.disabled = true
      }
    },
    surePay() {
      if (this.rechargeParams.totalAmt === '') {
        this.$message.warning('请输入金额')
      }
      // this.$router.push({path: '/code'});

      //   post('weixin/createOrderInfo', this.rechargeParams).then((res) => {
      //     var result = res.result
      //     if (res.code === 20000) {
      //       // 支付方式跳转

      //       if (this.rechargeParams.paymentType == '0') {
      //         this.$message.success('微信支付')
      //         this.wechatPay(result)
      //       } else if (this.rechargeParams.paymentType == '1') {
      //         this.$message.success('支付宝支付')
      //         const payDiv = document.getElementById('payDiv')
      //         if (payDiv) {
      //           document.body.removeChild(payDiv)
      //         }
      //         const div = document.createElement('div')
      //         div.id = 'payDiv'
      //         div.innerHTML = result
      //         document.body.appendChild(div)
      //         document
      //           .getElementById('payDiv')
      //           .getElementsByTagName('form')[0]
      //           .submit()
      //       }
      //     }
      //   })
    },
  },
  // },
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
