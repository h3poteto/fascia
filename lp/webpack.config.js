const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')

const production = process.env.NODE_ENV === 'production'

module.exports = {
  entry: {
    'css/lp': path.join(__dirname, './css/lp.scss'),
    'css/lp-webview': path.join(__dirname, './css/lp-webview.scss')
  },
  output: {
    path: path.resolve(__dirname, '../public/lp'),
    filename: '[name].js'
  },
  cache: true,
  resolve: {
    modules: [path.resolve(__dirname, './node_modules')],
    extensions: ['*', '.css', '.scss']
  },
  module: {
    rules: [
      {
        test: /\.(scss|css)$/,
        exclude: /node_modules/,
        use: [
          {
            loader: MiniCssExtractPlugin.loader,
            options: {
              publicPath: path.resolve(__dirname, '../public/lp')
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
    new CopyWebpackPlugin([{ from: path.resolve(__dirname, './images'), to: path.resolve(__dirname, '../public/lp/images') }])
  ]
}
