<template>
  <div
    ref="searchForm"
    v-click-outside="onBlur"
    class="searchFormDiv"
    :class="{ 'input-focus': inputFocus, 'show-histories': showHistories }"
  >
    <div class="search-input">
      <input
        v-model="keyword"
        name="q"
        class="input"
        type="text"
        maxlength="30"
        placeholder="输入你想查找的内容"
        autocomplete="off"
        @focus="onFocus"
        @input="onInput"
        @keyup.down="changeSelect(1)"
        @keyup.up="changeSelect(-1)"
        @keyup.enter="searchBoxOnEnter"
      />
      <span>
        <i class="iconfont icon-search" />
      </span>
    </div>
    <div class="histories">
      <ul>
        <li
          v-for="(item, index) in histories"
          :key="index"
          :class="{ selected: index === selectedIndex }"
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

<script>
import ClickOutside from 'vue-click-outside'
const localStorageKey = 'bbsgo.search.histories'
const maxHistoryLen = 10

export default {
  directives: {
    ClickOutside,
  },
  data() {
    return {
      keyword: '',
      inputFocus: false,
      selectedIndex: -1,
      allHistories: [],
    }
  },
  computed: {
    showHistories() {
      return this.inputFocus && this.histories && this.histories.length
    },
    histories() {
      if (this.keyword) {
        return this.allHistories.filter((history) => {
          return history.includes(this.keyword)
        })
      }
      return this.allHistories
    },
  },
  mounted() {
    this.keyword = this.$store.state.search.keyword
    this.loadAllHistories()
  },
  methods: {
    searchBoxOnEnter() {
      // 如果选中了历史搜索记录，那么使用历史搜索记录
      if (
        this.selectedIndex >= 0 &&
        this.histories &&
        this.histories.length > this.selectedIndex
      ) {
        this.keyword = this.histories[this.selectedIndex]
      }
      this.submitSearch()
    },
    historyItemClick(keyword) {
      this.keyword = keyword
      this.submitSearch()
    },
    submitSearch() {
      if (!this.keyword) {
        return
      }
      this.addHistories()
      window.location = '/search?q=' + encodeURIComponent(this.keyword)
    },
    onFocus() {
      this.inputFocus = true
    },
    onBlur() {
      this.inputFocus = false
    },
    onInput() {
      this.selectedIndex = -1
    },
    changeSelect(delta) {
      if (!this.histories || !this.histories.length) {
        return
      }
      let index = this.selectedIndex + delta
      if (index < 0) {
        // 选中熬第一个了，再往上取消选中
        index = -1
      } else if (index >= this.histories.length) {
        // 选中到最后了，再往下就回到第一个
        index = 0
      }
      this.selectedIndex = index
    },
    historyItemMouseOver(index) {
      this.selectedIndex = index
    },
    historyItemMouseOut() {
      this.selectedIndex = -1
    },
    loadAllHistories() {
      try {
        this.allHistories =
          JSON.parse(localStorage.getItem(localStorageKey)) || []
      } catch (error) {
        this.allHistories = []
      }
    },
    addHistories() {
      if (!this.keyword) {
        return
      }
      const newArray = []
      newArray.push(this.keyword)
      if (this.allHistories && this.allHistories.length) {
        for (let i = 0; i < this.allHistories.length; i++) {
          const element = this.allHistories[i]
          if (newArray.length < maxHistoryLen && !newArray.includes(element)) {
            newArray.push(element)
          }
        }
      }
      localStorage.setItem(localStorageKey, JSON.stringify(newArray))
      this.allHistories = newArray
    },
    deleteHistory(kw) {
      const newArray = []
      if (this.allHistories && this.allHistories.length) {
        for (let i = 0; i < this.allHistories.length; i++) {
          const element = this.allHistories[i]
          if (element !== kw && !newArray.includes(element)) {
            newArray.push(element)
          }
        }
      }
      localStorage.setItem(localStorageKey, JSON.stringify(newArray))
      this.allHistories = newArray
    },
  },
}
</script>

<style lang="scss" scoped>
.searchFormDiv {
  @media screen and (max-width: 768px) {
    & {
      display: none;
    }
  }

  $search-box-width: 380px;
  $border-color: #4e6ef2; // TODO

  &.input-focus {
    .search-input {
      background-color: var(--bg-color);
      border: 1px solid var(--border-color3);
    }
    .icon-search {
      color: #0065ff;
    }
  }

  &.show-histories {
    .histories {
      display: block;
    }
  }

  .search-input {
    width: $search-box-width;
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

    input {
      color: var(--text-color);
      font-weight: 400;
      font-size: 13px;
      height: 24px;
      line-height: 24px;
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
