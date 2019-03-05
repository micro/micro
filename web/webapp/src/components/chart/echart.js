/**
 * ECharts Vue Wrapper
 * Michael Wang
 */
import colors from "vuetify/es5/util/colors";
import _object from "lodash/object";

const ECharts = window.echarts || undefined;
if (ECharts === undefined) {
  console.error("ECharts is not defined");
}
// set color palette
const colorPalette = [];
Object.entries(colors).forEach(item => {
  if (item[1].base) {
    colorPalette.push(item[1].base);
  }
});
// default
// const colorPalette = ['#d87c7c', '#919e8b', '#d7ab82', '#6e7074', '#61a0a8', '#efa18d', '#787464', '#cc7e63', '#724e58', '#4b565b'];
// ECharts.registerTheme('material', {
//   color: colorPalette,
//   graph: {
//     color: colorPalette
//   }
//   textStyle: {

//   }
// });
(function() {
  const throttle = function(type, name, obj) {
    obj = obj || window;
    let running = false;
    let func = function() {
      if (running) {
        return;
      }
      running = true;
      requestAnimationFrame(function() {
        obj.dispatchEvent(new CustomEvent(name));
        running = false;
      });
    };
    obj.addEventListener(type, func);
  };
  /* init - you can init any event */
  throttle("resize", "optimizedResize");
})();
export default {
  name: "v-echart",

  render(h) {
    const data = {
      staticClass: "v-chart",
      style: this.canvasStyle,
      ref: "canvas",
      on: this.$listeners
    };
    return h("div", data);
  },

  props: {
    // args of  ECharts.init(dom, theme, opts)
    width: { type: String, default: "auto" },
    height: { type: String, default: "400px" },
    merged: {
      type: Boolean,
      default: true
    },
    // instace.setOption
    pathOption: [Object, Array],
    option: Object,
    // general config
    textStyle: Object,
    title: Object,
    legend: [Object, Array],
    tooltip: Object,
    grid: { type: [Object, Array] },
    xAxis: [Object, Array],
    yAxis: [Object, Array],
    series: [Object, Array],
    axisPointer: Object,
    dataset: {
      type: [Object, Array],
      default() {
        return {};
      }
    }, // option.dataSet
    colors: Array, // echarts.option.color
    backgroundColor: [Object, String],
    toolbox: { type: [Object, Array] },
    // resize delay
    widthChangeDelay: {
      type: Number,
      default: 450
    }
  },
  data: () => ({
    chartInstance: null,
    clientWidth: null,
    allowedOptions: [
      "textStyle",
      "title",
      "legend",
      "xAxis",
      "yAxis",
      "series",
      "tooltip",
      "axisPointer",
      "grid",
      "dataset",
      "colors",
      "backgroundColor"
    ],
    _defaultOption: {
      tooltip: {
        show: true
      },
      title: {
        show: true,
        textStyle: {
          color: "rgba(0, 0, 0 , .87)",
          fontFamily: "sans-serif"
        }
      },
      grid: {
        containLabel: true
      },
      xAxis: {
        show: true,
        type: "category",
        axisLine: {
          lineStyle: {
            color: "rgba(0, 0, 0 , .54)",
            type: "dashed"
          }
        },
        axisTick: {
          show: true,
          alignWithLabel: true,
          lineStyle: {
            show: true,
            color: "rgba(0, 0, 0 , .54)",
            type: "dashed"
          }
        },
        axisLabel: {
          show: false
        }
      },
      yAxis: {
        show: true,
        type: "value",
        axisLine: {
          lineStyle: {
            color: "rgba(0, 0, 0 , .54)",
            type: "dashed"
          }
        },
        axisLabel: {
          show: false
          // color: 'rgba(0, 0, 0 , .54)'
        },
        splitLine: {
          lineStyle: {
            type: "dashed"
          }
        },
        axisTick: {
          show: true,
          lineStyle: {
            show: true,
            color: "rgba(0, 0, 0 , .54)",
            type: "dashed"
          }
        }
      },
      series: [
        {
          type: "line"
        }
      ]
    }
  }),
  computed: {
    canvasStyle() {
      return {
        width: this.width,
        height: this.height
      };
    }
  },
  methods: {
    init() {
      const { widthChangeDelay } = this;
      // set
      if (this.pathOption) {
        this.pathOption.forEach(p => {
          _object.set(this.$data._defaultOption, p[0], p[1]);
        });
      }
      this.chartInstance = ECharts.init(this.$refs.canvas, "material");
      this.chartInstance.setOption(
        _object.merge(this.option, this.$data._defaultOption)
      );
      window.addEventListener("optimizedResize", e => {
        setTimeout(_ => {
          this.chartInstance.resize();
        }, this.widthChangeDelay);
      });
    },

    resize() {
      this.chartInstance.resize();
    },
    clean() {
      window.removeEventListener("resize", this.chartInstance.resize);
      this.chartInstance.clear();
    }
  },
  mounted() {
    this.init();
  },

  beforeDestroy() {
    this.clean();
  }
};
