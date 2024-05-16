<template>
  <section class="main">
    <div class="container">
      <div class="main-body redirect">
        <div>
          <img src="~/assets/images/logo.png" style="max-width: 100px" />
        </div>
        <div style="margin: 20px 0">
          <a :href="url" rel="nofollow"
            >即将跳往站外地址，点击该链接继续跳转&nbsp;&gt;&gt;</a
          >
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const route = useRoute();
const url = route.query.url || "";
const temp = url.toLowerCase();

if (!temp.startsWith("http://") && !temp.startsWith("https://")) {
  throw createError({
    statusCode: 500,
    message: "你访问的页面发生错误!",
  });
}
</script>

<style lang="scss" scoped>
.redirect {
  text-align: center;
  vertical-align: center;
  padding: 100px 0;
}
</style>
