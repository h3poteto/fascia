const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')

// eslint-disable-next-line no-undef
const production = process.env.NODE_ENV === 'production'
const filename = production ? '[name]-[hash]' : '[name]'
const devtool = production ? '' : '#eval-source-map'
const mode = production ? 'production' : 'development'

module.exports = {
  entry: {
    'javascripts/bundle': './frontend/javascripts/bundle.js',
    'stylesheets/application': './frontend/stylesheets/application.scss',
    'stylesheets/application-webview': './frontend/stylesheets/application-webview.scss'
  },
  output: {
    path: path.resolve(__dirname, './public/assets'),
    filename: '[name].js'
  },
  cache: true,
  devtool: devtool,
  mode: mode,
  watchOptions: {
    aggregateTimeout: 300,
    poll: 1000
  },
  resolve: {
    modules: [path.resolve(__dirname, './fronted/javascripts'), path.resolve(__dirname, './node_modules')],
    extensions: ['*', '.css', '.scss', '.js', '.jsx']
  },
  module: {
    rules: [
      {
        test: /\.js[x]?$/,
        exclude: /node_modules/,
        use: ['cache-loader', 'babel-loader']
      },
      {
        test: /\.(scss|css)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
            options: {
              publicPath: path.resolve(__dirname, './public/assets')
            }
          },
          'css-loader',
          'sass-loader'
        ]
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=application/octet-stream'
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'file-loader'
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader?limit=10000&mimetype=image/svg+xml'
      }
    ]
  },
  plugins: [
    new ManifestPlugin(),
    new MiniCssExtractPlugin({
      // Options similar to the same options in webpackOptions.output
      // both options are optional
      filename: production ? '[name].[hash].css' : '[name].css',
      chunkFilename: production ? '[id].[hash].css' : '[id].css'
    }),
    new CopyWebpackPlugin([{ from: './frontend/images', to: './images' }])
  ]
}
