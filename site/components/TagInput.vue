<template>
  <div class="select-tags">
    <input id="tags" v-model="tags" name="tags" type="hidden" />
    <div class="tags-selected">
      <span v-for="tag in tags" :key="tag" class="tag-item">
        <span class="text"
          >{{ tag
          }}<i
            :data-name="tag"
            class="iconfont icon-close"
            @click="clickRemoveTag"
        /></span>
      </span>
    </div>
    <input
      ref="tagInput"
      v-model="inputTag"
      :placeholder="
        '标签（请用逗号分隔每个标签，最多' +
        maxTagCount +
        '个，每个最长15字符）'
      "
      class="input"
      type="text"
      @input="autocomplete"
      @keydown.delete="removeTag"
      @keydown.enter="addTag"
      @keydown.32="addTag"
      @keydown.186="addTag"
      @keydown.188="addTag"
      @keydown.38="selectUp"
      @keydown.40="selectDown"
      @keydown.esc="close"
      @focus="openRecommendTags"
      @blur="closeRecommendTags"
      @click="openRecommendTags"
    />
    <transition name="el-zoom-in-bottom">
      <div v-show="autocompleteTags.length > 0" class="autocomplete-tags">
        <div class="tags-container">
          <section class="tag-section">
            <div
              v-for="(item, index) in autocompleteTags"
              :key="item"
              :class="{ active: index === selectIndex }"
              class="tag-item"
              @click="selectTag(index)"
              v-text="item"
            />
          </section>
        </div>
      </div>
    </transition>
    <transition name="el-zoom-in-bottom">
      <div v-show="showRecommendTags" class="recommend-tags">
        <div class="tags-container">
          <div class="header">
            <span>推荐标签</span>
            <span class="close-recommend"
              ><i class="iconfont icon-close" @click="closeRecommendTags"
            /></span>
          </div>
          <a
            v-for="tag in recommendTags"
            :key="tag"
            class="tag-item"
            @click="addRecommendTag(tag)"
            v-text="tag"
          />
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  props: {
    value: {
      type: Array,
      default() {
        return []
      },
    },
  },
  data() {
    return {
      tags: this.value || [],
      maxTagCount: 3, // 最多可以选择的标签数量
      maxWordCount: 15, // 每个标签最大字数
      showRecommendTags: false, // 是否显示推荐
      inputTag: '',
      autocompleteTags: [],
      selectIndex: -1,
    }
  },
  computed: {
    // 推荐标签
    recommendTags() {
      return this.$store.state.config.config.recommendTags
    },
  },
  methods: {
    removeTag(event, tag) {
      const selectionStart = this.$refs.tagInput.selectionStart
      if (!this.inputTag || selectionStart === 0) {
        // input框没内容，或者光标在首位的时候就删除最后一个标签
        this.tags.splice(this.tags.length - 1, 1)
        this.$emit('input', this.tags)
      }
    },

    clickRemoveTag(event) {
      const tag = event.target.dataset.name
      if (tag) {
        const index = this.tags.indexOf(tag)
        if (index !== -1) {
          this.tags.splice(index, 1)
          this.$emit('input', this.tags)
        }
      }
    },

    /**
     * 手动点击选择标签
     * @param index
     */
    selectTag(index) {
      this.selectIndex = index
      this.addTag()
    },

    /**
     * 添加标签
     * @param event
     */
    addTag(event) {
      if (event) {
        event.stopPropagation()
        event.preventDefault()
      }

      if (
        this.selectIndex >= 0 &&
        this.autocompleteTags.length > this.selectIndex
      ) {
        this.addTagName(this.autocompleteTags[this.selectIndex])
      } else {
        this.addTagName(this.inputTag)
      }
      this.autocompleteTags = []
      this.selectIndex = -1
    },

    /**
     * 添加推荐标签
     * @param tagName
     */
    addRecommendTag(tagName) {
      this.addTagName(tagName)
      this.closeRecommendTags()
    },

    /**
     * 添加标签
     * @param tagName 标签名称
     * @returns {boolean} 是否成功
     */
    addTagName(tagName) {
      if (!tagName) {
        return false
      }

      // 最多四个标签
      if (this.tags && this.tags.length >= this.maxTagCount) {
        return false
      }

      // 每个标签最多15个字符
      if (tagName.length > this.maxWordCount) {
        return false
      }

      // 标签已经存在
      if (this.tags && this.tags.includes(tagName)) {
        return false
      }

      this.tags.push(tagName)
      this.inputTag = ''
      this.$emit('input', this.tags)
      return true
    },

    async autocomplete() {
      this.closeRecommendTags()
      this.selectIndex = -1
      if (!this.inputTag) {
        this.autocompleteTags = []
      } else {
        const ret = await this.$axios.post('/api/tag/autocomplete', {
          input: this.inputTag,
        })
        this.autocompleteTags = []
        if (ret.length > 0) {
          for (let i = 0; i < ret.length; i++) {
            this.autocompleteTags.push(ret[i].name)
          }
        }
      }
    },

    selectUp(event) {
      event.stopPropagation()
      event.preventDefault()
      const curIndex = this.selectIndex
      const maxIndex = this.autocompleteTags.length - 1
      if (maxIndex < 0 || curIndex < 0) {
        // 已经在最顶部
        return
      }
      this.selectIndex--
    },

    selectDown(event) {
      event.stopPropagation()
      event.preventDefault()
      const curIndex = this.selectIndex
      const maxIndex = this.autocompleteTags.length - 1
      if (maxIndex < 0 || curIndex >= maxIndex) {
        // 已经在最底部
        return
      }
      this.selectIndex++
    },

    // 关闭推荐
    openRecommendTags() {
      this.showRecommendTags = true
    },

    // 开启推荐
    closeRecommendTags() {
      setTimeout(() => {
        this.showRecommendTags = false
      }, 300)
    },

    // 关闭自动补全
    close() {
      if (this.autocompleteTags && this.autocompleteTags.length > 0) {
        this.autocompleteTags = []
        this.selectIndex = -1
      }
      this.closeRecommendTags()
    },
  },
}
</script>

