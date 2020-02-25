const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')

// eslint-disable-next-line no-undef
const production = process.env.NODE_ENV === 'production'
// const filename = production ? '[name]-[hash]' : '[name]'

module.exports = {
  entry: {
    'js/main': path.join(__dirname, './js/main.tsx')
    //'css/application': path.join(__dirname, './css/application.scss'),
    //'css/application-webview': path.join(__dirname, './css/application-webview.scss')
  },
  output: {
    path: path.resolve(__dirname, '../public/assets'),
    filename: '[name]-[hash].js'
  },
  cache: true,
  watchOptions: {
    aggregateTimeout: 300,
    poll: 1000
  },
  resolve: {
    alias: {
      '@': path.join(__dirname, './js')
    },
    extensions: ['*', '.css', '.scss', '.ts', '.tsx', '.js', '.jsx']
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        use: 'babel-loader',
        exclude: /node_modules/
      },
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: 'babel-loader'
          },
          {
            loader: 'ts-loader'
          }
        ]
      },
      {
        test: /\.(scss|css)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
            options: {
              publicPath: path.resolve(__dirname, '../public/assets')
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
    new CopyWebpackPlugin([{ from: path.resolve(__dirname, './images'), to: path.resolve(__dirname, '../public/assets/images') }])
  ]
}
