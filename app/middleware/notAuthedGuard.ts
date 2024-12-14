const pb = usePocketbase();

export default defineNuxtRouteMiddleware((to, from) => {
  // console.log(
  //   "MIDDLEWARE::AUTH_CHECK",
  //   pb.authStore.model,
  //   pb.authStore.isValid,
  // );

  if (!pb.authStore.isValid) {
    if (to.path !== "/login") {
      return navigateTo("/login");
    }
  } else {
    if (to.path === "/login") {
      return navigateTo("/");
    }
  }

  return;
});
