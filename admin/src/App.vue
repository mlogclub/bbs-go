<template>
  <div id="app">
    <!--    <router-link to="/home">-->
    <!--      Home-->
    <!--    </router-link>-->
    <!--    <router-link to="/about">-->
    <!--      About-->
    <!--    </router-link>-->
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style lang="scss">
html, body {
  width: 100%;
  height: 100%;
  margin: 0;
}

#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 100%;
  height: 100%;
}
</style>
