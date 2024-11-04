<script lang="ts" setup>
const isDark = useDark();
const toggleDark = useToggle(isDark);

onMounted(() => {
  toggleDark(true);
});
</script>

<template>
  <div class="bg-background text-foreground">
    <NuxtLayout>
      <NuxtLoadingIndicator />
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
