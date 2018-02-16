const path = require('path')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')

// eslint-disable-next-line no-undef
const filename = process.env.NODE_ENV === 'production' ? '[name]-[hash]' : '[name]'

module.exports = {
  entry: {
    'javascripts/bundle': './frontend/javascripts/bundle.js',
    'stylesheets/application': './frontend/javascripts/application.js',
    'stylesheets/application-webview': './frontend/javascripts/application-webview.js',
  },
  output: {
    path: path.resolve(__dirname, './public/assets'),
    filename: `${filename}.js`,
  },
  cache: true,
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
    new ManifestPlugin(),
    new ExtractTextPlugin(`${filename}.css`),
    new CopyWebpackPlugin([{ from: './frontend/images', to: './images' }])
  ]
}
