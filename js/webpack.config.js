const path = require('path');

module.exports = env => {
    env.NODE_ENV = 'dev'
    console.log('NODE_ENV: ', env.NODE_ENV); // 'local'
    console.log('Production: ', env.production); // true
    return {
        mode: 'development',
        watch: true,
        entry: {
            wall: ['./wall.js'],
            login: ['./login.js'],
            join: ['./join.js'],
            'reset-password': ['./reset-password.js'],
            'forgot-password': ['./forgot-password.js'],
        },
        output: {
            path: path.resolve(__dirname + "/../static/js"),
            filename: '[name].bundle.js'
        },
        devtool: "source-map",
        module: {
            rules: [
                {
                    test: /\.m?js$/,
                    exclude: /(node_modules|bower_components)/,
                }
            ]
        }
    }
};