<style lang="scss" scoped>
.select-tags {
  display: flex;
  background-color: var(--bg-color);
  border: 1px solid var(--border-color2);
  border-radius: 4px;
  box-shadow: inset 0 1px 2px rgba(10, 10, 10, 0.1);
  color: var(--text-color);
  padding: 0 8px;

  .input {
    border: none;
    box-shadow: none;
    margin: 0;
    padding: 0;
  }

  .tags-selected {
    display: flex;

    span.tag-item {
      margin: 5px;
      padding: 0 10px;
      background: #eee;
      color: var(--text-color);
      line-height: 30px;
      border-radius: 5px;

      .text {
        text-align: center;
        vertical-align: middle;
        font-size: 12px;
        color: rgba(0, 0, 0, 0.5);
        white-space: nowrap;
        display: inline-block;

        i {
          font-size: 12px;
          margin-left: 4px;
        }

        i:hover {
          color: red;
          cursor: pointer;
        }
      }
    }
  }

  .autocomplete-tags {
    z-index: 2000;
    left: 0;
    right: 0;
    top: 42px;
    bottom: 0;
    position: absolute;

    .tags-container {
      scroll-behavior: smooth;
      position: relative;
      background: #f7f7f7;
      border-left: 1px solid var(--border-color2);
      border-right: 1px solid var(--border-color2);
      border-bottom: 1px solid var(--border-color2);

      .tag-section {
        font-size: 14px;
        line-height: 16px;

        .tag-item {
          padding: 8px 15px;
          cursor: pointer;

          &.active,
          &:hover {
            color: var(--text-color5);
            background: #006bde;
          }
        }
      }
    }
  }

  .recommend-tags {
    z-index: 2000;
    left: 0;
    right: 0;
    top: 42px;
    bottom: 0;
    position: absolute;

    .tags-container {
      scroll-behavior: smooth;
      position: relative;
      background: #f7f7f7;
      border-left: 1px solid var(--border-color2);
      border-right: 1px solid var(--border-color2);
      border-bottom: 1px solid var(--border-color2);
      padding: 0 10px 10px 10px;

      .header {
        font-weight: bold;
        font-size: 15px;
        color: #017e66;
        border-bottom: 1px solid var(--border-color2);
        margin-bottom: 5px;
        padding-top: 5px;
        padding-bottom: 5px;

        .close-recommend {
          float: right;
          cursor: pointer;
          &:hover {
            color: red;
          }
        }
      }

      .tag-item {
        padding: 0 11px;
        border-radius: 5px;
        display: inline-block;
        color: #017e66;
        background-color: rgba(1, 126, 102, 0.08);
        height: 22px;
        line-height: 22px;
        font-weight: normal;
        font-size: 13px;
        text-align: center;

        &:not(:last-child) {
          margin-right: 5px;
        }

        img {
          width: 16px;
          height: 16px;
          margin-right: 5px;
          margin-top: -1px;
          vertical-align: middle;
        }
      }

      .tag-item:hover,
      .tag-item:focus {
        background-color: #017e66;
        color: var(--text-color5);
        text-decoration: none;
      }
    }
  }
}
</style>
