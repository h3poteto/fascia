const path = require('path')
const { WebpackManifestPlugin } = require('webpack-manifest-plugin')

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
        use: ['babel-loader', 'ts-loader']
      },
      {
        test: /\.scss$/,
        exclude: /node_modules/,
        use: [
          'style-loader',
          {
            loader: 'css-loader',
            options: {
              modules: true
            }
          },
          'sass-loader'
        ]
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      },
      {
        test: /\.(woff|woff2)(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url-loader?limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url-loader?limit=10000&mimetype=application/octet-stream'
      },
      {
        test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
        use: 'file-loader'
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
        use: 'url-loader?limit=10000&mimetype=image/svg+xml'
      }
    ]
  },
  plugins: [new WebpackManifestPlugin({ publicPath: '' })]
}
