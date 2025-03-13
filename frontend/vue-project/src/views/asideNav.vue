<template>
    <div class="aside" style="height: 100%; width: 100%;">
        <el-row class="tac" style="height: 100%; width: 100%;">
            <el-col :span="12" style="height: 100%; width: 100%;">
                <el-menu default-active="this.$route.path" router class="el-menu-vertical-demo" @open="handleOpen"
                    @close="handleClose" background-color="#545c64" text-color="#fff" active-text-color="#ffd04b"
                    style="height: 100%; width: 100%;">
                    <el-submenu index="/PageOne">
                        <template slot="title">
                            <i class="el-icon-location"></i>
                            <span>导航一</span>
                        </template>
                        <el-menu-item-group>
                            <el-menu-item index="/server/${$route.params.id}/pageone">选项1</el-menu-item>
                            <el-menu-item index="/server/${$route.params.id}/pagetwo">选项2</el-menu-item>
                            <el-menu-item index="/">选项3</el-menu-item>
                        </el-menu-item-group>
                    </el-submenu>

                </el-menu>
            </el-col>
        </el-row>
    </div>
</template>

<script>

export default {
    methods: {
        handleOpen(key, keyPath) {
            console.log(key, keyPath);
        },
        handleClose(key, keyPath) {
            console.log(key, keyPath);
        }
    }
}
</script>
