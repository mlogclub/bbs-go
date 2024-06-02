<template>
  <div class="select-tags">
    <input id="tags" v-model="tags" name="tags" type="hidden" />
    <div class="tags-selected">
      <div v-for="tag in tags" :key="tag" class="tag-item">
        <span>{{ tag }}</span>
        <i
          :data-name="tag"
          class="iconfont icon-close"
          @click="clickRemoveTag"
        />
      </div>
    </div>
    <input
      ref="tagInput"
      v-model="inputTag"
      placeholder="标签"
      class="input"
      type="text"
      @input="autocomplete"
      @keydown.delete="removeTag"
      @keydown.enter="addTag"
      @keydown.;="addTag"
      @keydown.,="addTag"
      @keydown.up="selectUp"
      @keydown.down="selectDown"
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

<script setup>
const props = defineProps({
  modelValue: {
    type: Array,
    default() {
      return [];
    },
  },
});

const configStore = useConfigStore();

const maxTagCount = 3;
const maxWordCount = 15;

const tagInput = ref(null);
const tags = ref(props.modelValue || []);
const showRecommendTags = ref(false);
const inputTag = ref("");
const autocompleteTags = ref([]);
const selectIndex = ref(-1);

const emits = defineEmits(["update:modelValue"]);

const recommendTags = computed(() => {
  return configStore.config.recommendTags;
});

function removeTag(event, tag) {
  if (event.currentTarget.value) {
    return;
  }
  const selectionStart = tagInput.value.selectionStart;
  if (!inputTag.value || selectionStart === 0) {
    // input框没内容，或者光标在首位的时候就删除最后一个标签
    tags.value.splice(tags.value.length - 1, 1);
    emits("update:modelValue", tags.value);
  }
}

function clickRemoveTag(event) {
  const tag = event.target.dataset.name;
  if (tag) {
    const index = tags.value.indexOf(tag);
    if (index !== -1) {
      tags.value.splice(index, 1);
      emits("update:modelValue", tags.value);
    }
  }
}

/**
 * 手动点击选择标签
 * @param index
 */
function selectTag(index) {
  selectIndex.value = index;
  addTag();
}

/**
 * 添加标签
 * @param event
 */
function addTag(event) {
  if (event) {
    event.stopPropagation();
    event.preventDefault();
  }

  if (
    selectIndex.value >= 0 &&
    autocompleteTags.value.length > selectIndex.value
  ) {
    addTagName(autocompleteTags.value[selectIndex.value]);
  } else {
    addTagName(inputTag.value);
  }
  autocompleteTags.value = [];
  selectIndex.value = -1;
}

/**
 * 添加推荐标签
 * @param tagName
 */
function addRecommendTag(tagName) {
  addTagName(tagName);
  closeRecommendTags();
}

/**
 * 添加标签
 * @param tagName 标签名称
 * @returns {boolean} 是否成功
 */
function addTagName(tagName) {
  if (!tagName) {
    return false;
  }

  // 最多四个标签
  if (tags.value && tags.value.length >= maxTagCount) {
    return false;
  }

  // 每个标签最多15个字符
  if (tagName.length > maxWordCount) {
    return false;
  }

  // 标签已经存在
  if (tags.value && tags.value.includes(tagName)) {
    return false;
  }

  tags.value.push(tagName);
  inputTag.value = "";
  emits("update:modelValue", tags.value);
  return true;
}

async function autocomplete() {
  closeRecommendTags();
  selectIndex.value = -1;
  if (!inputTag.value) {
    autocompleteTags.value = [];
  } else {
    const ret = await useHttpPostForm("/api/tag/autocomplete", {
      body: {
        input: inputTag.value,
      },
    });
    autocompleteTags.value = [];
    if (ret.length > 0) {
      for (let i = 0; i < ret.length; i++) {
        autocompleteTags.value.push(ret[i].name);
      }
    }
  }
}

function selectUp(event) {
  event.stopPropagation();
  event.preventDefault();
  const curIndex = selectIndex.value;
  const maxIndex = autocompleteTags.value.length - 1;
  if (maxIndex < 0 || curIndex < 0) {
    // 已经在最顶部
    return;
  }
  selectIndex.value--;
}

function selectDown(event) {
  event.stopPropagation();
  event.preventDefault();
  const curIndex = selectIndex.value;
  const maxIndex = autocompleteTags.value.length - 1;
  if (maxIndex < 0 || curIndex >= maxIndex) {
    // 已经在最底部
    return;
  }
  selectIndex.value++;
}

// 关闭推荐
function openRecommendTags() {
  showRecommendTags.value = true;
}

// 开启推荐
function closeRecommendTags() {
  setTimeout(() => {
    showRecommendTags.value = false;
  }, 300);
}

// 关闭自动补全
function close() {
  if (autocompleteTags.value && autocompleteTags.value.length > 0) {
    autocompleteTags.value = [];
    selectIndex.value = -1;
  }
  closeRecommendTags();
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

    &:focus-visible {
      outline-width: 0;
    }
  }

  .tags-selected {
    display: flex;

    .tag-item {
      margin: 5px;
      padding: 0 10px;
      background: var(--bg-color3);
      color: var(--text-color);
      line-height: 30px;
      border-radius: 5px;

      text-align: center;
      font-size: 12px;
      white-space: nowrap;

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
      // background: #f7f7f7;
      background-color: var(--bg-color);
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
            // background-color: var(--bg-color2);
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
