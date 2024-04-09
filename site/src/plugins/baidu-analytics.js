export default defineNuxtPlugin((nuxtApp) => {
  if (process.client) {
    const config = useRuntimeConfig();
    const baiduAnalyticsID = config.baiduAnalyticsID;
    window._hmt = window._hmt || [];
    (function () {
      var hm = document.createElement("script");
      hm.src = `https://hm.baidu.com/hm.js?${baiduAnalyticsID}`;
      var s = document.getElementsByTagName("script")[0];
      s.parentNode.insertBefore(hm, s);
    })();
  }
});
