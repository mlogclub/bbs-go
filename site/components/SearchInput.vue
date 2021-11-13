<template>
  <div
    ref="searchForm"
    class="searchFormDiv"
    :class="{ 'input-focus': inputFocus, 'show-histories': showHistories }"
  >
    <div class="control has-icons-right">
      <input
        v-model="keyword"
        name="q"
        class="input"
        type="text"
        maxlength="30"
        placeholder="搜索"
        autocomplete="off"
        @focus="onFocus"
        @blur="onBlur"
        @input="onInput"
        @keyup.esc="onBlur"
        @keyup.down="changeSelect(1)"
        @keyup.up="changeSelect(-1)"
        @keyup.enter="searchBoxOnEnter"
      />
      <span class="icon is-medium is-right">
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
          <i class="iconfont icon-delete" @click="deleteHistory(item)" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
const localStorageKey = 'search.histories'
const maxHistoryLen = 10

export default {
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
      console.log('onFocus')
      this.inputFocus = true
    },
    onBlur() {
      console.log('onBlur')
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

  $search-box-width: 230px;
  $search-focus-box-width: 430px;
  $border-color: #ebebeb;
  //$border-focus-color: #4e6ef2;
  $border-focus-color: #e7672e;

  &.input-focus {
    .input {
      width: $search-focus-box-width;
      // border-color: $border-color !important;
      opacity: 1;
      filter: alpha(opacity=100) \9;
    }
  }

  &.show-histories {
    .input {
      border-radius: 6px 6px 0 0;
      border-bottom: 1px solid $border-color !important;
      border-top: 2px solid $border-focus-color;
      border-left: 2px solid $border-focus-color;
      border-right: 2px solid $border-focus-color;
    }
    .histories {
      width: $search-focus-box-width;
      display: block;
    }
  }

  .input {
    transition: width 0.4s;
    width: $search-box-width;
    box-shadow: none;
    background-color: #fff;
    float: right;
    position: relative;
    border-radius: 6px;
    border: 1px solid $border-color;
  }

  .histories {
    transition: all 0.4s;
    display: none;
    height: auto;
    width: $search-box-width;
    top: 48px;
    border-radius: 0 0 6px 6px;
    border: 2px solid $border-focus-color !important;
    border-top: 0 !important;
    box-shadow: none;
    font-family: Arial, 'PingFang SC', 'Microsoft YaHei', sans-serif;
    z-index: 1;
    position: absolute;
    background: #fff;

    ul {
      li {
        padding: 10px;
        color: #626675;
        display: flex;
        cursor: pointer;
        span {
          width: 100%;
          font-size: 13px;
        }
        &:hover {
          i.icon-delete {
            display: block;
          }
        }
        i.icon-delete {
          display: none;
          font-size: 13px;
          padding: 0px 6px;
          &:hover {
            color: #4e6ef2;
          }
        }
        &.selected {
          span {
            color: #4e6ef2;
          }
          background-color: #f5f5f6;
        }
      }
    }
  }
}
</style>
