const pkg = require('./package')
const path = require('path');
const tools = require('./gulpfile.js/tools')
const appConfig = require('./config').default;

// 开发用plugins，如mock等
let pluginsDev = [];

// 读取iconfont配置文件
const iconfontConfig = tools.parseConfig('./.iconfont')

// build的html文件最终路径
const output = path.resolve(__dirname, '../', 'examples/wwwroot/admin');

module.exports = {
    ssr: false,
    router: {
        mode: 'history',
        // base: './'
    },
    env: {
        env: process.env.env,
    },
    /*
     ** Headers of the page
     */
    head: {
        title: "Kaligo-admin",
        meta: [{
            charset: 'utf-8'
        },
        {
            name: 'viewport',
            content: 'width=device-width, initial-scale=1'
        },
        {
            hid: 'description',
            name: 'description',
            content: pkg.description
        }
        ],
        link: [{
            rel: 'icon',
            type: 'image/x-icon',
            href: '/favicon.ico'
        },{
            rel: 'stylesheet',
            href: 'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css'
        }
        ].concat(iconfontConfig.url ? [{
            rel: 'stylesheet',
            href: process.env.NODE_ENV === 'development' ? iconfontConfig.url : `${appConfig.client.cdn}/font/iconfont.css?v=${pkg.version}`
        }] : []),
        script: []
    },
    generate: {
        dir: output,
    },
    /*
     ** Customize the progress-bar color
     */
    loading: {
        color: '#fff'
    },
    css: [
        '~/assets/less/reset.less',
        '~/assets/less/index.less',
        '~/assets/less/animate.less',
        'element-ui/lib/theme-chalk/index.css',
    ],
    plugins: [
        '@/plugins/ui',
        '@/plugins/config',
        '@/plugins/axios',
        '@/plugins/api',
        '@/plugins/i18n',
        '@/plugins/day',
        '@/plugins/utils',
        '@/plugins/directive',
        '@/plugins/errorHandler',
        { src: "@/plugins/darkMode", ssr: false },
        ...pluginsDev
    ],
    modules: [
        // Doc: https://axios.nuxtjs.org/usage
        '@nuxtjs/style-resources',
    ],
    axios: {
        // See https://github.com/nuxt-community/axios-module#options
    },
    styleResources: {
        less: [
            '~/assets/less/var.less',
        ]
    },
    build: {
        babel: {
            plugins: [
                [
                    'component',
                    {
                        'libraryName': 'element-ui',
                        'styleLibraryName': 'theme-chalk'
                    }
                ]
            ]
        },
        /*
        ** You can extend webpack config here
        * loaders 是支持template 里面直接 :src="<文件路径>"
        * extend 是打包时候文件模块的命名处理
        */
        loaders: {
            vue: {
                transformAssetUrls: {
                    audio: 'src',
                },
            },
        },
        extend(config, ctx) {
            // loader配置读取音频相关
            config.module.rules.push({
                test: /\.(ogg|mp3|wav|mpe?g)$/i,
                loader: 'file-loader',
                options: {
                    name: '[path][name].[ext]',
                },
            })
            // loader配置读取svg相关
            config.module.rules.push({
                test: /\.svg$/,
                use: [
                {
                    loader: 'svgo-loader',
                    options: {
                        plugins: [
                            {removeTitle: true},
                            {convertColors: {shorthex: false}},
                            {convertPathData: false}
                        ],
                        name: '[path][name].[ext]',
                        limit: 4 * 1024,
                    }
                }
                ]
            })
            // 外网环境，需要做的一些自动化构建
            if (process.env.NODE_ENV !== 'development') {
                // 判断是否需要拉取iconfont
                if (iconfontConfig && iconfontConfig.url) {
                    const preTasksIconfont = require('./gulpfile.js/iconfont');
                    preTasksIconfont(iconfontConfig.url);
                }
            }

        },
        publicPath: `${appConfig.client.cdn}/static/`,
    }
}