const path = require('path')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')

module.exports = {
  entry: {
    'javascripts/bundle.js': './frontend/javascripts/bundle.js',
    'stylesheets/application.css': './frontend/stylesheets/application.scss',
    'stylesheets/application-webview.css': './frontend/stylesheets/application-webview.scss',
  },
  output: {
    path: path.resolve(__dirname, './public/assets'),
    filename: '[name]'
  },

  resolve: {
    modules: [
      path.resolve(__dirname, './fronted/javascripts'),
      path.resolve(__dirname, './node_modules'),
    ],
    extensions: ['*', '.css', '.scss', '.js', '.jsx']
  },
  module: {
    loaders: [
      {
        test: /\.js[x]?$/,
        exclude: /node_modules/,
        loader: 'babel-loader',
        query: {
          presets: ['es2015', 'react'],
        },
      },
      {
        test: /\.(scss|css)$/,
        loader: ExtractTextPlugin.extract({ fallback: 'style-loader', use: 'css-loader!sass-loader' }),
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=application/font-woff',
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=application/octet-stream',
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file-loader',
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=image/svg+xml',
},
      {
        test: /\.json/,
        loader: 'json-loader',
      }
    ]
  },
  plugins: [
    new ExtractTextPlugin('[name]'),
    new CopyWebpackPlugin([{ from: './frontend/images', to: './images' }])
  ]
}
