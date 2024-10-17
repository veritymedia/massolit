const pb = usePocketbase();

export default defineNuxtRouteMiddleware((to, from) => {
  if (!pb.authStore.isValid) {
    return navigateTo("/");
  }
  return;
});
