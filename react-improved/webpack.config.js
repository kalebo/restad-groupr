const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin');

const HtmlWebpackPluginConfig = new HtmlWebpackPlugin({
  template: './client/index.html',
  filename: 'index.html',
  inject: 'body'
})

module.exports = {
    entry: "./client/index.js",
    output: {
        path: path.resolve("./dist"),
        filename: "bundle.js"
    },
    module: {
        loaders: [
            { test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"},
            { test: /\.jsx$/, exclude: /node_modules/, loader: "babel-loader"}
        ]
    },

    plugins: [HtmlWebpackPluginConfig]
}
