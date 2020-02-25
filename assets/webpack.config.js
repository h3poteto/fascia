const path = require('path')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')

// eslint-disable-next-line no-undef
const production = process.env.NODE_ENV === 'production'

module.exports = {
  entry: {
    'js/main': path.join(__dirname, './js/main.tsx')
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
        use: ['style-loader', 'css-loader', 'sass-loader']
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
    new CopyWebpackPlugin([{ from: path.resolve(__dirname, './images'), to: path.resolve(__dirname, '../public/assets/images') }])
  ]
}
