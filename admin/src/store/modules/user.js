import { getToken, setToken, removeToken } from "@/utils/auth";
import router, { resetRouter } from "@/router";

const state = {
  token: getToken(),
  name: "",
  avatar: "",
  introduction: "",
  roles: [],
};

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token;
  },
  SET_INTRODUCTION: (state, introduction) => {
    state.introduction = introduction;
  },
  SET_USER_ID: (state, id) => {
    state.id = id;
  },
  SET_NAME: (state, name) => {
    state.name = name;
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar;
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles;
  },
};

const actions = {
  // user login
  login({ commit }, loginForm) {
    const { username, password, captchaId, captchaCode } = loginForm;
    return new Promise((resolve, reject) => {
      this._vm.axios
        .form("/api/login/signin", {
          username,
          password,
          captchaId,
          captchaCode,
        })
        .then((data) => {
          commit("SET_TOKEN", data.token);
          setToken(data.token);
          resolve();
        })
        .catch((error) => {
          reject(error);
        });
    });
  },

  // get user info
  async getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      this._vm.axios
        .get("/api/user/current")
        .then((data) => {
          if (!data || !data.id) {
            reject(new Error("Verification failed, please Login again."));
          }
          // roles must be a non-empty array
          if (!data.roles || data.roles.length <= 0) {
            reject(new Error("getInfo: roles must be a non-null array!"));
          }
          commit("SET_USER_ID", data.id);
          commit("SET_ROLES", data.roles);
          commit("SET_NAME", data.nickname);
          commit("SET_AVATAR", data.avatar);
          commit("SET_INTRODUCTION", data.description);
          resolve(data);
        })
        .catch((error) => {
          reject(error);
        });
    });
  },

  // user logout
  async logout({ commit, state, dispatch }) {
    await this._vm.axios.get("/api/login/signout");
    commit("SET_TOKEN", "");
    commit("SET_ROLES", []);
    removeToken();
    resetRouter();
    // reset visited views and cached views
    // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
    dispatch("tagsView/delAllViews", null, { root: true });
  },

  // remove token
  resetToken({ commit }) {
    return new Promise((resolve) => {
      commit("SET_TOKEN", "");
      commit("SET_ROLES", []);
      removeToken();
      resolve();
    });
  },

  // dynamically modify permissions
  changeRoles({ commit, dispatch }, role) {
    // eslint-disable-next-line no-async-promise-executor
    return new Promise(async (resolve) => {
      const token = `${role}-token`;

      commit("SET_TOKEN", token);
      setToken(token);

      const { roles } = await dispatch("getInfo");

      resetRouter();

      // generate accessible routes map based on roles
      const accessRoutes = await dispatch("permission/generateRoutes", roles, { root: true });

      // dynamically add accessible routes
      router.addRoutes(accessRoutes);

      // reset visited views and cached views
      dispatch("tagsView/delAllViews", null, { root: true });

      resolve();
    });
  },
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
};
