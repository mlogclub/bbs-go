<template>
  <div
    ref="searchForm"
    class="searchFormDiv"
    v-click-outside="onBlur"
    :class="{ 'input-focus': data.inputFocus, 'show-histories': showHistories }"
  >
    <div class="search-input">
      <input
        v-model="data.keyword"
        name="q"
        class="input"
        type="text"
        maxlength="30"
        placeholder="输入你想查找的内容"
        autocomplete="off"
        aria-autocomplete="off"
        @focus="onFocus"
        @input="onInput"
        @keyup.down="changeSelect(1)"
        @keyup.up="changeSelect(-1)"
        @keyup.enter="searchBoxOnEnter"
      />
      <span @click="submitSearch">
        <i class="iconfont icon-search" />
      </span>
    </div>
    <div class="histories">
      <ul>
        <li
          v-for="(item, index) in histories"
          :key="index"
          :class="{ selected: index === data.selectedIndex }"
          @mouseover="historyItemMouseOver(index)"
          @mouseout="historyItemMouseOut()"
        >
          <span @click="historyItemClick(item)">{{ item }}</span>
          <i class="iconfont icon-close" @click="deleteHistory(item)" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const localStorageKey = "bbsgo.search.histories";
const maxHistoryLen = 10;
const route = useRoute();

const data = reactive({
  keyword: route.query.q || "",
  inputFocus: false,
  selectedIndex: -1,
  allHistories: [],
});

const showHistories = computed(() => {
  return (
    data.inputFocus && histories && histories.value && histories.value.length
  );
});

const histories = computed(() => {
  if (data.keyword) {
    return data.allHistories.filter((history) => {
      return history.includes(data.keyword);
    });
  }
  return data.allHistories;
});

onMounted(() => {
  loadAllHistories();
});

const searchBoxOnEnter = () => {
  // 如果选中了历史搜索记录，那么使用历史搜索记录
  if (
    data.selectedIndex >= 0 &&
    histories.value &&
    histories.value.length > data.selectedIndex
  ) {
    data.keyword = histories.value[data.selectedIndex];
  }
  submitSearch();
};
const historyItemClick = (keyword) => {
  data.keyword = keyword;
  submitSearch();
};
const submitSearch = () => {
  if (!data.keyword) {
    return;
  }
  addHistories();
  window.location = "/search?q=" + encodeURIComponent(data.keyword);
};
const onFocus = () => {
  data.inputFocus = true;
};
const onBlur = () => {
  data.inputFocus = false;
};
const onInput = () => {
  data.selectedIndex = -1;
};
const changeSelect = (delta) => {
  if (!histories.value || !histories.value.length) {
    return;
  }
  let index = data.selectedIndex + delta;
  if (index < 0) {
    // 选中熬第一个了，再往上取消选中
    index = -1;
  } else if (index >= histories.value.length) {
    // 选中到最后了，再往下就回到第一个
    index = 0;
  }
  data.selectedIndex = index;
};
const historyItemMouseOver = (index) => {
  data.selectedIndex = index;
};
const historyItemMouseOut = () => {
  data.selectedIndex = -1;
};
const loadAllHistories = () => {
  try {
    data.allHistories = JSON.parse(localStorage.getItem(localStorageKey)) || [];
  } catch (error) {
    data.allHistories = [];
  }
};
const addHistories = () => {
  if (!data.keyword) {
    return;
  }
  const newArray = [];
  newArray.push(data.keyword);
  if (data.allHistories && data.allHistories.length) {
    for (let i = 0; i < data.allHistories.length; i++) {
      const element = data.allHistories[i];
      if (newArray.length < maxHistoryLen && !newArray.includes(element)) {
        newArray.push(element);
      }
    }
  }
  localStorage.setItem(localStorageKey, JSON.stringify(newArray));
  data.allHistories = newArray;
};
const deleteHistory = (kw) => {
  const newArray = [];
  if (data.allHistories && data.allHistories.length) {
    for (let i = 0; i < data.allHistories.length; i++) {
      const element = data.allHistories[i];
      if (element !== kw && !newArray.includes(element)) {
        newArray.push(element);
      }
    }
  }
  localStorage.setItem(localStorageKey, JSON.stringify(newArray));
  data.allHistories = newArray;
};
</script>

<style lang="scss" scoped>
.searchFormDiv {
  @media screen and (max-width: 768px) {
    & {
      display: none;
    }
  }

  $search-box-width: 280px;
  $search-box-unfocus-width: 180px;
  $focus-color: #0065ff; // TODO

  &.input-focus {
    .search-input {
      background-color: var(--bg-color);
      border: 1px solid var(--border-color3);
    }
    .icon-search {
      color: $focus-color;
    }
  }

  &.show-histories {
    .histories {
      display: block;
    }
  }

  .search-input {
    width: $search-box-unfocus-width;
    background-color: var(--bg-color2);
    border: 1px solid var(--border-color);
    height: 34px;
    border-radius: 17px;
    padding: 4px 10px;
    font-size: 14px;
    display: flex;
    transition-property: background-color, border, color;
    transition-duration: 0.25s;
    transition-timing-function: ease-in;
    transition: width 0.5s;

    &:has(.input:focus) {
      width: $search-box-width;
      border: 1px solid $focus-color;
    }

    input {
      color: var(--text-color);
      font-weight: 400;
      font-size: 13px;
      height: 100%;
      width: 100%;
      flex: 1 1;
      padding: 0;
      margin: 0 5px;
      overflow: hidden;
      background: transparent;
      border: none;
      resize: none;
      box-shadow: none;

      &:focus {
        outline: none;
      }
    }
  }

  span {
    padding: 0 4px 0 10px;
    cursor: pointer;
  }

  .histories {
    width: $search-box-width;
    display: none;
    height: auto;
    top: 48px;
    position: fixed;

    background-color: var(--bg-color);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    box-shadow: 0 5px 20px rgb(18 18 18 / 10%);
    z-index: 203;

    animation: showHistoriesAnimation 0.3s;
    -moz-animation: showHistoriesAnimation 0.3s; /* Firefox */
    -webkit-animation: showHistoriesAnimation 0.3s; /* Safari and Chrome */
    -o-animation: showHistoriesAnimation 0.3s; /* Opera */

    @keyframes showHistoriesAnimation {
      0% {
        opacity: 0;
        transform: translate(1px, 30px);
      }
      100% {
        opacity: 1;
        transform: translate(0, 0);
      }
    }
    ul {
      li {
        padding: 7px 16px;
        width: 100%;
        text-align: left;
        cursor: pointer;
        -webkit-box-sizing: border-box;
        box-sizing: border-box;

        display: flex;
        span {
          width: 100%;
          font-size: 13px;
          font-weight: 400;
          color: var(--text-color);
        }
        &:hover {
          i.iconfont {
            display: block;
          }
        }
        i.iconfont {
          display: none;
          font-size: 13px;
          padding: 0 6px;
          &:hover {
            color: #4e6ef2;
          }
        }
        &.selected {
          background-color: var(--bg-color2);
        }
      }
    }
  }
}
</style>
