<template>
  <div v-viewer="viewerOptions" class="image-uploads">
    <div
      v-for="(image, index) in previewFiles"
      :key="index"
      class="preview-item"
      :class="{ deleted: image.deleted }"
      :style="{ width: size, height: size, margin: previewItemMargin }"
    >
      <img :src="image.url" class="image-item" />
      <el-progress
        v-show="image.progress < 100"
        :percentage="image.progress"
        color="#25A9F6"
        :show-text="false"
        class="progress"
      />
      <div v-show="image.progress < 100" class="cover">上传中...</div>
      <div
        :class="{
          'upload-delete': true,
          'show-delete': image.progress === 100,
        }"
        @click="removeItem(index)"
      >
        <i class="iconfont icon-delete" />
      </div>
    </div>
    <div
      v-show="previewFiles.length < limit"
      class="add-image-btn"
      :style="{ width: size, height: size }"
      @click="onClick($event)"
    >
      <input
        ref="uploadImage"
        :accept="accept"
        type="file"
        multiple
        @input="onInput"
      />
      <slot name="add-image-button">
        <i class="iconfont icon-add" />
      </slot>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    accept: {
      type: String,
      default: 'image/*',
    },
    limit: {
      type: Number,
      default: 9,
    },
    sizeLimit: {
      type: Number,
      default: 1024 * 1024 * 20,
    },
    onUpload: {
      type: Boolean,
      default: false,
    },
    size: {
      type: String,
      default: '94px',
    },
  },
  data() {
    return {
      fileList: [],
      previewFiles: [],
      currentInput: null,
      viewerOptions: {
        zIndex: 9999,
        navbar: false,
        title: false,
        tooltip: false,
        movable: false,
        scalable: false,
        url: 'src',
        filter(image) {
          return (
            image.classList &&
            image.classList.length &&
            image.classList.contains('image-item')
          )
        },
      },
    }
  },
  computed: {
    previewItemMargin() {
      if (this.previewFiles.length > 1 || this.limit > 1) {
        // margin-right: 10px;
        // margin-bottom: 10px;
        return '0 10px 10px 0'
      }
      return '0'
    },
  },
  watch: {
    fileList: {
      handler() {
        if (
          this.fileList.length > this.previewFiles.length &&
          this.previewFiles.length === 0
        ) {
          // 初始试，回显
          this.previewFiles.push(...this.fileList)
          this.previewFiles.map((item) => {
            item.progress = 100
          }) // 增加 deleted progress 属性
        }
      },
      deep: true,
      immediate: true,
    },
  },
  methods: {
    onClick() {
      const currentObj = this.$refs.uploadImage
      this.currentInput = currentObj
      currentObj.dispatchEvent(new MouseEvent('click'))
    },
    onInput(e) {
      const files = e.target.files
      this.addFiles(files)
    },
    addFiles(files) {
      if (!files || !files.length) return // 没有文件
      if (!this.checkSizeLimit(files)) return // 文件大小检查
      if (!this.checkLengthLimit(files)) return // 文件数量检查

      const fileArray = []
      for (let i = 0; i < files.length; i++) {
        const url = this.getObjectURL(files[i])
        this.previewFiles.push({
          name: files[i].name,
          url,
          progress: 0,
          deleted: false,
          size: files[i].size,
        })
        fileArray.push(files[i])
      }
      const promiseList = fileArray.reduce((result, file, index, array) => {
        result.push(this.uploadFile(file, index, array.length))
        return result
      }, [])
      this.uploadFiles(promiseList)
    },
    uploadFile(file, index, length) {
      const me = this
      const config = {
        onUploadProgress(progressEvent) {
          if (!progressEvent.lengthComputable) {
            // 当进度不可估量,直接等于 100
            me.previewFiles[
              me.previewFiles.length - length + index
            ].progress = 100
            return
          }
          me.previewFiles[me.previewFiles.length - length + index].progress =
            parseInt(
              Math.round(
                (progressEvent.loaded / progressEvent.total) * 100
              ).toString()
            ) * 0.9
        },
      }
      const formData = new FormData()
      formData.append('image', file, file.name)
      return this.$axios.post('/api/upload', formData, config)
    },
    uploadFiles(promiseList) {
      this.$emit(`update:onUpload`, true)
      Promise.all(promiseList).then(
        (resList) => {
          this.previewFiles.map((item) => {
            item.progress = 100
          }) // 请求响应后，更新到 100%
          // const _fileList = [...this.fileList]
          // resList.forEach((item) => _fileList.push(item))
          resList.forEach((item) => this.fileList.push(item))
          if (this.currentInput) {
            this.currentInput.value = ''
          }
          this.$emit('input', this.fileList)
          this.$emit(`update:onUpload`, false)
        },
        (e) => {
          // 失败的时候取消对应的预览照片
          if (this.currentInput) {
            this.currentInput.value = ''
          }
          const length = promiseList.length
          this.$emit(`update:onUpload`, false)
          this.previewFiles.splice(this.previewFiles.length - length, length)
          // this.handleError(e).then(() => {})
          console.error(e)
        }
      )
    },
    removeItem(index) {
      this.$confirm('确定删除此内容吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(
        () => {
          this.previewFiles[index].deleted = true // 删除动画
          // const _fileList = [...this.fileList]
          // _fileList.splice(index, 1)
          this.fileList.splice(index, 1)
          this.$emit('input', this.fileList) // 避免和回显冲突，先修改 fileList
          setTimeout(() => {
            this.previewFiles.splice(index, 1)
            this.$message.success('删除成功')
          }, 900)
        },
        () => console.log('取消删除')
      )
    },
    checkSizeLimit(files) {
      let pass = true
      for (let i = 0; i < files.length; i++) {
        if (files[i].size > this.sizeLimit) {
          pass = false
        }
      }
      if (!pass)
        this.$message.error(
          `图片大小不可超过 ${this.sizeLimit / 1024 / 1024} MB`
        )
      return pass
    },
    checkLengthLimit(files) {
      if (this.previewFiles.length + files.length > this.limit) {
        this.$message.warning(`图片最多上传${this.limit}张`)
        this.$emit('exceed', files)
        return false
      } else {
        return true
      }
    },
    getObjectURL(file) {
      let url = null
      if (window.createObjectURL) {
        // basic
        url = window.createObjectURL(file)
      } else if (window.URL) {
        // mozilla(firefox)
        url = window.URL.createObjectURL(file)
      } else if (window.webkitURL) {
        // webkit or chrome
        url = window.webkitURL.createObjectURL(file)
      }
      return url
    },
    clear() {
      this.fileList = []
      this.previewFiles = []
    },
  },
}
</script>

