<script setup lang="ts">
const pocketbase = usePocketbase();

const user = pocketbase.authStore.model;

function logout() {
  pocketbase.authStore.clear();
  navigateTo("/auth/login");
}
</script>
<template>
  <div class="h-screen w-screen flex items-center justify-center p-10">
    <div class="outline rounded p-5 outline-muted">
      <h2>Login</h2>
    </div>
  </div>
</template>
