// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  // Must be false to be used served by pocketbase
  ssr: false,

  srcDir: "app",
  css: ["./app/assets/css/main.css", "./app/assets/css/shadcn.css"],

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
    componentDir: "./app/components/ui",
  },

  modules: ["shadcn-nuxt"],

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  imports: {
    autoImport: true,
    addons: {
      vueTemplate: true,
    },
  },

  compatibilityDate: "2024-10-15",
});
