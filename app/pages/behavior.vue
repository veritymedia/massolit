<script setup lang="ts">
definePageMeta({
  middleware: ["not-authed-guard"],
});

const pb = usePocketbase();

const isLoading = ref(false);
</script>
<template>
  <div class="flex flex-col justify-start">
    <div
      v-if="isLoading"
      class="fixed top-0 left-0 flex items-center justify-center w-screen h-screen bg-background"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    <AppAlert />
    <div>
      <div>
        <h2 class="mt-10 mb-5">Detention Tracker</h2>
        <BehaviorDetentionTable />
      </div>
    </div>
  </div>
</template>
