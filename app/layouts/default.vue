<script lang="ts" setup>
const appConfig = useAppConfig();
const pb = usePocketbase();
onMounted(() => {
  const url = window.location.host;
  if (!url.includes("localhost")) {
    pb.baseUrl = "https://" + url;
  }
});
</script>

<template>
  <div
    class="flex flex-col justify-start w-full h-full min-h-screen p-5 background-gradient"
  >
    <div class="flex items-center justify-between">
      <NuxtLink to="/">
        <img
          src="../public/images/logos/massolit-logo.png"
          class="w-24"
          alt=""
        />
      </NuxtLink>
      <p class="text-sm">{{ pb.authStore.model?.email }}</p>
    </div>
    <slot />
  </div>
</template>
