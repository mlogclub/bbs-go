<template>
  <div>
    <my-nav />
    <nuxt />
  </div>
</template>

<script>
import MyNav from '~/components/MyNav'

export default {
  middleware: 'user',
  components: {
    MyNav
  }
}
</script>
