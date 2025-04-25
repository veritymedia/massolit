<script lang="ts" setup>
const isLoading = ref(false);

const pb = usePocketbase();
</script>

<template>
  <div class="">
    <div
      v-if="isLoading"
      class="fixed top-0 left-0 flex items-center justify-center w-screen h-screen bg-background"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    <AppAlert />
    <div>
      <h2 class="mt-10 mb-5">Detention Tracker</h2>
    </div>
  </div>
</template>

<style></style>
