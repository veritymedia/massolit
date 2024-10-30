<template>
  <transition name="slide-down">
    <div
      v-if="activeAlert"
      class="w-full flex justify-center pt-4 left-0 fixed top-0"
    >
      <Alert class="w-11/12">
        <AlertTitle>{{ activeAlert.title }}</AlertTitle>
        <AlertDescription>{{ activeAlert.message }}</AlertDescription>
      </Alert>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { useAlert } from "#imports";
const { createAlert, activeAlert } = useAlert();
</script>

<style scoped>
/* Slide-down animation */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}
.slide-down-enter-from {
  transform: translateY(-100%);
  opacity: 1;
}
.slide-down-enter-to {
  transform: translateY(0);
  opacity: 1;
}
.slide-down-leave-from {
  transform: translateY(0);
  opacity: 1;
}
.slide-down-leave-to {
  transform: translateY(-100%);
  opacity: 1;
}
</style>
