<template>
  <div class="flex flex-col items-center">
    <Icon name="material-symbols:check-circle" class="w-12 h-12 opacity-70" />
    <p class="w-full text-center my-5 text-lg font-bold uppercase">
      Book available
    </p>
    <BookCard
      :book="{
        id: props.rentedBookStatus.bookId,
        isbn: props.rentedBookStatus.book?.isbn,
        title: props.rentedBookStatus.book?.title,
        cover_url: props.rentedBookStatus.book?.cover_url,
      }"
    />
  </div>
</template>

<script setup lang="ts">
import type { RentedBookStatus as BookStatus } from "../app.vue";

interface Props {
  rentedBookStatus: BookStatus;
}
const props = withDefaults(defineProps<Props>(), {});
</script>
