<template>
  <section class="main">
    <div class="container-wrapper">
      <div class="main-body">
        <div class="notice">
          <h1>收录规则：</h1>
          <ul>
            <li>只收录编程相关的技术博客</li>
            <li>拒绝SEO类型博客</li>
            <li>拒绝大量转载文章的博客</li>
            <li>独立博客优先</li>
            <li>提交后我们会进行审核，审核通过后即可显示</li>
          </ul>
        </div>

        <div class="widget">
          <div class="widget-header">
            提交博客地址
          </div>
          <div class="widget-content">
            <div class="field">
              <label class="label">博客链接</label>
              <div class="control">
                <input
                  v-model="url"
                  @blur="detect"
                  @keyup.enter="submitLink"
                  class="input is-success"
                  type="text"
                  placeholder="博客链接（必填）"
                />
              </div>
            </div>

            <div class="field">
              <label class="label">博客标题</label>
              <div class="control">
                <input
                  v-model="title"
                  @keyup.enter="submitLink"
                  class="input is-success"
                  type="text"
                  placeholder="博客标题（必填）"
                />
              </div>
            </div>

            <div class="field">
              <label class="label">博客简介</label>
              <div class="control">
                <input
                  v-model="summary"
                  @keyup.enter="submitLink"
                  class="input is-success"
                  type="text"
                  placeholder="博客简介（必填）"
                />
              </div>
            </div>

            <div class="field">
              <label class="label">博客Logo</label>
              <div class="control">
                <input
                  v-model="logo"
                  @keyup.enter="submitLink"
                  class="input is-success"
                  type="text"
                  placeholder="博客Logo（非必填请填写Logo链接）"
                />
              </div>
            </div>

            <div class="field">
              <div class="control">
                <button
                  :class="{ 'is-loading': publishing }"
                  :disabled="publishing"
                  @click="submitLink"
                  class="button is-success"
                >
                  提交
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
import utils from '~/common/utils'
export default {
  data() {
    return {
      publishing: false,
      url: '',
      title: '',
      summary: '',
      logo: ''
    }
  },
  methods: {
    async submitLink() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true

      try {
        await this.$axios.post('/api/link/create', {
          url: this.url,
          title: this.title,
          summary: this.summary,
          logo: this.logo
        })
        this.$toast.success('提交成功，请耐心等待审核。', {
          duration: 1000,
          onComplete() {
            utils.linkTo('/links')
          }
        })
      } catch (e) {
        me.publishing = false
        this.$toast.error(e.message || e)
      }
    },
    async detect() {
      if (!this.url) {
        return
      }
      try {
        const data = await this.$axios.post('/api/link/detect', {
          url: this.url
        })
        if (data) {
          if (!this.title && data.title) {
            this.title = data.title
          }
          if (!this.summary && data.description) {
            this.summary = data.description
          }
        }
      } catch (e) {
        console.error(e)
      }
    }
  },
  head() {
    return {
      title: this.$siteTitle('提交博客')
    }
  }
}
</script>

<style lang="scss" scoped>
.notice {
  padding: 7px 15px;
  margin-bottom: 20px;
  border: 1px solid transparent;
  border-radius: 4px;
  background-color: #fcf8e3;
  border-color: #faebcc;
  color: #8a6d3b;

  a {
    color: #3273dc;
    cursor: pointer;
  }

  h1 {
    font-weight: bold;
  }

  p:not(:last-child) {
    margin-bottom: 10px;
  }

  ul {
    list-style: disc;
    margin-left: 20px;
    margin-top: 10px;
  }
}
</style>
