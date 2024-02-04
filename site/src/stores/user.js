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
      this.user = await useHttpGet("/api/user/current");
      return this.user;
    },
    async signin(body) {
      const { user, token, redirect } = await useHttpPostForm(
        "/api/login/signin",
        {
          body: body,
        }
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
      const { user, token, redirect } = await useHttpPostForm(
        "/api/login/signup",
        {
          body: form,
        }
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
