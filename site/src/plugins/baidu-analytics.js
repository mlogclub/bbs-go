export default defineNuxtPlugin((nuxtApp) => {
  if (process.client) {
    const baiduAnalyticsID = "79b8ff82974d0769ef5c629e4cd46629";
    console.log(import.meta.env);
    window._hmt = window._hmt || [];
    (function () {
      var hm = document.createElement("script");
      hm.src = `https://hm.baidu.com/hm.js?${baiduAnalyticsID}`;
      var s = document.getElementsByTagName("script")[0];
      s.parentNode.insertBefore(hm, s);
    })();
  }
});
