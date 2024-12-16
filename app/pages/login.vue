<script setup lang="ts">
definePageMeta({
  middleware: ["not-authed-guard"],
});

const pb = usePocketbase();

const userModel = ref({
  email: "",
  pass: "",
});

const emailError = ref<string | null>(null);
const passwordError = ref<string | null>(null);

async function login() {
  try {
    // Reset errors
    emailError.value = null;
    passwordError.value = null;

    // Validate inputs
    if (!userModel.value.email) {
      emailError.value = "Email is required";
      return;
    }
    if (!userModel.value.pass) {
      passwordError.value = "Password is required";
      return;
    }

    // Attempt login
    let user = await pb
      .collection("users")
      .authWithPassword(userModel.value.email, userModel.value.pass);
    console.log(user);
    console.log(pb.authStore.model);
    await navigateTo("/");
  } catch (err: any) {
    console.log(err);
    if (err.response?.data?.email) {
      emailError.value = "Invalid email or password";
    } else {
      passwordError.value = "Invalid email or password";
    }
  }
}
</script>

<template>
  <div class="h-screen w-screen flex flex-col items-center justify-between p-5">
    <div class="flex items-center justify-between w-full">
      <img
        src="/images/logos/massolit-logo.png"
        class="w-24"
        alt="Massolit Logo"
      />
    </div>

    <Card>
      <CardHeader>Login to Massolit</CardHeader>
      <CardContent class="flex flex-col gap-4">
        <form @submit.prevent="login" class="flex flex-col gap-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700"
              >Email</label
            >
            <Input
              id="email"
              name="email"
              v-model="userModel.email"
              placeholder="you@massolit.app"
              type="email"
              autocomplete="email"
              aria-required="true"
              :aria-invalid="emailError ? 'true' : 'false'"
              required
              class="mt-1 block w-full"
            />
            <p
              v-if="emailError"
              id="email-error"
              class="mt-2 text-sm text-red-600"
            >
              {{ emailError }}
            </p>
          </div>

          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-700"
              >Password</label
            >
            <Input
              id="password"
              name="password"
              v-model="userModel.pass"
              placeholder="••••••••"
              type="password"
              autocomplete="current-password"
              aria-required="true"
              :aria-invalid="passwordError ? 'true' : 'false'"
              required
              class="mt-1 block w-full"
            />
            <p
              v-if="passwordError"
              id="password-error"
              class="mt-2 text-sm bg-destructive text-foreground py-0.5 px-2 rounded"
            >
              {{ passwordError }}
            </p>
          </div>

          <Button type="submit" class="mt-4 w-full">Login</Button>
        </form>
      </CardContent>
    </Card>
    <div></div>
  </div>
</template>
