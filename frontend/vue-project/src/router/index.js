import { createRouter, createWebHistory } from "vue-router";
import Login from "@/views/Login.vue";
//import Register from "@/views/Register.vue";
import Home from "@/views/Home.vue";
import ServerDetail from "@/views/ServerDetail.vue";
//import Profile from "@/views/Profile.vue";


const routes = [
  // 登录页（设为默认路由）
  {
    path: "/",
    name: "Login",
    component: Login,
  },

  // // 注册页
  // {
  //   path: "/register",
  //   name: "Register",
  //   component: Register,
  // },

  // 主界面
  {
    path: "/home",
    name: "Home",
    component: Home,
  },

  // 服务器详情页（动态路由）
  {
    path: "/server/:id",
    name: "ServerDetail",
    component: ServerDetail,
    props: true,
    children: [
      {
        path: 'pageone', // 相对路径，无需重复父级路径
        name: 'pageone',
        component: () => import('@/views/PageOne.vue'),
        props: true
      },
      {
        path: 'pagetwo',
        name: 'pagetwo',
        component: () => import('@/views/PageTwo.vue'),
        props: true
      },
    ]
  },

  // // 个人信息页
  // {
  //   path: "/profile",
  //   name: "Profile",
  //   component: Profile,
  // },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;