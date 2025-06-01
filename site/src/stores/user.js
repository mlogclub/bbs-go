import { defineStore } from "pinia";

export const useUserStore = defineStore("user", {
  state: () => ({
    user: null,
  }),
  getters: {
    isLogin() {
      return !!this.user;
    },
  },
  actions: {
    async fetchCurrent() {
      const { data } = await useMyFetch("/api/user/current");
      this.user = data.value;
      return this.user;
    },
    async signin(body) {
      const { user, token, redirect } = await useHttpPost(
        "/api/login/signin",
        useJsonToForm(body)
      );
      this.user = user;
      return {
        user,
        token,
        redirect,
      };
    },
    async signout() {
      await useHttpGet("/api/login/signout");
      this.user = null;
    },
    async signup(form) {
      const { user, token, redirect } = await useHttpPost(
        "/api/login/signup",
        useJsonToForm(form)
      );
      this.user = user;
      return {
        user,
        token,
        redirect,
      };
    },
  },
});
