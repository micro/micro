<template>
    <v-navigation-drawer id="appDrawer" :mini-variant.sync="mini" fixed :dark="$vuetify.dark" app v-model="drawer"
                         width="260">
        <v-toolbar color="primary darken-1" dark>
            <v-toolbar-title class="ml-0 pl-3">
                <span class="hidden-sm-and-down">Micro</span>
            </v-toolbar-title>
        </v-toolbar>
        <vue-perfect-scrollbar class="drawer-menu--scroll" :settings="scrollSettings">
            <v-list dense expand>
                <template v-for="(item, i) in menus">
                    <!--group with subitems-->
                    <v-list-group v-if="item.items" :key="item.name" :group="item.group" :prepend-icon="item.icon"
                                  no-action="no-action">
                        <v-list-tile slot="activator" ripple="ripple">
                            <v-list-tile-content>
                                <v-list-tile-title> {{ $t("menu."+item.title) }}</v-list-tile-title>
                            </v-list-tile-content>
                        </v-list-tile>
                        <template v-for="(subItem, i) in item.items">
                            <!--sub group-->
                            <v-list-group v-if="subItem.items" :key="subItem.name" :group="subItem.group"
                                          sub-group="sub-group">
                                <v-list-tile slot="activator" ripple="ripple">
                                    <v-list-tile-content>
                                        <v-list-tile-title>{{ $t("menu."+subItem.title) }}</v-list-tile-title>
                                    </v-list-tile-content>
                                </v-list-tile>
                                <v-list-tile v-for="(grand, i) in subItem.children" :key="i"
                                             :to="genChildTarget(item, grand)" :href="grand.href" ripple="ripple">
                                    <v-list-tile-content>
                                        <v-list-tile-title>{{ $t("menu."+grand.title) }}</v-list-tile-title>
                                    </v-list-tile-content>
                                </v-list-tile>
                            </v-list-group>
                            <!--child item-->
                            <v-list-tile v-else :key="i" :to="genChildTarget(item, subItem)" :href="subItem.href"
                                         :disabled="subItem.disabled" :target="subItem.target" ripple="ripple">
                                <v-list-tile-content>
                                    <v-list-tile-title><span>{{ $t("menu."+subItem.title) }} </span>
                                    </v-list-tile-title>
                                </v-list-tile-content>
                                <v-list-tile-action v-if="subItem.action">
                                    <v-icon :class="[subItem.actionClass || 'success--text']">{{ subItem.action }}
                                    </v-icon>
                                </v-list-tile-action>
                            </v-list-tile>
                        </template>
                    </v-list-group>
                    <v-subheader v-else-if="item.header" :key="i">{{ item.header }}</v-subheader>
                    <v-divider v-else-if="item.divider" :key="i"></v-divider>
                    <!--top-level link-->
                    <v-list-tile v-else :to="!item.href ? { name: item.name } : null" :href="item.href" ripple="ripple"
                                 :disabled="item.disabled" :target="item.target" rel="noopener" :key="item.name">
                        <v-list-tile-action v-if="item.icon">
                            <v-icon>{{ item.icon }}</v-icon>
                        </v-list-tile-action>
                        <v-list-tile-content>
                            <v-list-tile-title>{{ $t("menu."+item.title) }}</v-list-tile-title>
                        </v-list-tile-content>
                        <v-list-tile-action v-if="item.subAction">
                            <v-icon class="success--text">{{ item.subAction }}</v-icon>
                        </v-list-tile-action>
                    </v-list-tile>
                </template>
            </v-list>
        </vue-perfect-scrollbar>
    </v-navigation-drawer>
</template>
<script lang="ts">
    import {Component, Prop, Vue} from "vue-property-decorator";
    import VuePerfectScrollbar from "vue-perfect-scrollbar";
    import menu from "@/api/menu";

    @Component({
        components: {
            VuePerfectScrollbar
        }
    })
    export default class AppDrawer extends Vue {
        @Prop() private parent: any;

        private mini = false;
        private drawer = true;
        private menus = menu;
        private scrollSettings = {
            maxScrollbarLength: 160
        };

        created() {
            this.parent.$on("APP_DRAWER_TOGGLED", () => {
                this.drawer = !this.drawer;
            });
        }

        genChildTarget(item: any, subItem: any) {
            if (subItem.href) return;
            if (subItem.component) {
                return {
                    name: subItem.component
                };
            }
            return {name: `${item.group}/${subItem.name}`};
        }
    }
</script>


<style lang="stylus">
    #appDrawer {
        overflow: hidden;

        .drawer-menu--scroll {
            height: calc(100vh - 48px);
            overflow: auto;
        }
    }
</style>
