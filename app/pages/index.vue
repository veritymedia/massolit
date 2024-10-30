<script setup lang="ts">
import type { Label } from "radix-vue";

const pb = usePocketbase();

const user = pb.authStore.model;

const userModel = ref({
  email: "",
  pass: "",
});

async function login() {
  try {
    let user = await pb
      .collection("users")
      .authWithPassword(userModel.value.email, userModel.value.pass);
    console.log(user);
    console.log(pb.authStore.model);
    await navigateTo("/app");
  } catch (err) {
    console.log(err);
  }
}
</script>
<template>
  <div class="h-screen w-screen flex flex-col items-center justify-between p-5">
    <div class="flex items-center justify-between w-full">
      <img src="/images/logos/massolit-logo.png" class="w-24" alt="" />
    </div>

    <Card class="">
      <CardHeader>Login to Massolit</CardHeader>
      <CardContent class="flex flex-col gap-4">
        <form>
          <label for="email">Email</label>
          <Input
            id="email"
            :modelValue="userModel.email"
            @update:modelValue="(v: string) => (userModel.email = v)"
            placeholder="someone@email.com"
            type="email"
          />
        </form>
        <form>
          <label for="password">Password</label>
          <Input
            id="password"
            :modelValue="userModel.pass"
            @update:modelValue="(v: string) => (userModel.pass = v)"
            type="password"
          />
        </form>
        <Button @click="login">Login</Button>
      </CardContent>
    </Card>
    <div></div>
  </div>
</template>
