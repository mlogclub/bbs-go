module.exports = {
  root: true,

  env: {
    node: true
  },

  extends: [
    // 'plugin:vue/essential',
    // 'eslint:recommended',
    'plugin:vue/recommended',
    '@vue/standard'
  ],

  parserOptions: {
    parser: 'babel-eslint'
  },

  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'import/extensions': ['error', 'always', {
      js: 'never',
      vue: 'never'
    }],
    'no-plusplus': 'off',
    'no-continue': 'off',
    'vue/no-v-html': 'off'
  }
}
