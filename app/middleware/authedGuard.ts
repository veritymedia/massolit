const pb = usePocketbase();

export default defineNuxtRouteMiddleware((to, from) => {
  console.log("MIDDLEWARE: ", pb.authStore.model, pb.authStore.isValid);
  if (pb.authStore.isValid) {
    return navigateTo("/");
  }
  return;
});
