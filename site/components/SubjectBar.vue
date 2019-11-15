<template>
  <div class="left-container">
    <ul class="m-nav">
      <li>
        <a :class="{ active: subject && subject.id === 1 }" href="/subject/1">
          <i class="iconfont icon-topic" />
          <span>Go</span>
        </a>
      </li>
      <li>
        <a :class="{ active: subject && subject.id === 2 }" href="/subject/2">
          <i class="iconfont icon-topic" />
          <span>Java</span>
        </a>
      </li>
      <li>
        <a :class="{ active: subject && subject.id === 3 }" href="/subject/3">
          <i class="iconfont icon-topic" />
          <span>Python</span>
        </a>
      </li>
      <li>
        <a :class="{ active: subject && subject.id === 4 }" href="/subject/4">
          <i class="iconfont icon-topic" />
          <span>Android</span>
        </a>
      </li>
      <li>
        <a :class="{ active: subject && subject.id === 5 }" href="/subject/5">
          <i class="iconfont icon-topic" />
          <span>程序员</span>
        </a>
      </li>
    </ul>
    <!-- 展示广告190x90 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:90px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="9345305153"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>

    <!-- 展示广告190x480 -->
    <ins
      class="adsbygoogle"
      style="display:inline-block;width:190px;height:480px"
      data-ad-client="ca-pub-5683711753850351"
      data-ad-slot="3438372357"
    />
    <script>
      ;(adsbygoogle = window.adsbygoogle || []).push({})
    </script>
  </div>
</template>

<script>
export default {
  props: {
    subject: {
      type: Object,
      default: null
    }
  }
}
</script>

<style lang="scss" scoped>
.subject-aside {
  background: #fff;
  border: 1px solid #eaeaea;
  font-weight: bold;
  margin: 0 0 20px 0;
  list-style: none;
  padding: 0;
  border-radius: 4px;
  padding: 7px 0;

  li {
    a {
      display: block;
      padding: 7px 20px;
      color: #444;
      text-decoration: none;
    }

    a:hover {
      // background: #72d0eb;
      // color: #fff;
      // text-shadow: none;
      color: #72d0eb;
    }

    a.active {
      color: #428bca;
    }
  }
}
</style>