<style lang="scss" scoped>
.image-uploads {
  display: flex;
  flex-wrap: wrap;

  .preview-item {
    position: relative;
    border: 2px dashed var(--border-color);

    &.deleted {
      transition: 1s all;
      transform: translateY(-100%);
      opacity: 0;
    }

    .image-item {
      cursor: pointer;
      // border-radius: 5px;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .progress {
      position: absolute;
      top: 80px;
      width: 100%;
      height: 6px;
      padding: 0 10px;
    }

    .cover {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      color: var(--text-color2);
      background: rgba(255, 255, 255, 0.5);
      font-size: 12px;
      display: flex;
      justify-content: center;
      align-items: center;
    }

    .upload-delete {
      cursor: pointer;
      position: absolute;
      left: 0;
      bottom: 0;
      height: 20px;
      width: 100%;
      display: none;
      justify-content: center;
      align-items: center;
      background: rgba(0, 0, 0, 0.3);
      text-align: center;
      vertical-align: middle;
      line-height: 20px;

      i.iconfont {
        font-size: 14px;
        fill: white;
        color: var(--text-color5);
        font-weight: 700;
      }
    }

    &:hover {
      .upload-delete.show-delete {
        display: flex;
      }
    }
  }

  .add-image-btn {
    cursor: pointer;
    border: 2px dashed var(--border-color);
    position: relative;

    input[type='file'] {
      cursor: pointer;
      display: none;
    }

    .icon-add {
      font-size: 30px;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      color: #1878f3; // TODO
    }
  }
}
</style>
