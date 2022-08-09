let path = require('path')
let CopyPlugin = require('copy-webpack-plugin')

module.exports = {
  lintOnSave: false,
  configureWebpack: (config) => {


    config.plugins.push(new CopyPlugin([
      {
        from: path.resolve(__dirname, './config.json'),
        to: path.resolve(__dirname, './dist/config.json')
      }
    ]))
  }
};
