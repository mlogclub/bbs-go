<template>
  <div>
    <pc-nav class="pc-nav" />
    <mobile-nav class="mobile-nav" />
  </div>
</template>

<script>
export default {}
</script>

<style lang="scss">
// 当现实顶部导航栏的时候，顶部预留一段空白，防止重叠
body {
  @media screen and (max-width: 1024px) {
    padding-top: 46px;
  }
}
</style>

<style scoped>
@media screen and (max-width: 1024px) {
  .pc-nav {
    display: none !important;
  }
}

@media screen and (min-width: 1024px) {
  .mobile-nav {
    display: none !important;
  }
}
</style>
