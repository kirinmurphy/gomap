module.exports = {
  proxy: 'localhost:8080',
  files: [
    'src/templates/**/*.{html,js,css}',
    '!src/templates/css/**/*' 
  ],
  open: false,
  notify: false,
  watchOptions: {
    ignoreInitial: true,
    ignored: '*.txt'
  },
  reloadDebounce: 200
};
