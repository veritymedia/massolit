<template>
  <div>Add book</div>
</template>

<script setup lang="ts">
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});
</script>
