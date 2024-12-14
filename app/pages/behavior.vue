<script setup lang="ts">
definePageMeta({
  middleware: ["not-authed-guard"],
});

const isLoading = ref(false);
</script>
<template>
  <div
    class="h-screen w-screen background-gradient flex-col justify-between p-6 flex"
  >
    <div
      v-if="isLoading"
      class="w-screen h-screen fixed top-0 left-0 bg-background flex items-center justify-center"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    <AppAlert />
    <div>
      <div class="flex items-center justify-between">
        <NuxtLink to="/app">
          <img
            src="../public/images/logos/massolit-logo.png"
            class="w-24"
            alt=""
          />
        </NuxtLink>
        <!-- <p class="text-sm">{{ pb.authStore.model?.email }}</p> -->
      </div>
      <div class="mt-20">
        <div v-if="$route.name === 'app'">
          <StudentBooks />
        </div>
      </div>
    </div>
  </div>
</template>
