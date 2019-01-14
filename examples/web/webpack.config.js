var path = require('path');
module.exports = {
  mode: 'development',
  entry: "./src/main.js",
  devtool: 'inline-source-map',
  devServer: {
    contentBase: path.join(__dirname, 'dist'),
    compress: true,
    port: 8080
  },
};