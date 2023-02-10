import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

//https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    // 是否快起 $ref , vue3.2 的语法糖

    vue({
      refTransform: true,
      reactivityTransform: true,
    }),
  ],

  server: {
    //prot:8080

    proxy: {
      "/user": {
        target: "http://127.0.0.1:8989",
      },
    },
  },

  base: "./", //打包相对路径
});
