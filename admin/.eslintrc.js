module.exports = {
  root: true,

  env: {
    node: true,
  },

  extends: [
    // 'plugin:vue/essential',
    // '@vue/airbnb',
    'plugin:vue/recommended',
    '@vue/airbnb',
  ],

  parserOptions: {
    parser: 'babel-eslint',
  },

  rules: {
    'no-console': 'off',
    'no-debugger': 'off',
  },
};
