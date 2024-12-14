<script lang="ts" setup>
const isDark = useDark();
const toggleDark = useToggle(isDark);

onMounted(() => {
  toggleDark(true);
});
</script>

<template>
  <div class="bg- background-gradient text-foreground">
    <NuxtLayout>
      <NuxtLoadingIndicator />
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<style>
.background-gradient {
  background: rgb(2, 0, 36);
  background: linear-gradient(
    327deg,
    rgba(2, 0, 36, 1) 0%,
    rgba(35, 35, 99, 1) 100%
  );
}
</style>
