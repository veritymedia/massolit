<script lang="ts" setup>


definePageMeta({
    layout: "empty",
});

</script>

<template>

    <div class="justify-content-center align-content-center display: flex flex-wrap fill-height mt-8">

        <div class="surface-card p-4 shadow-2 border-round md:w-6 lg:w-3">
            <div class="text-center mb-5">
                <img src="/images/logos/pocketnuxt.svg" alt="Image" height="50" class="mb-3">
                <div class="text-900 text-3xl font-medium mb-3">Join Now!</div>
                <span class="text-600 font-medium line-height-3">Already have an account?</span>
                <a class="font-medium no-underline ml-2 text-blue-500 cursor-pointer" href="/auth/login">Back to
                    login</a>
            </div>

            <div>
                <label for="email" class="block text-900 font-medium mb-2">Email</label>
                <InputText id="email" type="text" class="w-full mb-3" v-model="email" />

                <label for="password" class="block text-900 font-medium mb-2">Password</label>
                <InputText id="password" type="password" class="w-full mb-3" v-model="password" />

                <label for="password" class="block text-900 font-medium mb-2">Confirm Password</label>
                <InputText id="confirmPassword" type="password" class="w-full mb-3" v-model="confirmPassword" />

                <Button label="Sign In" icon="pi pi-user" class="w-full" " :loading="isLoading"></Button>
            </div>
        </div>

    </div>

</template>
