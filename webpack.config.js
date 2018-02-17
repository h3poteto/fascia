const path = require('path')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const ManifestPlugin = require('webpack-manifest-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

// eslint-disable-next-line no-undef
const production = process.env.NODE_ENV === 'production'
const filename = production ? '[name]-[hash]' : '[name]'
const devtool = production ? '' : '#eval-source-map'

module.exports = {
  entry: {
    'javascripts/bundle': './frontend/javascripts/bundle.js',
    'stylesheets/application':  './frontend/javascripts/application.js',
    'stylesheets/application-webview': './frontend/javascripts/application-webview.js',
  },
  output: {
    path: path.resolve(__dirname, './public/assets'),
    filename: `${filename}.js`,
  },
  cache: true,
  devtool: devtool,
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
        loader: [
          'cache-loader',
          'babel-loader',
        ],
      },
      {
        test: /\.(scss|css)$/,
        exclude: /node_modules/,
        loader: ExtractTextPlugin.extract({ fallback: 'style-loader', use: 'cache-loader!css-loader!sass-loader' }),
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
    new CopyWebpackPlugin([{ from: './frontend/images', to: './images' }]),
    ...(
      production ? [
        new UglifyJsPlugin()
      ] : []
    ),
  ]
}
