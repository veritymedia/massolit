<script setup lang="ts">
import type { Label } from "radix-vue";

const pb = usePocketbase();

const user = pb.authStore.model;

const userModel = ref({
  email: "",
  pass: "",
});

async function login() {
  let user = await pb
    .collection("users")
    .authWithPassword(userModel.value.email, userModel.value.pass);
  console.log(user);
  console.log(pb.authStore.model);
}
</script>
<template>
  <div class="h-screen w-screen flex items-center justify-center p-10">
    <Card class="">
      <CardHeader>Access Library</CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div>
          <Label for="email">Email</Label>
          <Input
            id="email"
            :modelValue="userModel.email"
            @update:modelValue="(v: string) => (userModel.email = v)"
            placeholder="someone@email.com"
            type="email"
          />
        </div>
        <div>
          <Label for="password">Password</Label>
          <Input
            id="password"
            :modelValue="userModel.pass"
            @update:modelValue="(v: string) => (userModel.pass = v)"
            type="password"
          />
        </div>
        <Button @click="login">Login</Button>
      </CardContent>
    </Card>
  </div>
</template>
