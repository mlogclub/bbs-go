<template>
  <div class="post-tweets-wrapper">
    <ul class="tab-list">
      <li class="tab-item current">
        <div class="tab-name">发表动态</div>
      </li>
    </ul>
    <div class="tweets-box">
      <textarea
        v-model="content"
        placeholder="有什么新鲜事想告诉大家"
        class="title-input"
        @input="onInput"
        @paste="handleParse"
        @drop="handleDrag"
        @keydown.ctrl.enter="doSubmit"
        @keydown.meta.enter="doSubmit"
      />
      <p class="words-number">{{ wordCount }}/{{ maxWordCount }}字</p>
      <div class="box-footer">
        <div class="bui-left">
          <span class="action-btn" @click="showUploader = !showUploader">
            <i class="iconfont icon-image" />
            <span>图片</span>
          </span>
          <!--
          <span class="action-btn">
            <i class="iconfont icon-emoji" />
            <span>表情</span>
          </span>
          -->
        </div>
        <div class="bui-right">
          <span class="msg-tip">{{ message }}</span>
          <span class="tweets-help">Ctrl or ⌘ + Enter</span>
          <a
            :class="{ active: hasContent }"
            class="upload-publish"
            @click="doSubmit"
            >发布</a
          >
        </div>
      </div>

      <div v-show="showUploader" class="uploader-popup">
        <div class="imgUploadBox">
          <p class="uploader-title">本地上传</p>
          <p class="uploader-meta">
            共 {{ imageCount }} 张，还能上传 {{ maxImageCount }} 张
          </p>
          <i
            class="close-popup iconfont icon-close"
            @click="showUploader = false"
          />
          <div class="upload-box">
            <form style="display: none;">
              <input
                ref="imageInput"
                type="file"
                accept="image/*"
                multiple="multiple"
                @change="handleImageUploadChange"
              />
            </form>
            <ul class="upload-img-list">
              <li v-for="(image, i) in images" :key="i" class="upload-img-item">
                <img :src="image" />
                <i
                  class="iconfont icon-close remove"
                  @click="removeImg(image)"
                />
              </li>
              <li
                v-if="imageCount < maxImageCount"
                class="upload-img-item upload-img-add"
                @click="handleImageUploadClick"
              >
                <i class="iconfont icon-add" />
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    nodeId: {
      type: Number,
      default: 0,
    },
  },
  data() {
    return {
      content: '',
      images: [
        // 'https://file.mlog.club/images/2020/02/27/0aadf3d7c46dba756f4e228e8e8f8ed6.jpg',
        // 'https://file.mlog.club/images/2020/02/28/6819d3e0afb535594fb55c1108e1ad37.jpg'
      ],
      message: '',
      maxWordCount: 1000,
      showUploader: false,
      maxImageCount: 9,
    }
  },
  computed: {
    hasContent() {
      return this.content && this.content.length > 0
    },
    wordCount() {
      return this.content ? this.content.length : 0
    },
    imageCount() {
      return this.images ? this.images.length : 0
    },
    user() {
      return this.$store.state.user.current
    },
  },
  methods: {
    onInput() {},
    async doSubmit() {
      if (!this.user) {
        this.message = '发表失败，请登录后重试'
        return
      }
      if (!this.hasContent) {
        this.message = '发表失败，请输入内容'
        return
      }
      this.showUploader = false // 关闭图片上传框
      try {
        const ret = await this.$axios.post('/api/tweet/create', {
          content: this.content,
          imageList: JSON.stringify(this.images),
        })
        this.content = ''
        this.message = ''
        this.$emit('created', ret)
        this.$toast.success('发布成功')
      } catch (e) {
        this.message = e.message || e
      }
    },
    handleImageUploadClick() {
      this.$refs.imageInput.click()
    },
    handleParse(e) {
      const items = e.clipboardData && e.clipboardData.items
      let file = null
      if (items && items.length) {
        for (let i = 0; i < items.length; i++) {
          if (items[i].type.includes('image')) {
            file = items[i].getAsFile()
          }
        }
      }

      if (!file) {
        return
      }

      this.showUploader = true // 展开上传面板
      e.preventDefault() // 阻止默认行为即不让剪贴板内容在div中显示出来

      if (this.imageCount + 1 > this.maxImageCount) {
        this.message = '图片数量超过上限'
        return
      }

      this.upload(file) // 上传
    },
    handleDrag(e) {
      e.stopPropagation()
      e.preventDefault()

      const files = []
      const items = e.dataTransfer.items
      if (items && items.length) {
        if (items && items.length) {
          for (let i = 0; i < items.length; i++) {
            if (items[i].type.includes('image')) {
              files.push(items[i].getAsFile())
            }
          }
        }
      }

      if (files && files.length) {
        this.showUploader = true // 展开上传面板
        this.uploadFiles(files)
      }
    },
    async handleImageUploadChange(ev) {
      const files = ev.target.files
      if (!files) return

      await this.uploadFiles(files)

      // 清理文件输入框
      this.$refs.imageInput.value = null
    },
    async uploadFiles(files) {
      if (files.length === 0) {
        return
      }

      if (this.imageCount + files.length > this.maxImageCount) {
        this.message = '图片数量超过上限'
        return
      }

      for (let i = 0; i < files.length; i++) {
        await this.upload(files[i])
      }
    },
    async upload(file) {
      try {
        const formData = new FormData()
        formData.append('image', file, file.name)
        const ret = await this.$axios.post('/api/upload', formData)
        this.images.push(ret.url)
      } catch (e) {
        this.message = e.message || e
      }
    },
    removeImg(img) {
      const index = this.images.indexOf(img)
      if (index !== -1) {
        this.images.splice(index, 1)
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.post-tweets-wrapper {
  position: relative;
  /*border: 1px solid #e8e8e8;*/
  width: 100%;

  .tab-list {
    height: 36px;
    border-bottom: 1px solid #e8e8e8;
    display: block;
    zoom: 1;

    .tab-item {
      margin-left: 19px;
      font-size: 15px;
      color: #222;
      line-height: 34px;
      border-bottom: 2px solid transparent;
      cursor: pointer;
      float: left;

      .tab-name {
        padding: 0 3px;
        text-align: center;
      }

      &.current {
        border-bottom-color: #ed4040;
        color: #222;
      }
    }
  }

  .tweets-box {
    padding: 0;
    margin: 0;
    box-sizing: border-box;
    position: relative;

    .title-input {
      width: 100%;
      height: 100px;
      display: block;
      font-size: 14px;
      line-height: 1.4;
      padding: 10px;
      border: 0;
      outline: 0;
      resize: none;
      overflow: auto;
      background-color: #f4f5f6;
      transition: all 0.5s;
      animation-duration: 0.8s;
      animation-fill-mode: both;
    }

    .words-number {
      position: absolute;
      z-index: 3;
      bottom: 50px;
      right: 10px;
      display: inline-block;
      background-color: rgba(0, 0, 0, 0.5);
      border-radius: 50px;
      padding: 0 8px;
      color: #fff;
      font-size: 13px;
    }

    .box-footer {
      border-top: 1px solid #e8e8e8;
      height: 36px;
      display: block;
      zoom: 1;

      .bui-left {
        float: left;
        user-select: none;

        .action-btn {
          color: #222;
          font-size: 14px;
          line-height: 32px;
          display: inline-block;
          vertical-align: middle;
          margin: 0 0 0 20px;
          cursor: pointer;
          user-select: none;

          .iconfont {
            font-size: 20px;
            color: #ed4040;
            top: 2px;
            position: relative;
          }
        }
      }

      .bui-right {
        float: right;
        text-align: right;
        user-select: none;

        .msg-tip {
          color: #ed4040;
          font-size: 12px;
          margin-right: 10px;
        }

        .tweets-help {
          color: #c2c2c2;
          font-size: 13px;
          user-select: none;
        }

        .upload-publish {
          display: inline-block;
          width: 120px;
          line-height: 36px;
          text-align: center;
          font-size: 14px;
          font-weight: 700;
          background-color: #ed4040;
          color: #fff;
          opacity: 0.6;
          user-select: none;

          &.active {
            opacity: 1;
          }
        }
      }
    }

    .uploader-popup {
      position: absolute;
      bottom: -245px;
      left: -1px;
      width: 420px;
      height: 246px;
      background-color: #fff;
      border: 1px solid #e8e8e8;
      box-shadow: 0 2px 15px rgba(0, 0, 0, 0.12);
      border-radius: 5px;
      z-index: 300;

      &:after {
        content: '';
        width: 0;
        height: 0;
        border: 6px solid transparent;
        border-bottom-color: #fff;
        position: absolute;
        top: -12px;
        left: 36px;
      }

      .imgUploadBox {
        padding: 12px;

        .uploader-title {
          margin-bottom: 5px;
          font-size: 14px;
          color: #222;
        }

        .uploader-meta {
          color: #999;
          font-size: 14px;
          margin-bottom: 10px;
        }

        .close-popup {
          position: absolute;
          top: 8px;
          right: 6px;
          color: #cacaca;
          font-size: 18px;
          font-weight: 700;
          cursor: pointer;
        }

        .upload-box {
          padding: 0;
          margin: 0;
          box-sizing: border-box;
          position: relative;
          display: block;
          zoom: 1;

          .upload-img-list {
            font-size: 0;
            margin-right: -8px;

            .upload-img-add {
              background-color: #fff !important;
              text-align: center;

              i.icon-add {
                font-size: 30px;
                color: rgb(221, 221, 221);
              }
            }

            .upload-img-item {
              cursor: pointer;
              border: 2px dashed #ddd;
              line-height: 72px;
              text-align: center;

              display: inline-block;
              vertical-align: middle;
              width: 72px;
              height: 72px;
              margin: 0 8px 8px 0;
              background-color: #e8e8e8;
              background-size: 32px 32px;
              background-position: 50%;
              background-repeat: no-repeat;
              overflow: hidden;
              position: relative;

              &:hover {
                i.remove {
                  font-size: 14px;
                  font-weight: 700;
                  color: #fff;
                  position: absolute;
                  top: 3px;
                  right: 3px;
                  width: 16px;
                  height: 16px;
                  line-height: 16px;
                  text-align: center;
                  background-color: #ed4040;
                  border-radius: 50%;
                }
              }
            }
          }
        }
      }
    }
  }
}
</style>
