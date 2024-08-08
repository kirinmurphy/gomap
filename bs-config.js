module.exports = {
  proxy: 'localhost:8080',
  files: [
    'templates/**/*.{html,js,css}'
  ],
  open: false,
  notify: false,
  watchOptions: {
    ignoreInitial: true,
    ignored: '*.txt'
  },
  reloadDebounce: 200
};
