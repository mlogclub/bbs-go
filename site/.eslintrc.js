module.exports = {
  root: true,
  env: {
    browser: true,
    node: true
  },
  parserOptions: {
    parser: 'babel-eslint'
  },
  extends: ["@nuxtjs", "plugin:nuxt/recommended"],
  // 'extends': [
  //   'eslint:recommended',
  //   'plugin:vue/recommended'
  // ],
  // 'extends': [
  //   'plugin:vue/essential',
  //   'eslint:recommended'
  // ],
  // add your custom rules here
  rules: {
    'no-unused-vars': [
      'warn',
      {
        'vars': 'all',
        'args': 'none'
      }
    ],
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'vue/no-v-html': 'off'
  }
}
