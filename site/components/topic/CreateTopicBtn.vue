<template>
  <el-popover
    placement="bottom"
    width="120"
    min-width="120"
    trigger="click"
    popper-class="custom-popover-width"
  >
    <ul class="type-container">
      <li>帖子</li>
      <li>动态</li>
    </ul>
    <el-button slot="reference" type="primary" class="new-post">
      <span class="add-icon">+</span>发帖
    </el-button>
  </el-popover>
</template>
<script>
export default {
  data() {
    return {
      visible: false,
    }
  },
  computed: {},
  methods: {},
}
</script>
<style lang="scss" scoped>
@import '@/assets/styles/variable.scss';

// 重置element
.new-post {
  background-color: $color-blue-base !important;
  border-color: $color-blue-base !important;
  padding: 9.5px 13.5px !important;
  border-radius: 2px;
  &:hover {
    background-color: $color-blue-deep !important;
    border-color: $color-blue-deep !important;
  }
  .add-icon {
    display: inline-block;
    margin-right: 5px;
  }
}
.type-container {
  padding: 5px 0;
  li {
    position: relative;
    list-style: none;
    width: 120px;
    height: 40px;
    line-height: 40px;
    text-align: center;
    cursor: pointer;
    outline: none;
    &:hover {
      background-color: #e5f2ff;
      color: $color-blue-base;
    }
    &:before {
      content: '';
      position: absolute;
      left: 0;
      bottom: 0;
      width: 100px;
      height: 1px;
      background: #ededed;
      margin: 0 10px;
    }
    &:last-child:before {
      height: 0;
    }
  }
}
</style>
