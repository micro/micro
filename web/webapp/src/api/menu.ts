const Menu = [
    {
        title: "homePage",
        group: "apps",
        icon: "home",
        name: "home"
    },
    {
        title: "cliTerminal",
        group: "apps",
        icon: "tune",
        name: "cli"
    },
    {
        title: "registryInfo",
        component: "apps",
        icon: "cloud",
        name: "registry"
    },
    {
        title: "callService",
        component: "apps",
        icon: "train",
        name: "call"
    },
    {
        title: "stats",
        component: "stats",
        icon: "bar_chart",
        name: "statistics",
        items: [
            {name: 'apiStatistics', title: 'statsAPI', component: 'apiStatistics'},
            {name: 'serviceStatistics', title: 'statsService', component: 'serviceStatistics'},
        ]
    },
    {divider: true}
];
// reorder menu
Menu.forEach((item: any) => {
    if (item.items) {
        item.items.sort((x: any, y: any) => {
            let textA = x.title.toUpperCase();
            let textB = y.title.toUpperCase();
            return textA < textB ? -1 : textA > textB ? 1 : 0;
        });
    }
});

export default Menu;
