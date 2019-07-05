const path = require('path')

module.exports =  {
    mode: 'development',
    entry: path.join(__dirname, 'services/app/index.js'),
    module: {
      rules: [
        {
          test: /\.(js)$/,
          exclude: /node_modules/,
          loaders: ["babel-loader"],
        }
      ],
    },
    output: {
        path: path.resolve(__dirname, 'services/nginx/public/assets/js'),
        filename: "app.js",
        publicPath: "/assets/js"
    },
    devServer: {
        contentBase: path.join(__dirname, 'services/nginx/public/assets/js'),
        compress: true,
        port: 8082
    }
}
