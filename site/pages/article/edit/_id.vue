<template>
  <section class="main">
    <div class="container">
      <div class="article-create-form">
        <h1 class="title">修改文章</h1>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input"
              type="text"
              placeholder="标题"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              placeholder="请输入内容，将图片复制或拖入编辑器可上传"
            />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <image-upload v-model="postForm.cover" :limit="1" size="120px">
              <template #add-image-button>
                <div class="cover-add-btn">
                  <i class="iconfont icon-add" />
                  <span>封面</span>
                </div>
              </template>
            </image-upload>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              :class="{ 'is-loading': publishing }"
              :disabled="publishing"
              class="button is-success"
              @click="submitCreate"
              >提交更改</a
            >
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  middleware: 'authenticated',
  async asyncData({ $axios, params, error }) {
    try {
      const [article] = await Promise.all([
        $axios.get('/api/article/edit/' + params.id),
      ])
      const cover = []
      if (article.cover) {
        cover.push(article.cover)
      }
      return {
        article,
        postForm: {
          title: article.title,
          content: article.content,
          tags: article.tags,
          cover,
        },
      }
    } catch (e) {
      error({
        statusCode: 403,
        message: e.message || '403',
      })
    }
  },
  data() {
    return {
      publishing: false, // 当前是否正处于发布中...
      postForm: {
        title: '',
        content: '',
        tags: [],
        cover: [],
      },
    }
  },
  head() {
    return {
      title: this.$siteTitle('修改文章'),
    }
  },
  computed: {
    currentUser() {
      return this.$store.state.user.current
    },
  },
  methods: {
    async submitCreate() {
      const me = this
      if (me.publishing) {
        return
      }
      me.publishing = true

      try {
        const article = await this.$axios.post(
          '/api/article/edit/' + this.article.articleId,
          {
            title: this.postForm.title,
            content: this.postForm.content,
            tags: this.postForm.tags ? this.postForm.tags.join(',') : '',
            cover:
              me.postForm.cover && me.postForm.cover.length
                ? JSON.stringify(me.postForm.cover[0])
                : null,
          }
        )
        this.$msg({
          message: '删除成功',
          onClose() {
            me.$linkTo('/article/' + article.articleId)
          },
        })
      } catch (e) {
        me.publishing = false
        this.$message.error('提交失败：' + (e.message || e))
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.article-create-form {
  background-color: var(--bg-color);
  padding: 30px;
}
.cover-add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  i {
    font-size: 24px;
    color: #1878f3;
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
  }
}
</style>
