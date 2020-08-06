<template>
  <div class="home-icons">
    <div class="item">
      <div class="icon">
        <a href="/topics"
          ><i class="iconfont icon-topic" style="color: #ef7853;"></i
        ></a>
      </div>
      <div class="text">
        <a href="/topics">码农社区 <i class="iconfont icon-arrowright"></i></a>
      </div>
    </div>
    <!--
    <div class="item">
      <div class="icon">
        <a href="/articles"
          ><i class="iconfont icon-article" style="color:#fee83c;"></i
        ></a>
      </div>
      <div class="text">
        <a href="/articles"
          >码农资讯 <i class="iconfont icon-arrowright"></i
        ></a>
      </div>
    </div>
    -->
    <div class="item">
      <div class="icon">
        <a href="/topic/40"
          ><i class="iconfont icon-book" style="color: #fee83c;"></i
        ></a>
      </div>
      <div class="text">
        <a href="/topic/40"
          >在线书籍 <i class="iconfont icon-arrowright"></i
        ></a>
      </div>
    </div>
    <div class="item">
      <div class="icon">
        <a href="/projects"
          ><i class="iconfont icon-open-source" style="color: #4c90e3;"></i
        ></a>
      </div>
      <div class="text">
        <a href="/projects"
          >开源项目 <i class="iconfont icon-arrowright"></i
        ></a>
      </div>
    </div>
    <div class="item">
      <div class="icon">
        <a href="/links"
          ><i class="iconfont icon-blog" style="color: #21d503;"></i
        ></a>
      </div>
      <div class="text">
        <a href="/links">博客导航 <i class="iconfont icon-arrowright"></i></a>
      </div>
    </div>
  </div>
</template>

<script>
export default {}
</script>

<style lang="scss" scoped>
.home-icons {
  display: flex;

  @media screen and (max-width: 768px) {
    & {
      display: none;
    }
  }

  .item {
    display: block;
    background: #fff;
    border-radius: 3px;
    width: 100%;
    text-align: center;

    &:not(:first-child) {
      margin-left: 20px;
    }

    .icon {
      padding: 10px;
      height: 80px;
      i {
        font-size: 60px;
      }
    }

    .text {
      border-top: 1px solid #e9e9e9;
      font-size: 14px;
      font-weight: bold;
      display: block;

      a {
        display: block;
        padding: 6px 15px;
        color: #666;

        i {
          float: right;
          font-size: 12px;
        }
      }
    }
  }
}
</style>
