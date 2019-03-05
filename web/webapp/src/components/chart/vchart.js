const G2 = window.G2;
import colors from "vuetify/es5/util/colors";
const chartColors = [];
Object.entries(colors).forEach(item => {
  chartColors.push(item[1].base);
});

export default {
  name: "v-chart",

  render(h) {
    const data = {
      staticClass: "v-chart",
      ref: "canvas",
      on: this.$listeners
    };
    return h("div", data);
  },

  props: {
    option: Object,
    height: Number
  },
  data: () => ({
    chartInstance: null
  }),

  methods: {
    init() {
      this.chartInstance = new G2.Chart({
        container: this.$refs.canvas,
        forceFit: true,
        height: window.innerWidth
      });
    },
    resize() {
      this.chartInstance.resize();
    },
    clean() {
      window.removeEventListener("resize", this.chartInstance.resize);
      this.chartInstance.dispose();
    }
  },
  mounted() {
    this.init();
    window.addEventListener("resize", () => {
      this.resize();
    });
  },
  beforeDestroy() {
    this.clean();
  }
};
