<template>
    <div>
        <!-- <myHeader></myHeader> -->

        <el-container style="height: 100vh; 
      border: 1px solid #eee; 
      margin: 0;
      border: 0;
      padding: 0;"> <!-- style： 隐藏滚动条  去掉白边-->
            <!-- 侧边栏 -->
            <el-aside style="width: 200px;">
                <asideNav></asideNav>
            </el-aside>
            <!-- 主页面 -->
            <el-main>
                <router-view></router-view>
            </el-main>
        </el-container>
    </div>
</template>

<script>
import asideNav from "@/views/asideNav.vue"
export default {
    name: 'ServerDetail',
    components: { asideNav }  // 注册组件
}

</script>

<style scoped></style>
