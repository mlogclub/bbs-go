<template>
  <div />
</template>

<script>
export default {
  async asyncData({ $axios, params, redirect, error }) {
    try {
      const ret = await $axios.get('/api/article/redirect/' + params.id)
      redirect(302, ret.url)
    } catch (e) {
      error({
        statusCode: 404,
        message: e.message || '404',
      })
    }
  },
}
</script>
