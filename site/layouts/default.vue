<template>
  <div>
    <my-nav />

    <nuxt />

    <my-footer />
  </div>
</template>

<script>
import MyNav from '~/components/MyNav'
import MyFooter from '~/components/MyFooter'

export default {
  middleware: 'global',
  components: {
    MyNav,
    MyFooter
  }
}
</script>
<script>
(function(){
    var bp = document.createElement('script');
    var curProtocol = window.location.protocol.split(':')[0];
    if (curProtocol === 'https') {
        bp.src = 'https://zz.bdstatic.com/linksubmit/push.js';
    }
    else {
        bp.src = 'http://push.zhanzhang.baidu.com/push.js';
    }
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(bp, s);
})();
</script>
