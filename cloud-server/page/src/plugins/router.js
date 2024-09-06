import {
    createRouter,
    createWebHistory
} from 'vue-router'

const Home = () => import('../views/home.vue')
const Auth = () => import('../views/auth.vue')
const List = () => import('../views/list.vue')
const OAuth = () => import('../views/oauth2.vue')

const router = createRouter({
    history: createWebHistory('/app'),
    routes: [
        {
            name: 'Home',
            path: '/',
            component: Home,
            meta: {
                title: '首页'
            }
        },
        {
            name: 'Auth',
            path: '/auth',
            component: Auth,
            meta: {
                title: '登录'
            }
        }, 
        {
            name: 'OAuth',
            path: '/oauth2',
            component: OAuth,
            meta: {
                title: '授权登录'
            }
        }, 
        {
            name: 'List',
            path: '/list',
            component: List,
            meta: {
                title: '设备列表'
            }
        }
    ],
})

router.beforeEach((to, _, next) => {
    document.title = (to.meta.title === undefined ? '未知页面 - ' : to.meta.title + ' - ') + 'Obsidian NAS Sync'
    next()
})

export default router