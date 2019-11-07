<template>
  <div>
    <my-nav />
    <nuxt />
  </div>
</template>

<script>
import MyNav from '~/components/MyNav'

export default {
  middleware: 'global',
  components: {
    MyNav
  }
}
</script>
