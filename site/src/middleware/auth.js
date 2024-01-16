export default defineNuxtRouteMiddleware((to) => {
  const userStore = useUserStore();
  if (!userStore.user) {
    useMsgSignIn();
    return;
  }
});
