const path = require('path')

module.exports =  {
    mode: 'production',
    entry: path.join(__dirname, 'services/app/index.js'),
    module: {
      rules: [
        {
          test: /\.(js)$/,
          exclude: /node_modules/,
          use: [
            {
              loader: "babel-loader"
            }
          ]
        },
        {
          test: /\.(css)$/,
          use: [
            'style-loader',
            'css-loader'
          ]
        },
        {
          test: /\.(jpg)$/,
          use: [
            {
              loader: 'file-loader',
              options: {
                name: '[name].[ext]',
                publicPath: '/img',
                outputPath: '/img',
              }
            }
          ]
        }
      ]
    },
    output: {
        path: path.resolve(__dirname, 'services/nginx/public/assets/'),
        filename: "js/app.js",
        publicPath: "/"
    },
    devServer: {
        contentBase: [
          path.join(__dirname, 'services/nginx/public/'),
          path.join(__dirname, 'services/nginx/public/assets')
        ],
        historyApiFallback: true,
        clientLogLevel: 'debug',
        compress: true,
        watchOptions: {
          poll: true
        },
        port: 8082
    }
}
