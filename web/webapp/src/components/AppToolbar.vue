<template>
    <v-toolbar color="primary" fixed dark app>
        <v-toolbar-title class="ml-0 pl-3">
            <v-toolbar-side-icon
                    @click.stop="handleDrawerToggle"
            ></v-toolbar-side-icon>
        </v-toolbar-title>
        <!-- <v-text-field
               flat
               solo-inverted
               prepend-icon="search"
               label="Search"
               class="hidden-sm-and-down"
               >
             </v-text-field>-->
        <v-spacer></v-spacer>
        <v-toolbar-items class="hidden-sm-and-down">
            <v-btn icon @click="handleFullScreen()">
                <v-icon>fullscreen</v-icon>
            </v-btn>
            <v-menu
                    bottom
                    origin="center center"
                    offset-y
                    transition="scale-transition"
                    :position-y="23"
            >
                <template v-slot:activator="{ on }">
                    <v-btn icon v-on="on">
                        <img :src="currentLanFlag">
                    </v-btn>
                </template>

                <v-list>
                    <v-list-tile
                            v-for="(item, i) in lanItems"
                            :key="i"
                            @click="changeLanguage(item.lan)"
                    >
                        <v-list-tile-avatar :size="25" :tile="true">
                            <img :src="item.flag">
                        </v-list-tile-avatar>
                        <v-list-tile-title>{{ item.title }}</v-list-tile-title>
                    </v-list-tile>

                </v-list>
            </v-menu>
        </v-toolbar-items>
    </v-toolbar>
</template>
<script lang="ts">
    import Util from "@/utils";
    import {Vue, Component, Watch, Prop} from "vue-property-decorator";

    @Component({
        components: {}
    })
    export default class AppToolbar extends Vue {

        @Prop() private parent: any;

        private currentLanFlag: string = '';

        private lanItems = {
            en: {flag: 'https://cdn.vuetifyjs.com/images/flags/us.png', title: 'English', lan: 'en'},
            cn: {flag: 'https://cdn.vuetifyjs.com/images/flags/cn.png', title: '简体中文', lan: 'cn'}
        }

        mounted() {
            this.setDefaultLanguage(Util.getCookieValue('locale'))
        }

        setDefaultLanguage(lan?: string) {
            if (lan && this.lanItems[lan]) {
                this.currentLanFlag = this.lanItems[lan].flag
            } else {
                this.currentLanFlag = this.lanItems.en.flag
            }
        }

        changeLanguage(lan?: string) {
            this.$root.$emit('localeChange', lan)
            this.setDefaultLanguage(lan)
        }


        handleDrawerToggle() {
            this.parent.$emit("APP_DRAWER_TOGGLED");
        }

        handleFullScreen() {
            Util.toggleFullScreen();
        }


    }
</script>
