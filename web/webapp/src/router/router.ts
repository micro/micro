import Vue from "vue";
import Router from "vue-router";

let NProgress = require("nprogress");

import "nprogress/nprogress.css";

Vue.use(Router);
const router = new Router({
    base: "/",
    mode: "hash",
    linkActiveClass: "active",
    routes: [
        {
            path: "*",
            meta: {
                public: true
            },
            redirect: {
                path: "/404"
            }
        },
        {
            path: "/404",
            meta: {
                public: true
            },
            name: "NotFound",
            component: () => import(`@/views/NotFound.vue`)
        },
        {
            path: "/403",
            meta: {
                public: true
            },
            name: "AccessDenied",
            component: () => import(  `@/views/Deny.vue`)
        },
        {
            path: "/500",
            meta: {
                public: true
            },
            name: "ServerError",
            component: () => import(`@/views/Error.vue`)
        },
        {
            path: "/",
            meta: {},
            name: "Root",
            redirect: {
                name: "Home"
            }
        },
        {
            path: "/home",
            meta: {breadcrumb: true},
            name: "Home",
            component: () => import(`@/views/home/Home.vue`)
        },
        {
            path: "/cli",
            meta: {breadcrumb: true},
            name: "cli",
            component: () => import(`@/views/cli/Cli.vue`)
        }
    ]
});
// router gards
router.beforeEach((to, from, next) => {
    NProgress.start();
    next();
});

router.afterEach((to, from) => {
    // ...
    NProgress.done();
});

export default router;
