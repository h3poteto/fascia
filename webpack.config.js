var ExtractTextPlugin = require("extract-text-webpack-plugin");
var CopyWebpackPlugin = require("copy-webpack-plugin");
var subscript = require('markdown-it-sub');
var superscript = require('markdown-it-sup');

module.exports = {
  devtool: "source-map",
  entry: {
    'javascripts/bundle.js': "./frontend/javascripts/bundle.js",
    'stylesheets/application.css': "./frontend/stylesheets/application.scss",
    'stylesheets/application-webview.css': "./frontend/stylesheets/application-webview.scss",
  },
  output: {
    path: "./public/assets",
    filename: "[name]"
  },

  resolve: {
    modulesDirectories: [
      __dirname + "/fronted/javascripts",
      __dirname + "/node_modules"
    ],
    extensions: ["", ".js", ".jsx"]
  },

  module: {
    loaders: [
      {
        test: /\.js[x]?$/,
        exclude: /node_modules/,
        loader: ["babel"],
        query: {}
      },
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract("style", "css")
      },
      {
        test: /\.scss$/,
        loader: ExtractTextPlugin.extract("style", "css!sass")
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: "url?limit=10000&mimetype=application/font-woff"
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: "url?limit=10000&mimetype=application/octet-stream"
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: "file"
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: "url?limit=10000&mimetype=image/svg+xml"
      },
      {
        test: /\.json/,
        loader: 'json'
      }
    ]
  },
  plugins: [
    new ExtractTextPlugin("[name]"),
    new CopyWebpackPlugin([{ from: "./frontend/images", to: "./images" }])
  ]
}
