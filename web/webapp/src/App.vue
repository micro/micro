<template>
    <div id="appRoot">
        <template>
            <v-app id="inspire" class="app">
                <app-drawer class="app--drawer" :parent="this"></app-drawer>
                <app-toolbar class="app--toolbar" :parent="this"></app-toolbar>
                <v-content style="height: 100%;">
                    <div class="page-wrapper">
                        <page-header v-if="$route.meta.breadcrumb"></page-header>
                        <keep-alive>
                            <router-view></router-view>
                        </keep-alive>
                    </div>
                </v-content>
                <!-- Go to top -->
                <app-fab></app-fab>
                <!-- theme setting
                <v-btn small fab dark falt fixed top="top" right="right" class="setting-fab" color="red"
                       @click="openThemeSettings">
                    <v-icon>settings</v-icon>
                </v-btn>
                <v-navigation-drawer
                        class="setting-drawer"
                        temporary
                        right
                        v-model="rightDrawer"
                        hide-overlay
                        fixed
                >
                    <theme-settings></theme-settings>
                </v-navigation-drawer> -->
            </v-app>
        </template>
    </div>
</template>

<script lang="ts">
    import AppDrawer from "@/views/home/AppDrawer.vue";
    import AppToolbar from "@/components/AppToolbar.vue";
    import AppFab from "@/components/AppFab.vue";
    import PageHeader from "@/components/PageHeader.vue";
    // import ThemeSettings from "@/components/ThemeSettings.vue";

    import {Vue, Component, Watch} from "vue-property-decorator";

    @Component({
        components: {
            AppDrawer,
            AppToolbar,
            AppFab,
            PageHeader,
            // ThemeSettings
        }
    })
    export default class Home extends Vue {
        private expanded = true;
        private rightDrawer = false;
        private snackbar = {
            show: false,
            text: "",
            color: ""
        };

        computed() {
        }

        created() {
        }

        openThemeSettings() {
            this.$vuetify.goTo(0);
            this.rightDrawer = !this.rightDrawer;
        }
    }
</script>


<style scoped>
    .page-wrapper {
        padding-left: 20px;
    }
</style>
