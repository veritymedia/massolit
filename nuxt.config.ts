// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false, // Must be false to be used served by pocketbase
  srcDir: "app",
  css: [],
  build: {
    transpile: ["primevue"],
  },
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: "",
    /**
     * Directory that the component lives in.
     * @default "./components/ui"
     */
    componentDir: "./components/ui",
  },
  modules: ["@nuxtjs/tailwindcss", "shadcn-nuxt"],
  imports: {
    autoImport: true,
    addons: {
      vueTemplate: true,
    },
  },
});
